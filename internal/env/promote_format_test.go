package env

import (
	"strings"
	"testing"
)

func TestFprintPromoteResultNoPromoted(t *testing.T) {
	r := PromoteResult{FromEnv: "staging", ToEnv: "production"}
	var sb strings.Builder
	FprintPromoteResult(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected from env in output, got: %s", out)
	}
	if !strings.Contains(out, "No keys promoted") {
		t.Errorf("expected 'No keys promoted' in output, got: %s", out)
	}
}

func TestFprintPromoteResultWithPromoted(t *testing.T) {
	r := PromoteResult{
		FromEnv:  "staging",
		ToEnv:    "production",
		Promoted: []string{"DB_HOST", "API_KEY"},
	}
	var sb strings.Builder
	FprintPromoteResult(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "Promoted (2)") {
		t.Errorf("expected promoted count in output, got: %s", out)
	}
}

func TestFprintPromoteResultWithSkipped(t *testing.T) {
	r := PromoteResult{
		FromEnv:  "staging",
		ToEnv:    "production",
		Promoted: []string{"LOG_LEVEL"},
		Skipped:  []string{"DB_HOST"},
	}
	var sb strings.Builder
	FprintPromoteResult(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "Skipped (1") {
		t.Errorf("expected skipped count in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in skipped section, got: %s", out)
	}
}

func TestPromoteSummaryNothingToPromote(t *testing.T) {
	r := PromoteResult{}
	got := PromoteSummary(r)
	if got != "nothing to promote" {
		t.Errorf("expected 'nothing to promote', got %q", got)
	}
}

func TestPromoteSummaryWithCounts(t *testing.T) {
	r := PromoteResult{
		Promoted: []string{"A", "B"},
		Skipped:  []string{"C"},
	}
	got := PromoteSummary(r)
	if !strings.Contains(got, "2 promoted") {
		t.Errorf("expected '2 promoted' in summary, got %q", got)
	}
	if !strings.Contains(got, "1 skipped") {
		t.Errorf("expected '1 skipped' in summary, got %q", got)
	}
}
