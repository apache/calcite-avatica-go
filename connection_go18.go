// +build go1.8

package avatica

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"context"
)

func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {

	if opts.ReadOnly {
		return nil, errors.New("Read-only transactions are not supported")
	}

	var isolation isoLevel

	switch sql.IsolationLevel(opts.Isolation) {
	case sql.LevelDefault:
		isolation = isolationUseCurrent
	case sql.LevelReadUncommitted:
		isolation = isolationReadUncommitted
	case sql.LevelReadCommitted:
		isolation = isolationReadComitted
	case sql.LevelWriteCommitted:
		return nil, errors.New("LevelWriteCommitted isolation level is not supported")
	case sql.LevelRepeatableRead:
		isolation = isolationRepeatableRead
	case sql.LevelSnapshot:
		return nil, errors.New("LevelSnapshot isolation level is not supported")
	case sql.LevelSerializable:
		isolation = isolationSerializable
	case sql.LevelLinearizable:
		return nil, errors.New("LevelLinearizable isolation level is not supported")
	default:
		return nil, fmt.Errorf("Unsupported transaction isolation level: %d", opts.Isolation)
	}

	return c.begin(ctx, isolation)
}

func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return c.prepare(ctx, query)
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	list, err := driverNamedValueToNamedValue(args)

	if err != nil {
		return nil, fmt.Errorf("could not execute statement: %s", err)
	}

	return c.exec(ctx, query, list)
}

func (c *conn) Ping(ctx context.Context) error {

	_, err := c.ExecContext(ctx, "SELECT 1", []driver.NamedValue{})

	if err != nil {
		return fmt.Errorf("Error pinging database: %s", err)
	}

	return nil
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	list, err := driverNamedValueToNamedValue(args)

	if err != nil {
		return nil, fmt.Errorf("could not execute query: %s", err)
	}

	return c.query(ctx, query, list)
}
