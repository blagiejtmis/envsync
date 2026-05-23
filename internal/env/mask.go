// Package env provides utilities for parsing, serializing, and manipulating
// .env files.
package env

import (
	"strings"
)

// MaskOptions controls how values are masked when displayed.
type MaskOptions struct {
	// VisiblePrefix is the number of leading characters to keep visible.
	VisiblePrefix int
	// VisibleSuffix is the number of trailing characters to keep visible.
	VisibleSuffix int
	// MaskChar is the character used for masking. Defaults to '*'.
	MaskChar rune
}

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		VisiblePrefix: 2,
		VisibleSuffix: 0,
		MaskChar:      '*',
	}
}

// MaskValue masks a single string value according to the given options.
// If the value is shorter than or equal to VisiblePrefix+VisibleSuffix,
// the entire value is replaced with mask characters of the same length.
func MaskValue(value string, opts MaskOptions) string {
	if opts.MaskChar == 0 {
		opts.MaskChar = '*'
	}
	n := len(value)
	if n == 0 {
		return ""
	}
	total := opts.VisiblePrefix + opts.VisibleSuffix
	if total >= n {
		return strings.Repeat(string(opts.MaskChar), n)
	}
	var sb strings.Builder
	sb.WriteString(value[:opts.VisiblePrefix])
	sb.WriteString(strings.Repeat(string(opts.MaskChar), n-total))
	if opts.VisibleSuffix > 0 {
		sb.WriteString(value[n-opts.VisibleSuffix:])
	}
	return sb.String()
}

// MaskEntries returns a new slice of entries where sensitive values are masked.
// Non-sensitive values are left unchanged. Uses IsSensitive to determine
// which keys require masking.
func MaskEntries(entries []Entry, opts MaskOptions) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if IsSensitive(e.Key) {
			out[i] = Entry{Key: e.Key, Value: MaskValue(e.Value, opts)}
		} else {
			out[i] = e
		}
	}
	return out
}
