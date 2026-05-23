package env

import (
	"testing"
)

var groupEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "APP_DEBUG", Value: "false"},
	{Key: "STANDALONE", Value: "yes"},
}

func TestGroupByPrefixKeepsPrefix(t *testing.T) {
	opts := DefaultGroupOptions()
	groups := GroupByPrefix(groupEntries, opts)

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}

	// Order should match insertion order.
	if groups[0].Name != "DB" {
		t.Errorf("expected first group DB, got %s", groups[0].Name)
	}
	if len(groups[0].Entries) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groups[0].Entries))
	}
	// Keys retained with prefix.
	if groups[0].Entries[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", groups[0].Entries[0].Key)
	}
}

func TestGroupByPrefixStripsPrefix(t *testing.T) {
	opts := DefaultGroupOptions()
	opts.KeepPrefix = false
	groups := GroupByPrefix(groupEntries, opts)

	var dbGroup *Group
	for i := range groups {
		if groups[i].Name == "DB" {
			dbGroup = &groups[i]
			break
		}
	}
	if dbGroup == nil {
		t.Fatal("DB group not found")
	}
	for _, e := range dbGroup.Entries {
		if e.Key == "DB_HOST" || e.Key == "DB_PORT" {
			t.Errorf("prefix should be stripped, got key %s", e.Key)
		}
	}
	keys := map[string]bool{}
	for _, e := range dbGroup.Entries {
		keys[e.Key] = true
	}
	if !keys["HOST"] || !keys["PORT"] {
		t.Errorf("expected stripped keys HOST and PORT, got %v", keys)
	}
}

func TestGroupByPrefixNoSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "NOSEP", Value: "1"},
		{Key: "ALSONONE", Value: "2"},
	}
	opts := DefaultGroupOptions()
	groups := GroupByPrefix(entries, opts)

	// Each key becomes its own group (no underscore).
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupNames(t *testing.T) {
	names := GroupNames(groupEntries, "_")
	expected := []string{"APP", "DB", "STANDALONE"}
	if len(names) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, names)
	}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], n)
		}
	}
}

func TestGroupByPrefixPreservesValues(t *testing.T) {
	opts := DefaultGroupOptions()
	groups := GroupByPrefix(groupEntries, opts)
	for _, g := range groups {
		for _, e := range g.Entries {
			if e.Value == "" {
				t.Errorf("entry %s lost its value", e.Key)
			}
		}
	}
}
