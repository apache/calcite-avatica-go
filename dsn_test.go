package avatica

import (
	"testing"
	"time"
)

func TestParseDSN(t *testing.T) {

	config, err := ParseDSN("http://localhost:8765?maxRowsTotal=1&frameMaxSize=1&location=Australia/Melbourne&schema=myschema")

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if config.endpoint != "http://localhost:8765" {
		t.Errorf("Expected endpoint to be %s, got %s", "http://localhost:8765", config.endpoint)
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
