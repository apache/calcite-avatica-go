// +build go1.8

package avatica

import (
	"database/sql/driver"
	"fmt"
)

func driverNamedValueToNamedValue(values []driver.NamedValue) ([]namedValue, error) {
	list := make([]namedValue, len(values))

	for i, nv := range values {
		list[i] = namedValue(nv)

		if nv.Name != "" {
			return list, fmt.Errorf("named parameters are not supported: %s given", nv.Name)
		}
	}

	return list, nil
}
