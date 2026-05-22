// Package sync provides push/pull synchronisation of .env files against the
// secret store backend.
package sync

import (
	"fmt"

	"github.com/yourorg/envsync/internal/env"
)

// MergePullOptions configures the behaviour of MergePull.
type MergePullOptions struct {
	// Strategy controls how key conflicts are resolved.
	Strategy env.MergeStrategy
	// DryRun reports what would change without writing any files.
	DryRun bool
}

// MergePullResult summarises what happened during a MergePull.
type MergePullResult struct {
	Conflicts  []string
	Written    bool
	Merged     []env.Entry
}

// MergePull fetches the remote .env from the store, merges it with the local
// file, and writes the result back to disk (unless DryRun is set).
func (s *Syncer) MergePull(opts MergePullOptions) (MergePullResult, error) {
	// Load local file (may not exist yet).
	var localEntries []env.Entry
	if pairs, err := env.LoadFile(s.envPath); err == nil {
		for k, v := range pairs {
			localEntries = append(localEntries, env.Entry{Key: k, Value: v})
		}
	}

	// Pull remote entries via normal pull logic.
	remotePairs, err := s.fetchRemote()
	if err != nil {
		return MergePullResult{}, fmt.Errorf("merge pull: fetch remote: %w", err)
	}

	var remoteEntries []env.Entry
	for k, v := range remotePairs {
		remoteEntries = append(remoteEntries, env.Entry{Key: k, Value: v})
	}

	res := env.Merge(localEntries, remoteEntries, opts.Strategy)

	if !opts.DryRun {
		mergedMap := make(map[string]string, len(res.Merged))
		for _, e := range res.Merged {
			mergedMap[e.Key] = e.Value
		}
		if err := env.WriteFile(s.envPath, mergedMap); err != nil {
			return MergePullResult{}, fmt.Errorf("merge pull: write file: %w", err)
		}
	}

	return MergePullResult{
		Conflicts: res.Conflicts,
		Written:   !opts.DryRun,
		Merged:    res.Merged,
	}, nil
}
