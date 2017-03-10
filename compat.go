package avatica

import (
	"database/sql/driver"
)

type namedValue struct {
	Name    string
	Ordinal int
	Value   driver.Value
}

func driverValueToNamedValue(values []driver.Value) []namedValue {
	list := make([]namedValue, len(values))

	for i, v := range values {
		list[i] = namedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}

	return list
}

type isoLevel int32

const (
	isolationUseCurrent      isoLevel = -1
	isolationNone            isoLevel = 0
	isolationReadUncommitted isoLevel = 1
	isolationReadComitted    isoLevel = 2
	isolationRepeatableRead  isoLevel = 4
	isolationSerializable    isoLevel = 8
)
