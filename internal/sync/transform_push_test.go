package sync

import (
	"testing"

	"github.com/yourorg/envsync/internal/env"
)

func transformPushEntries() []env.Entry {
	return []env.Entry{
		{Key: "  db_host", Value: "  localhost  "},
		{Key: "api_key", Value: `"topsecret"`},
		{Key: "app_env", Value: "production"},
	}
}

func TestApplyTransformPushDefaultTrimValues(t *testing.T) {
	ents := transformPushEntries()
	out := applyTransformPush(ents, DefaultTransformPushOptions())

	if out[0].Value != "localhost" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
	if out[0].Key != "db_host" {
		t.Errorf("expected trimmed key, got %q", out[0].Key)
	}
}

func TestApplyTransformPushUppercaseKeys(t *testing.T) {
	ents := transformPushEntries()
	opts := DefaultTransformPushOptions()
	opts.UppercaseKeys = true
	out := applyTransformPush(ents, opts)

	if out[1].Key != "API_KEY" {
		t.Errorf("expected uppercase key, got %q", out[1].Key)
	}
	if out[2].Key != "APP_ENV" {
		t.Errorf("expected uppercase key, got %q", out[2].Key)
	}
}

func TestApplyTransformPushStripQuotes(t *testing.T) {
	ents := transformPushEntries()
	opts := DefaultTransformPushOptions()
	opts.StripQuotes = true
	out := applyTransformPush(ents, opts)

	if out[1].Value != "topsecret" {
		t.Errorf("expected unquoted value, got %q", out[1].Value)
	}
}

func TestApplyTransformPushNoOptsPreservesValues(t *testing.T) {
	ents := []env.Entry{{Key: "KEY", Value: "value"}}
	opts := DefaultTransformPushOptions()
	opts.TrimValues = false
	out := applyTransformPush(ents, opts)

	if out[0].Value != "value" {
		t.Errorf("unexpected value change: %q", out[0].Value)
	}
}

func TestApplyTransformPushDoesNotMutate(t *testing.T) {
	ents := []env.Entry{{Key: "  key", Value: "  val  "}}
	opts := DefaultTransformPushOptions()
	_ = applyTransformPush(ents, opts)

	if ents[0].Key != "  key" {
		t.Error("original slice was mutated")
	}
}
