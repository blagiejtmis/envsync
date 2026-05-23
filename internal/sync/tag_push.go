// Package sync provides push/pull synchronisation of .env files.
package sync

import (
	"fmt"

	"github.com/yourorg/envsync/internal/env"
)

// TagPushOptions controls which tag to isolate before pushing.
type TagPushOptions struct {
	// TagName restricts the push to entries belonging to this tag.
	// When empty, all entries are pushed.
	TagName string
}

// DefaultTagPushOptions returns options that push all entries.
func DefaultTagPushOptions() TagPushOptions {
	return TagPushOptions{}
}

// applyTagPush filters entries to those belonging to the requested tag.
// If opts.TagName is empty the original slice is returned unchanged.
func applyTagPush(entries []env.Entry, opts TagPushOptions) ([]env.Entry, error) {
	if opts.TagName == "" {
		return entries, nil
	}
	tags := env.ParseTags(entries)
	filtered, err := env.FilterByTag(entries, tags, opts.TagName)
	if err != nil {
		return nil, fmt.Errorf("tag push: %w", err)
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("tag push: no entries found for tag %q", opts.TagName)
	}
	return filtered, nil
}

// AvailableTags returns the set of distinct tag names present in entries.
// It is a convenience wrapper around env.ParseTags for callers that only
// need the tag names without the full tag metadata.
func AvailableTags(entries []env.Entry) []string {
	tags := env.ParseTags(entries)
	names := make([]string, 0, len(tags))
	for name := range tags {
		names = append(names, name)
	}
	return names
}
