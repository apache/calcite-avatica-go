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
	"context"
	"database/sql/driver"
	"errors"

	avaticaErrors "github.com/apache/calcite-avatica-go/v5/errors"
	"github.com/apache/calcite-avatica-go/v5/message"
)

type conn struct {
	connectionId  string
	config        *Config
	httpClient    *httpClient
	adapter       Adapter
	connectorInfo map[string]string
}

// Prepare returns a prepared statement, bound to this connection.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return c.prepare(context.Background(), query)
}

func (c *conn) prepare(ctx context.Context, query string) (driver.Stmt, error) {
	if c.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	response, err := c.httpClient.post(ctx, message.PrepareRequest_builder{
		ConnectionId: c.connectionId,
		Sql:          query,
		MaxRowsTotal: c.config.maxRowsTotal,
	}.Build())

	if err != nil {
		return nil, c.avaticaErrorToResponseErrorOrError(err)
	}

	prepareResponse := response.(*message.PrepareResponse)

	return &stmt{
		statementID:  prepareResponse.GetStatement().GetId(),
		conn:         c,
		parameters:   prepareResponse.GetStatement().GetSignature().GetParameters(),
		handle:       prepareResponse.GetStatement(),
		batchUpdates: make([]*message.UpdateBatch, 0),
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

	_, err := c.httpClient.post(context.Background(), message.CloseConnectionRequest_builder{
		ConnectionId: c.connectionId,
	}.Build())

	c.connectionId = ""

	if err != nil {
		return c.avaticaErrorToResponseErrorOrError(err)
	}

	return nil
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

	_, err := c.httpClient.post(ctx, message.ConnectionSyncRequest_builder{
		ConnectionId: c.connectionId,
		ConnProps: message.ConnectionProperties_builder{
			AutoCommit:           false,
			HasAutoCommit:        true,
			TransactionIsolation: uint32(isolationLevel),
		}.Build(),
	}.Build())

	if err != nil {
		return nil, c.avaticaErrorToResponseErrorOrError(err)
	}

	return &tx{conn: c}, nil
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

	st, err := c.httpClient.post(ctx, message.CreateStatementRequest_builder{
		ConnectionId: c.connectionId,
	}.Build())

	if err != nil {
		return nil, c.avaticaErrorToResponseErrorOrError(err)
	}

	statementID := st.(*message.CreateStatementResponse).GetStatementId()
	defer c.closeStatement(ctx, statementID)

	res, err := c.httpClient.post(ctx, message.PrepareAndExecuteRequest_builder{
		ConnectionId:      c.connectionId,
		StatementId:       statementID,
		Sql:               query,
		MaxRowsTotal:      c.config.maxRowsTotal,
		FirstFrameMaxSize: c.config.frameMaxSize,
	}.Build())

	if err != nil {
		return nil, c.avaticaErrorToResponseErrorOrError(err)
	}

	// Currently there is only 1 ResultSet per response for exec
	changed := int64(res.(*message.ExecuteResponse).GetResults()[0].GetUpdateCount())

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

	st, err := c.httpClient.post(ctx, message.CreateStatementRequest_builder{
		ConnectionId: c.connectionId,
	}.Build())

	if err != nil {
		return nil, c.avaticaErrorToResponseErrorOrError(err)
	}

	statementID := st.(*message.CreateStatementResponse).GetStatementId()

	res, err := c.httpClient.post(ctx, message.PrepareAndExecuteRequest_builder{
		ConnectionId:      c.connectionId,
		StatementId:       statementID,
		Sql:               query,
		MaxRowsTotal:      c.config.maxRowsTotal,
		FirstFrameMaxSize: c.config.frameMaxSize,
	}.Build())

	if err != nil {
		_ = c.closeStatement(ctx, statementID)
		return nil, c.avaticaErrorToResponseErrorOrError(err)
	}

	resultSets := res.(*message.ExecuteResponse).GetResults()

	return newRows(c, statementID, true, resultSets), nil
}

func (c *conn) avaticaErrorToResponseErrorOrError(err error) error {

	var avaticaErr avaticaError

	ok := errors.As(err, &avaticaErr)

	if !ok {
		return err
	}

	if c.adapter != nil {
		return c.adapter.ErrorResponseToResponseError(avaticaErr.message)
	}

	return avaticaErrors.ResponseError{
		Exceptions:   avaticaErr.message.GetExceptions(),
		ErrorMessage: avaticaErr.message.GetErrorMessage(),
		Severity:     int8(avaticaErr.message.GetSeverity()),
		ErrorCode:    avaticaErrors.ErrorCode(avaticaErr.message.GetErrorCode()),
		SqlState:     avaticaErrors.SQLState(avaticaErr.message.GetSqlState()),
		Metadata: &avaticaErrors.RPCMetadata{
			ServerAddress: message.ServerAddressFromMetadata(avaticaErr.message),
		},
	}
}

// ResetSession implements driver.SessionResetter.
// (From Go 1.10)
func (c *conn) ResetSession(_ context.Context) error {
	if c.connectionId == "" {
		return driver.ErrBadConn
	}
	return nil
}

func (c *conn) closeStatement(ctx context.Context, statementID uint32) error {
	_, err := c.httpClient.post(ctx, message.CloseStatementRequest_builder{
		ConnectionId: c.connectionId,
		StatementId:  statementID,
	}.Build())
	return c.avaticaErrorToResponseErrorOrError(err)
}
