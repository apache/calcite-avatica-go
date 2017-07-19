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
http://[username:password@]address:port[/schema][?parameter1=value&...parameterN=value]
```

In other words, the scheme (http), address and port is mandatory, but the schema and parameters are optional.

#### username
This is the JDBC username that is passed directly to the backing database. It is *NOT* used for authenticating
against Avatica.

#### password
This is the JDBC password that is passed directly to the backing database. It is *NOT* used for authenticating
against Avatica.

#### schema
The `schema` path sets the default schema to use for this connection. For example, if you set it to `myschema`,
then executing the query `SELECT * FROM my_table` will have the equivalence of `SELECT * FROM myschema.my_table`.
If schema is set, you can still work on tables in other schemas by supplying a schema prefix:
`SELECT * FROM myotherschema.my_other_table`.

The following parameters are supported:

#### authentication
The authentication type to use when authenticating against Avatica. Valid values are `BASIC` for HTTP Basic authentication,
`DIGEST` for HTTP Digest authentication, and `SPNEGO` for Kerberos with SPNEGO authentication.

#### avaticaUser
The user to use when authenticating against Avatica. This parameter is required if `authentication` is `BASIC` or `DIGEST`.

#### avaticaPassword
The password to use when authenticating against Avatica. This parameter is required if `authentication` is `BASIC` or `DIGEST`.

#### principal
The Kerberos principal to use when authenticating against Avatica. It should be in the form `primary/instance@realm`, where
the instance is optional. This parameter is required if `authentication` is `SPNEGO` and you want the driver to perform the
Kerberos login.

#### keytab
The path to the Kerberos keytab to use when authenticating against Avatica. This parameter is required if `authentication`
is `SPNEGO` and you want the driver to perform the Kerberos login.

#### krb5Conf
The path to the Kerberos configuration to use when authenticating against Avatica. This parameter is required if `authentication`
is `SPNEGO` and you want the driver to perform the Kerberos login.

#### krb5CredentialsCache
The path to the Kerberos credential cache file to use when authenticating against Avatica. This parameter is required if
`authentication` is `SPNEGO` and you have logged into Kerberos already and want the driver to use the existing credentials.

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

#### transactionIsolation

Setting `transactionIsolation` allows you to set the isolation level for transactions using the connection. The value
should be a positive integer analogous to the transaction levels defined by the JDBC specification. The default value
is `0`, which means transactions are not supported. This is to deal with the fact that Calcite/Avatica works with
many types of backends, with some backends having no transaction support. If you are using Apache Phoenix 4.7 onwards,
we recommend setting it to `4`, which is the maximum isolation level supported.

The supported values for `transactionIsolation` are:

| Value | JDBC Constant                  | Description                                                                      |
| ----- | ------------------------------ | -------------------------------------------------------------------------------- |
| 0     | none                           | Transactions are not supported                                                   |
| 1     | `TRANSACTION_READ_UNCOMMITTED` | Dirty reads, non-repeatable reads and phantom reads may occur.                   |
| 2     | `TRANSACTION_READ_COMMITTED`   | Dirty reads are prevented, but non-repeatable reads and phantom reads may occur. |
| 4     | `TRANSACTION_REPEATABLE_READ`  | Dirty reads and non-repeatable reads are prevented, but phantom reads may occur. |
| 8     | `TRANSACTION_SERIALIZABLE`     | Dirty reads, non-repeatable reads, and phantom reads are all prevented.          |

### time.Time support

The following Phoenix/Avatica datatypes are automatically converted to and from `time.Time`:
`TIME`, `DATE` and `TIMESTAMP`.

It is important to understand that avatica and the underlying database ignores the timezone. If you save a `time.Time`
to the database, the timezone is ignored and vice-versa. This is why you need to make sure the `location` parameter
in your DSN is set to the same value as the location of the `time.Time` values you are inserting into the database.

We recommend using `UTC`, which is the default value of `location`.

## Version compatibility
| Driver Version  | Phoenix Version   | Calcite-Avatica Version |
| --------------- | ----------------- | ----------------------- |
| 1.x.x           | >= 4.8.0          | >= 1.8.0                |
| 2.x.x           | >= 4.8.0          | >= 1.8.0                |

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