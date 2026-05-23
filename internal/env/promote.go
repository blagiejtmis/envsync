package env

import "fmt"

// PromoteOptions controls how entries are promoted between environments.
type PromoteOptions struct {
	// FromEnv is the source environment label (e.g. "staging").
	FromEnv string
	// ToEnv is the destination environment label (e.g. "production").
	ToEnv string
	// AllowOverwrite controls whether existing keys in the destination are overwritten.
	AllowOverwrite bool
	// Keys restricts promotion to specific keys; empty means all keys.
	Keys []string
}

// DefaultPromoteOptions returns sensible defaults for promotion.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		AllowOverwrite: false,
	}
}

// Promote copies selected entries from src into dst according to opts.
// It returns the resulting merged slice and a list of keys that were promoted.
func Promote(src, dst []Entry, opts PromoteOptions) ([]Entry, []string, error) {
	srcMap := make(map[string]Entry, len(src))
	for _, e := range src {
		srcMap[e.Key] = e
	}

	dstMap := make(map[string]Entry, len(dst))
	dstOrder := make([]string, 0, len(dst))
	for _, e := range dst {
		dstMap[e.Key] = e
		dstOrder = append(dstOrder, e.Key)
	}

	promoteKeys := opts.Keys
	if len(promoteKeys) == 0 {
		for _, e := range src {
			promoteKeys = append(promoteKeys, e.Key)
		}
	}

	var promoted []string
	for _, k := range promoteKeys {
		srcEntry, ok := srcMap[k]
		if !ok {
			return nil, nil, fmt.Errorf("promote: key %q not found in source environment %q", k, opts.FromEnv)
		}
		if _, exists := dstMap[k]; exists && !opts.AllowOverwrite {
			continue
		}
		if _, exists := dstMap[k]; !exists {
			dstOrder = append(dstOrder, k)
		}
		dstMap[k] = srcEntry
		promoted = append(promoted, k)
	}

	result := make([]Entry, 0, len(dstOrder))
	for _, k := range dstOrder {
		result = append(result, dstMap[k])
	}
	return result, promoted, nil
}
