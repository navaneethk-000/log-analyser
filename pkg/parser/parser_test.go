package parser

import (
	"os"
	"path/filepath"
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

func TestParseLogFiles(t *testing.T) {
	entries, err := ParseLogFiles("/home/navaneeth/project_go/logs")

	if err != nil {
		t.Fatalf("Expected no error, got:%v", err)
	}

	if len(entries) == 0 {
		t.Fatalf("Expected some log entries, but got 0")
	}
}

func TestParseLogFilesValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	logContent := `2025-10-23 15:17:08.636 | INFO | api-server | host=worker01 | request_id=req-xyz | msg="Cache cleared"`
	tmpFile := filepath.Join(tmpDir, "valid.log")

	err := os.WriteFile(tmpFile, []byte(logContent+"\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Host != "worker01" {
		t.Errorf("Expected host=worker01, got %s", entries[0].Host)
	}
}

func TestParseLogFilesInvalidDirectoryPath(t *testing.T) {
	path := "../logg"
	_, err := ParseLogFiles(path)
	if err == nil {
		t.Errorf("Expected 'no such directory' error but got none.")
	}
}

func TestParseLogFiles_FailedToOpen(t *testing.T) {
	tmpDir := t.TempDir()
	errorFile := filepath.Join(tmpDir, "failopen.log")

	err := os.WriteFile(errorFile, []byte("data"), 0000)
	if err != nil {
		t.Fatalf("Unable to create file: %v", err)
	}
	defer os.Chmod(errorFile, 0644)

	_, err = ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
}

func TestParseLogFiles_IgnoresSubdir(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "logs")
	os.Mkdir(subDir, 0755)

	sampleLogEntry := `2025-10-29 15:41:14.814 | DEBUG | cache | host=db01 | request_id=req-z2t4p5-1421 | msg="Deadlock detected, retrying"`
	tmpFile := filepath.Join(tmpDir, "log1.log")
	os.WriteFile(tmpFile, []byte(sampleLogEntry+"\n"), 0644)

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 file entry, but got %d", len(entries))
	}
}
