package sync

import (
	"testing"

	"github.com/example/envsync/internal/env"
)

var sortPushEntries = []env.Entry{
	{Key: "Z_LAST", Value: "z"},
	{Key: "A_FIRST", Value: "a"},
	{Key: "M_MIDDLE", Value: "m"},
}

func TestApplySortPushDisabled(t *testing.T) {
	opts := DefaultSortPushOptions()
	out := applySortPush(sortPushEntries, opts)

	if out[0].Key != "Z_LAST" {
		t.Errorf("disabled sort should preserve order; got %q first", out[0].Key)
	}
}

func TestApplySortPushAlpha(t *testing.T) {
	opts := DefaultSortPushOptions()
	opts.Enabled = true
	out := applySortPush(sortPushEntries, opts)

	if out[0].Key != "A_FIRST" {
		t.Errorf("expected A_FIRST first, got %q", out[0].Key)
	}
	if out[2].Key != "Z_LAST" {
		t.Errorf("expected Z_LAST last, got %q", out[2].Key)
	}
}

func TestApplySortPushByPrefix(t *testing.T) {
	entries := []env.Entry{
		{Key: "DB_URL", Value: "postgres"},
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "dev"},
	}
	opts := DefaultSortPushOptions()
	opts.Enabled = true
	opts.Opts.Order = env.SortByPrefix
	out := applySortPush(entries, opts)

	if out[0].Key != "APP_ENV" && out[0].Key != "APP_NAME" {
		t.Errorf("expected APP_ prefix first, got %q", out[0].Key)
	}
}

func TestApplySortPushDoesNotMutate(t *testing.T) {
	original := []env.Entry{
		{Key: "Z", Value: "1"},
		{Key: "A", Value: "2"},
	}
	opts := DefaultSortPushOptions()
	opts.Enabled = true
	applySortPush(original, opts)

	if original[0].Key != "Z" {
		t.Error("applySortPush must not mutate the original slice")
	}
}

func TestApplySortPushEmptyEntries(t *testing.T) {
	opts := DefaultSortPushOptions()
	opts.Enabled = true
	out := applySortPush([]env.Entry{}, opts)

	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}
