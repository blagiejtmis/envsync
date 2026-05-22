package env

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting problem found in an env file.
type LintIssue struct {
	Key     string
	Message string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// LintOptions controls which lint rules are applied.
type LintOptions struct {
	RequireUppercase  bool // keys should be ALL_CAPS
	ForbidLeadingUnderscores bool // keys should not start with _
	MaxValueLength   int  // 0 means no limit
}

// DefaultLintOptions returns sensible defaults.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		RequireUppercase:         true,
		ForbidLeadingUnderscores: false,
		MaxValueLength:           1024,
	}
}

// Lint checks a slice of Entry values against the given options and returns
// any issues found. An empty slice means the entries are clean.
func Lint(entries []Entry, opts LintOptions) []LintIssue {
	var issues []LintIssue

	for _, e := range entries {
		if opts.RequireUppercase {
			if e.Key != strings.ToUpper(e.Key) {
				issues = append(issues, LintIssue{
					Key:     e.Key,
					Message: "key should be uppercase",
				})
			}
		}

		if opts.ForbidLeadingUnderscores && strings.HasPrefix(e.Key, "_") {
			issues = append(issues, LintIssue{
				Key:     e.Key,
				Message: "key must not start with an underscore",
			})
		}

		if opts.MaxValueLength > 0 && len(e.Value) > opts.MaxValueLength {
			issues = append(issues, LintIssue{
				Key:     e.Key,
				Message: fmt.Sprintf("value exceeds maximum length of %d", opts.MaxValueLength),
			})
		}
	}

	return issues
}

// LintClean returns true when Lint produces no issues.
func LintClean(entries []Entry, opts LintOptions) bool {
	return len(Lint(entries, opts)) == 0
}
