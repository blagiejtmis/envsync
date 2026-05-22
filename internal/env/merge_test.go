package env

import (
	"testing"
)

func entries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMergeNoConflicts(t *testing.T) {
	local := entries("FOO", "1", "BAR", "2")
	remote := entries("BAZ", "3")

	res := Merge(local, remote, PreferLocal)

	if len(res.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", res.Conflicts)
	}
	if len(res.Merged) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Merged))
	}
}

func TestMergePreferLocal(t *testing.T) {
	local := entries("FOO", "local")
	remote := entries("FOO", "remote")

	res := Merge(local, remote, PreferLocal)

	if len(res.Conflicts) != 1 || res.Conflicts[0] != "FOO" {
		t.Fatalf("expected conflict on FOO, got %v", res.Conflicts)
	}
	if res.Merged[0].Value != "local" {
		t.Errorf("expected local value, got %q", res.Merged[0].Value)
	}
}

func TestMergePreferRemote(t *testing.T) {
	local := entries("FOO", "local")
	remote := entries("FOO", "remote")

	res := Merge(local, remote, PreferRemote)

	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %v", res.Conflicts)
	}
	if res.Merged[0].Value != "remote" {
		t.Errorf("expected remote value, got %q", res.Merged[0].Value)
	}
}

func TestMergeRemoteOnlyKeysAdded(t *testing.T) {
	local := entries("A", "1")
	remote := entries("A", "1", "B", "2", "C", "3")

	res := Merge(local, remote, PreferLocal)

	if len(res.Conflicts) != 0 {
		t.Fatalf("unexpected conflicts: %v", res.Conflicts)
	}
	if len(res.Merged) != 3 {
		t.Fatalf("expected 3 merged entries, got %d", len(res.Merged))
	}
}

func TestMergePreservesOrder(t *testing.T) {
	local := entries("Z", "26", "A", "1")
	remote := entries("M", "13")

	res := Merge(local, remote, PreferLocal)

	if res.Merged[0].Key != "Z" || res.Merged[1].Key != "A" || res.Merged[2].Key != "M" {
		t.Errorf("order not preserved: %v", res.Merged)
	}
}

func TestMergeEmptyLocal(t *testing.T) {
	remote := entries("X", "42")
	res := Merge(nil, remote, PreferLocal)

	if len(res.Merged) != 1 || res.Merged[0].Key != "X" {
		t.Errorf("expected remote entry in merged, got %v", res.Merged)
	}
}
