// Package env provides utilities for parsing, manipulating, and serializing
// environment variable files.
package env

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested key segments are flattened.
type FlattenOptions struct {
	// Separator is the string used to join key segments (default "_").
	Separator string
	// Uppercase forces all resulting keys to uppercase.
	Uppercase bool
	// Prefix is prepended to every resulting key.
	Prefix string
}

// DefaultFlattenOptions returns sensible defaults for FlattenOptions.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "_",
		Uppercase: true,
	}
}

// Flatten takes a map of dot-notation or slash-notation keys (e.g. "db.host"
// or "db/host") and converts them into flat ENV-style entries using the
// configured separator. Existing Entry comments are preserved when the source
// key matches exactly.
//
// Example:
//
//	{"db.host": "localhost", "db.port": "5432"}
//	→ [{Key:"DB_HOST", Value:"localhost"}, {Key:"DB_PORT", Value:"5432"}]
func Flatten(pairs map[string]string, opts FlattenOptions) ([]Entry, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	seen := make(map[string]struct{}, len(pairs))
	entries := make([]Entry, 0, len(pairs))

	for rawKey, val := range pairs {
		flat, err := flattenKey(rawKey, opts)
		if err != nil {
			return nil, fmt.Errorf("flatten: %w", err)
		}
		if _, dup := seen[flat]; dup {
			return nil, fmt.Errorf("flatten: duplicate key %q after flattening", flat)
		}
		seen[flat] = struct{}{}
		entries = append(entries, Entry{Key: flat, Value: val})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries, nil
}

// flattenKey converts a single raw key into a flat ENV key.
func flattenKey(raw string, opts FlattenOptions) (string, error) {
	// Normalise separators: replace dots and slashes with the target separator.
	replacer := strings.NewReplacer(".", opts.Separator, "/", opts.Separator)
	key := replacer.Replace(raw)

	if opts.Prefix != "" {
		key = opts.Prefix + opts.Separator + key
	}
	if opts.Uppercase {
		key = toUpperASCII(key)
	}
	if key == "" {
		return "", fmt.Errorf("empty key produced from %q", raw)
	}
	return key, nil
}
