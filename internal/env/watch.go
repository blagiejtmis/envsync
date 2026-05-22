// Package env provides utilities for parsing, serializing, and monitoring .env files.
package env

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"time"
)

// ChangeEvent is emitted when a watched .env file is modified.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
}

// WatchOptions controls the behaviour of Watch.
type WatchOptions struct {
	// PollInterval is how often the file is stat-checked. Defaults to 2s.
	PollInterval time.Duration
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{PollInterval: 2 * time.Second}
}

// Watch polls the file at path and sends a ChangeEvent on the returned channel
// whenever the file content changes. The channel is closed when ctx is done.
func Watch(ctx context.Context, path string, opts WatchOptions) (<-chan ChangeEvent, error) {
	if opts.PollInterval <= 0 {
		opts.PollInterval = DefaultWatchOptions().PollInterval
	}

	initialHash, err := hashFile(path)
	if err != nil {
		return nil, fmt.Errorf("watch: initial hash: %w", err)
	}

	ch := make(chan ChangeEvent, 1)

	go func() {
		defer close(ch)
		current := initialHash
		ticker := time.NewTicker(opts.PollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				next, err := hashFile(path)
				if err != nil {
					continue // file may be temporarily unavailable
				}
				if next != current {
					ch <- ChangeEvent{Path: path, OldHash: current, NewHash: next}
					current = next
				}
			}
		}
	}()

	return ch, nil
}

// hashFile returns a hex SHA-256 digest of the file at path.
func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum), nil
}
