package env

import (
	"testing"
)

func promoteEntries(kvs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestPromoteAllKeys(t *testing.T) {
	src := promoteEntries("DB_HOST", "staging-db", "API_KEY", "stg-secret")
	dst := promoteEntries("APP_ENV", "production")
	opts := DefaultPromoteOptions()
	opts.AllowOverwrite = true

	result, promoted, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(promoted) != 2 {
		t.Fatalf("expected 2 promoted keys, got %d", len(promoted))
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestPromoteNoOverwrite(t *testing.T) {
	src := promoteEntries("DB_HOST", "staging-db")
	dst := promoteEntries("DB_HOST", "prod-db")
	opts := DefaultPromoteOptions() // AllowOverwrite = false

	result, promoted, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(promoted) != 0 {
		t.Fatalf("expected 0 promoted keys, got %d", len(promoted))
	}
	if result[0].Value != "prod-db" {
		t.Errorf("expected destination value to be preserved, got %q", result[0].Value)
	}
}

func TestPromoteSpecificKeys(t *testing.T) {
	src := promoteEntries("DB_HOST", "staging-db", "API_KEY", "stg-secret", "LOG_LEVEL", "debug")
	dst := promoteEntries("APP_ENV", "production")
	opts := DefaultPromoteOptions()
	opts.Keys = []string{"LOG_LEVEL"}
	opts.AllowOverwrite = true

	_, promoted, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(promoted) != 1 || promoted[0] != "LOG_LEVEL" {
		t.Errorf("expected only LOG_LEVEL promoted, got %v", promoted)
	}
}

func TestPromoteMissingKeyError(t *testing.T) {
	src := promoteEntries("DB_HOST", "staging-db")
	dst := promoteEntries("APP_ENV", "production")
	opts := DefaultPromoteOptions()
	opts.Keys = []string{"MISSING_KEY"}

	_, _, err := Promote(src, dst, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestPromoteDoesNotMutateSrc(t *testing.T) {
	src := promoteEntries("DB_HOST", "staging-db")
	dst := promoteEntries("APP_ENV", "production")
	opts := DefaultPromoteOptions()
	opts.AllowOverwrite = true

	origLen := len(src)
	_, _, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(src) != origLen {
		t.Errorf("Promote mutated src slice")
	}
}
