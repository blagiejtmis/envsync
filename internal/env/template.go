// Package env provides utilities for working with .env files.
package env

import (
	"fmt"
	"os"
	"strings"
)

// TemplateOptions controls how a template is generated or applied.
type TemplateOptions struct {
	// RedactValues replaces all values with empty strings in the output.
	RedactValues bool
	// IncludeComments preserves comment lines from the source.
	IncludeComments bool
}

// DefaultTemplateOptions returns sensible defaults for template generation.
func DefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{
		RedactValues:    true,
		IncludeComments: true,
	}
}

// GenerateTemplate produces a .env.template from a slice of Entry values.
// Values are blanked out so the template can be safely committed to VCS.
func GenerateTemplate(entries []Entry, opts TemplateOptions) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		out[i] = Entry{Key: e.Key, Value: e.Value}
		if opts.RedactValues {
			out[i].Value = ""
		}
	}
	return out
}

// ApplyTemplate fills in missing keys in dst from template, leaving existing
// values untouched. It returns the merged result and a list of keys that were
// added from the template.
func ApplyTemplate(dst []Entry, template []Entry) ([]Entry, []string) {
	existing := make(map[string]bool, len(dst))
	for _, e := range dst {
		existing[e.Key] = true
	}

	var added []string
	result := make([]Entry, len(dst))
	copy(result, dst)

	for _, t := range template {
		if !existing[t.Key] {
			result = append(result, Entry{Key: t.Key, Value: ""})
			added = append(added, t.Key)
		}
	}
	return result, added
}

// CheckMissingKeys returns keys present in template but absent (or blank) in env.
func CheckMissingKeys(env []Entry, template []Entry) []string {
	values := make(map[string]string, len(env))
	for _, e := range env {
		values[e.Key] = e.Value
	}

	var missing []string
	for _, t := range template {
		v, ok := values[t.Key]
		if !ok || strings.TrimSpace(v) == "" {
			missing = append(missing, t.Key)
		}
	}
	return missing
}

// WriteTemplateFile serialises entries (with values blanked) to path.
func WriteTemplateFile(path string, entries []Entry, opts TemplateOptions) error {
	tpl := GenerateTemplate(entries, opts)
	data := Serialize(tpl)
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		return fmt.Errorf("write template %s: %w", path, err)
	}
	return nil
}
