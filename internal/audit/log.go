// Package audit provides a simple audit log for tracking env sync operations.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventPush   EventKind = "push"
	EventPull   EventKind = "pull"
	EventDelete EventKind = "delete"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	Key       string    `json:"key"`
	User      string    `json:"user,omitempty"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit entries to a file in JSON-lines format.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the file at path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an Entry to the audit log file.
func (l *Logger) Record(kind EventKind, key, user, message string) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Key:       key,
		User:      user,
		Message:   message,
	}
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// ReadAll reads and returns all audit entries from the log file.
func (l *Logger) ReadAll() ([]Entry, error) {
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
