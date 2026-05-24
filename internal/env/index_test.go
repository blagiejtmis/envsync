package env

import (
	"testing"
)

var indexEntries = []Entry{
	{Key: "HOST", Value: "localhost"},
	{Key: "PORT", Value: "8080"},
	{Key: "DEBUG", Value: "true"},
	{Key: "SECRET", Value: "abc123"},
}

func TestIndexBuildsMap(t *testing.T) {
	idx := Index(indexEntries)
	if len(idx) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(idx))
	}
	if idx["PORT"].Value != "8080" {
		t.Errorf("expected PORT=8080, got %s", idx["PORT"].Value)
	}
}

func TestLookupFound(t *testing.T) {
	v, ok := Lookup(indexEntries, "HOST")
	if !ok {
		t.Fatal("expected HOST to be found")
	}
	if v != "localhost" {
		t.Errorf("expected localhost, got %s", v)
	}
}

func TestLookupNotFound(t *testing.T) {
	_, ok := Lookup(indexEntries, "MISSING")
	if ok {
		t.Fatal("expected MISSING to not be found")
	}
}

func TestContainsTrue(t *testing.T) {
	if !Contains(indexEntries, "DEBUG") {
		t.Error("expected Contains to return true for DEBUG")
	}
}

func TestContainsFalse(t *testing.T) {
	if Contains(indexEntries, "NOPE") {
		t.Error("expected Contains to return false for NOPE")
	}
}

func TestKeys(t *testing.T) {
	keys := Keys(indexEntries)
	expected := []string{"HOST", "PORT", "DEBUG", "SECRET"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("index %d: expected %s, got %s", i, k, keys[i])
		}
	}
}

func TestValues(t *testing.T) {
	vals := Values(indexEntries)
	if len(vals) != 4 {
		t.Fatalf("expected 4 values, got %d", len(vals))
	}
	if vals[1] != "8080" {
		t.Errorf("expected 8080 at index 1, got %s", vals[1])
	}
}

func TestPickReturnsSubset(t *testing.T) {
	picked := Pick(indexEntries, "HOST", "SECRET")
	if len(picked) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(picked))
	}
	if picked[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", picked[0].Key)
	}
	if picked[1].Key != "SECRET" {
		t.Errorf("expected SECRET, got %s", picked[1].Key)
	}
}

func TestPickMissingKeyOmitted(t *testing.T) {
	picked := Pick(indexEntries, "HOST", "MISSING")
	if len(picked) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(picked))
	}
	if picked[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", picked[0].Key)
	}
}
