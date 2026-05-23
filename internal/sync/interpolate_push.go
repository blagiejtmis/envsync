// Package sync provides push/pull synchronisation of .env files.
package sync

import (
	"fmt"

	"github.com/yourorg/envsync/internal/env"
)

// InterpolatePushOptions controls pre-push interpolation behaviour.
type InterpolatePushOptions struct {
	// Interpolate expands variable references before encrypting and pushing.
	Interpolate bool
	// FailOnMissing aborts the push when a referenced variable cannot be resolved.
	FailOnMissing bool
}

// DefaultInterpolatePushOptions returns sensible defaults for push.
func DefaultInterpolatePushOptions() InterpolatePushOptions {
	return InterpolatePushOptions{
		Interpolate:   false,
		FailOnMissing: true,
	}
}

// prepareForPush optionally interpolates entries before they are pushed to the
// secret store. When opts.Interpolate is false the entries are returned as-is.
func prepareForPush(entries []env.Entry, opts InterpolatePushOptions) ([]env.Entry, error) {
	if !opts.Interpolate {
		return entries, nil
	}

	iopts := env.DefaultInterpolateOptions()
	iopts.FallbackToOS = true
	iopts.FailOnMissing = opts.FailOnMissing

	resolved, err := env.Interpolate(entries, iopts)
	if err != nil {
		return nil, fmt.Errorf("pre-push interpolation failed: %w", err)
	}
	return resolved, nil
}
