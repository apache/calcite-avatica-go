package avatica

import (
	"database/sql/driver"
	"github.com/Boostport/avatica/message"
	"golang.org/x/net/context"
	"io"
	"time"
)

type rows struct {
	conn        *conn
	statementID uint32
	columnNames []string
	columnTypes []message.Rep
	done        bool
	offset      uint64
	data        [][]*message.TypedValue
	currentRow  int
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice.  If a particular column name isn't known, an empty
// string should be returned for that entry.
func (r *rows) Columns() []string {
	return r.columnNames
}

// Close closes the rows iterator.
func (r *rows) Close() error {

	r.conn = nil
	return nil
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// The dest slice may be populated only with
// a driver Value type, but excluding string.
// All string values must be converted to []byte.
//
// Next should return io.EOF when there are no more rows.
func (r *rows) Next(dest []driver.Value) error {

	if r.currentRow >= len(r.data) {

		if r.done {
			// Finished iterating through all results
			return io.EOF
		}

		// Fetch more results from the server
		res, err := r.conn.httpClient.post(context.Background(), &message.FetchRequest{
			ConnectionId: r.conn.connectionId,
			StatementId:  r.statementID,
			Offset:       r.offset,
			FrameMaxSize: r.conn.config.frameMaxSize,
		})

		if err != nil {
			return err
		}

		frame := res.(*message.FetchResponse).Frame

		data := [][]*message.TypedValue{}

		// In some cases the server does not return done as true
		// until it returns a result with no rows
		if len(frame.Rows) == 0 {
			return io.EOF
		}

		for _, row := range frame.Rows {
			rowData := []*message.TypedValue{}

			for _, col := range row.Value {
				rowData = append(rowData, col.ScalarValue)
			}

			data = append(data, rowData)
		}

		r.done = frame.Done
		r.data = data
		r.currentRow = 0

	}

	for i, val := range r.data[r.currentRow] {
		dest[i] = typedValueToNative(r.columnTypes[i], val, r.conn.config)
	}

	r.currentRow++

	return nil
}

// newRows create a new set of rows from a result set.
func newRows(conn *conn, statementID uint32, resultSet *message.ResultSetResponse) *rows {

	columnNames := []string{}
	columnTypes := []message.Rep{}

	for _, col := range resultSet.Signature.Columns {
		columnNames = append(columnNames, col.ColumnName)

		// Special case for floats, date, time and timestamp
		switch col.Type.Name {
		case "FLOAT":
			columnTypes = append(columnTypes, message.Rep_FLOAT)
		case "UNSIGNED_FLOAT":
			columnTypes = append(columnTypes, message.Rep_FLOAT)
		case "TIME":
			columnTypes = append(columnTypes, message.Rep_JAVA_SQL_TIME)
		case "DATE":
			columnTypes = append(columnTypes, message.Rep_JAVA_SQL_DATE)
		case "TIMESTAMP":
			columnTypes = append(columnTypes, message.Rep_JAVA_SQL_TIMESTAMP)
		default:
			columnTypes = append(columnTypes, col.Type.Rep)
		}
	}

	frame := resultSet.FirstFrame

	data := [][]*message.TypedValue{}

	for _, row := range frame.Rows {
		rowData := []*message.TypedValue{}

		for _, col := range row.Value {
			rowData = append(rowData, col.ScalarValue)
		}

		data = append(data, rowData)
	}

	return &rows{
		conn:        conn,
		statementID: statementID,
		columnNames: columnNames,
		columnTypes: columnTypes,
		done:        frame.Done,
		offset:      frame.Offset,
		data:        data,
	}
}

// typedValueToNative converts values from avatica's types to Go's native types
func typedValueToNative(rep message.Rep, v *message.TypedValue, config *Config) interface{} {

	switch rep {
	case message.Rep_BOOLEAN, message.Rep_PRIMITIVE_BOOLEAN:
		return v.BoolValue

	case message.Rep_STRING, message.Rep_PRIMITIVE_CHAR, message.Rep_CHARACTER, message.Rep_BIG_DECIMAL:
		return v.StringValue

	case message.Rep_FLOAT, message.Rep_PRIMITIVE_FLOAT:
		return float32(v.DoubleValue)

	case message.Rep_LONG,
		message.Rep_PRIMITIVE_LONG,
		message.Rep_INTEGER,
		message.Rep_PRIMITIVE_INT,
		message.Rep_BIG_INTEGER,
		message.Rep_NUMBER,
		message.Rep_BYTE,
		message.Rep_PRIMITIVE_BYTE,
		message.Rep_SHORT,
		message.Rep_PRIMITIVE_SHORT:
		return v.NumberValue

	case message.Rep_BYTE_STRING:
		return v.BytesValue

	case message.Rep_DOUBLE, message.Rep_PRIMITIVE_DOUBLE:
		return v.DoubleValue

	case message.Rep_JAVA_SQL_DATE, message.Rep_JAVA_UTIL_DATE:

		// We receive the number of days since 1970/1/1 from the server
		// Because a location can have multiple time zones due to daylight savings,
		// we first do all our calculations in UTC and then force the timezone to
		// the one the user has chosen.
		t, _ := time.ParseInLocation("2006-Jan-02", "1970-Jan-01", time.UTC)
		days := time.Hour * 24 * time.Duration(v.NumberValue)
		t = t.Add(days)

		return forceTimezone(t, config.location)

	case message.Rep_JAVA_SQL_TIME:

		// We receive the number of milliseconds since 00:00:00.000 from the server
		// Because a location can have multiple time zones due to daylight savings,
		// we first do all our calculations in UTC and then force the timezone to
		// the one the user has chosen.
		t, _ := time.ParseInLocation("15:04:05", "00:00:00", time.UTC)
		ms := time.Millisecond * time.Duration(v.NumberValue)
		t = t.Add(ms)
		return forceTimezone(t, config.location)

	case message.Rep_JAVA_SQL_TIMESTAMP:

		// We receive the number of milliseconds since 1970-01-01 00:00:00.000 from the server
		// Force to UTC for consistency because time.Unix uses the local timezone
		t := time.Unix(0, v.NumberValue*int64(time.Millisecond)).In(time.UTC)
		return forceTimezone(t, config.location)

	default:
		return nil
	}
}

// forceTimezone takes a time.Time and changes its location without shifting the timezone.
func forceTimezone(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}
