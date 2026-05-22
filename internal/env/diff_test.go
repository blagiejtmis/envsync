package env

import (
	"testing"
)

func TestDiffNoChanges(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1", "B": "2"}
	changes := Diff(old, new)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestDiffAdded(t *testing.T) {
	old := map[string]string{"A": "1"}
	new := map[string]string{"A": "1", "B": "2"}
	changes := Diff(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != ChangeAdded || changes[0].Key != "B" || changes[0].NewValue != "2" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffRemoved(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1"}
	changes := Diff(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != ChangeRemoved || changes[0].Key != "B" || changes[0].OldValue != "2" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffUpdated(t *testing.T) {
	old := map[string]string{"A": "old"}
	new := map[string]string{"A": "new"}
	changes := Diff(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != ChangeUpdated || c.Key != "A" || c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiffSortedByKey(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"Z": "1", "A": "2", "M": "3"}
	changes := Diff(old, new)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}
	if changes[0].Key != "A" || changes[1].Key != "M" || changes[2].Key != "Z" {
		t.Errorf("changes not sorted: %v", changes)
	}
}

func TestHasChanges(t *testing.T) {
	if HasChanges(map[string]string{"A": "1"}, map[string]string{"A": "1"}) {
		t.Error("expected no changes")
	}
	if !HasChanges(map[string]string{"A": "1"}, map[string]string{"A": "2"}) {
		t.Error("expected changes")
	}
}
