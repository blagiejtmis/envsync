// Package sync provides push/pull operations for syncing .env files.
package sync

import "github.com/yourorg/envsync/internal/env"

// PrefixPushOptions controls how key prefixes are handled before a push.
type PrefixPushOptions struct {
	// AddPrefix, when non-empty, prepends this prefix to all keys before
	// storing them in the secret store.
	AddPrefix string

	// StripPrefix, when non-empty, removes this prefix from all keys before
	// storing them in the secret store. Keys without the prefix are kept
	// unless StripStrict is true.
	StripPrefix string

	// StripStrict causes keys that do not carry StripPrefix to be omitted.
	StripStrict bool

	env.PrefixOptions
}

// DefaultPrefixPushOptions returns sensible defaults.
func DefaultPrefixPushOptions() PrefixPushOptions {
	return PrefixPushOptions{
		PrefixOptions: env.DefaultPrefixOptions(),
	}
}

// applyPrefixPush transforms entries according to the supplied options.
// It is a pure function and does not mutate the input slice.
func applyPrefixPush(entries []env.Entry, opts PrefixPushOptions) []env.Entry {
	out := entries
	if opts.StripPrefix != "" {
		out = env.StripPrefix(out, opts.StripPrefix, opts.PrefixOptions, opts.StripStrict)
	}
	if opts.AddPrefix != "" {
		out = env.AddPrefix(out, opts.AddPrefix, opts.PrefixOptions)
	}
	return out
}
