// Package env provides utilities for parsing, serializing, and managing
// .env file contents.
package env

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Label     string    `json:"label,omitempty"`
	Entries   []Entry   `json:"entries"`
	Checksum  string    `json:"checksum"`
}

// TakeSnapshot creates a Snapshot from a slice of entries with an optional label.
func TakeSnapshot(entries []Entry, label string) Snapshot {
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
		Checksum:  checksumEntries(entries),
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at the given path.
func SaveSnapshot(path string, snap Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from a JSON file at the given path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return snap, nil
}

// Verify checks whether the snapshot's checksum matches its entries.
func (s Snapshot) Verify() bool {
	return checksumEntries(s.Entries) == s.Checksum
}

// checksumEntries produces a deterministic hex checksum over sorted key=value pairs.
func checksumEntries(entries []Entry) string {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	h := sha256.New()
	for _, e := range sorted {
		fmt.Fprintf(h, "%s=%s\n", e.Key, e.Value)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
