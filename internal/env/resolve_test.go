package env

import (
	"testing"
)

func resolveEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_URL", Value: "postgres://localhost/dev"},
	}
}

func TestResolveAppliesOverrides(t *testing.T) {
	overrides := map[string]string{
		"APP_HOST": "prod.example.com",
		"APP_PORT": "443",
	}
	opts := DefaultResolveOptions()
	opts.Overrides = overrides

	out, err := Resolve(resolveEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["APP_HOST"] != "prod.example.com" {
		t.Errorf("APP_HOST = %q, want prod.example.com", m["APP_HOST"])
	}
	if m["APP_PORT"] != "443" {
		t.Errorf("APP_PORT = %q, want 443", m["APP_PORT"])
	}
	if m["DB_URL"] != "postgres://localhost/dev" {
		t.Errorf("DB_URL should be unchanged")
	}
}

func TestResolveFallbackToEntry(t *testing.T) {
	opts := DefaultResolveOptions()
	opts.Overrides = map[string]string{"APP_HOST": "newhost"}

	out, err := Resolve(resolveEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestResolveNoFallbackError(t *testing.T) {
	opts := DefaultResolveOptions()
	opts.FallbackToEntry = false
	opts.Overrides = map[string]string{"APP_HOST": "x"}

	_, err := Resolve(resolveEntries(), opts)
	if err == nil {
		t.Fatal("expected error when FallbackToEntry=false and key missing from overrides")
	}
}

func TestResolveUnknownKeyError(t *testing.T) {
	opts := DefaultResolveOptions()
	opts.IgnoreUnknownKeys = false
	opts.Overrides = map[string]string{"NONEXISTENT": "val"}

	_, err := Resolve(resolveEntries(), opts)
	if err == nil {
		t.Fatal("expected error for unknown override key")
	}
}

func TestResolveIgnoreUnknownKeys(t *testing.T) {
	opts := DefaultResolveOptions()
	opts.IgnoreUnknownKeys = true
	opts.Overrides = map[string]string{"NONEXISTENT": "val"}

	out, err := Resolve(resolveEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestResolveFromString(t *testing.T) {
	raw := "APP_HOST=staging.example.com\nAPP_PORT=9090\n"
	opts := DefaultResolveOptions()

	out, err := ResolveFromString(resolveEntries(), raw, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["APP_HOST"] != "staging.example.com" {
		t.Errorf("APP_HOST = %q, want staging.example.com", m["APP_HOST"])
	}
	if m["APP_PORT"] != "9090" {
		t.Errorf("APP_PORT = %q, want 9090", m["APP_PORT"])
	}
}

func TestResolveFromStringInvalidLine(t *testing.T) {
	raw := "BADLINE\n"
	opts := DefaultResolveOptions()

	_, err := ResolveFromString(resolveEntries(), raw, opts)
	if err == nil {
		t.Fatal("expected error for invalid override line")
	}
}
