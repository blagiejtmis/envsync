package env

import (
	"sort"
	"strings"
)

// SortOrder defines how entries should be sorted.
type SortOrder int

const (
	// SortAlpha sorts entries alphabetically by key (A→Z).
	SortAlpha SortOrder = iota
	// SortAlphaDesc sorts entries reverse-alphabetically by key (Z→A).
	SortAlphaDesc
	// SortByPrefix groups entries by their prefix (before the first '_').
	SortByPrefix
)

// DefaultSortOptions returns the default sorting options.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Order:          SortAlpha,
		Stable:         true,
		CaseSensitive:  false,
	}
}

// SortOptions controls how entries are sorted.
type SortOptions struct {
	// Order specifies the sort order.
	Order SortOrder
	// Stable preserves the relative order of equal keys.
	Stable bool
	// CaseSensitive enables case-sensitive key comparison.
	CaseSensitive bool
}

// Sort returns a sorted copy of entries according to the given options.
// The original slice is not modified.
func Sort(entries []Entry, opts SortOptions) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)

	less := buildLessFunc(opts)

	if opts.Stable {
		sort.SliceStable(out, less)
	} else {
		sort.Slice(out, less)
	}
	return out
}

func buildLessFunc(opts SortOptions) func(i, j int) bool {
	normKey := func(k string) string {
		if opts.CaseSensitive {
			return k
		}
		return strings.ToLower(k)
	}

	switch opts.Order {
	case SortAlphaDesc:
		return func(i, j int) bool {
			return normKey(entries[i].Key) > normKey(entries[j].Key)
		}
	case SortByPrefix:
		return func(i, j int) bool {
			pi := prefixOf(normKey(entries[i].Key))
			pj := prefixOf(normKey(entries[j].Key))
			if pi != pj {
				return pi < pj
			}
			return normKey(entries[i].Key) < normKey(entries[j].Key)
		}
	default: // SortAlpha
		return func(i, j int) bool {
			return normKey(entries[i].Key) < normKey(entries[j].Key)
		}
	}
}

func prefixOf(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
