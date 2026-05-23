package sync

import (
	"github.com/example/envsync/internal/env"
)

// DefaultSortPushOptions returns conservative defaults: sort alphabetically
// before pushing so the stored payload has a deterministic key order.
func DefaultSortPushOptions() SortPushOptions {
	return SortPushOptions{
		Enabled: false,
		Opts:    env.DefaultSortOptions(),
	}
}

// SortPushOptions controls whether and how entries are sorted before a push.
type SortPushOptions struct {
	// Enabled activates sorting. When false the entries are unchanged.
	Enabled bool
	// Opts are passed directly to env.Sort.
	Opts env.SortOptions
}

// applySortPush returns a sorted copy of entries when Enabled is true,
// otherwise it returns a shallow copy of the original slice unchanged.
func applySortPush(entries []env.Entry, opts SortPushOptions) []env.Entry {
	if !opts.Enabled {
		out := make([]env.Entry, len(entries))
		copy(out, entries)
		return out
	}
	return env.Sort(entries, opts.Opts)
}
