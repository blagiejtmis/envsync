package env

import (
	"strings"
	"testing"
)

var truncateEntries = []Entry{
	{Key: "SHORT", Value: "hi"},
	{Key: "LONG_VALUE", Value: strings.Repeat("x", 300)},
	{Key: "EXACT", Value: strings.Repeat("a", 256)},
	{Key: "SKIP_ME", Value: strings.Repeat("z", 300)},
}

func TestTruncateShortValuesUnchanged(t *testing.T) {
	out, err := Truncate(truncateEntries, DefaultTruncateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "hi" {
		t.Errorf("expected 'hi', got %q", out[0].Value)
	}
}

func TestTruncateLongValueClipped(t *testing.T) {
	opts := DefaultTruncateOptions()
	out, err := Truncate(truncateEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out[1].Value) != opts.MaxLength {
		t.Errorf("expected length %d, got %d", opts.MaxLength, len(out[1].Value))
	}
	if !strings.HasSuffix(out[1].Value, opts.Suffix) {
		t.Errorf("expected suffix %q in %q", opts.Suffix, out[1].Value)
	}
}

func TestTruncateExactLengthUnchanged(t *testing.T) {
	out, err := Truncate(truncateEntries, DefaultTruncateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out[2].Value) != 256 {
		t.Errorf("expected 256, got %d", len(out[2].Value))
	}
}

func TestTruncateSpecificKeysOnly(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.Keys = []string{"LONG_VALUE"}
	out, err := Truncate(truncateEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out[3].Value) != 300 {
		t.Errorf("SKIP_ME should be untouched, got length %d", len(out[3].Value))
	}
	if len(out[1].Value) != opts.MaxLength {
		t.Errorf("LONG_VALUE should be clipped to %d, got %d", opts.MaxLength, len(out[1].Value))
	}
}

func TestTruncateDoesNotMutateOriginal(t *testing.T) {
	original := []Entry{{Key: "K", Value: strings.Repeat("v", 400)}}
	_, err := Truncate(original, DefaultTruncateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(original[0].Value) != 400 {
		t.Error("original entry was mutated")
	}
}

func TestTruncateInvalidMaxLength(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.MaxLength = 0
	_, err := Truncate(truncateEntries, opts)
	if err == nil {
		t.Error("expected error for MaxLength=0")
	}
}

func TestTruncateSuffixTooLong(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.MaxLength = 3
	opts.Suffix = "..."
	_, err := Truncate(truncateEntries, opts)
	if err == nil {
		t.Error("expected error when suffix >= MaxLength")
	}
}
