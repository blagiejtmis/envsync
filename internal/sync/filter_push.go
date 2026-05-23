package sync

import "github.com/yourorg/envsync/internal/env"

// FilterPushOptions configures which entries are included when pushing.
type FilterPushOptions struct {
	// Prefix restricts the push to keys with this prefix.
	Prefix string
	// Keys restricts the push to exactly these keys (takes precedence over Prefix).
	Keys []string
	// ExcludeKeys removes these keys before pushing.
	ExcludeKeys []string
}

// DefaultFilterPushOptions returns a FilterPushOptions that passes all entries.
func DefaultFilterPushOptions() FilterPushOptions {
	return FilterPushOptions{}
}

// applyFilterPush filters entries according to opts before they are pushed to
// the store. It returns a new slice and never mutates the input.
func applyFilterPush(entries []env.Entry, opts FilterPushOptions) []env.Entry {
	if opts.Prefix == "" && len(opts.Keys) == 0 && len(opts.ExcludeKeys) == 0 {
		return entries
	}
	return env.Filter(entries, env.FilterOptions{
		Prefix:      opts.Prefix,
		Keys:        opts.Keys,
		ExcludeKeys: opts.ExcludeKeys,
	})
}
