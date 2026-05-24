package env

import (
	"testing"
)

var setBase = []Entry{
	{Key: "HOST", Value: "localhost"},
	{Key: "PORT", Value: "5432"},
	{Key: "DEBUG", Value: "false"},
}

func TestSetOverwritesExistingKey(t *testing.T) {
	out, err := Set(setBase, map[string]string{"PORT": "9999"}, DefaultSetOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := mustGet(out, "PORT"); v != "9999" {
		t.Errorf("expected PORT=9999, got %q", v)
	}
}

func TestSetAddsNewKey(t *testing.T) {
	out, err := Set(setBase, map[string]string{"TIMEOUT": "30"}, DefaultSetOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := mustGet(out, "TIMEOUT"); v != "30" {
		t.Errorf("expected TIMEOUT=30, got %q", v)
	}
}

func TestSetDoesNotMutateOriginal(t *testing.T) {
	_, err := Set(setBase, map[string]string{"HOST": "changed"}, DefaultSetOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if setBase[0].Value != "localhost" {
		t.Error("original slice was mutated")
	}
}

func TestSetDisallowNew(t *testing.T) {
	opts := DefaultSetOptions()
	opts.AllowNew = false
	_, err := Set(setBase, map[string]string{"NEW_KEY": "val"}, opts)
	if err == nil {
		t.Error("expected error when AllowNew=false and key is new")
	}
}

func TestSetDisallowOverwrite(t *testing.T) {
	opts := DefaultSetOptions()
	opts.OverwriteExisting = false
	_, err := Set(setBase, map[string]string{"HOST": "newhost"}, opts)
	if err == nil {
		t.Error("expected error when OverwriteExisting=false and key exists")
	}
}

func TestSetEmptyKeyError(t *testing.T) {
	_, err := Set(setBase, map[string]string{"": "val"}, DefaultSetOptions())
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestSetMultiplePairs(t *testing.T) {
	out, err := Set(setBase, map[string]string{"HOST": "remotehost", "NEWVAR": "1"}, DefaultSetOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := mustGet(out, "HOST"); v != "remotehost" {
		t.Errorf("HOST: got %q, want remotehost", v)
	}
	if v := mustGet(out, "NEWVAR"); v != "1" {
		t.Errorf("NEWVAR: got %q, want 1", v)
	}
}

// mustGet is a test helper that returns the value for key or empty string.
func mustGet(entries []Entry, key string) string {
	for _, e := range entries {
		if e.Key == key {
			return e.Value
		}
	}
	return ""
}
