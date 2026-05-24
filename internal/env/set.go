package env

import "fmt"

// SetOptions controls behaviour of Set.
type SetOptions struct {
	// AllowNew permits adding keys that do not already exist.
	AllowNew bool
	// OverwriteExisting permits changing the value of an existing key.
	OverwriteExisting bool
}

// DefaultSetOptions returns permissive defaults: allow both new and existing keys.
func DefaultSetOptions() SetOptions {
	return SetOptions{
		AllowNew:          true,
		OverwriteExisting: true,
	}
}

// Set applies one or more key=value pairs to entries according to opts.
// It returns the updated slice (original is not mutated) and an error if any
// constraint in opts is violated.
func Set(entries []Entry, pairs map[string]string, opts SetOptions) ([]Entry, error) {
	out := Clone(entries)

	for key, value := range pairs {
		if key == "" {
			return nil, fmt.Errorf("set: empty key is not allowed")
		}

		idx := indexOfEntry(out, key)
		if idx >= 0 {
			if !opts.OverwriteExisting {
				return nil, fmt.Errorf("set: key %q already exists and OverwriteExisting is false", key)
			}
			out[idx].Value = value
		} else {
			if !opts.AllowNew {
				return nil, fmt.Errorf("set: key %q does not exist and AllowNew is false", key)
			}
			out = append(out, Entry{Key: key, Value: value})
		}
	}
	return out, nil
}

// indexOfEntry returns the index of the first Entry with the given key, or -1.
func indexOfEntry(entries []Entry, key string) int {
	for i, e := range entries {
		if e.Key == key {
			return i
		}
	}
	return -1
}
