package env

import "fmt"

// PatchOp represents a single patch operation on an env entry.
type PatchOp struct {
	Key    string
	Value  string
	Delete bool
}

// PatchOptions controls the behaviour of Patch.
type PatchOptions struct {
	// AllowNewKeys permits patch ops that introduce keys not present in base.
	AllowNewKeys bool
	// IgnoreMissingDeletes silently skips delete ops for keys that do not exist.
	IgnoreMissingDeletes bool
}

// DefaultPatchOptions returns sensible defaults for Patch.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		AllowNewKeys:         true,
		IgnoreMissingDeletes: false,
	}
}

// Patch applies a slice of PatchOp operations to base and returns a new
// slice without mutating the original. Set operations overwrite existing
// values or append new entries when AllowNewKeys is true. Delete operations
// remove the matching key.
func Patch(base []Entry, ops []PatchOp, opts PatchOptions) ([]Entry, error) {
	// Work on a deep copy so callers are not surprised.
	out := Clone(base)

	for _, op := range ops {
		if op.Key == "" {
			return nil, fmt.Errorf("patch: empty key in operation")
		}

		if op.Delete {
			idx := indexOfKey(out, op.Key)
			if idx == -1 {
				if opts.IgnoreMissingDeletes {
					continue
				}
				return nil, fmt.Errorf("patch: delete key %q not found", op.Key)
			}
			out = append(out[:idx], out[idx+1:]...)
			continue
		}

		// Set operation.
		idx := indexOfKey(out, op.Key)
		if idx == -1 {
			if !opts.AllowNewKeys {
				return nil, fmt.Errorf("patch: key %q not found and AllowNewKeys is false", op.Key)
			}
			out = append(out, Entry{Key: op.Key, Value: op.Value})
		} else {
			out[idx].Value = op.Value
		}
	}

	return out, nil
}

// indexOfKey returns the index of the first Entry with the given key, or -1.
func indexOfKey(entries []Entry, key string) int {
	for i, e := range entries {
		if e.Key == key {
			return i
		}
	}
	return -1
}
