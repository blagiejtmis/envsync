package env

import (
	"testing"
)

func renameEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "APP_SECRET", Value: "s3cr3t"},
	}
}

func TestRenameBasic(t *testing.T) {
	got, err := Rename(renameEntries(), "APP_PORT", "APP_LISTEN_PORT", DefaultRenameOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[1].Key != "APP_LISTEN_PORT" {
		t.Errorf("expected APP_LISTEN_PORT at index 1, got %q", got[1].Key)
	}
	if got[1].Value != "8080" {
		t.Errorf("expected value 8080, got %q", got[1].Value)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 entries, got %d", len(got))
	}
}

func TestRenameDoesNotMutateOriginal(t *testing.T) {
	orig := renameEntries()
	_, err := Rename(orig, "APP_HOST", "HOST", DefaultRenameOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig[0].Key != "APP_HOST" {
		t.Error("Rename mutated the original slice")
	}
}

func TestRenameMissingKeyError(t *testing.T) {
	opts := DefaultRenameOptions()
	_, err := Rename(renameEntries(), "NONEXISTENT", "NEW_KEY", opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRenameMissingKeyNoError(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.FailOnMissing = false
	got, err := Rename(renameEntries(), "NONEXISTENT", "NEW_KEY", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 entries unchanged, got %d", len(got))
	}
}

func TestRenameCollisionError(t *testing.T) {
	opts := DefaultRenameOptions()
	_, err := Rename(renameEntries(), "APP_HOST", "APP_PORT", opts)
	if err == nil {
		t.Fatal("expected collision error, got nil")
	}
}

func TestRenameCollisionOverwrite(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.FailOnCollision = false
	got, err := Rename(renameEntries(), "APP_HOST", "APP_PORT", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Original APP_PORT entry is dropped; APP_HOST is renamed in place.
	if len(got) != 2 {
		t.Errorf("expected 2 entries after overwrite, got %d", len(got))
	}
	if got[0].Key != "APP_PORT" || got[0].Value != "localhost" {
		t.Errorf("unexpected first entry: %+v", got[0])
	}
}
