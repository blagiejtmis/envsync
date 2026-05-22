package env_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envsync/internal/env"
)

func TestTakeSnapshotFields(t *testing.T) {
	entries := []env.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	}
	before := time.Now().UTC()
	snap := env.TakeSnapshot(entries, "v1")
	after := time.Now().UTC()

	if snap.Label != "v1" {
		t.Errorf("expected label v1, got %q", snap.Label)
	}
	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Errorf("timestamp out of expected range: %v", snap.Timestamp)
	}
	if snap.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
}

func TestSnapshotVerify(t *testing.T) {
	entries := []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
	}
	snap := env.TakeSnapshot(entries, "")

	if !snap.Verify() {
		t.Error("expected Verify() to return true for untampered snapshot")
	}

	snap.Entries[0].Value = "tampered"
	if snap.Verify() {
		t.Error("expected Verify() to return false after tampering")
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := env.TakeSnapshot([]env.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}, "test-label")

	if err := env.SaveSnapshot(path, orig); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := env.LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if loaded.Label != orig.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, orig.Label)
	}
	if loaded.Checksum != orig.Checksum {
		t.Errorf("checksum mismatch: got %q, want %q", loaded.Checksum, orig.Checksum)
	}
	if !loaded.Verify() {
		t.Error("loaded snapshot failed Verify()")
	}
}

func TestLoadSnapshotFileNotFound(t *testing.T) {
	_, err := env.LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestChecksumIsDeterministic(t *testing.T) {
	entries := []env.Entry{
		{Key: "Z", Value: "last"},
		{Key: "A", Value: "first"},
	}
	s1 := env.TakeSnapshot(entries, "")

	reversed := []env.Entry{entries[1], entries[0]}
	s2 := env.TakeSnapshot(reversed, "")

	if s1.Checksum != s2.Checksum {
		t.Errorf("checksum should be order-independent: %q vs %q", s1.Checksum, s2.Checksum)
	}
}

func TestSaveSnapshotBadPath(t *testing.T) {
	snap := env.TakeSnapshot(nil, "")
	err := env.SaveSnapshot("/no/such/dir/snap.json", snap)
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
	_ = os.Remove("/no/such/dir/snap.json")
}
