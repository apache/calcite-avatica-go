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

/*
Package avatica provides an Apache Phoenix Query Server/Avatica driver for Go's database/sql package.

Quickstart

Import the database/sql package along with the avatica driver.

	import "database/sql"
	import _ "github.com/apache/calcite-avatica-go"

	db, err := sql.Open("avatica", "http://phoenix-query-server:8765")

See https://github.com/apache/calcite-avatica-go#usage for more details
*/
package avatica

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/apache/calcite-avatica-go/message"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// Driver is exported to allow it to be used directly.
type Driver struct{}

// Open a Connection to the server.
// See https://github.com/apache/calcite-avatica-go#dsn for more information
// on how the DSN is formatted.
func (a *Driver) Open(dsn string) (driver.Conn, error) {

	config, err := ParseDSN(dsn)

	if err != nil {
		return nil, fmt.Errorf("Unable to open connection: %s", err)
	}

	httpClient, err := NewHTTPClient(config.endpoint, httpClientAuthConfig{
		authenticationType:  config.authentication,
		username:            config.avaticaUser,
		password:            config.avaticaPassword,
		principal:           config.principal,
		keytab:              config.keytab,
		krb5Conf:            config.krb5Conf,
		krb5CredentialCache: config.krb5CredentialCache,
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to create HTTP client: %s", err)
	}

	connectionId, err := uuid.NewV4()

	if err != nil {
		return nil, fmt.Errorf("Error generating connection id: %s", err)
	}

	info := map[string]string{
		"AutoCommit":  "true",
		"Consistency": "8",
	}

	if config.user != "" {
		info["user"] = config.user
	}

	if config.password != "" {
		info["password"] = config.password
	}

	// Open a connection to the server
	req := &message.OpenConnectionRequest{
		ConnectionId: connectionId.String(),
		Info:         info,
	}

	if config.schema != "" {
		req.Info["schema"] = config.schema
	}

	_, err = httpClient.post(context.Background(), req)

	if err != nil {
		return nil, err
	}

	conn := &conn{
		connectionId: connectionId.String(),
		httpClient:   httpClient,
		config:       config,
	}

	return conn, nil
}

func init() {
	sql.Register("avatica", &Driver{})
}
