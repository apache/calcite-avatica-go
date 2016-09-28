# Apache Phoenix/Avatica SQL Driver
[![GoDoc](https://godoc.org/github.com/Boostport/avatica?status.png)](https://godoc.org/github.com/Boostport/Avatica)
[![wercker status](https://app.wercker.com/status/1abd1c7014e780ba7754decadb212451/s/master "wercker status")](https://app.wercker.com/project/byKey/1abd1c7014e780ba7754decadb212451)
[![Coverage Status](https://coveralls.io/repos/github/Boostport/avatica/badge.svg?branch=master)](https://coveralls.io/github/Boostport/avatica?branch=master)

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

#### maxRowsTotal

The `maxRowsTotal` parameter sets the maximum number of rows to return for a given query. By default, this is set to
`-1`, so that there is no limit on the number of rows returned.

#### frameMaxSize

The `frameMaxSize` parameter sets the maximum number of rows to return in a frame. Depending on the number of rows
returned and subject to the limits of `maxRowsTotal`, a query result set can contain rows in multiple frames. These
additional frames are then fetched on a as-needed basis. `frameMaxSize` allows you to control the number of rows
in each frame to suit your application's performance profile. By default this is set to `-1`, so that there is no limit
on the number of rows in a frame.

### time.Time support

The following Phoenix/Avatica datatypes are automatically converted to and from `time.Time`:
`TIME`, `DATE` and `TIMESTAMP`.

It is important to understand that avatica and the underlying database ignores the timezone. If you save a `time.Time`
to the database, the timezone is ignored and vice-versa. This is why you need to make sure the `location` parameter
in your DSN is set to the same value as the location of the `time.Time` values you are inserting into the database.

We recommend using `UTC`, which is the default value of `location`.

## Version compatibility
| Driver Version  | Phoenix Version | Calcite/Avatica Version |
| --------------- | -------------   | ----------------------- |
| 1.0.0           | 4.8.0           | 1.8.0                   |

## Development

To run tests, but skip tests in the vendor directory, run:

```
go test $(go list ./... | grep -v /vendor/)
```

The driver is not feature-complete yet, so contributions are very appreciated.

#### Updating protocol buffer definitions
To update the procotol buffer definitions, update `CALCITE_VER` in `gen-protobuf.bat` and `gen-protobuf.sh` to match
the version included by Phoenix and then run the appropriate script for your platform.

#### About the moby.yml file
The moby.yml file is used by our internal tool to automatically reload and test the code during development.
We hope to have this tool open-sourced soon.

## License
The driver is licensed under the Apache 2 license.