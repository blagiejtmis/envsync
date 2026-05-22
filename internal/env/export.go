// Package env provides utilities for parsing, serializing, and manipulating
// .env files.
package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat controls how exported variables are formatted.
type ExportFormat int

const (
	// FormatRaw writes KEY=VALUE lines with no shell export prefix.
	FormatRaw ExportFormat = iota
	// FormatExport writes "export KEY=VALUE" lines suitable for sourcing.
	FormatExport
	// FormatDocker writes "--env KEY=VALUE" flags for use with docker run.
	FormatDocker
)

// Fprint writes the entries to w using the specified format.
// Keys are written in sorted order for deterministic output.
func Fprint(w io.Writer, entries []Entry, format ExportFormat) error {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, e := range sorted {
		var line string
		switch format {
		case FormatExport:
			line = fmt.Sprintf("export %s=%s\n", e.Key, quoteIfNeeded(e.Value))
		case FormatDocker:
			line = fmt.Sprintf("--env %s=%s\n", e.Key, e.Value)
		default:
			line = fmt.Sprintf("%s=%s\n", e.Key, e.Value)
		}
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("export: write: %w", err)
		}
	}
	return nil
}

// quoteIfNeeded wraps value in double quotes if it contains spaces or special
// shell characters that would require quoting when sourced.
func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n\r'\"\\$`!#&;|<>(){}") {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return v
}
