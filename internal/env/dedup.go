package env

// DedupStrategy controls which entry is kept when duplicate keys are found.
type DedupStrategy int

const (
	// DedupKeepFirst retains the first occurrence of a duplicate key.
	DedupKeepFirst DedupStrategy = iota
	// DedupKeepLast retains the last occurrence of a duplicate key.
	DedupKeepLast
)

// DefaultDedupOptions returns a DedupOptions with sensible defaults.
func DefaultDedupOptions() DedupOptions {
	return DedupOptions{
		Strategy: DedupKeepLast,
	}
}

// DedupOptions configures the Dedup function.
type DedupOptions struct {
	// Strategy determines which duplicate entry survives.
	Strategy DedupStrategy
	// Report, if non-nil, is called with each key that was removed as a duplicate.
	Report func(key string)
}

// Dedup removes duplicate keys from entries according to opts.
// The relative order of surviving entries is preserved.
func Dedup(entries []Entry, opts DedupOptions) []Entry {
	if len(entries) == 0 {
		return nil
	}

	seen := make(map[string]int, len(entries)) // key -> index in result
	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			switch opts.Strategy {
			case DedupKeepFirst:
				// discard current entry
				if opts.Report != nil {
					opts.Report(e.Key)
				}
			case DedupKeepLast:
				// replace the earlier entry in-place
				if opts.Report != nil {
					opts.Report(result[idx].Key)
				}
				result[idx] = e
			}
		} else {
			seen[e.Key] = len(result)
			result = append(result, e)
		}
	}

	return result
}

// DupKeys returns a slice of keys that appear more than once in entries.
// The returned slice is in first-occurrence order.
func DupKeys(entries []Entry) []string {
	counts := make(map[string]int, len(entries))
	order := make([]string, 0)
	for _, e := range entries {
		if counts[e.Key] == 0 {
			order = append(order, e.Key)
		}
		counts[e.Key]++
	}
	var dups []string
	for _, k := range order {
		if counts[k] > 1 {
			dups = append(dups, k)
		}
	}
	return dups
}
