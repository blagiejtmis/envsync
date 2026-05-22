package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendHistoryCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	entries := []Entry{{Key: "FOO", Value: "bar"}}
	if err := AppendHistory(path, "initial", entries, 10); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestAppendHistoryAccumulates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 3; i++ {
		e := []Entry{{Key: "K", Value: "v"}}
		if err := AppendHistory(path, "", e, 10); err != nil {
			t.Fatalf("AppendHistory iteration %d: %v", i, err)
		}
	}

	hf, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(hf.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(hf.Entries))
	}
}

func TestAppendHistoryPrunesOldEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 7; i++ {
		e := []Entry{{Key: "K", Value: "v"}}
		if err := AppendHistory(path, "", e, 5); err != nil {
			t.Fatalf("AppendHistory: %v", err)
		}
	}

	hf, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(hf.Entries) != 5 {
		t.Errorf("expected 5 entries after pruning, got %d", len(hf.Entries))
	}
}

func TestLoadHistoryNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")

	_, err := LoadHistory(path)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestHistoryEntryLabel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	e := []Entry{{Key: "A", Value: "1"}}
	if err := AppendHistory(path, "deploy-v1.2", e, 10); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}

	hf, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if hf.Entries[0].Label != "deploy-v1.2" {
		t.Errorf("expected label 'deploy-v1.2', got %q", hf.Entries[0].Label)
	}
	if hf.Entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
