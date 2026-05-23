// Package env provides utilities for loading, parsing, and manipulating .env files.
package env

import "strings"

// TransformFunc is a function that transforms a single Entry.
type TransformFunc func(e Entry) Entry

// TransformOptions controls which transformations are applied.
type TransformOptions struct {
	// TrimKeys removes leading/trailing whitespace from keys.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// LowercaseKeys converts all keys to lowercase.
	LowercaseKeys bool
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// StripQuotes removes surrounding single or double quotes from values.
	StripQuotes bool
	// Custom is an optional user-supplied transform applied last.
	Custom TransformFunc
}

// DefaultTransformOptions returns a TransformOptions with safe defaults.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{
		TrimKeys:   true,
		TrimValues: true,
	}
}

// Transform applies the configured transformations to each Entry in entries.
// It does not mutate the original slice.
func Transform(entries []Entry, opts TransformOptions) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if opts.TrimKeys {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.LowercaseKeys {
			e.Key = strings.ToLower(e.Key)
		}
		if opts.UppercaseKeys {
			e.Key = strings.ToUpper(e.Key)
		}
		if opts.StripQuotes {
			e.Value = stripSurroundingQuotes(e.Value)
		}
		if opts.Custom != nil {
			e = opts.Custom(e)
		}
		out[i] = e
	}
	return out
}

// stripSurroundingQuotes removes a matching pair of surrounding single or
// double quotes from s. If no matching pair is found, s is returned unchanged.
func stripSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
