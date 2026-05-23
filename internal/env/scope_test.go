package env

import (
	"testing"
)

func scopeEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost", Comment: "# scope=development"},
		{Key: "DB_HOST_PROD", Value: "db.prod.example.com", Comment: "# scope=production"},
		{Key: "API_KEY", Value: "secret", Comment: "# scope=staging"},
		{Key: "APP_NAME", Value: "envsync", Comment: ""},
	}
}

func TestScopeFilterNoTargets(t *testing.T) {
	entries := scopeEntries()
	out, err := ScopeFilter(entries, nil, DefaultScopeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(out))
	}
}

func TestScopeFilterSingleScope(t *testing.T) {
	out, err := ScopeFilter(scopeEntries(), []string{"development"}, DefaultScopeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// expects DB_HOST (development) + APP_NAME (no scope)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out[0].Key)
	}
	if out[1].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", out[1].Key)
	}
}

func TestScopeFilterMultipleScopes(t *testing.T) {
	out, err := ScopeFilter(scopeEntries(), []string{"development", "staging"}, DefaultScopeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DB_HOST, API_KEY, APP_NAME
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestScopeFilterUnknownScopeError(t *testing.T) {
	_, err := ScopeFilter(scopeEntries(), []string{"canary"}, DefaultScopeOptions())
	if err == nil {
		t.Fatal("expected error for unknown scope, got nil")
	}
}

func TestScopeFilterDoesNotMutateOriginal(t *testing.T) {
	original := scopeEntries()
	copy := make([]Entry, len(original))
	for i, e := range original {
		copy[i] = e
	}
	_, _ = ScopeFilter(original, []string{"production"}, DefaultScopeOptions())
	for i, e := range original {
		if e != copy[i] {
			t.Errorf("entry %d was mutated", i)
		}
	}
}

func TestExtractAnnotationFound(t *testing.T) {
	got := extractAnnotation("# scope=production env=cloud", "scope")
	if got != "production" {
		t.Errorf("expected 'production', got %q", got)
	}
}

func TestExtractAnnotationNotFound(t *testing.T) {
	got := extractAnnotation("# some random comment", "scope")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestExtractAnnotationEmptyComment(t *testing.T) {
	got := extractAnnotation("", "scope")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
