package avatica

import (
	"testing"
	"time"
)

func TestParseDSN(t *testing.T) {

	config, err := ParseDSN("http://localhost:8765?maxRowCount=1&location=Australia/Melbourne")

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if config.endpoint != "http://localhost:8765" {
		t.Errorf("Expected endpoint to be %s, got %s", "http://localhost:8765", config.endpoint)
	}

	if config.fetchMaxRowCount != 1 {
		t.Errorf("Expected fetchMaxRowCount to be %d, got %d", 1, config.fetchMaxRowCount)
	}

	if config.maxRowCount != 1 {
		t.Errorf("Expected maxRowCount to be %d, got %d", 1, config.maxRowCount)
	}

	if config.location.String() != "Australia/Melbourne" {
		t.Errorf("Expected timezone to be %s, got %s", "Australia/Melbourne", config.location)
	}
}

func TestParseEmptyDSN(t *testing.T) {

	_, err := ParseDSN("")

	if err == nil {
		t.Fatalf("Expected error due to empty DSN, but received nothing")
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

	if config.maxRowCount == 0 {
		t.Error("There was no maxRowCount set.")
	}

	if config.fetchMaxRowCount == 0 {
		t.Error("There was no fetchMaxRowCount set.")
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

	_, err = ParseDSN("http://localhost:8765?maxRowCount=abc")

	if err == nil {
		t.Fatal("Expected error due to invalid maxRowCount, but did not receive any.")
	}

	_, err = ParseDSN("http://localhost:8765?maxRowCount=0")

	if err == nil {
		t.Fatal("Expected error due to invalid maxRowCount, but did not receive any.")
	}

	_, err = ParseDSN("http://localhost:8765?maxRowCount=-1")

	if err == nil {
		t.Fatal("Expected error due to invalid maxRowCount, but did not receive any.")
	}
}
