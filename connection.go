package avatica

import (
	"database/sql/driver"
	"github.com/Boostport/avatica/message"
	"golang.org/x/net/context"
)

type conn struct {
	connectionId string
	config       *Config
	httpClient   *httpClient
}

// Prepare returns a prepared statement, bound to this connection.
func (c *conn) Prepare(query string) (driver.Stmt, error) {

	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	response, err := c.httpClient.post(context.Background(), &message.PrepareRequest{
		ConnectionId: c.connectionId,
		Sql:          query,
		MaxRowCount:  c.config.maxRowCount,
	})

	if err != nil {
		return nil, err
	}

	prepareResponse := response.(*message.PrepareResponse)

	return &stmt{
		statementID: prepareResponse.Statement.Id,
		conn:        c,
		parameters:  prepareResponse.Statement.Signature.Parameters,
		handle:      *prepareResponse.Statement,
	}, nil
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
func (c *conn) Close() error {

	if c.connectionId == "" {
		return driver.ErrBadConn
	}

	_, err := c.httpClient.post(context.Background(), &message.CloseConnectionRequest{
		ConnectionId: c.connectionId,
	})

	c.connectionId = ""

	return err
}

// Begin starts and returns a new transaction.
func (c *conn) Begin() (driver.Tx, error) {

	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	_, err := c.httpClient.post(context.Background(), &message.ConnectionSyncRequest{
		ConnectionId: c.connectionId,
		ConnProps: &message.ConnectionProperties{
			AutoCommit:           false,
			HasAutoCommit:        true,
			TransactionIsolation: 8,
		},
	})

	if err != nil {
		return nil, err
	}

	return &tx{
		conn: c,
	}, nil
}

// Exec prepares and executes a query and returns the result directly.
func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {

	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip
	}

	st, err := c.httpClient.post(context.Background(), &message.CreateStatementRequest{
		ConnectionId: c.connectionId,
	})

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.post(context.Background(), &message.PrepareAndExecuteRequest{
		ConnectionId: c.connectionId,
		StatementId:  st.(*message.CreateStatementResponse).StatementId,
		Sql:          query,
		MaxRowCount:  c.config.maxRowCount,
	})

	if err != nil {
		return nil, err
	}

	// Currently there is only 1 ResultSet per response
	changed := int64(res.(*message.ExecuteResponse).Results[0].UpdateCount)

	return &result{
		affectedRows: changed,
	}, nil

}

// Query prepares and executes a query and returns the result directly.
// Query's optimizations are currently disabled due to CALCITE-1181.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {

	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	return nil, driver.ErrSkip

	// Disabled due to CALCITE-1181

	/*if len(args) != 0 {
		return nil, driver.ErrSkip
	}

	st, err := c.httpClient.post(context.Background(), &message.CreateStatementRequest{
		ConnectionId: c.connectionId,
	})

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.post(context.Background(), &message.PrepareAndExecuteRequest{
		ConnectionId: c.connectionId,
		StatementId:  st.(*message.CreateStatementResponse).StatementId,
		Sql:          query,
		MaxRowCount:  maxRowCount,
	})

	if err != nil {
		return nil, err
	}

	// Currently there is only 1 ResultSet per response
	resultSet := res.(*message.ExecuteResponse).Results[0]

	return NewRows(c, st.(*message.CreateStatementResponse).StatementId, resultSet), nil*/
}
