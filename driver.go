/*
Package avatica provides an Apache Phoenix Query Server/Avatica driver for Go's database/sql package.

Quickstart

Import the database/sql package along with the avatica driver.

	import "database/sql"
	import _ "github.com/Boostport/avatica"

	db, err := sql.Open("avatica", "http://phoenix-query-server:8765")

See https://github.com/Boostport/avatica#usage for more details
*/
package avatica

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/Boostport/avatica/message"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// Driver is exported to allow it to be used directly.
type Driver struct{}

// Open a Connection to the server.
// See https://github.com/Boostport/avatica#dsn for more information
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

	connectionId := uuid.NewV4().String()

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
		ConnectionId: connectionId,
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
		connectionId: connectionId,
		httpClient:   httpClient,
		config:       config,
	}

	return conn, nil
}

func init() {
	sql.Register("avatica", &Driver{})
}
