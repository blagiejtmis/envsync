package env

import "sort"

// Group holds a named collection of env entries.
type Group struct {
	Name    string
	Entries []Entry
}

// GroupOptions controls how entries are grouped.
type GroupOptions struct {
	// PrefixSep is the separator used to split a key into group + remainder.
	// Defaults to "_".
	PrefixSep string

	// KeepPrefix, when true, retains the full original key inside each group.
	// When false, the group prefix and separator are stripped from the key.
	KeepPrefix bool
}

// DefaultGroupOptions returns sensible defaults.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		PrefixSep:  "_",
		KeepPrefix: true,
	}
}

// GroupByPrefix splits entries into named groups based on their key prefix.
// Keys that have no separator are placed in a group named "".
func GroupByPrefix(entries []Entry, opts GroupOptions) []Group {
	if opts.PrefixSep == "" {
		opts.PrefixSep = "_"
	}

	index := make(map[string][]Entry)
	var order []string
	seen := make(map[string]bool)

	for _, e := range entries {
		groupName, member := splitPrefix(e.Key, opts.PrefixSep)

		entry := e
		if !opts.KeepPrefix && groupName != "" {
			entry.Key = member
		}

		if !seen[groupName] {
			seen[groupName] = true
			order = append(order, groupName)
		}
		index[groupName] = append(index[groupName], entry)
	}

	groups := make([]Group, 0, len(order))
	for _, name := range order {
		groups = append(groups, Group{Name: name, Entries: index[name]})
	}
	return groups
}

// GroupNames returns a sorted slice of unique group names for the given entries.
func GroupNames(entries []Entry, sep string) []string {
	if sep == "" {
		sep = "_"
	}
	seen := make(map[string]bool)
	for _, e := range entries {
		name, _ := splitPrefix(e.Key, sep)
		seen[name] = true
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// splitPrefix returns the prefix and the remainder of key split on the first
// occurrence of sep. If sep is not found both prefix and remainder equal key.
func splitPrefix(key, sep string) (prefix, remainder string) {
	for i := 0; i <= len(key)-len(sep); i++ {
		if key[i:i+len(sep)] == sep {
			return key[:i], key[i+len(sep):]
		}
	}
	return key, key
}
