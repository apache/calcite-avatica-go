/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package avatica

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"time"

	avaticaMessage "github.com/apache/calcite-avatica-go/v4/message"
	"github.com/golang/protobuf/proto"
	digest_auth_client "github.com/xinsnake/go-http-digest-auth-client"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/credentials"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	gokrbSPNEGO "gopkg.in/jcmturner/gokrb5.v7/spnego"
)

var (
	badConnRe = regexp.MustCompile(`org\.apache\.calcite\.avatica\.NoSuchConnectionException`)
)

type httpClientAuthConfig struct {
	authenticationType authentication

	username string
	password string
	token    string

	principal           krb5Principal
	keytab              string
	krb5Conf            string
	krb5CredentialCache string
}

// httpClient wraps the default http.Client to communicate with the Avatica server.
type httpClient struct {
	host           string
	authConfig     httpClientAuthConfig
	httpClient     *http.Client
	kerberosClient *client.Client
}

type avaticaError struct {
	message *avaticaMessage.ErrorResponse
}

func (e avaticaError) Error() string {
	return fmt.Sprintf("avatica encountered an error: %s", e.message.ErrorMessage)
}

// NewHTTPClient creates a new httpClient from a host.
func NewHTTPClient(host string, authenticationConf httpClientAuthConfig) (*httpClient, error) {

	c := &httpClient{
		host:       host,
		authConfig: authenticationConf,

		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
		},
	}

	if authenticationConf.authenticationType == digest {
		rt := digest_auth_client.NewTransport(authenticationConf.username, authenticationConf.password)
		c.httpClient.Transport = &rt

	} else if authenticationConf.authenticationType == spnego {

		if authenticationConf.krb5CredentialCache != "" {

			tc, err := credentials.LoadCCache(authenticationConf.krb5CredentialCache)

			if err != nil {
				return nil, fmt.Errorf("error reading kerberos ticket cache: %s", err)
			}

			kc, err := client.NewClientFromCCache(tc, config.NewConfig())
			if err != nil {
				return nil, fmt.Errorf("error creating kerberos client: %s", err)
			}

			c.kerberosClient = kc

		} else {

			cfg, err := config.Load(authenticationConf.krb5Conf)

			if err != nil {
				return nil, fmt.Errorf("error reading kerberos config: %s", err)
			}

			kt, err := keytab.Load(authenticationConf.keytab)

			if err != nil {
				return nil, fmt.Errorf("error reading kerberos keytab: %s", err)
			}

			kc := client.NewClientWithKeytab(authenticationConf.principal.username, authenticationConf.principal.realm, kt, cfg)

			err = kc.Login()

			if err != nil {
				return nil, fmt.Errorf("error performing kerberos login with keytab: %s", err)
			}

			c.kerberosClient = kc
		}
	}

	return c, nil
}

// post posts a protocol buffer message to the Avatica server.
func (c *httpClient) post(ctx context.Context, message proto.Message) (proto.Message, error) {

	wrapped, err := proto.Marshal(message)

	if err != nil {
		return nil, err
	}

	wire := &avaticaMessage.WireMessage{
		Name:           classNameFromRequest(message),
		WrappedMessage: wrapped,
	}

	body, err := proto.Marshal(wire)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.host, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-google-protobuf")

	switch c.authConfig.authenticationType {
	case basic:
		req.SetBasicAuth(c.authConfig.username, c.authConfig.password)
	case spnego:
		if err := gokrbSPNEGO.SetSPNEGOHeader(c.kerberosClient, req, ""); err != nil {
			return nil, err
		}
	case token:
		req.Header.Set("Authorization", "Bearer "+c.authConfig.token)
	}

	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	result := &avaticaMessage.WireMessage{}

	err = proto.Unmarshal(response, result)

	if err != nil {
		return nil, err
	}

	inner, err := responseFromClassName(result.Name)

	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(result.WrappedMessage, inner)

	if err != nil {
		return nil, err
	}

	if v, ok := inner.(*avaticaMessage.ErrorResponse); ok {

		for _, exception := range v.Exceptions {
			if badConnRe.MatchString(exception) {
				return nil, driver.ErrBadConn
			}
		}

		return nil, avaticaError{v}
	}

	return inner, nil
}
