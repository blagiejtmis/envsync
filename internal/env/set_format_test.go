package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestSetDiffAdded(t *testing.T) {
	_, result, err := SetDiff(setBase, map[string]string{"NEWKEY": "v"}, DefaultSetOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Added) != 1 || result.Added[0] != "NEWKEY" {
		t.Errorf("expected Added=[NEWKEY], got %v", result.Added)
	}
}

func TestSetDiffUpdated(t *testing.T) {
	_, result, err := SetDiff(setBase, map[string]string{"HOST": "newhost"}, DefaultSetOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Updated) != 1 || result.Updated[0] != "HOST" {
		t.Errorf("expected Updated=[HOST], got %v", result.Updated)
	}
}

func TestSetDiffUnchanged(t *testing.T) {
	_, result, err := SetDiff(setBase, map[string]string{"HOST": "localhost"}, DefaultSetOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Unchanged) != 1 || result.Unchanged[0] != "HOST" {
		t.Errorf("expected Unchanged=[HOST], got %v", result.Unchanged)
	}
}

func TestFprintSetResultNoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintSetResult(&buf, SetResult{})
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got %q", buf.String())
	}
}

func TestFprintSetResultShowsSymbols(t *testing.T) {
	var buf bytes.Buffer
	r := SetResult{
		Added:    []string{"FOO"},
		Updated:  []string{"BAR"},
		Unchanged: []string{"BAZ"},
	}
	FprintSetResult(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "+ FOO") {
		t.Errorf("expected '+ FOO' in output: %q", out)
	}
	if !strings.Contains(out, "~ BAR") {
		t.Errorf("expected '~ BAR' in output: %q", out)
	}
	if !strings.Contains(out, "  BAZ") {
		t.Errorf("expected '  BAZ' in output: %q", out)
	}
}

func TestSetSummaryString(t *testing.T) {
	r := SetResult{
		Added:    []string{"A", "B"},
		Updated:  []string{"C"},
		Unchanged: []string{},
	}
	s := SetSummary(r)
	if s != "2 added, 1 updated, 0 unchanged" {
		t.Errorf("unexpected summary: %q", s)
	}
}
