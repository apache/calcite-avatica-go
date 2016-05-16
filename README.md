# Avatica SQL Driver
An Apache Phoenix/Avatica driver for Go's [database/sql](http://golang.org/pkg/database/sql) package

## Getting started
Install using the go tool or your dependency management tool:

```
$ go get github.com/Boostport/avatica
```

## Usage

The Phoenix/Avatica driver implements Go's `database/sql/driver` interface, so, import Go's
`database/sql` package and the driver:

```
import "database/sql"
import _ "github.com/Boostport/avatica"

db, err := sql.Open("avatica", "http://localhost:8765")
```

Then simply use the database connection to query some data, for example:

```
rows := db.Query("SELECT COUNT(*) FROM test")
```

### DSN (Data Source Name)

The DSN has the following format (optional parts are marked by square brackets):

```
http://address:ports[?parameter1=value&...parameterN=value]
```

In other words, the scheme (http), address and port is mandatory, but the parameters are optional.

The following parameters are supported:

#### location

The `location` will be set as the location of unserialized `time.Time` values. It must be a valid timezone.
If you want to use the local timezone, use `Local`. By default, this is set to `UTC`.

#### maxRowCount

The `maxRowCount` sets the number of rows to fetch from the avatica server. Setting a high number causes avatica
to return a large number of rows and setting a low number will cause avatica to return a small amount of rows
and fetch more from the server if required. By default, this is set to `100`.

### time.Time support

The following Phoenix/Avatica datatypes are automatically converted to and from `time.Time`:
`TIME`, `DATE` and `TIMESTAMP`.

It is important to understand that avatica and the underlying database ignores the timezone. If you save a `time.Time`
to the database, the timezone is ignored and vice-versa. This is why you need to make sure the `location` parameter
in your DSN is set to the same value as the location of the `time.Time` values you are inserting into the database.

We recommend using `UTC`, which is the default value of `location`.

## Development

To run tests, but skip tests in the vendor directory, run:

```
go test $(go list ./... | grep -v /vendor/)
```

The driver is not feature-complete yet, so contributions are very appreciated.

#### About the moby.yml file
The moby.yml file is used by our internal tool to automatically reload and test the code during development.
We hope to have this tool open-sourced soon.