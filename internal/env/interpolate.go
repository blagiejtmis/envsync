// Package env provides utilities for working with .env files.
package env

import (
	"fmt"
	"os"
	"strings"
)

// InterpolateOptions controls the behaviour of Interpolate.
type InterpolateOptions struct {
	// FallbackToOS allows unresolved references to be looked up in os.Environ.
	FallbackToOS bool
	// FailOnMissing returns an error if a referenced key cannot be resolved.
	FailOnMissing bool
}

// DefaultInterpolateOptions returns sensible defaults.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		FallbackToOS:  true,
		FailOnMissing: false,
	}
}

// Interpolate expands ${VAR} and $VAR references within entry values.
// Resolution order: entries map → os.Environ (when FallbackToOS is set).
func Interpolate(entries []Entry, opts InterpolateOptions) ([]Entry, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		expanded, err := expandValue(e.Value, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolate %q: %w", e.Key, err)
		}
		result[i] = Entry{Key: e.Key, Value: expanded}
	}
	return result, nil
}

// expandValue replaces all $VAR / ${VAR} tokens in s.
func expandValue(s string, lookup map[string]string, opts InterpolateOptions) (string, error) {
	var errs []string
	result := os.Expand(s, func(key string) string {
		if v, ok := lookup[key]; ok {
			return v
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(key); ok {
				return v
			}
		}
		if opts.FailOnMissing {
			errs = append(errs, key)
		}
		return ""
	})
	if len(errs) > 0 {
		return "", fmt.Errorf("unresolved variables: %s", strings.Join(errs, ", "))
	}
	return result, nil
}
