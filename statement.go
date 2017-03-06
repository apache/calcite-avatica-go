package avatica

import (
	"database/sql/driver"
	"time"

	"github.com/Boostport/avatica/message"
	"golang.org/x/net/context"
)

type stmt struct {
	statementID uint32
	conn        *conn
	parameters  []*message.AvaticaParameter
	handle      message.StatementHandle
}

// Close closes a statement
func (s *stmt) Close() error {

	if s.conn.connectionId == "" {
		return driver.ErrBadConn
	}

	_, err := s.conn.httpClient.post(context.Background(), &message.CloseStatementRequest{
		ConnectionId: s.conn.connectionId,
		StatementId:  s.statementID,
	})

	return err
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (s *stmt) NumInput() int {
	return len(s.parameters)
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	list := driverValueToNamedValue(args)
	return s.exec(context.Background(), list)
}

func (s *stmt) exec(ctx context.Context, args []namedValue) (driver.Result, error) {

	if s.conn.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	res, err := s.conn.httpClient.post(ctx, &message.ExecuteRequest{
		StatementHandle:    &s.handle,
		ParameterValues:    s.parametersToTypedValues(args),
		FirstFrameMaxSize:  uint64(s.conn.config.frameMaxSize), //TODO: Due to CALCITE-1353, if frameMaxSize == -1, it overflows to 18446744073709551615 due to the conversion to uint64, which is basically all rows.
		HasParameterValues: true,
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

// Query executes a query that may return rows, such as a
// SELECT.
func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	list := driverValueToNamedValue(args)
	return s.query(context.Background(), list)
}

func (s *stmt) query(ctx context.Context, args []namedValue) (driver.Rows, error) {
	if s.conn.connectionId == "" {
		return nil, driver.ErrBadConn
	}

	res, err := s.conn.httpClient.post(ctx, &message.ExecuteRequest{
		StatementHandle:    &s.handle,
		ParameterValues:    s.parametersToTypedValues(args),
		FirstFrameMaxSize:  uint64(s.conn.config.frameMaxSize), //TODO: Due to CALCITE-1353, if frameMaxSize == -1, it overflows to 18446744073709551615 due to the conversion to uint64, which is basically all rows.
		HasParameterValues: true,
	})

	if err != nil {
		return nil, err
	}

	// Currently there is only 1 ResultSet per response
	resultSet := res.(*message.ExecuteResponse).Results[0]

	return newRows(s.conn, s.statementID, resultSet), nil
}

func (s *stmt) parametersToTypedValues(vals []namedValue) []*message.TypedValue {

	result := []*message.TypedValue{}

	for i, val := range vals {
		typed := message.TypedValue{}

		if val.Value == nil {
			typed.Null = true
		} else {

			switch v := val.Value.(type) {
			case int64:
				typed.Type = message.Rep_LONG
				typed.NumberValue = v
			case float64:
				typed.Type = message.Rep_DOUBLE
				typed.DoubleValue = v
			case bool:
				typed.Type = message.Rep_BOOLEAN
				typed.BoolValue = v
			case []byte:
				typed.Type = message.Rep_BYTE_STRING
				typed.BytesValue = v
			case string:
				typed.Type = message.Rep_STRING
				typed.StringValue = v
			case time.Time:
				avaticaParameter := s.parameters[i]

				switch avaticaParameter.TypeName {
				case "TIME":
					typed.Type = message.Rep_JAVA_SQL_TIME

					// Because a location can have multiple time zones due to daylight savings,
					// we need to be explicit and get the offset
					zone, offset := v.Zone()

					// Calculate milliseconds since 00:00:00.000
					base := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.FixedZone(zone, offset))
					typed.NumberValue = int64(v.Sub(base).Nanoseconds() / int64(time.Millisecond))

				case "DATE":
					typed.Type = message.Rep_JAVA_SQL_DATE

					// Because a location can have multiple time zones due to daylight savings,
					// we need to be explicit and get the offset
					zone, offset := v.Zone()

					// Calculate number of days since 1970/1/1
					base := time.Date(1970, 1, 1, 0, 0, 0, 0, time.FixedZone(zone, offset))
					typed.NumberValue = int64(v.Sub(base) / (24 * time.Hour))

				case "TIMESTAMP":
					typed.Type = message.Rep_JAVA_SQL_TIMESTAMP

					// Because a location can have multiple time zones due to daylight savings,
					// we need to be explicit and get the offset
					zone, offset := v.Zone()

					// Calculate number of milliseconds since 1970-01-01 00:00:00.000
					base := time.Date(1970, 1, 1, 0, 0, 0, 0, time.FixedZone(zone, offset))
					typed.NumberValue = int64(v.Sub(base).Nanoseconds() / int64(time.Millisecond))
				}
			}
		}

		result = append(result, &typed)
	}

	return result
}
