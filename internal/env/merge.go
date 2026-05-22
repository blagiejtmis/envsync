// Package env provides utilities for parsing, serializing, and manipulating
// .env files.
package env

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// PreferLocal keeps the local value when a key exists in both sources.
	PreferLocal MergeStrategy = iota
	// PreferRemote overwrites local values with remote values on conflict.
	PreferRemote
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	// Merged is the final set of key-value pairs after merging.
	Merged []Entry
	// Conflicts lists keys whose values differed between local and remote.
	Conflicts []string
}

// Entry is a single key-value pair from a .env file.
type Entry struct {
	Key   string
	Value string
}

// Merge combines local and remote entries according to the given strategy.
// Keys present only in remote are always added. Keys present only in local
// are always kept. Conflicting keys are resolved by strategy.
func Merge(local, remote []Entry, strategy MergeStrategy) MergeResult {
	localMap := make(map[string]string, len(local))
	for _, e := range local {
		localMap[e.Key] = e.Value
	}

	remoteMap := make(map[string]string, len(remote))
	for _, e := range remote {
		remoteMap[e.Key] = e.Value
	}

	var conflicts []string
	result := make(map[string]string)

	// Start with all local entries.
	for k, v := range localMap {
		result[k] = v
	}

	// Apply remote entries.
	for k, rv := range remoteMap {
		lv, exists := localMap[k]
		if !exists {
			result[k] = rv
			continue
		}
		if lv != rv {
			conflicts = append(conflicts, k)
			if strategy == PreferRemote {
				result[k] = rv
			}
		}
	}

	// Preserve insertion order: local first, then new remote keys.
	seen := make(map[string]bool)
	var merged []Entry
	for _, e := range local {
		merged = append(merged, Entry{Key: e.Key, Value: result[e.Key]})
		seen[e.Key] = true
	}
	for _, e := range remote {
		if !seen[e.Key] {
			merged = append(merged, Entry{Key: e.Key, Value: e.Value})
			seen[e.Key] = true
		}
	}

	return MergeResult{Merged: merged, Conflicts: conflicts}
}
