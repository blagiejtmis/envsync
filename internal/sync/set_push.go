package sync

import (
	"fmt"

	"github.com/example/envsync/internal/env"
)

// SetPushOptions configures how key=value overrides are applied before a push.
type SetPushOptions struct {
	// Pairs is the map of key=value overrides to apply.
	Pairs map[string]string
	// AllowNew controls whether new keys may be introduced via Pairs.
	AllowNew bool
	// OverwriteExisting controls whether existing keys may be overwritten.
	OverwriteExisting bool
}

// DefaultSetPushOptions returns permissive defaults.
func DefaultSetPushOptions() SetPushOptions {
	return SetPushOptions{
		AllowNew:          true,
		OverwriteExisting: true,
	}
}

// applySetPush applies SetPushOptions to entries before they are pushed to the
// store. If Pairs is empty the original slice is returned unchanged.
func applySetPush(entries []env.Entry, opts SetPushOptions) ([]env.Entry, error) {
	if len(opts.Pairs) == 0 {
		return entries, nil
	}

	setOpts := env.SetOptions{
		AllowNew:          opts.AllowNew,
		OverwriteExisting: opts.OverwriteExisting,
	}

	out, err := env.Set(entries, opts.Pairs, setOpts)
	if err != nil {
		return nil, fmt.Errorf("set_push: %w", err)
	}
	return out, nil
}
