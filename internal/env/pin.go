package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PinEntry records a pinned version of a key-value pair at a point in time.
type PinEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by"`
	Comment   string    `json:"comment,omitempty"`
}

// PinFile holds all pinned entries.
type PinFile struct {
	Version int        `json:"version"`
	Pins    []PinEntry `json:"pins"`
}

// Pin adds or updates a pin for the given key in entries.
// Returns an error if the key does not exist in entries.
func Pin(entries []Entry, key, pinnedBy, comment string) (PinEntry, error) {
	for _, e := range entries {
		if e.Key == key {
			return PinEntry{
				Key:      key,
				Value:    e.Value,
				PinnedAt: time.Now().UTC(),
				PinnedBy: pinnedBy,
				Comment:  comment,
			}, nil
		}
	}
	return PinEntry{}, fmt.Errorf("pin: key %q not found in entries", key)
}

// SavePins writes a PinFile to disk at path.
func SavePins(path string, pf PinFile) error {
	pf.Version = 1
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadPins reads a PinFile from disk. Returns an empty PinFile if not found.
func LoadPins(path string) (PinFile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return PinFile{Version: 1}, nil
	}
	if err != nil {
		return PinFile{}, fmt.Errorf("pin: read: %w", err)
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return PinFile{}, fmt.Errorf("pin: unmarshal: %w", err)
	}
	return pf, nil
}

// Unpin removes the pin for key from pf. Returns false if key was not pinned.
func Unpin(pf *PinFile, key string) bool {
	for i, p := range pf.Pins {
		if p.Key == key {
			pf.Pins = append(pf.Pins[:i], pf.Pins[i+1:]...)
			return true
		}
	}
	return false
}

// PinnedKeys returns a map of key -> PinEntry for quick lookup.
func PinnedKeys(pf PinFile) map[string]PinEntry {
	m := make(map[string]PinEntry, len(pf.Pins))
	for _, p := range pf.Pins {
		m[p.Key] = p
	}
	return m
}
