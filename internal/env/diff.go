// Package env provides utilities for parsing, serializing, and comparing
// .env files.
package env

import "sort"

// ChangeKind describes the type of change between two env maps.
type ChangeKind string

const (
	ChangeAdded   ChangeKind = "added"
	ChangeRemoved ChangeKind = "removed"
	ChangeUpdated ChangeKind = "updated"
)

// Change represents a single key-level difference between two env maps.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string // empty for added
	NewValue string // empty for removed
}

// Diff compares two env maps (old, new) and returns an ordered slice of
// Changes. Keys are compared case-sensitively. The result is sorted
// alphabetically by key for deterministic output.
func Diff(oldEnv, newEnv map[string]string) []Change {
	var changes []Change

	// Detect removed and updated keys.
	for k, oldVal := range oldEnv {
		newVal, exists := newEnv[k]
		if !exists {
			changes = append(changes, Change{
				Key:      k,
				Kind:     ChangeRemoved,
				OldValue: oldVal,
			})
		} else if oldVal != newVal {
			changes = append(changes, Change{
				Key:      k,
				Kind:     ChangeUpdated,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	// Detect added keys.
	for k, newVal := range newEnv {
		if _, exists := oldEnv[k]; !exists {
			changes = append(changes, Change{
				Key:      k,
				Kind:     ChangeAdded,
				NewValue: newVal,
			})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})
	return changes
}

// HasChanges returns true when Diff produces at least one Change.
func HasChanges(oldEnv, newEnv map[string]string) bool {
	return len(Diff(oldEnv, newEnv)) > 0
}
