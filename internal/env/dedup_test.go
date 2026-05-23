package env

import (
	"testing"
)

var dedupEntries = []Entry{
	{Key: "HOST", Value: "localhost"},
	{Key: "PORT", Value: "5432"},
	{Key: "HOST", Value: "remotehost"},
	{Key: "DEBUG", Value: "true"},
	{Key: "PORT", Value: "9999"},
}

func TestDedupKeepLast(t *testing.T) {
	result := Dedup(dedupEntries, DefaultDedupOptions())
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Value != "remotehost" {
		t.Errorf("expected HOST=remotehost, got %s", result[0].Value)
	}
	if result[1].Value != "9999" {
		t.Errorf("expected PORT=9999, got %s", result[1].Value)
	}
}

func TestDedupKeepFirst(t *testing.T) {
	opts := DedupOptions{Strategy: DedupKeepFirst}
	result := Dedup(dedupEntries, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Value != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", result[0].Value)
	}
	if result[1].Value != "5432" {
		t.Errorf("expected PORT=5432, got %s", result[1].Value)
	}
}

func TestDedupReportsRemovedKeys(t *testing.T) {
	var reported []string
	opts := DedupOptions{
		Strategy: DedupKeepLast,
		Report:   func(k string) { reported = append(reported, k) },
	}
	Dedup(dedupEntries, opts)
	if len(reported) != 2 {
		t.Fatalf("expected 2 reported removals, got %d", len(reported))
	}
}

func TestDedupNoDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	result := Dedup(entries, DefaultDedupOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestDedupEmptyInput(t *testing.T) {
	result := Dedup(nil, DefaultDedupOptions())
	if result != nil {
		t.Errorf("expected nil result for nil input")
	}
}

func TestDupKeys(t *testing.T) {
	dups := DupKeys(dedupEntries)
	if len(dups) != 2 {
		t.Fatalf("expected 2 dup keys, got %d", len(dups))
	}
	if dups[0] != "HOST" || dups[1] != "PORT" {
		t.Errorf("unexpected dup keys: %v", dups)
	}
}

func TestDupKeysNoDups(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "1"}, {Key: "Y", Value: "2"}}
	if got := DupKeys(entries); len(got) != 0 {
		t.Errorf("expected no dups, got %v", got)
	}
}
