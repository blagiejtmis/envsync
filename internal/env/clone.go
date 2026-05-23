package env

// Clone returns a deep copy of a slice of Entry values.
// Modifications to the returned slice do not affect the original.
func Clone(entries []Entry) []Entry {
	if entries == nil {
		return nil
	}
	out := make([]Entry, len(entries))
	copy(out, entries)
	return out
}

// CloneMap returns a deep copy of a map[string]string.
func CloneMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// Unique returns entries with duplicate keys removed, keeping the last
// occurrence of each key (consistent with shell semantics where the last
// assignment wins).
func Unique(entries []Entry) []Entry {
	seen := make(map[string]int, len(entries))
	for i, e := range entries {
		seen[e.Key] = i
	}
	out := make([]Entry, 0, len(seen))
	for i, e := range entries {
		if seen[e.Key] == i {
			out = append(out, e)
		}
	}
	return out
}

// Reorder returns a new slice containing only the entries whose keys appear in
// keys, in the order given by keys. Entries not present in keys are appended
// at the end in their original relative order.
func Reorder(entries []Entry, keys []string) []Entry {
	index := make(map[string]Entry, len(entries))
	for _, e := range entries {
		index[e.Key] = e
	}
	out := make([]Entry, 0, len(entries))
	added := make(map[string]bool, len(keys))
	for _, k := range keys {
		if e, ok := index[k]; ok {
			out = append(out, e)
			added[k] = true
		}
	}
	for _, e := range entries {
		if !added[e.Key] {
			out = append(out, e)
		}
	}
	return out
}
