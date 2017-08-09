package avatica

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"fmt"

	avaticaMessage "github.com/Boostport/avatica/message"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/jcmturner/gokrb5/client"
	"github.com/jcmturner/gokrb5/config"
	"github.com/jcmturner/gokrb5/credentials"
	"github.com/jcmturner/gokrb5/keytab"
	"github.com/xinsnake/go-http-digest-auth-client"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type httpClientAuthConfig struct {
	authenticationType authentication

	username string
	password string

	principal           krb5Principal
	keytab              string
	krb5Conf            string
	krb5CredentialCache string
}

// httpClient wraps the default http.Client to communicate with the Avatica server.
type httpClient struct {
	host       string
	authConfig httpClientAuthConfig

	httpClient *http.Client

	kerberosClient client.Client
}

// NewHTTPClient creates a new httpClient from a host.
func NewHTTPClient(host string, authenticationConf httpClientAuthConfig) (*httpClient, error) {

	hc := cleanhttp.DefaultPooledClient()

	c := &httpClient{
		host:       host,
		authConfig: authenticationConf,

		httpClient: hc,
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

			kc, err := client.NewClientFromCCache(tc)

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

			kc := client.NewClientWithKeytab(authenticationConf.principal.username, authenticationConf.principal.realm, kt)
			kc.WithConfig(cfg)

			err = kc.Login()

			if err != nil {
				return nil, fmt.Errorf("error performing kerberos login with keytab: %s", err)
			}

			kc.EnableAutoSessionRenewal()

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

	if c.authConfig.authenticationType == basic {
		req.SetBasicAuth(c.authConfig.username, c.authConfig.password)
	} else if c.authConfig.authenticationType == spnego {
		c.kerberosClient.SetSPNEGOHeader(req, "")
	}

	res, err := ctxhttp.Do(ctx, c.httpClient, req)

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
		return nil, errorResponseToResponseError(v)
	}

	return inner, nil
}
