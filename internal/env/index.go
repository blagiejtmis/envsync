package env

// Index builds a map from key to Entry for fast lookups.
// The returned map holds pointers into the original slice;
// callers should not modify entries through the map.
func Index(entries []Entry) map[string]*Entry {
	idx := make(map[string]*Entry, len(entries))
	for i := range entries {
		idx[entries[i].Key] = &entries[i]
	}
	return idx
}

// Lookup returns the value for key and whether it was found.
func Lookup(entries []Entry, key string) (string, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

// Contains reports whether key exists in entries.
func Contains(entries []Entry, key string) bool {
	_, ok := Lookup(entries, key)
	return ok
}

// Keys returns the ordered list of keys from entries.
func Keys(entries []Entry) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return keys
}

// Values returns the ordered list of values from entries.
func Values(entries []Entry) []string {
	vals := make([]string, len(entries))
	for i, e := range entries {
		vals[i] = e.Value
	}
	return vals
}

// Pick returns only the entries whose keys are in the provided set.
// Order follows the original entries slice.
func Pick(entries []Entry, keys ...string) []Entry {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	out := make([]Entry, 0, len(keys))
	for _, e := range entries {
		if _, ok := set[e.Key]; ok {
			out = append(out, e)
		}
	}
	return out
}
