package sync

import (
	"testing"

	"github.com/yourorg/envsync/internal/env"
)

func tagPushEntries() []env.Entry {
	return []env.Entry{
		{Key: "# @tag:db DB_HOST,DB_PORT"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_KEY", Value: "abc"},
	}
}

func TestApplyTagPushNoOpts(t *testing.T) {
	entries := tagPushEntries()
	out, err := applyTagPush(entries, DefaultTagPushOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(out))
	}
}

func TestApplyTagPushFiltersToTag(t *testing.T) {
	out, err := applyTagPush(tagPushEntries(), TagPushOptions{TagName: "db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Key != "DB_HOST" && e.Key != "DB_PORT" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestApplyTagPushUnknownTagError(t *testing.T) {
	_, err := applyTagPush(tagPushEntries(), TagPushOptions{TagName: "nope"})
	if err == nil {
		t.Error("expected error for unknown tag, got nil")
	}
}

func TestApplyTagPushDoesNotMutate(t *testing.T) {
	original := tagPushEntries()
	copy := make([]env.Entry, len(original))
	for i, e := range original {
		copy[i] = e
	}
	_, _ = applyTagPush(original, TagPushOptions{TagName: "db"})
	for i, e := range original {
		if e != copy[i] {
			t.Errorf("entry %d mutated", i)
		}
	}
}
