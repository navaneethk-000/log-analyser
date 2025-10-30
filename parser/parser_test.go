package parser

import (
	"testing"
	"time"
)

func TestParserLogEntry(t *testing.T) {
	entry := `2025-10-23 16:53:00.033 | WARN | auth | host=worker01 | request_id=req-nkc44l-7700 | msg="Connection to cache server established"`

	got, _ := ParseLogEntry(entry)
	expectedTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 16:53:00.033")

	if got.Time != expectedTime {
		t.Errorf("Expected: %v, but got: %v", expectedTime, got.Time)
	}

	if got.Level != "WARN" {
		t.Errorf("Expected: 'WARN' but got: '%s'", got.Level)
	}

	if got.Component != "auth" {
		t.Errorf("Expected: 'auth' but got: '%s'", got.Component)
	}

	if got.Host != "worker01" {
		t.Errorf("Expected: 'worker01' but got: '%s'", got.Host)

	}

	if got.Request_id != "req-nkc44l-7700" {
		t.Errorf("Expected: 'eq-nkc44l-7700' but got: '%s'", got.Request_id)
	}

	if got.Message != "Connection to cache server established" {
		t.Errorf("Expected: 'Connection to cache server established' but got: '%s'", got.Message)

	}

	if got.Log != entry {
		t.Errorf("Expected: %v, but got: '%v'", entry, got.Log)
	}

}
