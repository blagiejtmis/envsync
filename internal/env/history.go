// Package env provides utilities for working with .env files.
package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a snapshot of env entries at a point in time.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Label     string    `json:"label,omitempty"`
	Entries   []Entry   `json:"entries"`
}

// HistoryFile is the collection of history entries stored on disk.
type HistoryFile struct {
	Entries []HistoryEntry `json:"entries"`
}

// AppendHistory appends a new history entry to the given history file path.
// The file is created if it does not exist. At most maxEntries are retained
// (oldest are pruned first).
func AppendHistory(path string, label string, entries []Entry, maxEntries int) error {
	hf, err := LoadHistory(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("history: load: %w", err)
	}

	hf.Entries = append(hf.Entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	})

	if maxEntries > 0 && len(hf.Entries) > maxEntries {
		hf.Entries = hf.Entries[len(hf.Entries)-maxEntries:]
	}

	return saveHistory(path, hf)
}

// LoadHistory loads the history file from path. Returns an empty HistoryFile
// and a wrapped os.ErrNotExist if the file does not exist.
func LoadHistory(path string) (HistoryFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return HistoryFile{}, err
	}
	var hf HistoryFile
	if err := json.Unmarshal(data, &hf); err != nil {
		return HistoryFile{}, fmt.Errorf("history: unmarshal: %w", err)
	}
	return hf, nil
}

func saveHistory(path string, hf HistoryFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(hf, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}
