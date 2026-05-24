package sync

import (
	"fmt"

	"github.com/yourorg/envsync/internal/env"
)

// DefaultTruncatePushOptions returns sensible defaults.
func DefaultTruncatePushOptions() TruncatePushOptions {
	return TruncatePushOptions{
		Enabled:   false,
		Truncate:  env.DefaultTruncateOptions(),
	}
}

// TruncatePushOptions controls value truncation applied before a push.
type TruncatePushOptions struct {
	// Enabled must be true for truncation to be applied.
	Enabled bool
	// Truncate holds the underlying truncation parameters.
	Truncate env.TruncateOptions
}

// applyTruncatePush truncates long values in entries before they are pushed to
// the secret store. This is useful when the backend imposes a maximum secret
// size. When Enabled is false the original slice is returned unchanged.
func applyTruncatePush(entries []env.Entry, opts TruncatePushOptions) ([]env.Entry, error) {
	if !opts.Enabled {
		return entries, nil
	}
	result, err := env.Truncate(entries, opts.Truncate)
	if err != nil {
		return nil, fmt.Errorf("truncate push: %w", err)
	}
	return result, nil
}
