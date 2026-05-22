package env

import (
	"fmt"
	"io"
	"strings"
)

// FprintSchema writes a human-readable schema summary to w.
func FprintSchema(w io.Writer, schema Schema) {
	for _, e := range schema {
		if e.Comment != "" {
			fmt.Fprintf(w, "  # %s\n", e.Comment)
		}
		attrs := []string{}
		if e.Required {
			attrs = append(attrs, "required")
		} else {
			attrs = append(attrs, "optional")
		}
		if e.Default != "" {
			attrs = append(attrs, fmt.Sprintf("default=%s", e.Default))
		}
		fmt.Fprintf(w, "  %-30s [%s]\n", e.Key, strings.Join(attrs, ", "))
	}
}

// FprintSchemaViolations writes violations produced by ValidateAgainstSchema
// to w. Returns the number of violations written.
func FprintSchemaViolations(w io.Writer, violations []string) int {
	for _, v := range violations {
		fmt.Fprintf(w, "  ✗ %s\n", v)
	}
	return len(violations)
}

// SchemaDefaults returns a map of keys to their default values for entries
// that have a non-empty default defined in the schema.
func SchemaDefaults(schema Schema) map[string]string {
	defaults := make(map[string]string)
	for _, e := range schema {
		if e.Default != "" {
			defaults[e.Key] = e.Default
		}
	}
	return defaults
}

// ApplySchemaDefaults fills in default values from the schema for any keys
// not already present in entries. Returns the augmented slice.
func ApplySchemaDefaults(entries []Entry, schema Schema) []Entry {
	present := make(map[string]bool, len(entries))
	for _, e := range entries {
		present[e.Key] = true
	}
	for _, s := range schema {
		if !present[s.Key] && s.Default != "" {
			entries = append(entries, Entry{Key: s.Key, Value: s.Default})
		}
	}
	return entries
}
