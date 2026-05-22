// Package env provides utilities for working with .env files.
package env

import (
	"strings"
	"unicode"
)

// NormalizeOptions controls how entries are normalized.
type NormalizeOptions struct {
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
	// RemoveDuplicates keeps only the last occurrence of a duplicate key.
	RemoveDuplicates bool
}

// DefaultNormalizeOptions returns sensible defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		UppercaseKeys:    true,
		TrimValues:       true,
		RemoveDuplicates: true,
	}
}

// Normalize applies normalization rules to a slice of Entry values and
// returns a new slice. The original slice is never mutated.
func Normalize(entries []Entry, opts NormalizeOptions) []Entry {
	seen := make(map[string]int) // key -> index in result
	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		key := e.Key
		val := e.Value

		if opts.UppercaseKeys {
			key = toUpperASCII(key)
		}
		if opts.TrimValues {
			val = strings.TrimSpace(val)
		}

		normalized := Entry{Key: key, Value: val}

		if opts.RemoveDuplicates {
			if idx, exists := seen[key]; exists {
				result[idx] = normalized
				continue
			}
			seen[key] = len(result)
		}

		result = append(result, normalized)
	}

	return result
}

// toUpperASCII uppercases only ASCII letters, leaving non-ASCII runes intact.
func toUpperASCII(s string) string {
	return strings.Map(func(r rune) rune {
		if r <= unicode.MaxASCII {
			return unicode.ToUpper(r)
		}
		return r
	}, s)
}
