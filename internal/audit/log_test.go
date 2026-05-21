package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/audit"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.log")
}

func TestRecordAndReadAll(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	if err := logger.Record(audit.EventPush, "DB_URL", "alice", "pushed secret"); err != nil {
		t.Fatalf("Record push: %v", err)
	}
	if err := logger.Record(audit.EventPull, "DB_URL", "bob", "pulled secret"); err != nil {
		t.Fatalf("Record pull: %v", err)
	}

	entries, err := logger.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Kind != audit.EventPush {
		t.Errorf("expected push, got %s", entries[0].Kind)
	}
	if entries[1].User != "bob" {
		t.Errorf("expected user bob, got %s", entries[1].User)
	}
}

func TestReadAllEmptyFile(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	entries, err := logger.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecordCreatesFile(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	if err := logger.Record(audit.EventDelete, "API_KEY", "", ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("log file not created: %v", err)
	}
}

func TestRecordMultipleKinds(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	kinds := []audit.EventKind{audit.EventPush, audit.EventPull, audit.EventDelete}
	for _, k := range kinds {
		if err := logger.Record(k, "KEY", "user", ""); err != nil {
			t.Fatalf("Record %s: %v", k, err)
		}
	}

	entries, err := logger.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, e := range entries {
		if e.Kind != kinds[i] {
			t.Errorf("entry %d: expected %s, got %s", i, kinds[i], e.Kind)
		}
	}
}
