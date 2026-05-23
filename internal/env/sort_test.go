package env

import (
	"testing"
)

var sortEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "APP_NAME", Value: "envsync"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "CACHE_URL", Value: "redis://localhost"},
}

func TestSortAlpha(t *testing.T) {
	opts := DefaultSortOptions()
	out := Sort(sortEntries, opts)

	expected := []string{"APP_ENV", "APP_NAME", "CACHE_URL", "DB_HOST", "DB_PORT"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSortAlphaDesc(t *testing.T) {
	opts := DefaultSortOptions()
	opts.Order = SortAlphaDesc
	out := Sort(sortEntries, opts)

	expected := []string{"DB_PORT", "DB_HOST", "CACHE_URL", "APP_NAME", "APP_ENV"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSortByPrefix(t *testing.T) {
	opts := DefaultSortOptions()
	opts.Order = SortByPrefix
	out := Sort(sortEntries, opts)

	// Prefixes: app, app, cache, db, db
	if out[0].Key != "APP_ENV" && out[0].Key != "APP_NAME" {
		t.Errorf("expected APP_ prefix first, got %q", out[0].Key)
	}
	if out[2].Key != "CACHE_URL" {
		t.Errorf("expected CACHE_URL third, got %q", out[2].Key)
	}
	if out[3].Key != "DB_HOST" && out[3].Key != "DB_PORT" {
		t.Errorf("expected DB_ prefix last, got %q", out[3].Key)
	}
}

func TestSortDoesNotMutateOriginal(t *testing.T) {
	original := []Entry{
		{Key: "Z_KEY", Value: "z"},
		{Key: "A_KEY", Value: "a"},
	}
	Sort(original, DefaultSortOptions())

	if original[0].Key != "Z_KEY" {
		t.Error("Sort mutated the original slice")
	}
}

func TestSortCaseSensitive(t *testing.T) {
	entries := []Entry{
		{Key: "b_key", Value: "1"},
		{Key: "A_KEY", Value: "2"},
	}
	opts := DefaultSortOptions()
	opts.CaseSensitive = true
	out := Sort(entries, opts)

	// ASCII: uppercase letters come before lowercase
	if out[0].Key != "A_KEY" {
		t.Errorf("case-sensitive sort: expected A_KEY first, got %q", out[0].Key)
	}
}

func TestSortStable(t *testing.T) {
	entries := []Entry{
		{Key: "SAME", Value: "first"},
		{Key: "SAME", Value: "second"},
	}
	opts := DefaultSortOptions()
	opts.Stable = true
	out := Sort(entries, opts)

	if out[0].Value != "first" {
		t.Error("stable sort should preserve original order for equal keys")
	}
}
