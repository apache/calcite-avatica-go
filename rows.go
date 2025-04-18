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
	"io"
	"time"

	"github.com/apache/calcite-avatica-go/v5/internal"
	"github.com/apache/calcite-avatica-go/v5/message"
)

type resultSet struct {
	columns    []*internal.Column
	done       bool
	offset     uint64
	data       [][]*message.TypedValue
	currentRow int
}

type rows struct {
	conn             *conn
	statementID      uint32
	closeStatement   bool
	resultSets       []*resultSet
	currentResultSet int
	columnNames      []string
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice.  If a particular column name isn't known, an empty
// string should be returned for that entry.
func (r *rows) Columns() []string {
	if r.columnNames != nil {
		return r.columnNames
	}
	var cols []string
	if len(r.resultSets) == 0 {
		return cols
	}
	for _, column := range r.resultSets[r.currentResultSet].columns {
		cols = append(cols, column.Name)
	}
	r.columnNames = cols
	return cols
}

// Close closes the rows iterator.
func (r *rows) Close() error {
	var err error
	if r.closeStatement {
		err = r.conn.closeStatement(context.Background(), r.statementID)
	}
	r.conn = nil
	return err
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
	if len(r.resultSets) == 0 {
		return io.EOF
	}
	resultSet := r.resultSets[r.currentResultSet]

	if resultSet.currentRow >= len(resultSet.data) {

		if resultSet.done {
			// Finished iterating through all results
			return io.EOF
		}

		// Fetch more results from the server
		res, err := r.conn.httpClient.post(context.Background(), message.FetchRequest_builder{
			ConnectionId: r.conn.connectionId,
			StatementId:  r.statementID,
			Offset:       resultSet.offset + uint64(len(resultSet.data)),
			FrameMaxSize: r.conn.config.frameMaxSize,
		}.Build())

		if err != nil {
			return r.conn.avaticaErrorToResponseErrorOrError(err)
		}

		frame := res.(*message.FetchResponse).GetFrame()

		var data [][]*message.TypedValue

		// In some cases the server does not return done as true
		// until it returns a result with no rows
		if len(frame.GetRows()) == 0 {
			return io.EOF
		}

		for _, row := range frame.GetRows() {
			var rowData []*message.TypedValue

			for _, col := range row.GetValue() {
				rowData = append(rowData, col.GetScalarValue())
			}

			data = append(data, rowData)
		}

		resultSet.done = frame.GetDone()
		resultSet.data = data
		resultSet.currentRow = 0
		resultSet.offset = frame.GetOffset()

	}

	for i, val := range resultSet.data[resultSet.currentRow] {
		dest[i] = typedValueToNative(resultSet.columns[i].Rep, val, r.conn.config)
	}

	resultSet.currentRow++

	return nil
}

// newRows create a new set of rows from a result set.
func newRows(conn *conn, statementID uint32, closeStatement bool, resultSets []*message.ResultSetResponse) *rows {

	var rsets []*resultSet

	for _, result := range resultSets {
		if result.GetSignature() == nil {
			break
		}

		var columns []*internal.Column

		for _, col := range result.GetSignature().GetColumns() {
			column := conn.adapter.GetColumnTypeDefinition(col)
			columns = append(columns, column)
		}

		frame := result.GetFirstFrame()

		var data [][]*message.TypedValue

		for _, row := range frame.GetRows() {
			var rowData []*message.TypedValue

			for _, col := range row.GetValue() {
				rowData = append(rowData, col.GetScalarValue())
			}

			data = append(data, rowData)
		}

		rsets = append(rsets, &resultSet{
			columns: columns,
			done:    frame.GetDone(),
			offset:  frame.GetOffset(),
			data:    data,
		})
	}

	return &rows{
		conn:             conn,
		statementID:      statementID,
		closeStatement:   closeStatement,
		resultSets:       rsets,
		currentResultSet: 0,
	}
}

// typedValueToNative converts values from avatica's types to Go's native types
func typedValueToNative(rep message.Rep, v *message.TypedValue, config *Config) interface{} {

	if v.GetType() == message.Rep_NULL {
		return nil
	}

	switch rep {

	case message.Rep_BOOLEAN, message.Rep_PRIMITIVE_BOOLEAN:
		return v.GetBoolValue()

	case message.Rep_STRING, message.Rep_PRIMITIVE_CHAR, message.Rep_CHARACTER, message.Rep_BIG_DECIMAL:
		return v.GetStringValue()

	case message.Rep_FLOAT, message.Rep_PRIMITIVE_FLOAT:
		return float32(v.GetDoubleValue())

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
		return v.GetNumberValue()

	case message.Rep_BYTE_STRING:
		return v.GetBytesValue()

	case message.Rep_DOUBLE, message.Rep_PRIMITIVE_DOUBLE:
		return v.GetDoubleValue()

	case message.Rep_JAVA_SQL_DATE, message.Rep_JAVA_UTIL_DATE:

		// We receive the number of days since 1970/1/1 from the server
		// Because a location can have multiple time zones due to daylight savings,
		// we first do all our calculations in UTC and then force the timezone to
		// the one the user has chosen.
		t, _ := time.ParseInLocation("2006-Jan-02", "1970-Jan-01", time.UTC)
		days := time.Hour * 24 * time.Duration(v.GetNumberValue())
		t = t.Add(days)

		return forceTimezone(t, config.location)

	case message.Rep_JAVA_SQL_TIME:

		// We receive the number of milliseconds since 00:00:00.000 from the server
		// Because a location can have multiple time zones due to daylight savings,
		// we first do all our calculations in UTC and then force the timezone to
		// the one the user has chosen.
		t, _ := time.ParseInLocation("15:04:05", "00:00:00", time.UTC)
		ms := time.Millisecond * time.Duration(v.GetNumberValue())
		t = t.Add(ms)
		return forceTimezone(t, config.location)

	case message.Rep_JAVA_SQL_TIMESTAMP:

		// We receive the number of milliseconds since 1970-01-01 00:00:00.000 from the server
		// Force to UTC for consistency because time.Unix uses the local timezone
		t := time.Unix(0, v.GetNumberValue()*int64(time.Millisecond)).In(time.UTC)
		return forceTimezone(t, config.location)

	default:
		return nil
	}
}

// forceTimezone takes a time.Time and changes its location without shifting the timezone.
func forceTimezone(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}
