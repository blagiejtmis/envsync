package env

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a named label that can be applied to a set of env keys.
type Tag struct {
	Name string
	Keys []string
}

// TagMap maps tag names to their associated keys.
type TagMap map[string][]string

// ParseTags parses tag annotations from env entries.
// Lines of the form: # @tag:<name> KEY1,KEY2 are recognised.
func ParseTags(entries []Entry) TagMap {
	tags := make(TagMap)
	for _, e := range entries {
		if !strings.HasPrefix(e.Key, "#") {
			continue
		}
		line := strings.TrimPrefix(e.Key, "#")
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "@tag:") {
			continue
		}
		rest := strings.TrimPrefix(line, "@tag:")
		parts := strings.Fields(rest)
		if len(parts) < 2 {
			continue
		}
		name := parts[0]
		keys := strings.Split(parts[1], ",")
		tags[name] = append(tags[name], keys...)
	}
	return tags
}

// FilterByTag returns entries whose keys belong to the given tag.
func FilterByTag(entries []Entry, tags TagMap, tagName string) ([]Entry, error) {
	keys, ok := tags[tagName]
	if !ok {
		return nil, fmt.Errorf("tag %q not found", tagName)
	}
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[strings.TrimSpace(k)] = struct{}{}
	}
	var out []Entry
	for _, e := range entries {
		if _, ok := set[e.Key]; ok {
			out = append(out, e)
		}
	}
	return out, nil
}

// TagNames returns sorted tag names from a TagMap.
func TagNames(tags TagMap) []string {
	names := make([]string, 0, len(tags))
	for n := range tags {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
