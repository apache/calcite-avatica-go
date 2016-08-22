package avatica

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Config is a configuration parsed from a DSN string
type Config struct {
	endpoint     string
	maxRowsTotal int64
	frameMaxSize int32
	location     *time.Location
}

// ParseDSN parses a DSN string to a Config
func ParseDSN(dsn string) (*Config, error) {

	conf := &Config{
		maxRowsTotal: -1,
		frameMaxSize: -1,
		location:     time.UTC,
	}

	parsed, err := url.ParseRequestURI(dsn)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse DSN: %s", err)
	}

	queries := parsed.Query()

	if v := queries.Get("maxRowsTotal"); v != "" {

		maxRowTotal, err := strconv.Atoi(v)

		if err != nil {
			return nil, fmt.Errorf("Invalid value for maxRowsTotal: %s", err)
		}

		conf.maxRowsTotal = int64(maxRowTotal)
	}

	if v := queries.Get("frameMaxSize"); v != "" {

		maxRowTotal, err := strconv.Atoi(v)

		if err != nil {
			return nil, fmt.Errorf("Invalid value for frameMaxSize: %s", err)
		}

		conf.frameMaxSize = int32(maxRowTotal)
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
