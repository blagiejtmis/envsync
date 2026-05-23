package env

import (
	"testing"
)

var prefixEntries = []Entry{
	{Key: "HOST", Value: "localhost"},
	{Key: "PORT", Value: "5432"},
	{Key: "NAME", Value: "mydb"},
}

func TestAddPrefixBasic(t *testing.T) {
	opts := DefaultPrefixOptions()
	got := AddPrefix(prefixEntries, "DB", opts)
	expected := []string{"DB_HOST", "DB_PORT", "DB_NAME"}
	for i, e := range got {
		if e.Key != expected[i] {
			t.Errorf("entry %d: got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestAddPrefixPreservesValues(t *testing.T) {
	opts := DefaultPrefixOptions()
	got := AddPrefix(prefixEntries, "DB", opts)
	if got[0].Value != "localhost" {
		t.Errorf("expected value %q, got %q", "localhost", got[0].Value)
	}
}

func TestAddPrefixEmptyPrefixIsNoop(t *testing.T) {
	opts := DefaultPrefixOptions()
	got := AddPrefix(prefixEntries, "", opts)
	if len(got) != len(prefixEntries) {
		t.Fatalf("expected %d entries, got %d", len(prefixEntries), len(got))
	}
	for i, e := range got {
		if e.Key != prefixEntries[i].Key {
			t.Errorf("entry %d: key changed from %q to %q", i, prefixEntries[i].Key, e.Key)
		}
	}
}

func TestAddPrefixDoesNotMutateOriginal(t *testing.T) {
	opts := DefaultPrefixOptions()
	origKey := prefixEntries[0].Key
	_ = AddPrefix(prefixEntries, "X", opts)
	if prefixEntries[0].Key != origKey {
		t.Errorf("original entry mutated")
	}
}

func TestStripPrefixBasic(t *testing.T) {
	opts := DefaultPrefixOptions()
	prefixed := AddPrefix(prefixEntries, "DB", opts)
	got := StripPrefix(prefixed, "DB", opts, true)
	for i, e := range got {
		if e.Key != prefixEntries[i].Key {
			t.Errorf("entry %d: got key %q, want %q", i, e.Key, prefixEntries[i].Key)
		}
	}
}

func TestStripPrefixStrictOmitsNonMatching(t *testing.T) {
	opts := DefaultPrefixOptions()
	mixed := []Entry{
		{Key: "DB_HOST", Value: "a"},
		{Key: "APP_NAME", Value: "b"},
	}
	got := StripPrefix(mixed, "DB", opts, true)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Key != "HOST" {
		t.Errorf("expected key HOST, got %q", got[0].Key)
	}
}

func TestStripPrefixNonStrictKeepsNonMatching(t *testing.T) {
	opts := DefaultPrefixOptions()
	mixed := []Entry{
		{Key: "DB_HOST", Value: "a"},
		{Key: "APP_NAME", Value: "b"},
	}
	got := StripPrefix(mixed, "DB", opts, false)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestCommonPrefix(t *testing.T) {
	opts := DefaultPrefixOptions()
	entries := []Entry{
		{Key: "DB_HOST"},
		{Key: "DB_PORT"},
		{Key: "DB_NAME"},
	}
	got := CommonPrefix(entries, opts)
	if got != "DB" {
		t.Errorf("expected common prefix %q, got %q", "DB", got)
	}
}

func TestCommonPrefixNoCommon(t *testing.T) {
	opts := DefaultPrefixOptions()
	entries := []Entry{
		{Key: "DB_HOST"},
		{Key: "APP_NAME"},
	}
	got := CommonPrefix(entries, opts)
	if got != "" {
		t.Errorf("expected empty common prefix, got %q", got)
	}
}
