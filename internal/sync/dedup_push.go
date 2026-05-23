// Package sync provides push/pull synchronisation helpers.
package sync

import (
	"fmt"
	"strings"

	"github.com/example/envsync/internal/env"
)

// DefaultDedupPushOptions returns DedupPushOptions with sensible defaults.
func DefaultDedupPushOptions() DedupPushOptions {
	return DedupPushOptions{
		Enabled:  true,
		Strategy: env.DedupKeepLast,
	}
}

// DedupPushOptions configures duplicate-key removal during a push.
type DedupPushOptions struct {
	// Enabled toggles deduplication; when false entries are returned unchanged.
	Enabled bool
	// Strategy selects which duplicate entry survives.
	Strategy env.DedupStrategy
}

// applyDedupPush removes duplicate keys from entries before they are pushed
// to the remote store. It returns the deduplicated slice and a human-readable
// summary of any keys that were collapsed.
func applyDedupPush(entries []env.Entry, opts DedupPushOptions) ([]env.Entry, string, error) {
	if !opts.Enabled {
		return entries, "", nil
	}

	dups := env.DupKeys(entries)
	if len(dups) == 0 {
		return entries, "", nil
	}

	var removed []string
	dedupOpts := env.DedupOptions{
		Strategy: opts.Strategy,
		Report:   func(k string) { removed = append(removed, k) },
	}

	result := env.Dedup(entries, dedupOpts)

	summary := fmt.Sprintf(
		"dedup: collapsed %d duplicate key(s): %s",
		len(dups),
		strings.Join(dups, ", "),
	)
	_ = removed // available for structured logging if needed

	return result, summary, nil
}
