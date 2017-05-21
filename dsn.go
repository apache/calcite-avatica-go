package avatica

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Config is a configuration parsed from a DSN string
type Config struct {
	endpoint             string
	maxRowsTotal         int64
	frameMaxSize         int32
	location             *time.Location
	schema               string
	transactionIsolation uint32
	user                 string
	password             string
}

// ParseDSN parses a DSN string to a Config
func ParseDSN(dsn string) (*Config, error) {

	conf := &Config{
		maxRowsTotal:         -1,
		frameMaxSize:         -1,
		location:             time.UTC,
		transactionIsolation: 0,
	}

	parsed, err := url.ParseRequestURI(dsn)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse DSN: %s", err)
	}

	userInfo := parsed.User

	if userInfo != nil {
		if userInfo.Username() != "" {
			conf.user = userInfo.Username()
		}

		if pass, ok := userInfo.Password(); ok {
			conf.password = pass
		}
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

	if parsed.Path != "" {
		conf.schema = strings.TrimPrefix(parsed.Path, "/")
	}

	if v := queries.Get("transactionIsolation"); v != "" {

		isolation, err := strconv.Atoi(v)

		if err != nil {
			return nil, fmt.Errorf("Invalid value for transactionIsolation: %s", err)
		}

		if isolation < 0 || isolation > 8 || isolation&(isolation-1) != 0 {
			return nil, fmt.Errorf("transactionIsolation must be 0, 1, 2, 4 or 8, %d given", isolation)
		}

		conf.transactionIsolation = uint32(isolation)
	}

	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""

	conf.endpoint = parsed.String()

	return conf, nil
}
