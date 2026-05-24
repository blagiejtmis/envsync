package env

import (
	"strings"
	"unicode"
)

// SanitizeOptions controls how entries are sanitized.
type SanitizeOptions struct {
	// StripControlChars removes non-printable control characters from values.
	StripControlChars bool
	// TrimQuotes removes surrounding single or double quotes from values.
	TrimQuotes bool
	// NormalizeNewlines replaces \r\n and bare \r with \n in values.
	NormalizeNewlines bool
	// MaxValueLen truncates values longer than this length. 0 means no limit.
	MaxValueLen int
}

// DefaultSanitizeOptions returns a SanitizeOptions with safe defaults.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		StripControlChars: true,
		TrimQuotes:        false,
		NormalizeNewlines: true,
		MaxValueLen:       0,
	}
}

// Sanitize applies sanitization rules to a slice of Entry values.
// It does not mutate the original slice.
func Sanitize(entries []Entry, opts SanitizeOptions) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		v := e.Value

		if opts.NormalizeNewlines {
			v = strings.ReplaceAll(v, "\r\n", "\n")
			v = strings.ReplaceAll(v, "\r", "\n")
		}

		if opts.StripControlChars {
			v = stripControl(v)
		}

		if opts.TrimQuotes {
			v = trimMatchingQuotes(v)
		}

		if opts.MaxValueLen > 0 && len(v) > opts.MaxValueLen {
			v = v[:opts.MaxValueLen]
		}

		out[i] = Entry{Key: e.Key, Value: v}
	}
	return out
}

// stripControl removes non-printable, non-space control characters.
func stripControl(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// trimMatchingQuotes removes a single layer of surrounding quotes if they match.
func trimMatchingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
