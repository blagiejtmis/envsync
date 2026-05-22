package sync

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/env"
)

func tmplEntries(keys ...string) []env.Entry {
	out := make([]env.Entry, len(keys))
	for i, k := range keys {
		out[i] = env.Entry{Key: k}
	}
	return out
}

func envEntries(pairs ...string) []env.Entry {
	out := make([]env.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCheckAgainstTemplateOK(t *testing.T) {
	result := CheckAgainstTemplate(
		envEntries("A", "1", "B", "2"),
		tmplEntries("A", "B"),
	)
	if !result.OK() {
		t.Fatalf("expected OK, got missing=%v extra=%v", result.MissingKeys, result.ExtraKeys)
	}
}

func TestCheckAgainstTemplateMissingKey(t *testing.T) {
	result := CheckAgainstTemplate(
		envEntries("A", "1"),
		tmplEntries("A", "B"),
	)
	if result.OK() {
		t.Fatal("expected not OK")
	}
	if len(result.MissingKeys) != 1 || result.MissingKeys[0] != "B" {
		t.Errorf("expected missing [B], got %v", result.MissingKeys)
	}
}

func TestCheckAgainstTemplateExtraKey(t *testing.T) {
	result := CheckAgainstTemplate(
		envEntries("A", "1", "EXTRA", "x"),
		tmplEntries("A"),
	)
	if len(result.ExtraKeys) != 1 || result.ExtraKeys[0] != "EXTRA" {
		t.Errorf("expected extra [EXTRA], got %v", result.ExtraKeys)
	}
}

func TestCheckAgainstTemplateBlankCountsAsMissing(t *testing.T) {
	result := CheckAgainstTemplate(
		envEntries("A", ""),
		tmplEntries("A"),
	)
	if len(result.MissingKeys) != 1 {
		t.Errorf("blank value should count as missing, got %v", result.MissingKeys)
	}
}

func TestFprintCheckResultOK(t *testing.T) {
	var buf bytes.Buffer
	FprintCheckResult(&buf, TemplateCheckResult{})
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected 'passed' in output, got %q", buf.String())
	}
}

func TestFprintCheckResultWithIssues(t *testing.T) {
	var buf bytes.Buffer
	FprintCheckResult(&buf, TemplateCheckResult{
		MissingKeys: []string{"SECRET"},
		ExtraKeys:   []string{"LEGACY"},
	})
	out := buf.String()
	if !strings.Contains(out, "SECRET") {
		t.Error("expected SECRET in output")
	}
	if !strings.Contains(out, "LEGACY") {
		t.Error("expected LEGACY in output")
	}
}
