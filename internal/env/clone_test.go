package env

import (
	"testing"
)

func TestCloneIsDeepCopy(t *testing.T) {
	orig := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	got := Clone(orig)
	if len(got) != len(orig) {
		t.Fatalf("expected len %d, got %d", len(orig), len(got))
	}
	got[0].Value = "mutated"
	if orig[0].Value == "mutated" {
		t.Error("Clone mutated original slice")
	}
}

func TestCloneNil(t *testing.T) {
	if Clone(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestCloneMapIsDeepCopy(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	got := CloneMap(orig)
	got["A"] = "mutated"
	if orig["A"] == "mutated" {
		t.Error("CloneMap mutated original map")
	}
}

func TestCloneMapNil(t *testing.T) {
	if CloneMap(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestUniqueKeepsLast(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "first"},
		{Key: "B", Value: "b"},
		{Key: "A", Value: "last"},
	}
	got := Unique(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Value != "last" {
		t.Errorf("expected last value 'last', got %q", got[0].Value)
	}
}

func TestUniqueNoDuplicates(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "1"}, {Key: "Y", Value: "2"}}
	got := Unique(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestReorderKnownKeysFirst(t *testing.T) {
	entries := []Entry{
		{Key: "C", Value: "3"},
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	got := Reorder(entries, []string{"B", "A"})
	if got[0].Key != "B" || got[1].Key != "A" || got[2].Key != "C" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestReorderMissingKeyIgnored(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}}
	got := Reorder(entries, []string{"Z", "A"})
	if len(got) != 1 || got[0].Key != "A" {
		t.Errorf("unexpected result: %v", got)
	}
}
