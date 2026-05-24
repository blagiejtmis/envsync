package env

// OmitOptions controls which entries are removed.
type OmitOptions struct {
	// Keys lists exact keys to omit.
	Keys []string
	// Prefixes removes any entry whose key starts with one of these prefixes.
	Prefixes []string
	// Blank removes entries whose value is the empty string.
	Blank bool
	// Sensitive removes entries detected as sensitive by IsSensitive.
	Sensitive bool
}

// DefaultOmitOptions returns an OmitOptions with no filtering enabled.
func DefaultOmitOptions() OmitOptions {
	return OmitOptions{}
}

// Omit returns a new slice of entries with the specified entries removed
// according to opts. The original slice is never mutated.
func Omit(entries []Entry, opts OmitOptions) []Entry {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if _, drop := keySet[e.Key]; drop {
			continue
		}
		if opts.Blank && e.Value == "" {
			continue
		}
		if opts.Sensitive && IsSensitive(e.Key) {
			continue
		}
		if hasAnyPrefix(e.Key, opts.Prefixes) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// OmittedKeys returns the keys that would be removed from entries by Omit.
func OmittedKeys(entries []Entry, opts OmitOptions) []string {
	full := Omit(entries, opts)
	kept := make(map[string]struct{}, len(full))
	for _, e := range full {
		kept[e.Key] = struct{}{}
	}
	var removed []string
	for _, e := range entries {
		if _, ok := kept[e.Key]; !ok {
			removed = append(removed, e.Key)
		}
	}
	return removed
}

func hasAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if p != "" && len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
