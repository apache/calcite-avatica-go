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
	"database/sql/driver"
	"io"
	"time"

	"reflect"

	"math"

	"fmt"

	"github.com/apache/calcite-avatica-go/message"
	"golang.org/x/net/context"
)

type precisionScale struct {
	precision int64
	scale     int64
}

type column struct {
	name           string
	typeName       string
	rep            message.Rep
	length         int64
	nullable       bool
	precisionScale *precisionScale
	scanType       reflect.Type
}

type resultSet struct {
	columns    []*column
	done       bool
	offset     uint64
	data       [][]*message.TypedValue
	currentRow int
}

type rows struct {
	conn             *conn
	statementID      uint32
	resultSets       []*resultSet
	currentResultSet int
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice.  If a particular column name isn't known, an empty
// string should be returned for that entry.
func (r *rows) Columns() []string {

	cols := []string{}

	for _, column := range r.resultSets[r.currentResultSet].columns {
		cols = append(cols, column.name)
	}

	return cols
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

	resultSet := r.resultSets[r.currentResultSet]

	if resultSet.currentRow >= len(resultSet.data) {

		if resultSet.done {
			// Finished iterating through all results
			return io.EOF
		}

		// Fetch more results from the server
		res, err := r.conn.httpClient.post(context.Background(), &message.FetchRequest{
			ConnectionId: r.conn.connectionId,
			StatementId:  r.statementID,
			Offset:       resultSet.offset,
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

		resultSet.done = frame.Done
		resultSet.data = data
		resultSet.currentRow = 0

	}

	for i, val := range resultSet.data[resultSet.currentRow] {
		dest[i] = typedValueToNative(resultSet.columns[i].rep, val, r.conn.config)
	}

	resultSet.currentRow++

	return nil
}

// newRows create a new set of rows from a result set.
func newRows(conn *conn, statementID uint32, resultSets []*message.ResultSetResponse) *rows {

	rsets := []*resultSet{}

	for _, result := range resultSets {
		if result.Signature == nil {
			break
		}
		columns := []*column{}

		for _, col := range result.Signature.Columns {

			column := &column{
				name:     col.ColumnName,
				typeName: col.Type.Name,
				nullable: col.Nullable != 0,
			}

			// Handle precision and length
			switch col.Type.Name {
			case "DECIMAL":

				precision := int64(col.Precision)

				if precision == 0 {
					precision = math.MaxInt64
				}

				scale := int64(col.Scale)

				if scale == 0 {
					scale = math.MaxInt64
				}

				column.precisionScale = &precisionScale{
					precision: precision,
					scale:     scale,
				}
			case "VARCHAR", "CHAR", "BINARY":
				column.length = int64(col.Precision)
			case "VARBINARY":
				column.length = math.MaxInt64
			}

			// Handle scan types
			switch col.Type.Name {
			case "INTEGER", "UNSIGNED_INT", "BIGINT", "UNSIGNED_LONG", "TINYINT", "UNSIGNED_TINYINT", "SMALLINT", "UNSIGNED_SMALLINT":
				column.scanType = reflect.TypeOf(int64(0))

			case "FLOAT", "UNSIGNED_FLOAT", "DOUBLE", "UNSIGNED_DOUBLE":
				column.scanType = reflect.TypeOf(float64(0))

			case "DECIMAL", "VARCHAR", "CHAR":
				column.scanType = reflect.TypeOf("")

			case "BOOLEAN":
				column.scanType = reflect.TypeOf(bool(false))

			case "TIME", "DATE", "TIMESTAMP", "UNSIGNED_TIME", "UNSIGNED_DATE", "UNSIGNED_TIMESTAMP":
				column.scanType = reflect.TypeOf(time.Time{})

			case "BINARY", "VARBINARY":
				column.scanType = reflect.TypeOf([]byte{})

			default:
				panic(fmt.Sprintf("scantype for %s is not implemented", col.Type.Name))
			}

			// Handle rep type special cases for decimals, floats, date, time and timestamp
			switch col.Type.Name {
			case "DECIMAL":
				column.rep = message.Rep_BIG_DECIMAL
			case "FLOAT":
				column.rep = message.Rep_FLOAT
			case "UNSIGNED_FLOAT":
				column.rep = message.Rep_FLOAT
			case "TIME", "UNSIGNED_TIME":
				column.rep = message.Rep_JAVA_SQL_TIME
			case "DATE", "UNSIGNED_DATE":
				column.rep = message.Rep_JAVA_SQL_DATE
			case "TIMESTAMP", "UNSIGNED_TIMESTAMP":
				column.rep = message.Rep_JAVA_SQL_TIMESTAMP
			default:
				column.rep = col.Type.Rep
			}

			columns = append(columns, column)
		}

		frame := result.FirstFrame

		data := [][]*message.TypedValue{}

		for _, row := range frame.Rows {
			rowData := []*message.TypedValue{}

			for _, col := range row.Value {
				rowData = append(rowData, col.ScalarValue)
			}

			data = append(data, rowData)
		}

		rsets = append(rsets, &resultSet{
			columns: columns,
			done:    frame.Done,
			offset:  frame.Offset,
			data:    data,
		})
	}

	return &rows{
		conn:             conn,
		statementID:      statementID,
		resultSets:       rsets,
		currentResultSet: 0,
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
