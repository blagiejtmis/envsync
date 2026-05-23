package env

import "strings"

// PrefixOptions controls how key prefixes are added or stripped.
type PrefixOptions struct {
	// Separator placed between the prefix and the original key.
	// Defaults to "_".
	Separator string
}

// DefaultPrefixOptions returns sensible defaults.
func DefaultPrefixOptions() PrefixOptions {
	return PrefixOptions{Separator: "_"}
}

// AddPrefix returns a new slice of entries with prefix prepended to every key.
// Entries with empty keys are left unchanged.
func AddPrefix(entries []Entry, prefix string, opts PrefixOptions) []Entry {
	if prefix == "" {
		return Clone(entries)
	}
	sep := opts.Separator
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Key == "" {
			out[i] = e
			continue
		}
		out[i] = Entry{Key: prefix + sep + e.Key, Value: e.Value, Comment: e.Comment}
	}
	return out
}

// StripPrefix removes a leading prefix (plus separator) from every key.
// Keys that do not carry the prefix are left unchanged unless strict is true,
// in which case they are omitted from the result.
func StripPrefix(entries []Entry, prefix string, opts PrefixOptions, strict bool) []Entry {
	if prefix == "" {
		return Clone(entries)
	}
	head := prefix + opts.Separator
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if strings.HasPrefix(e.Key, head) {
			out = append(out, Entry{Key: e.Key[len(head):], Value: e.Value, Comment: e.Comment})
		} else if !strict {
			out = append(out, e)
		}
	}
	return out
}

// CommonPrefix returns the longest prefix shared by all keys,
// up to and including the separator. Returns "" when no common
// prefix exists or entries is empty.
func CommonPrefix(entries []Entry, opts PrefixOptions) string {
	if len(entries) == 0 {
		return ""
	}
	sep := opts.Separator
	ref := entries[0].Key
	for _, e := range entries[1:] {
		for !strings.HasPrefix(e.Key, ref) {
			idx := strings.LastIndex(ref, sep)
			if idx < 0 {
				return ""
			}
			ref = ref[:idx]
		}
	}
	// ref must end at a separator boundary
	if idx := strings.LastIndex(ref, sep); idx >= 0 {
		return ref[:idx]
	}
	return ""
}
