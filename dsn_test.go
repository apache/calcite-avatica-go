package avatica

import (
	"strconv"
	"testing"
	"time"
)

func TestParseDSN(t *testing.T) {

	config, err := ParseDSN("http://localhost:8765/myschema?maxRowsTotal=1&frameMaxSize=1&location=Australia/Melbourne&transactionIsolation=8")

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if config.endpoint != "http://localhost:8765/myschema" {
		t.Errorf("Expected endpoint to be %s, got %s", "http://localhost:8765/myschema", config.endpoint)
	}

	if config.frameMaxSize != 1 {
		t.Errorf("Expected frameMaxSize to be %d, got %d", 1, config.frameMaxSize)
	}

	if config.maxRowsTotal != 1 {
		t.Errorf("Expected maxRowsTotal to be %d, got %d", 1, config.maxRowsTotal)
	}

	if config.location.String() != "Australia/Melbourne" {
		t.Errorf("Expected timezone to be %s, got %s", "Australia/Melbourne", config.location)
	}

	if config.schema != "myschema" {
		t.Errorf("Expected schema to be %s, got %s", "myschema", config.schema)
	}

	if config.transactionIsolation != 8 {
		t.Errorf("Expected transactionIsolation to be %d, got %d", 8, config.transactionIsolation)
	}
}

func TestParseEmptyDSN(t *testing.T) {

	_, err := ParseDSN("")

	if err == nil {
		t.Fatal("Expected error due to empty DSN, but received nothing")
	}
}

func TestDSNDefaults(t *testing.T) {

	config, err := ParseDSN("http://localhost:8765")

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if config.location.String() == "" {
		t.Error("There was no timezone set.")
	}

	if config.maxRowsTotal == 0 {
		t.Error("There was no maxRowsTotal set.")
	}

	if config.frameMaxSize == 0 {
		t.Error("There was no fetchMaxSize set.")
	}

	if config.schema != "" {
		t.Errorf("Unexpected schema set: %s", config.schema)
	}

	if config.transactionIsolation != 0 {
		t.Errorf("Default transaction level should be %d, got %d.", 0, config.transactionIsolation)
	}
}

func TestLocallocation(t *testing.T) {

	config, err := ParseDSN("http://localhost:8765?location=Local")

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if config.location != time.Local {
		t.Fatal("DSN has location set to 'local', but configuration did not set location to 'local'.")
	}
}

func TestBadInput(t *testing.T) {

	_, err := ParseDSN("http://localhost:8765?location=asdfasdf")

	if err == nil {
		t.Fatal("Expected error due to invalid location, but did not receive any.")
	}

	_, err = ParseDSN("http://localhost:8765?maxRowsTotal=abc")

	if err == nil {
		t.Fatal("Expected error due to invalid maxRowsTotal, but did not receive any.")
	}

	_, err = ParseDSN("http://localhost:8765?frameMaxSize=abc")

	if err == nil {
		t.Fatal("Expected error due to invalid frameMaxSize, but did not receive any.")
	}
}

func TestInvalidTransactionIsolation(t *testing.T) {

	badIsolationLevels := []int{-1, 3, 5, 6, 7, 9, 10, 11, 100}

	for _, isolationLevel := range badIsolationLevels {

		_, err := ParseDSN("http://localhost:8765?transactionIsolation=" + strconv.Itoa(isolationLevel))

		if err == nil {
			t.Fatal("Expected error due to invalid transactionIsolation, but did not receive any.")
		}
	}
}

func TestValidTransactionIsolation(t *testing.T) {

	validIsolationLevels := []int{0, 1, 2, 4, 8}

	for _, isolationLevel := range validIsolationLevels {

		_, err := ParseDSN("http://localhost:8765?transactionIsolation=" + strconv.Itoa(isolationLevel))

		if err != nil {
			t.Fatalf("Unexpected error when %d is set as the isolation level: %s", isolationLevel, err)
		}
	}
}
