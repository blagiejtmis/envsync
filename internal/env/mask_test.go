package env

import (
	"strings"
	"testing"
)

func TestMaskValueEmpty(t *testing.T) {
	opts := DefaultMaskOptions()
	if got := MaskValue("", opts); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMaskValueShortString(t *testing.T) {
	opts := DefaultMaskOptions() // prefix=2, suffix=0
	// value shorter than or equal to prefix: fully masked
	got := MaskValue("ab", opts)
	if got != "**" {
		t.Fatalf("expected '**', got %q", got)
	}
}

func TestMaskValueDefaultOptions(t *testing.T) {
	opts := DefaultMaskOptions() // prefix=2, suffix=0
	got := MaskValue("supersecret", opts)
	if !strings.HasPrefix(got, "su") {
		t.Fatalf("expected prefix 'su', got %q", got)
	}
	if !strings.Contains(got, "*") {
		t.Fatalf("expected mask chars in %q", got)
	}
	if len(got) != len("supersecret") {
		t.Fatalf("expected same length, got %d", len(got))
	}
}

func TestMaskValueWithSuffix(t *testing.T) {
	opts := MaskOptions{VisiblePrefix: 2, VisibleSuffix: 2, MaskChar: '#'}
	got := MaskValue("abcdefgh", opts)
	// expect "ab####gh"
	if got != "ab####gh" {
		t.Fatalf("expected 'ab####gh', got %q", got)
	}
}

func TestMaskValueCustomChar(t *testing.T) {
	opts := MaskOptions{VisiblePrefix: 1, VisibleSuffix: 0, MaskChar: '-'}
	got := MaskValue("hello", opts)
	if got != "h----" {
		t.Fatalf("expected 'h----', got %q", got)
	}
}

func TestMaskEntriesSensitiveOnly(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t!"},
		{Key: "API_KEY", Value: "abc123xyz"},
	}
	opts := DefaultMaskOptions()
	result := MaskEntries(entries, opts)

	if result[0].Value != "myapp" {
		t.Errorf("APP_NAME should not be masked, got %q", result[0].Value)
	}
	if result[1].Value == "s3cr3t!" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if result[2].Value == "abc123xyz" {
		t.Errorf("API_KEY should be masked")
	}
}

func TestMaskEntriesDoesNotMutateOriginal(t *testing.T) {
	original := []Entry{
		{Key: "SECRET_TOKEN", Value: "original-value"},
	}
	opts := DefaultMaskOptions()
	_ = MaskEntries(original, opts)
	if original[0].Value != "original-value" {
		t.Errorf("original slice was mutated")
	}
}
