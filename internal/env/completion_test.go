package env

import (
	"strings"
	"testing"
)

var completionEntries = []Entry{
	{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "SECRET_KEY", Value: "hunter2"},
	{Key: "EMPTY_VAL", Value: ""},
}

func TestFprintCompletionBash(t *testing.T) {
	var sb strings.Builder
	opts := DefaultCompletionOptions()
	if err := FprintCompletion(&sb, completionEntries, ShellBash, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DATABASE_URL=") {
		t.Errorf("expected DATABASE_URL in output")
	}
}

func TestFprintCompletionBashExport(t *testing.T) {
	var sb strings.Builder
	opts := CompletionOptions{ExportPrefix: true}
	if err := FprintCompletion(&sb, completionEntries, ShellBash, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected 'export APP_ENV=production', got:\n%s", out)
	}
}

func TestFprintCompletionFish(t *testing.T) {
	var sb strings.Builder
	opts := DefaultCompletionOptions()
	if err := FprintCompletion(&sb, completionEntries, ShellFish, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "set -x APP_ENV production") {
		t.Errorf("expected fish set statement, got:\n%s", out)
	}
}

func TestFprintCompletionOnlyKeys(t *testing.T) {
	var sb strings.Builder
	opts := CompletionOptions{OnlyKeys: true}
	if err := FprintCompletion(&sb, completionEntries, ShellZsh, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(completionEntries) {
		t.Errorf("expected %d lines, got %d", len(completionEntries), len(lines))
	}
	for _, l := range lines {
		if strings.Contains(l, "=") {
			t.Errorf("OnlyKeys should not emit values, got: %s", l)
		}
	}
}

func TestFprintCompletionSortedOrder(t *testing.T) {
	var sb strings.Builder
	opts := CompletionOptions{OnlyKeys: true}
	if err := FprintCompletion(&sb, completionEntries, ShellBash, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] < lines[i-1] {
			t.Errorf("output not sorted: %q before %q", lines[i-1], lines[i])
		}
	}
}

func TestShellQuoteSpecialChars(t *testing.T) {
	cases := []struct {
		input string
		wantQuoted bool
	}{
		{"simple", false},
		{"has space", true},
		{"has$dollar", true},
		{"plain123", false},
	}
	for _, tc := range cases {
		result := shellQuote(tc.input)
		hasQuote := strings.HasPrefix(result, "'")
		if hasQuote != tc.wantQuoted {
			t.Errorf("shellQuote(%q) = %q, wantQuoted=%v", tc.input, result, tc.wantQuoted)
		}
	}
}
