package env

import (
	"strings"
	"testing"
)

func lintEntries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestLintCleanEntries(t *testing.T) {
	entries := lintEntries("DATABASE_URL", "postgres://localhost/db", "PORT", "5432")
	issues := Lint(entries, DefaultLintOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestLintRequiresUppercase(t *testing.T) {
	entries := lintEntries("database_url", "postgres://localhost", "PORT", "5432")
	opts := DefaultLintOptions()
	issues := Lint(entries, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "database_url" {
		t.Errorf("unexpected key %q in issue", issues[0].Key)
	}
	if !strings.Contains(issues[0].Message, "uppercase") {
		t.Errorf("expected uppercase message, got %q", issues[0].Message)
	}
}

func TestLintForbidLeadingUnderscores(t *testing.T) {
	entries := lintEntries("_PRIVATE", "secret", "PUBLIC", "value")
	opts := DefaultLintOptions()
	opts.ForbidLeadingUnderscores = true
	issues := Lint(entries, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "underscore") {
		t.Errorf("unexpected message: %q", issues[0].Message)
	}
}

func TestLintMaxValueLength(t *testing.T) {
	long := strings.Repeat("x", 2048)
	entries := lintEntries("TOKEN", long)
	opts := DefaultLintOptions()
	issues := Lint(entries, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
	if !strings.Contains(issues[0].Message, "maximum length") {
		t.Errorf("unexpected message: %q", issues[0].Message)
	}
}

func TestLintCleanHelper(t *testing.T) {
	entries := lintEntries("GOOD_KEY", "value")
	if !LintClean(entries, DefaultLintOptions()) {
		t.Error("expected LintClean to return true")
	}

	bad := lintEntries("bad_key", "value")
	if LintClean(bad, DefaultLintOptions()) {
		t.Error("expected LintClean to return false for bad key")
	}
}

func TestLintIssueString(t *testing.T) {
	issue := LintIssue{Key: "MY_KEY", Message: "some problem"}
	got := issue.String()
	if got != "MY_KEY: some problem" {
		t.Errorf("unexpected String() output: %q", got)
	}
}
