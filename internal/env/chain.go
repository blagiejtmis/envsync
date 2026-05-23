// Package env provides utilities for working with .env files.
package env

import "fmt"

// ChainOptions controls how multiple env entry slices are merged in order.
type ChainOptions struct {
	// LastWins: if true, later sources override earlier ones for the same key.
	// If false, the first value found for a key is kept.
	LastWins bool
}

// DefaultChainOptions returns the default ChainOptions (first-wins).
func DefaultChainOptions() ChainOptions {
	return ChainOptions{LastWins: false}
}

// Chain merges multiple slices of Entry in order according to opts.
// Duplicate keys are resolved by ChainOptions.LastWins:
//   - false (default): the first occurrence of a key is kept.
//   - true: the last occurrence of a key wins.
//
// The relative order of keys from the winning source is preserved in output.
func Chain(opts ChainOptions, sources ...[]Entry) ([]Entry, error) {
	if len(sources) == 0 {
		return nil, nil
	}

	seen := make(map[string]int) // key -> index in result
	result := []Entry{}

	for srcIdx, src := range sources {
		_ = srcIdx
		for _, e := range src {
			if e.Key == "" {
				return nil, fmt.Errorf("chain: empty key in source %d", srcIdx)
			}
			if idx, exists := seen[e.Key]; exists {
				if opts.LastWins {
					result[idx] = e
				}
				// first-wins: do nothing
			} else {
				seen[e.Key] = len(result)
				result = append(result, e)
			}
		}
	}

	return result, nil
}
