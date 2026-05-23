// Package sync provides functionality for pushing and pulling .env files
// to and from the shared secret store.
package sync

import (
	"fmt"

	"github.com/yourorg/envsync/internal/env"
)

// ScopePushOptions controls which scopes are included when pushing entries
// to the secret store.
type ScopePushOptions struct {
	// Scopes is the list of scope names to include. An empty list includes all entries.
	Scopes []string
	// ScopeOpts configures scope annotation parsing.
	ScopeOpts env.ScopeOptions
}

// DefaultScopePushOptions returns ScopePushOptions with sensible defaults.
func DefaultScopePushOptions() ScopePushOptions {
	return ScopePushOptions{
		ScopeOpts: env.DefaultScopeOptions(),
	}
}

// applyScopePush filters entries by scope before they are pushed to the store.
// If no scopes are specified all entries pass through unchanged.
func applyScopePush(entries []env.Entry, opts ScopePushOptions) ([]env.Entry, error) {
	if len(opts.Scopes) == 0 {
		return env.Clone(entries), nil
	}
	out, err := env.ScopeFilter(entries, opts.Scopes, opts.ScopeOpts)
	if err != nil {
		return nil, fmt.Errorf("scope push: %w", err)
	}
	return out, nil
}
