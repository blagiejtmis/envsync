// Package env provides utilities for working with .env files.
package env

import (
	"fmt"
	"strings"
)

// ResolveOptions controls how environment variable references are resolved
// when merging a set of entries with an override map.
type ResolveOptions struct {
	// Overrides is a map of key->value that takes precedence over entry values.
	Overrides map[string]string

	// FallbackToEntry uses the existing entry value when an override is not
	// present. When false, missing overrides cause an error.
	FallbackToEntry bool

	// IgnoreUnknownKeys silently skips override keys that do not exist in entries.
	IgnoreUnknownKeys bool
}

// DefaultResolveOptions returns a ResolveOptions with sensible defaults.
func DefaultResolveOptions() ResolveOptions {
	return ResolveOptions{
		FallbackToEntry:   true,
		IgnoreUnknownKeys: true,
	}
}

// Resolve applies the override map onto entries, returning a new slice.
// Keys present in overrides replace matching entry values. Keys in entries
// that are absent from overrides are kept when FallbackToEntry is true,
// otherwise an error is returned.
func Resolve(entries []Entry, opts ResolveOptions) ([]Entry, error) {
	if opts.Overrides == nil {
		opts.Overrides = map[string]string{}
	}

	// Validate that all override keys exist in entries unless IgnoreUnknownKeys.
	if !opts.IgnoreUnknownKeys {
		keySet := make(map[string]struct{}, len(entries))
		for _, e := range entries {
			keySet[e.Key] = struct{}{}
		}
		for k := range opts.Overrides {
			if _, ok := keySet[k]; !ok {
				return nil, fmt.Errorf("resolve: unknown override key %q", k)
			}
		}
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		resolved := e
		if v, ok := opts.Overrides[e.Key]; ok {
			resolved.Value = v
		} else if !opts.FallbackToEntry {
			return nil, fmt.Errorf("resolve: no override provided for key %q", e.Key)
		}
		out = append(out, resolved)
	}
	return out, nil
}

// ResolveFromString parses a simple KEY=VALUE string list (newline-separated)
// and uses the resulting map as overrides.
func ResolveFromString(entries []Entry, raw string, opts ResolveOptions) ([]Entry, error) {
	overrides := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("resolve: invalid override line %q", line)
		}
		overrides[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	opts.Overrides = overrides
	return Resolve(entries, opts)
}
