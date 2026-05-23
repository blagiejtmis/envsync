package env

import (
	"strings"
	"testing"
)

func transformEntries() []Entry {
	return []Entry{
		{Key: "  db_host  ", Value: "  localhost  "},
		{Key: "API_KEY", Value: `"secret123"`},
		{Key: "app_name", Value: "'my app'"},
	}
}

func TestTransformTrimKeysAndValues(t *testing.T) {
	ents := transformEntries()
	opts := DefaultTransformOptions()
	out := Transform(ents, opts)

	if out[0].Key != "db_host" {
		t.Errorf("expected trimmed key, got %q", out[0].Key)
	}
	if out[0].Value != "localhost" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
}

func TestTransformUppercaseKeys(t *testing.T) {
	ents := []Entry{{Key: "db_host", Value: "localhost"}}
	opts := DefaultTransformOptions()
	opts.UppercaseKeys = true
	out := Transform(ents, opts)

	if out[0].Key != "DB_HOST" {
		t.Errorf("expected uppercase key, got %q", out[0].Key)
	}
}

func TestTransformLowercaseKeys(t *testing.T) {
	ents := []Entry{{Key: "API_KEY", Value: "val"}}
	opts := DefaultTransformOptions()
	opts.LowercaseKeys = true
	out := Transform(ents, opts)

	if out[0].Key != "api_key" {
		t.Errorf("expected lowercase key, got %q", out[0].Key)
	}
}

func TestTransformStripQuotes(t *testing.T) {
	ents := transformEntries()
	opts := DefaultTransformOptions()
	opts.StripQuotes = true
	out := Transform(ents, opts)

	if out[1].Value != "secret123" {
		t.Errorf("expected unquoted value, got %q", out[1].Value)
	}
	if out[2].Value != "my app" {
		t.Errorf("expected single-unquoted value, got %q", out[2].Value)
	}
}

func TestTransformCustomFunc(t *testing.T) {
	ents := []Entry{{Key: "greeting", Value: "hello"}}
	opts := DefaultTransformOptions()
	opts.Custom = func(e Entry) Entry {
		e.Value = strings.ToUpper(e.Value)
		return e
	}
	out := Transform(ents, opts)

	if out[0].Value != "HELLO" {
		t.Errorf("expected custom transform result, got %q", out[0].Value)
	}
}

func TestTransformDoesNotMutateOriginal(t *testing.T) {
	ents := []Entry{{Key: "  key  ", Value: "  val  "}}
	opts := DefaultTransformOptions()
	_ = Transform(ents, opts)

	if ents[0].Key != "  key  " {
		t.Error("original slice was mutated")
	}
}

func TestTransformEmptySlice(t *testing.T) {
	out := Transform([]Entry{}, DefaultTransformOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
