package env

import (
	"fmt"
	"os"
	"sort"
)

// LoadFile reads and parses a .env file from the given path.
func LoadFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	entries, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parsing env file %q: %w", path, err)
	}
	return entries, nil
}

// WriteFile serializes entries and writes them to path, creating or truncating the file.
func WriteFile(path string, entries []Entry) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("opening env file %q for writing: %w", path, err)
	}
	defer f.Close()

	_, err = fmt.Fprint(f, Serialize(entries))
	if err != nil {
		return fmt.Errorf("writing env file %q: %w", path, err)
	}
	return nil
}

// ToMap converts a slice of Entries into a key-value map.
// If duplicate keys exist, the last value wins.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// FromMap converts a map into a slice of Entries sorted by key for deterministic output.
func FromMap(m map[string]string) []Entry {
	entries := make([]Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}
