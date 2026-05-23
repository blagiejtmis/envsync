package env

import "strings"

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	// Prefix keeps only entries whose key starts with the given prefix.
	Prefix string
	// Keys keeps only entries whose key is in the set (ignored when empty).
	Keys []string
	// ExcludeKeys removes entries whose key is in the set.
	ExcludeKeys []string
}

// Filter returns a new slice containing only entries that match opts.
func Filter(entries []Entry, opts FilterOptions) []Entry {
	exclude := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		exclude[k] = true
	}

	allow := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allow[k] = true
	}

	var out []Entry
	for _, e := range entries {
		if exclude[e.Key] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if len(allow) > 0 && !allow[e.Key] {
			continue
		}
		out = append(out, e)
	}
	return out
}

// FilterByPrefix is a convenience wrapper around Filter for prefix-only filtering.
func FilterByPrefix(entries []Entry, prefix string) []Entry {
	return Filter(entries, FilterOptions{Prefix: prefix})
}

// FilterByKeys is a convenience wrapper that keeps only the specified keys.
func FilterByKeys(entries []Entry, keys ...string) []Entry {
	return Filter(entries, FilterOptions{Keys: keys})
}

// Reject removes entries whose key is in the exclusion list.
func Reject(entries []Entry, keys ...string) []Entry {
	return Filter(entries, FilterOptions{ExcludeKeys: keys})
}
