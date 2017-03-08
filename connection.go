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
	return c.prepare(context.Background(), query)
}

func (c *conn) prepare(ctx context.Context, query string) (driver.Stmt, error) {
	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	response, err := c.httpClient.post(ctx, &message.PrepareRequest{
		ConnectionId: c.connectionId,
		Sql:          query,
		MaxRowsTotal: c.config.maxRowsTotal,
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
	return c.begin(context.Background(), isolationUseCurrent)
}

func (c *conn) begin(ctx context.Context, isolationLevel isoLevel) (driver.Tx, error) {
	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	if isolationLevel == isolationUseCurrent {
		isolationLevel = isoLevel(c.config.transactionIsolation)
	}

	_, err := c.httpClient.post(ctx, &message.ConnectionSyncRequest{
		ConnectionId: c.connectionId,
		ConnProps: &message.ConnectionProperties{
			AutoCommit:           false,
			HasAutoCommit:        true,
			TransactionIsolation: uint32(isolationLevel),
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
	list := driverValueToNamedValue(args)
	return c.exec(context.Background(), query, list)
}

func (c *conn) exec(ctx context.Context, query string, args []namedValue) (driver.Result, error) {
	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip
	}

	st, err := c.httpClient.post(ctx, &message.CreateStatementRequest{
		ConnectionId: c.connectionId,
	})

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.post(ctx, &message.PrepareAndExecuteRequest{
		ConnectionId:      c.connectionId,
		StatementId:       st.(*message.CreateStatementResponse).StatementId,
		Sql:               query,
		MaxRowsTotal:      c.config.maxRowsTotal,
		FirstFrameMaxSize: c.config.frameMaxSize,
	})

	if err != nil {
		return nil, err
	}

	// Currently there is only 1 ResultSet per response for exec
	changed := int64(res.(*message.ExecuteResponse).Results[0].UpdateCount)

	return &result{
		affectedRows: changed,
	}, nil
}

// Query prepares and executes a query and returns the result directly.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	list := driverValueToNamedValue(args)
	return c.query(context.Background(), query, list)
}

func (c *conn) query(ctx context.Context, query string, args []namedValue) (driver.Rows, error) {
	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip
	}

	st, err := c.httpClient.post(ctx, &message.CreateStatementRequest{
		ConnectionId: c.connectionId,
	})

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.post(ctx, &message.PrepareAndExecuteRequest{
		ConnectionId:      c.connectionId,
		StatementId:       st.(*message.CreateStatementResponse).StatementId,
		Sql:               query,
		MaxRowsTotal:      c.config.maxRowsTotal,
		FirstFrameMaxSize: c.config.frameMaxSize,
	})

	if err != nil {
		return nil, err
	}

	resultSets := res.(*message.ExecuteResponse).Results

	return newRows(c, st.(*message.CreateStatementResponse).StatementId, resultSets), nil
}
