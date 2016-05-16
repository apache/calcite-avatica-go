package avatica

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Config is a configuration parsed from a DSN string
type Config struct {
	endpoint         string
	maxRowCount      uint64
	fetchMaxRowCount uint32
	location         *time.Location
}

// ParseDSN parses a DSN string to a Config
func ParseDSN(dsn string) (*Config, error) {

	conf := &Config{
		maxRowCount:      100,
		fetchMaxRowCount: 100,
		location:         time.UTC,
	}

	parsed, err := url.ParseRequestURI(dsn)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse DSN: %s", err)
	}

	queries := parsed.Query()

	if v := queries.Get("maxRowCount"); v != "" {

		maxRowCount, err := strconv.Atoi(v)

		if err != nil {
			return nil, fmt.Errorf("Invalid value for maxRowCount: %s", err)
		}

		if maxRowCount <= 0 {
			return nil, fmt.Errorf("maxRowCount must be greater than 0")
		}

		conf.fetchMaxRowCount = uint32(maxRowCount)
		conf.maxRowCount = uint64(maxRowCount)
	}

	if v := queries.Get("location"); v != "" {

		loc, err := time.LoadLocation(v)

		if err != nil {
			return nil, fmt.Errorf("Invalid value for location: %s", err)
		}

		conf.location = loc
	}

	parsed.RawQuery = ""
	parsed.Fragment = ""

	conf.endpoint = parsed.String()

	return conf, nil
}
