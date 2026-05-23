package env

import (
	"testing"
)

func TestFlattenDotNotation(t *testing.T) {
	pairs := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	opts := DefaultFlattenOptions()
	entries, err := Flatten(pairs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	m := entriesToMap(entries)
	if m["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want %q", m["DB_HOST"], "localhost")
	}
	if m["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT = %q, want %q", m["DB_PORT"], "5432")
	}
}

func TestFlattenSlashNotation(t *testing.T) {
	pairs := map[string]string{"aws/region": "us-east-1"}
	opts := DefaultFlattenOptions()
	entries, err := Flatten(pairs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "AWS_REGION" {
		t.Errorf("key = %q, want AWS_REGION", entries[0].Key)
	}
}

func TestFlattenWithPrefix(t *testing.T) {
	pairs := map[string]string{"host": "localhost"}
	opts := DefaultFlattenOptions()
	opts.Prefix = "APP"
	entries, err := Flatten(pairs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "APP_HOST" {
		t.Errorf("key = %q, want APP_HOST", entries[0].Key)
	}
}

func TestFlattenPreservesLowercaseWhenDisabled(t *testing.T) {
	pairs := map[string]string{"my.key": "value"}
	opts := DefaultFlattenOptions()
	opts.Uppercase = false
	entries, err := Flatten(pairs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "my_key" {
		t.Errorf("key = %q, want my_key", entries[0].Key)
	}
}

func TestFlattenDuplicateKeyError(t *testing.T) {
	// "a.b" and "a/b" both flatten to "A_B" — should error.
	pairs := map[string]string{
		"a.b": "first",
		"a/b": "second",
	}
	opts := DefaultFlattenOptions()
	_, err := Flatten(pairs, opts)
	if err == nil {
		t.Fatal("expected duplicate key error, got nil")
	}
}

func TestFlattenSortedOutput(t *testing.T) {
	pairs := map[string]string{
		"z.key": "last",
		"a.key": "first",
		"m.key": "middle",
	}
	opts := DefaultFlattenOptions()
	entries, err := Flatten(pairs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "A_KEY" || entries[1].Key != "M_KEY" || entries[2].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v", entries)
	}
}

// entriesToMap is a helper used across tests.
func entriesToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
