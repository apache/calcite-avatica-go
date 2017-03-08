// +build go1.8

package avatica

import (
	"database/sql/driver"
	"fmt"
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

func driverNamedValueToNamedValue(values []driver.NamedValue) ([]namedValue,error ) {
	list := make([]namedValue, len(values))

	for i, nv := range values {
		list[i] = namedValue(nv)

		if nv.Name != ""{
			return list,fmt.Errorf("named paramters are not supported: %s given", nv.Name)
		}
	}

	return list, nil
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
