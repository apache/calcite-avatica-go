// +build go1.8

package avatica

import (
	"io"
	"reflect"
)

func (r *rows) HasNextResultSet() bool {
	lastResultSetID := len(r.resultSets) - 1
	return lastResultSetID > r.currentResultSet
}

func (r *rows) NextResultSet() error {

	lastResultSetID := len(r.resultSets) - 1

	if r.currentResultSet+1 > lastResultSetID {
		return io.EOF
	}

	r.currentResultSet++

	return nil
}

func (r *rows) ColumnTypeDatabaseTypeName(index int) string {

	return r.resultSets[r.currentResultSet].columns[index].typeName
}

func (r *rows) ColumnTypeLength(index int) (length int64, ok bool) {
	l := r.resultSets[r.currentResultSet].columns[index].length

	if l == 0 {
		return 0, false
	}

	return l, true
}

func (r *rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	return r.resultSets[r.currentResultSet].columns[index].nullable, true
}

func (r *rows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {

	ps := r.resultSets[r.currentResultSet].columns[index].precisionScale

	if ps != nil {
		return ps.precision, ps.scale, true
	}

	return 0, 0, false
}

func (r *rows) ColumnTypeScanType(index int) reflect.Type {
	return r.resultSets[r.currentResultSet].columns[index].scanType
}
