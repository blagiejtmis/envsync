package env

import (
	"os"
	"testing"
)

func interpEntries(pairs ...string) []Entry {
	out := make([]Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestInterpolateNoRefs(t *testing.T) {
	entries := interpEntries("FOO", "bar", "BAZ", "qux")
	got, err := Interpolate(entries, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0].Value != "bar" || got[1].Value != "qux" {
		t.Fatalf("values changed unexpectedly: %+v", got)
	}
}

func TestInterpolateIntraEntry(t *testing.T) {
	entries := interpEntries("BASE", "/opt/app", "BIN", "${BASE}/bin")
	got, err := Interpolate(entries, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "/opt/app/bin"; got[1].Value != want {
		t.Fatalf("want %q, got %q", want, got[1].Value)
	}
}

func TestInterpolateFallbackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	entries := interpEntries("GREETING", "hello-$OS_VAR")
	opts := DefaultInterpolateOptions()
	opts.FallbackToOS = true
	got, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "hello-from-os"; got[0].Value != want {
		t.Fatalf("want %q, got %q", want, got[0].Value)
	}
}

func TestInterpolateNoFallbackToOS(t *testing.T) {
	os.Setenv("SHOULD_NOT_USE", "leaked")
	defer os.Unsetenv("SHOULD_NOT_USE")
	entries := interpEntries("VAL", "$SHOULD_NOT_USE")
	opts := DefaultInterpolateOptions()
	opts.FallbackToOS = false
	got, err := Interpolate(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0].Value != "" {
		t.Fatalf("expected empty string, got %q", got[0].Value)
	}
}

func TestInterpolateFailOnMissing(t *testing.T) {
	entries := interpEntries("VAL", "${UNDEFINED_XYZ_123}")
	opts := DefaultInterpolateOptions()
	opts.FallbackToOS = false
	opts.FailOnMissing = true
	_, err := Interpolate(entries, opts)
	if err == nil {
		t.Fatal("expected error for unresolved variable, got nil")
	}
}

func TestInterpolateDoesNotMutateOriginal(t *testing.T) {
	original := interpEntries("A", "$B", "B", "hello")
	_, err := Interpolate(original, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if original[0].Value != "$B" {
		t.Fatal("original entries were mutated")
	}
}
