// +build go1.8
package avatica

import (
	"database/sql/driver"

	"fmt"

	"golang.org/x/net/context"
)

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {

	list, err := driverNamedValueToNamedValue(args)

	if err != nil {
		return nil, fmt.Errorf("Error executing statement: %s", err)
	}

	return s.exec(ctx, list)
}

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {

	list, err := driverNamedValueToNamedValue(args)

	if err != nil {
		return nil, fmt.Errorf("Error executing query: %s", err)
	}

	return s.query(ctx, list)
}
