package env

import (
	"testing"
)

func sanitizeEntries() []Entry {
	return []Entry{
		{Key: "A", Value: "hello\x00world"},
		{Key: "B", Value: "\"quoted\""},
		{Key: "C", Value: "line\r\nend"},
		{Key: "D", Value: "normal"},
		{Key: "E", Value: "'single'"},
	}
}

func TestSanitizeStripControlChars(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.TrimQuotes = false
	out := Sanitize(sanitizeEntries(), opts)
	for _, e := range out {
		if e.Key == "A" && e.Value != "helloworld" {
			t.Errorf("expected control char stripped, got %q", e.Value)
		}
	}
}

func TestSanitizeTrimDoubleQuotes(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.TrimQuotes = true
	out := Sanitize(sanitizeEntries(), opts)
	for _, e := range out {
		if e.Key == "B" && e.Value != "quoted" {
			t.Errorf("expected quotes trimmed, got %q", e.Value)
		}
	}
}

func TestSanitizeTrimSingleQuotes(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.TrimQuotes = true
	out := Sanitize(sanitizeEntries(), opts)
	for _, e := range out {
		if e.Key == "E" && e.Value != "single" {
			t.Errorf("expected single quotes trimmed, got %q", e.Value)
		}
	}
}

func TestSanitizeNormalizeNewlines(t *testing.T) {
	opts := DefaultSanitizeOptions()
	out := Sanitize(sanitizeEntries(), opts)
	for _, e := range out {
		if e.Key == "C" && e.Value != "line\nend" {
			t.Errorf("expected CRLF normalized, got %q", e.Value)
		}
	}
}

func TestSanitizeMaxValueLen(t *testing.T) {
	entries := []Entry{{Key: "LONG", Value: "abcdefghij"}}
	opts := DefaultSanitizeOptions()
	opts.MaxValueLen = 5
	out := Sanitize(entries, opts)
	if out[0].Value != "abcde" {
		t.Errorf("expected truncation to 5, got %q", out[0].Value)
	}
}

func TestSanitizeDoesNotMutateOriginal(t *testing.T) {
	original := []Entry{{Key: "X", Value: "\x01dirty"}}
	opts := DefaultSanitizeOptions()
	_ = Sanitize(original, opts)
	if original[0].Value != "\x01dirty" {
		t.Error("original slice was mutated")
	}
}

func TestSanitizeNormalValueUnchanged(t *testing.T) {
	entries := []Entry{{Key: "D", Value: "normal"}}
	opts := DefaultSanitizeOptions()
	out := Sanitize(entries, opts)
	if out[0].Value != "normal" {
		t.Errorf("expected value unchanged, got %q", out[0].Value)
	}
}
