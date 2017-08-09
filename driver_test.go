package avatica

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

var (
	dsn string
)

func init() {

	// get environment variables
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	dsn = env("AVATICA_HOST", "http://phoenix-server:8765")
}

func generateTableName() string {
	return "test" + strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

type DBTest struct {
	*testing.T
	db        *sql.DB
	tableName string
}

func (dbt *DBTest) fail(method, query string, err error) {

	if len(query) > 300 {
		query = "[query too large to print]"
	}

	dbt.Fatalf("error on %s %s: %s", method, query, err.Error())
}

func (dbt *DBTest) mustExec(query string, args ...interface{}) (res sql.Result) {
	res, err := dbt.db.Exec(query, args...)

	if err != nil {
		dbt.fail("exec", query, err)
	}

	return res
}

func (dbt *DBTest) mustQuery(query string, args ...interface{}) (rows *sql.Rows) {
	rows, err := dbt.db.Query(query, args...)

	if err != nil {
		dbt.fail("query", query, err)
	}

	return rows
}

func runTests(t *testing.T, dsn string, tests ...func(dbt *DBTest)) {

	db, err := sql.Open("avatica", dsn)

	if err != nil {
		t.Fatalf("error connecting: %s", err.Error())
	}

	defer db.Close()

	table := generateTableName()

	db.Exec("DROP TABLE IF EXISTS " + table)

	dbt := &DBTest{t, db, table}

	for _, test := range tests {
		test(dbt)
		dbt.db.Exec("DROP TABLE IF EXISTS " + table)
	}
}

func TestConnectionMustBeOpenedWithAutoCommitTrue(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec("CREATE TABLE " + dbt.tableName + " (id BIGINT PRIMARY KEY, val VARCHAR) TRANSACTIONAL=false")

		dbt.mustExec("UPSERT INTO " + dbt.tableName + " VALUES (1,'A')")

		dbt.mustExec("UPSERT INTO " + dbt.tableName + " VALUES (2,'B')")

		rows := dbt.mustQuery("SELECT COUNT(*) FROM " + dbt.tableName)
		defer rows.Close()

		for rows.Next() {

			var count int

			err := rows.Scan(&count)

			if err != nil {
				dbt.Fatal(err)
			}

			if count != 2 {
				dbt.Fatalf("There should be 2 rows, got %d", count)
			}
		}

	})
}

func TestZeroValues(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec("CREATE TABLE " + dbt.tableName + " (int INTEGER PRIMARY KEY, flt FLOAT, bool BOOLEAN, str VARCHAR) TRANSACTIONAL=false")

		dbt.mustExec("UPSERT INTO " + dbt.tableName + " VALUES (0, 0.0, false, '')")

		rows := dbt.mustQuery("SELECT * FROM " + dbt.tableName)
		defer rows.Close()

		for rows.Next() {

			var i int
			var flt float64
			var b bool
			var s string

			err := rows.Scan(&i, &flt, &b, &s)

			if err != nil {
				dbt.Fatal(err)
			}

			if i != 0 {
				dbt.Fatalf("Integer should be 0, got %v", i)
			}

			if flt != 0.0 {
				dbt.Fatalf("Float should be 0.0, got %v", flt)
			}

			if b != false {
				dbt.Fatalf("Boolean should be false, got %v", b)
			}

			if s != "" {
				dbt.Fatalf("String should be \"\", got %v", s)
			}
		}

	})
}

func TestDataTypes(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY,
				uint UNSIGNED_INT,
				bint BIGINT,
				ulong UNSIGNED_LONG,
				tint TINYINT,
				utint UNSIGNED_TINYINT,
				sint SMALLINT,
				usint UNSIGNED_SMALLINT,
				flt FLOAT,
				uflt UNSIGNED_FLOAT,
				dbl DOUBLE,
				udbl UNSIGNED_DOUBLE,
				dec DECIMAL,
				bool BOOLEAN,
				tm TIME,
				dt DATE,
				tmstmp TIMESTAMP,
				utm UNSIGNED_TIME,
				udt UNSIGNED_DATE,
				utmstmp UNSIGNED_TIMESTAMP,
				var VARCHAR,
				ch CHAR(3),
				bin BINARY(20),
				varbin VARBINARY
			    ) TRANSACTIONAL=false`)

		var (
			integerValue  int       = -20
			uintegerValue int       = 5
			bintValue     int       = -9223372036854775807
			ulongValue    int       = 9223372036854775807
			tintValue     int       = -128
			utintValue    int       = 126
			sintValue     int       = -32768
			usintValue    int       = 32767
			fltValue      float64   = -3.555
			ufltValue     float64   = 3.555
			dblValue      float64   = -9.555
			udblValue     float64   = 9.555
			decValue      string    = "1.333"
			booleanValue  bool      = true
			tmValue       time.Time = time.Date(0, 1, 1, 21, 21, 21, 222000000, time.UTC)
			dtValue       time.Time = time.Date(2100, 2, 1, 0, 0, 0, 0, time.UTC)
			tmstmpValue   time.Time = time.Date(2100, 2, 1, 21, 21, 21, 222000000, time.UTC)
			utmValue      time.Time = time.Date(0, 1, 1, 21, 21, 21, 222000000, time.UTC)
			udtValue      time.Time = time.Date(2100, 2, 1, 0, 0, 0, 0, time.UTC)
			utmstmpValue  time.Time = time.Date(2100, 2, 1, 21, 21, 21, 222000000, time.UTC)
			varcharValue  string    = "test string"
			chValue       string    = "a"
			binValue      []byte    = make([]byte, 20, 20)
			varbinValue   []byte    = []byte("testtesttest")
		)

		copy(binValue[:], []byte("test"))

		dbt.mustExec(`UPSERT INTO `+dbt.tableName+` VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			integerValue,
			uintegerValue,
			bintValue,
			ulongValue,
			tintValue,
			utintValue,
			sintValue,
			usintValue,
			fltValue,
			ufltValue,
			dblValue,
			udblValue,
			decValue,
			booleanValue,
			tmValue,
			dtValue,
			tmstmpValue,
			utmValue,
			udtValue,
			utmstmpValue,
			varcharValue,
			chValue,
			binValue,
			varbinValue,
		)

		rows := dbt.mustQuery("SELECT * FROM " + dbt.tableName)
		defer rows.Close()

		var (
			integer  int
			uinteger int
			bint     int
			ulong    int
			tint     int
			utint    int
			sint     int
			usint    int
			flt      float64
			uflt     float64
			dbl      float64
			udbl     float64
			dec      string
			boolean  bool
			tm       time.Time
			dt       time.Time
			tmstmp   time.Time
			utm      time.Time
			udt      time.Time
			utmstmp  time.Time
			varchar  string
			ch       string
			bin      []byte
			varbin   []byte
		)

		for rows.Next() {

			err := rows.Scan(&integer, &uinteger, &bint, &ulong, &tint, &utint, &sint, &usint, &flt, &uflt, &dbl, &udbl, &dec, &boolean, &tm, &dt, &tmstmp, &utm, &udt, &utmstmp, &varchar, &ch, &bin, &varbin)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		comparisons := []struct {
			result   interface{}
			expected interface{}
		}{
			{integer, integerValue},
			{uinteger, uintegerValue},
			{bint, bintValue},
			{ulong, ulongValue},
			{tint, tintValue},
			{utint, utintValue},
			{sint, sintValue},
			{usint, usintValue},
			{flt, fltValue},
			{uflt, ufltValue},
			{dbl, dblValue},
			{udbl, udblValue},
			{dec, decValue},
			{boolean, booleanValue},
			{tm, tmValue},
			{dt, dtValue},
			{tmstmp, tmstmpValue},
			{utm, utmValue},
			{udt, udtValue},
			{utmstmp, utmstmpValue},
			{varchar, varcharValue},
			{ch, chValue},
			{bin, binValue},
			{varbin, varbinValue},
		}

		for _, tt := range comparisons {

			if v, ok := tt.expected.(time.Time); ok {

				if !v.Equal(tt.result.(time.Time)) {
					dbt.Fatalf("Expected %v, got %v.", tt.expected, tt.result)
				}

			} else if v, ok := tt.expected.([]byte); ok {

				if !bytes.Equal(v, tt.result.([]byte)) {
					dbt.Fatalf("Expected %v, got %v.", tt.expected, tt.result)
				}

			} else if tt.expected != tt.result {

				dbt.Fatalf("Expected %v, got %v.", tt.expected, tt.result)
			}
		}
	})
}

func TestLocations(t *testing.T) {

	query := "?location=Australia/Melbourne"

	runTests(t, dsn+query, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				tm TIME PRIMARY KEY,
				dt DATE,
				tmstmp TIMESTAMP
			    ) TRANSACTIONAL=false`)

		loc, err := time.LoadLocation("Australia/Melbourne")

		if err != nil {
			dbt.Fatalf("Unexpected error: %s", err)
		}

		var (
			tmValue     time.Time = time.Date(0, 1, 1, 21, 21, 21, 222000000, loc)
			dtValue     time.Time = time.Date(2100, 2, 1, 0, 0, 0, 0, loc)
			tmstmpValue time.Time = time.Date(2100, 2, 1, 21, 21, 21, 222000000, loc)
		)

		dbt.mustExec(`UPSERT INTO `+dbt.tableName+` VALUES (?, ?, ?)`,
			tmValue,
			dtValue,
			tmstmpValue,
		)

		rows := dbt.mustQuery("SELECT * FROM " + dbt.tableName)
		defer rows.Close()

		var (
			tm     time.Time
			dt     time.Time
			tmstmp time.Time
		)

		for rows.Next() {

			err := rows.Scan(&tm, &dt, &tmstmp)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		comparisons := []struct {
			result   time.Time
			expected time.Time
		}{
			{tm, tmValue},
			{dt, dtValue},
			{tmstmp, tmstmpValue},
		}

		for _, tt := range comparisons {

			if !tt.result.Equal(tt.expected) {
				dbt.Fatalf("Expected %v, got %v.", tt.expected, tt.result)
			}
		}
	})
}

func TestDateAndTimestampsBefore1970(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY,
				dt DATE,
				tmstmp TIMESTAMP
			    ) TRANSACTIONAL=false`)

		var (
			integerValue int       = 1
			dtValue      time.Time = time.Date(1945, 5, 20, 0, 0, 0, 0, time.UTC)
			tmstmpValue  time.Time = time.Date(1911, 5, 20, 21, 21, 21, 222000000, time.UTC)
		)

		dbt.mustExec(`UPSERT INTO `+dbt.tableName+` VALUES (?, ?, ?)`,
			integerValue,
			dtValue,
			tmstmpValue,
		)

		rows := dbt.mustQuery("SELECT dt, tmstmp FROM " + dbt.tableName)
		defer rows.Close()

		var (
			dt     time.Time
			tmstmp time.Time
		)

		for rows.Next() {
			err := rows.Scan(&dt, &tmstmp)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		comparisons := []struct {
			result   time.Time
			expected time.Time
		}{
			{dt, dtValue},
			{tmstmp, tmstmpValue},
		}

		for _, tt := range comparisons {
			if !tt.expected.Equal(tt.result) {
				dbt.Fatalf("Expected %v, got %v.", tt.expected, tt.result)
			}
		}
	})
}

func TestStoreAndRetrieveBinaryData(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY,
				bin VARBINARY
			    ) TRANSACTIONAL=false`)

		filePath := filepath.Join("test-fixtures", "gopher.png")

		file, err := ioutil.ReadFile(filePath)

		if err != nil {
			t.Fatalf("Unable to read text-fixture: %s", filePath)
		}

		hash := sha256.Sum256(file)

		dbt.mustExec(`UPSERT INTO `+dbt.tableName+` VALUES (?, ?)`,
			1,
			file,
		)

		rows := dbt.mustQuery("SELECT bin FROM " + dbt.tableName)
		defer rows.Close()

		var receivedFile []byte

		for rows.Next() {

			err := rows.Scan(&receivedFile)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		ioutil.WriteFile("test-fixtures/gopher.png", receivedFile, os.ModePerm)

		receivedHash := sha256.Sum256(receivedFile)

		if !bytes.Equal(hash[:], receivedHash[:]) {
			t.Fatalf("Hash of stored file (%x) does not equal hash of retrieved file (%x).", hash[:], receivedHash[:])
		}
	})
}

func TestCommittingTransactions(t *testing.T) {

	query := "?transactionIsolation=4"

	runTests(t, dsn+query, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=true`)

		tx, err := dbt.db.Begin()

		if err != nil {
			t.Fatalf("Unable to create transaction: %s", err)
		}

		stmt, err := tx.Prepare(`UPSERT INTO ` + dbt.tableName + ` VALUES(?)`)

		if err != nil {
			t.Fatalf("Could not prepare statement: %s", err)
		}

		totalRows := 6

		for i := 1; i <= totalRows; i++ {
			_, err := stmt.Exec(i)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		r := tx.QueryRow("SELECT COUNT(*) FROM " + dbt.tableName)

		var count int

		err = r.Scan(&count)

		if err != nil {
			t.Fatalf("Unable to scan row result: %s", err)
		}

		if count != totalRows {
			t.Fatalf("Expected %d rows, got %d", totalRows, count)
		}

		// Commit the transaction
		tx.Commit()

		rows := dbt.mustQuery("SELECT COUNT(*) FROM " + dbt.tableName)

		var countAfterRollback int

		for rows.Next() {
			err := rows.Scan(&countAfterRollback)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		if countAfterRollback != totalRows {
			t.Fatalf("Expected %d rows, got %d", totalRows, countAfterRollback)
		}
	})
}

func TestRollingBackTransactions(t *testing.T) {

	query := "?transactionIsolation=4"

	runTests(t, dsn+query, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=true`)

		tx, err := dbt.db.Begin()

		if err != nil {
			t.Fatalf("Unable to create transaction: %s", err)
		}

		stmt, err := tx.Prepare(`UPSERT INTO ` + dbt.tableName + ` VALUES(?)`)

		if err != nil {
			t.Fatalf("Could not prepare statement: %s", err)
		}

		totalRows := 6

		for i := 1; i <= totalRows; i++ {
			_, err := stmt.Exec(i)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		r := tx.QueryRow(`SELECT COUNT(*) FROM ` + dbt.tableName)

		var count int

		err = r.Scan(&count)

		if err != nil {
			t.Fatalf("Unable to scan row result: %s", err)
		}

		if count != totalRows {
			t.Fatalf("Expected %d rows, got %d", totalRows, count)
		}

		// Rollback the transaction
		tx.Rollback()

		rows := dbt.mustQuery(`SELECT COUNT(*) FROM ` + dbt.tableName)

		var countAfterRollback int

		for rows.Next() {
			err := rows.Scan(&countAfterRollback)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		if countAfterRollback != 0 {
			t.Fatalf("Expected %d rows, got %d", 0, countAfterRollback)
		}
	})
}

func TestPreparedStatements(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=false`)

		stmt, err := dbt.db.Prepare(`UPSERT INTO ` + dbt.tableName + ` VALUES(?)`)

		if err != nil {
			dbt.Fatal(err)
		}

		totalRows := 6

		for i := 1; i <= totalRows; i++ {
			_, err := stmt.Exec(i)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		queryStmt, err := dbt.db.Prepare(`SELECT * FROM ` + dbt.tableName + ` WHERE int = ?`)

		if err != nil{
			dbt.Fatal(err)
		}

		var res int

		for i := 1; i <= totalRows; i++ {

			err := queryStmt.QueryRow(i).Scan(&res)

			if err != nil {
				dbt.Fatal(err)
			}

			if res != i {
				dbt.Fatalf("Unexpected query result. Expected %d, got %d.", i, res)
			}
		}
	})
}

func TestFetchingMoreRows(t *testing.T) {

	query := "?maxRowsTotal=-1&frameMaxSize=1"

	runTests(t, dsn+query, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=false`)

		stmt, err := dbt.db.Prepare(`UPSERT INTO ` + dbt.tableName + ` VALUES(?)`)

		if err != nil {
			dbt.Fatal(err)
		}

		totalRows := 6

		for i := 1; i <= totalRows; i++ {
			_, err := stmt.Exec(i)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		rows := dbt.mustQuery(`SELECT * FROM ` + dbt.tableName)
		defer rows.Close()

		count := 0

		for rows.Next() {
			count++
		}

		if count != totalRows {
			dbt.Fatalf("Expected %d rows to be retrieved, retrieved %d", totalRows, count)
		}
	})
}

func TestExecuteShortcut(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=false`)

		res, err := dbt.db.Exec(`UPSERT INTO ` + dbt.tableName + ` VALUES(1)`)

		if err != nil {
			dbt.Fatal(err)
		}

		affected, err := res.RowsAffected()

		if err != nil {
			dbt.Fatal(err)
		}

		if affected != 1 {
			dbt.Fatalf("Expected 1 row to be affected, %d affected", affected)
		}
	})
}

func TestQueryShortcut(t *testing.T) {

	query := "?maxRowsTotal=-1&frameMaxSize=1"

	runTests(t, dsn+query, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=false`)

		stmt, err := dbt.db.Prepare(`UPSERT INTO ` + dbt.tableName + ` VALUES(?)`)

		if err != nil {
			dbt.Fatal(err)
		}

		totalRows := 6

		for i := 1; i <= totalRows; i++ {
			_, err := stmt.Exec(i)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		rows := dbt.mustQuery(`SELECT * FROM ` + dbt.tableName)
		defer rows.Close()

		count := 0

		for rows.Next() {
			count++
		}

		if count != totalRows {
			dbt.Fatalf("Expected %d rows to be retrieved, retrieved %d", totalRows, count)
		}
	})
}

func TestOptimisticConcurrency(t *testing.T) {

	query := "?transactionIsolation=4"

	runTests(t, dsn+query, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				id INTEGER PRIMARY KEY,
				msg VARCHAR,
				version INTEGER
			    ) TRANSACTIONAL=true`)

		stmt, err := dbt.db.Prepare(`UPSERT INTO ` + dbt.tableName + ` VALUES(?, ?, ?)`)

		if err != nil {
			dbt.Fatal(err)
		}

		totalRows := 6

		for i := 1; i <= totalRows; i++ {
			_, err := stmt.Exec(i, fmt.Sprintf("message version %d", i), i)

			if err != nil {
				dbt.Fatal(err)
			}
		}

		// Start the transactions
		tx1, err := dbt.db.Begin()

		if err != nil {
			dbt.Fatal(err)
		}

		tx2, err := dbt.db.Begin()

		if err != nil {
			dbt.Fatal(err)
		}

		// Select from first transaction
		_ = tx1.QueryRow(`SELECT MAX(version) FROM ` + dbt.tableName)

		// Modify using second transaction
		_, err = tx2.Exec(`UPSERT INTO `+dbt.tableName+` VALUES(?, ?, ?)`, 7, "message value 7", 7)

		if err != nil {
			dbt.Fatal(err)
		}

		err = tx2.Commit()

		if err != nil {
			dbt.Fatal(err)
		}

		// Modify using tx1
		_, err = tx1.Exec(`UPSERT INTO `+dbt.tableName+` VALUES(?, ?, ?)`, 7, "message value 7", 7)

		if err != nil {
			dbt.Fatal(err)
		}

		err = tx1.Commit()

		if err == nil {
			dbt.Fatal("Expected an error, but did not receive any.")
		}

		errName := err.(ResponseError).Name()

		if errName != "transaction_conflict_exception" {
			dbt.Fatal("Expected transaction_conflict")
		}
	})
}

func TestLastInsertIDShouldReturnError(t *testing.T) {

	runTests(t, dsn, func(dbt *DBTest) {

		dbt.mustExec(`DROP SEQUENCE IF EXISTS test_sequence`)

		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				id INTEGER PRIMARY KEY,
				msg VARCHAR,
				version INTEGER
			    ) TRANSACTIONAL=false`)

		dbt.mustExec(`CREATE SEQUENCE test_sequence`)

		res, err := dbt.db.Exec(`UPSERT INTO ` + dbt.tableName + ` VALUES(NEXT VALUE FOR test_sequence, 'abc', 1)`)

		dbt.mustExec(`DROP SEQUENCE test_sequence`)

		if err != nil {
			dbt.Fatal(err)
		}

		_, err = res.LastInsertId()

		if err == nil {
			dbt.Fatal("Expected an error as Avatica does not support LastInsertId(), but there was no error.")
		}
	})
}

func TestConnectionToInvalidServerShouldReturnError(t *testing.T) {

	runTests(t, "http://invalid-server:8765", func(dbt *DBTest) {

		_, err := dbt.db.Exec(`CREATE TABLE ` + dbt.tableName + ` (
					id INTEGER PRIMARY KEY,
					msg VARCHAR,
			    	      ) TRANSACTIONAL=false`)

		if err == nil {
			dbt.Fatal("Expected an error due to connection to invalid server, but got nothing.")
		}
	})
}

func TestSchemaSupport(t *testing.T) {

	db, err := sql.Open("avatica", dsn)

	if err != nil {
		t.Fatalf("error connecting: %s", err.Error())
	}

	defer db.Close()

	db.Exec("CREATE SCHEMA IF NOT EXISTS avaticatest")
	defer db.Exec("DROP SCHEMA IF EXISTS avaticatest")

	path := "/avaticatest"

	runTests(t, dsn+path, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE ` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=false`)

		defer dbt.mustExec(`DROP TABLE IF EXISTS ` + dbt.tableName)

		_, err := dbt.db.Exec(`UPSERT INTO ` + dbt.tableName + ` VALUES(1)`)

		if err != nil {
			dbt.Fatal(err)
		}

		rows := dbt.mustQuery(`SELECT * FROM avaticatest.` + dbt.tableName)
		defer rows.Close()

		count := 0

		for rows.Next() {
			count++
		}

		if count != 1 {
			dbt.Errorf("Expected 1 row, got %d rows back,", count)
		}
	})
}

func TestMultipleSchemaSupport(t *testing.T) {

	db, err := sql.Open("avatica", dsn)

	if err != nil {
		t.Fatalf("error connecting: %s", err.Error())
	}

	defer db.Close()

	db.Exec("CREATE SCHEMA IF NOT EXISTS avaticatest1")
	defer db.Exec("DROP SCHEMA IF EXISTS avaticatest1")

	db.Exec("CREATE SCHEMA IF NOT EXISTS avaticatest2")
	defer db.Exec("DROP SCHEMA IF EXISTS avaticatest2")

	path := "/avaticatest1"

	runTests(t, dsn+path, func(dbt *DBTest) {

		// Create and seed table
		dbt.mustExec(`CREATE TABLE avaticatest2.` + dbt.tableName + ` (
				int INTEGER PRIMARY KEY
			    ) TRANSACTIONAL=false`)

		defer dbt.mustExec(`DROP TABLE IF EXISTS avaticatest2.` + dbt.tableName)

		_, err := dbt.db.Exec(`UPSERT INTO avaticatest2.` + dbt.tableName + ` VALUES(1)`)

		if err != nil {
			dbt.Fatal(err)
		}

		rows := dbt.mustQuery(`SELECT * FROM avaticatest2.` + dbt.tableName)
		defer rows.Close()

		count := 0

		for rows.Next() {
			count++
		}

		if count != 1 {
			dbt.Errorf("Expected 1 row, got %d rows back,", count)
		}
	})
}

func TestErrorCodeParsing(t *testing.T) {

	db, err := sql.Open("avatica", dsn)

	if err != nil {
		t.Fatalf("error connecting: %s", err.Error())
	}

	defer db.Close()

	_, err = db.Query("SELECT * FROM table_that_does_not_exist")

	if err == nil {
		t.Error("Expected error due to selecting from non-existent table, but there was no error.")
	}

	resErr, ok := err.(ResponseError)

	if !ok {
		t.Fatalf("Error type was not ResponseError")
	}

	if resErr.ErrorCode != 1012 {
		t.Errorf("Expected error code to be %d, got %d.", 1012, resErr.ErrorCode)
	}

	if resErr.SqlState != "42M03" {
		t.Errorf("Expected SQL state to be %s, got %s.", "42M03", resErr.SqlState)
	}
}
