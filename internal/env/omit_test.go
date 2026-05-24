package env

import (
	"testing"
)

var omitEntries = []Entry{
	{Key: "APP_HOST", Value: "localhost"},
	{Key: "APP_SECRET", Value: "s3cr3t"},
	{Key: "DB_PASSWORD", Value: "pass"},
	{Key: "DEBUG", Value: ""},
	{Key: "LOG_LEVEL", Value: "info"},
	{Key: "INTERNAL_FLAG", Value: "1"},
}

func TestOmitByKey(t *testing.T) {
	result := Omit(omitEntries, OmitOptions{Keys: []string{"DEBUG", "LOG_LEVEL"}})
	for _, e := range result {
		if e.Key == "DEBUG" || e.Key == "LOG_LEVEL" {
			t.Errorf("expected key %q to be omitted", e.Key)
		}
	}
	if len(result) != 4 {
		t.Errorf("expected 4 entries, got %d", len(result))
	}
}

func TestOmitBlankValues(t *testing.T) {
	result := Omit(omitEntries, OmitOptions{Blank: true})
	for _, e := range result {
		if e.Value == "" {
			t.Errorf("expected blank entry %q to be omitted", e.Key)
		}
	}
}

func TestOmitSensitive(t *testing.T) {
	result := Omit(omitEntries, OmitOptions{Sensitive: true})
	for _, e := range result {
		if IsSensitive(e.Key) {
			t.Errorf("expected sensitive key %q to be omitted", e.Key)
		}
	}
}

func TestOmitByPrefix(t *testing.T) {
	result := Omit(omitEntries, OmitOptions{Prefixes: []string{"APP_", "INTERNAL_"}})
	for _, e := range result {
		if e.Key == "APP_HOST" || e.Key == "APP_SECRET" || e.Key == "INTERNAL_FLAG" {
			t.Errorf("expected key %q to be omitted by prefix", e.Key)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestOmitDoesNotMutateOriginal(t *testing.T) {
	original := make([]Entry, len(omitEntries))
	copy(original, omitEntries)
	Omit(omitEntries, OmitOptions{Keys: []string{"DEBUG"}, Blank: true, Sensitive: true})
	for i, e := range omitEntries {
		if e != original[i] {
			t.Errorf("original mutated at index %d", i)
		}
	}
}

func TestOmittedKeysReturnsRemovedKeys(t *testing.T) {
	keys := OmittedKeys(omitEntries, OmitOptions{Blank: true, Keys: []string{"LOG_LEVEL"}})
	expected := map[string]bool{"DEBUG": true, "LOG_LEVEL": true}
	if len(keys) != 2 {
		t.Fatalf("expected 2 omitted keys, got %d: %v", len(keys), keys)
	}
	for _, k := range keys {
		if !expected[k] {
			t.Errorf("unexpected omitted key %q", k)
		}
	}
}

func TestOmitEmptyOptions(t *testing.T) {
	result := Omit(omitEntries, DefaultOmitOptions())
	if len(result) != len(omitEntries) {
		t.Errorf("expected all %d entries, got %d", len(omitEntries), len(result))
	}
}
