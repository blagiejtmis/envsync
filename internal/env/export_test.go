package env

import (
	"strings"
	"testing"
)

func TestFprintRaw(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	var sb strings.Builder
	if err := Fprint(&sb, entries, FormatRaw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := sb.String()
	if !strings.Contains(got, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux line, got:\n%s", got)
	}
	if !strings.Contains(got, "FOO=bar\n") {
		t.Errorf("expected FOO=bar line, got:\n%s", got)
	}
}

func TestFprintExport(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "secret"},
	}
	var sb strings.Builder
	if err := Fprint(&sb, entries, FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := sb.String()
	if got != "export API_KEY=secret\n" {
		t.Errorf("expected 'export API_KEY=secret', got %q", got)
	}
}

func TestFprintDocker(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Value: "8080"},
	}
	var sb strings.Builder
	if err := Fprint(&sb, entries, FormatDocker); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := sb.String()
	if got != "--env PORT=8080\n" {
		t.Errorf("expected '--env PORT=8080', got %q", got)
	}
}

func TestFprintSortedOrder(t *testing.T) {
	entries := []Entry{
		{Key: "ZEBRA", Value: "1"},
		{Key: "ALPHA", Value: "2"},
		{Key: "MIDDLE", Value: "3"},
	}
	var sb strings.Builder
	if err := Fprint(&sb, entries, FormatRaw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first line to start with ALPHA, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line to start with ZEBRA, got %q", lines[2])
	}
}

func TestQuoteIfNeeded(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"with space", `"with space"`},
		{"has\"quote", `"has\"quote"`},
		{"dollar$sign", `"dollar$sign"`},
	}
	for _, tc := range cases {
		got := quoteIfNeeded(tc.input)
		if got != tc.want {
			t.Errorf("quoteIfNeeded(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
