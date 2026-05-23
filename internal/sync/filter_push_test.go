package sync

import (
	"testing"

	"github.com/yourorg/envsync/internal/env"
)

var pushEntries = []env.Entry{
	{Key: "APP_HOST", Value: "localhost"},
	{Key: "APP_PORT", Value: "8080"},
	{Key: "DB_PASS", Value: "secret"},
	{Key: "DEBUG", Value: "true"},
}

func TestApplyFilterPushNoOpts(t *testing.T) {
	got := applyFilterPush(pushEntries, DefaultFilterPushOptions())
	if len(got) != len(pushEntries) {
		t.Errorf("expected %d entries, got %d", len(pushEntries), len(got))
	}
}

func TestApplyFilterPushPrefix(t *testing.T) {
	got := applyFilterPush(pushEntries, FilterPushOptions{Prefix: "APP_"})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key != "APP_HOST" && e.Key != "APP_PORT" {
			t.Errorf("unexpected key %q in result", e.Key)
		}
	}
}

func TestApplyFilterPushExplicitKeys(t *testing.T) {
	got := applyFilterPush(pushEntries, FilterPushOptions{Keys: []string{"DEBUG", "DB_PASS"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestApplyFilterPushExclude(t *testing.T) {
	got := applyFilterPush(pushEntries, FilterPushOptions{ExcludeKeys: []string{"DB_PASS"}})
	for _, e := range got {
		if e.Key == "DB_PASS" {
			t.Error("DB_PASS should have been excluded")
		}
	}
	if len(got) != 3 {
		t.Errorf("expected 3 entries, got %d", len(got))
	}
}

func TestApplyFilterPushDoesNotMutate(t *testing.T) {
	origLen := len(pushEntries)
	applyFilterPush(pushEntries, FilterPushOptions{Prefix: "APP_"})
	if len(pushEntries) != origLen {
		t.Error("original slice was mutated")
	}
}
