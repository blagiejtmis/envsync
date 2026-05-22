// Package sync provides push/pull synchronisation of .env files.
package sync

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envsync/internal/env"
)

// TemplateCheckResult holds the outcome of comparing an env file against its
// template.
type TemplateCheckResult struct {
	// MissingKeys are keys declared in the template but absent/blank in env.
	MissingKeys []string
	// ExtraKeys are keys present in env but not declared in the template.
	ExtraKeys []string
}

// OK returns true when there are no missing or extra keys.
func (r TemplateCheckResult) OK() bool {
	return len(r.MissingKeys) == 0 && len(r.ExtraKeys) == 0
}

// CheckAgainstTemplate compares envEntries against templateEntries and returns
// a TemplateCheckResult describing any discrepancies.
func CheckAgainstTemplate(envEntries, templateEntries []env.Entry) TemplateCheckResult {
	tmplKeys := make(map[string]bool, len(templateEntries))
	for _, e := range templateEntries {
		tmplKeys[e.Key] = true
	}

	envKeys := make(map[string]bool, len(envEntries))
	for _, e := range envEntries {
		envKeys[e.Key] = true
	}

	missing := env.CheckMissingKeys(envEntries, templateEntries)

	var extra []string
	for _, e := range envEntries {
		if !tmplKeys[e.Key] {
			extra = append(extra, e.Key)
		}
	}

	return TemplateCheckResult{MissingKeys: missing, ExtraKeys: extra}
}

// FprintCheckResult writes a human-readable summary of the check result to w.
func FprintCheckResult(w io.Writer, r TemplateCheckResult) {
	if r.OK() {
		fmt.Fprintln(w, "template check passed: all required keys are present")
		return
	}
	if len(r.MissingKeys) > 0 {
		fmt.Fprintf(w, "missing keys (%d): %s\n",
			len(r.MissingKeys), strings.Join(r.MissingKeys, ", "))
	}
	if len(r.ExtraKeys) > 0 {
		fmt.Fprintf(w, "undeclared keys (%d): %s\n",
			len(r.ExtraKeys), strings.Join(r.ExtraKeys, ", "))
	}
}
