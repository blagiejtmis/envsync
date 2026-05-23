package env

import (
	"testing"
)

var filterEntries = []Entry{
	{Key: "APP_HOST", Value: "localhost"},
	{Key: "APP_PORT", Value: "8080"},
	{Key: "DB_HOST", Value: "db"},
	{Key: "DB_PASS", Value: "secret"},
	{Key: "DEBUG", Value: "true"},
}

func TestFilterByPrefix(t *testing.T) {
	got := FilterByPrefix(filterEntries, "APP_")
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key != "APP_HOST" && e.Key != "APP_PORT" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestFilterByKeys(t *testing.T) {
	got := FilterByKeys(filterEntries, "DB_HOST", "DEBUG")
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilterExcludeKeys(t *testing.T) {
	got := Reject(filterEntries, "DB_PASS", "DEBUG")
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key == "DB_PASS" || e.Key == "DEBUG" {
			t.Errorf("excluded key %q still present", e.Key)
		}
	}
}

func TestFilterCombinedPrefixAndExclude(t *testing.T) {
	got := Filter(filterEntries, FilterOptions{
		Prefix:      "DB_",
		ExcludeKeys: []string{"DB_PASS"},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", got[0].Key)
	}
}

func TestFilterEmptyOptions(t *testing.T) {
	got := Filter(filterEntries, FilterOptions{})
	if len(got) != len(filterEntries) {
		t.Errorf("expected all %d entries, got %d", len(filterEntries), len(got))
	}
}

func TestFilterDoesNotMutateOriginal(t *testing.T) {
	copy := make([]Entry, len(filterEntries))
	for i, e := range filterEntries {
		copy[i] = e
	}
	FilterByPrefix(filterEntries, "APP_")
	for i, e := range filterEntries {
		if e != copy[i] {
			t.Errorf("original mutated at index %d", i)
		}
	}
}
