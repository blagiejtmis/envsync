package sync

import (
	"testing"

	"github.com/yourorg/envsync/internal/env"
)

var prefixPushEntries = []env.Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "APP_NAME", Value: "myapp"},
}

func TestApplyPrefixPushNoOpts(t *testing.T) {
	opts := DefaultPrefixPushOptions()
	got := applyPrefixPush(prefixPushEntries, opts)
	if len(got) != len(prefixPushEntries) {
		t.Fatalf("expected %d entries, got %d", len(prefixPushEntries), len(got))
	}
	for i, e := range got {
		if e.Key != prefixPushEntries[i].Key {
			t.Errorf("entry %d: key changed from %q to %q", i, prefixPushEntries[i].Key, e.Key)
		}
	}
}

func TestApplyPrefixPushAddsPrefix(t *testing.T) {
	opts := DefaultPrefixPushOptions()
	opts.AddPrefix = "PROD"
	got := applyPrefixPush(prefixPushEntries, opts)
	for _, e := range got {
		if len(e.Key) < 5 || e.Key[:5] != "PROD_" {
			t.Errorf("expected key to start with PROD_, got %q", e.Key)
		}
	}
}

func TestApplyPrefixPushStripsPrefix(t *testing.T) {
	opts := DefaultPrefixPushOptions()
	opts.StripPrefix = "DB"
	got := applyPrefixPush(prefixPushEntries, opts)
	// DB_HOST -> HOST, DB_PORT -> PORT, APP_NAME kept (non-strict)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if got[0].Key != "HOST" {
		t.Errorf("expected HOST, got %q", got[0].Key)
	}
	if got[1].Key != "PORT" {
		t.Errorf("expected PORT, got %q", got[1].Key)
	}
}

func TestApplyPrefixPushStripStrictOmitsNonMatching(t *testing.T) {
	opts := DefaultPrefixPushOptions()
	opts.StripPrefix = "DB"
	opts.StripStrict = true
	got := applyPrefixPush(prefixPushEntries, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries after strict strip, got %d", len(got))
	}
}

func TestApplyPrefixPushDoesNotMutate(t *testing.T) {
	opts := DefaultPrefixPushOptions()
	opts.AddPrefix = "X"
	origKey := prefixPushEntries[0].Key
	_ = applyPrefixPush(prefixPushEntries, opts)
	if prefixPushEntries[0].Key != origKey {
		t.Errorf("original entry mutated")
	}
}

func TestApplyPrefixPushStripThenAdd(t *testing.T) {
	opts := DefaultPrefixPushOptions()
	opts.StripPrefix = "DB"
	opts.StripStrict = true
	opts.AddPrefix = "PROD"
	got := applyPrefixPush(prefixPushEntries, opts)
	// DB_HOST -> HOST -> PROD_HOST
	// DB_PORT -> PORT -> PROD_PORT
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Key != "PROD_HOST" {
		t.Errorf("expected PROD_HOST, got %q", got[0].Key)
	}
	if got[1].Key != "PROD_PORT" {
		t.Errorf("expected PROD_PORT, got %q", got[1].Key)
	}
}
