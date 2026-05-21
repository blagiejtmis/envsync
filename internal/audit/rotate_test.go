package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func writeLogFile(t *testing.T, path string, size int) {
	t.Helper()
	data := make([]byte, size)
	for i := range data {
		data[i] = 'x'
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("writeLogFile: %v", err)
	}
}

func TestRotateNotNeeded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	writeLogFile(t, path, 100)

	opts := DefaultRotateOptions() // MaxBytes = 1 MiB
	rotated, err := Rotate(path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Fatal("expected no rotation for small file")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("original file should still exist: %v", err)
	}
}

func TestRotatePerformed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	writeLogFile(t, path, 200)

	opts := RotateOptions{MaxBytes: 100, MaxBackups: 0}
	rotated, err := Rotate(path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Fatal("expected rotation to occur")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("original file should have been renamed away")
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "audit.*.log"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(matches))
	}
}

func TestRotatePrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	// Create pre-existing fake backups.
	for _, ts := range []string{"20240101T000000Z", "20240102T000000Z", "20240103T000000Z"} {
		fakePath := filepath.Join(dir, "audit."+ts+".log")
		if err := os.WriteFile(fakePath, []byte("old"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	writeLogFile(t, path, 200)
	opts := RotateOptions{MaxBytes: 100, MaxBackups: 2}
	rotated, err := Rotate(path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Fatal("expected rotation")
	}

	// After rotation we added 1 backup → 4 total; MaxBackups=2 should prune to 2.
	matches, _ := filepath.Glob(filepath.Join(dir, "audit.*.log"))
	if len(matches) != 2 {
		t.Fatalf("expected 2 backups after pruning, got %d: %v", len(matches), matches)
	}
}

func TestRotateNonExistentFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.log")
	rotated, err := Rotate(path, DefaultRotateOptions())
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if rotated {
		t.Fatal("should not rotate a non-existent file")
	}
}
