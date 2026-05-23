package env

import (
	"testing"
)

func chainEntries(kvs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestChainFirstWins(t *testing.T) {
	a := chainEntries("FOO", "from_a", "BAR", "bar_a")
	b := chainEntries("FOO", "from_b", "BAZ", "baz_b")

	result, err := Chain(DefaultChainOptions(), a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Value != "from_a" {
		t.Errorf("FOO should be from_a, got %q", result[0].Value)
	}
}

func TestChainLastWins(t *testing.T) {
	a := chainEntries("FOO", "from_a", "BAR", "bar_a")
	b := chainEntries("FOO", "from_b", "BAZ", "baz_b")

	opts := ChainOptions{LastWins: true}
	result, err := Chain(opts, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key == "FOO" && e.Value != "from_b" {
			t.Errorf("FOO should be from_b (last-wins), got %q", e.Value)
		}
	}
}

func TestChainNoSources(t *testing.T) {
	result, err := Chain(DefaultChainOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for no sources")
	}
}

func TestChainEmptyKeyError(t *testing.T) {
	a := []Entry{{Key: "", Value: "oops"}}
	_, err := Chain(DefaultChainOptions(), a)
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestChainPreservesOrder(t *testing.T) {
	a := chainEntries("C", "c", "A", "a")
	b := chainEntries("B", "b", "A", "override")

	result, err := Chain(DefaultChainOptions(), a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Key != "C" || result[1].Key != "A" || result[2].Key != "B" {
		t.Errorf("unexpected order: %v", result)
	}
}
