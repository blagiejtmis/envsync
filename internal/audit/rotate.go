package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions configures log rotation behaviour.
type RotateOptions struct {
	// MaxBytes is the maximum size of the log file before rotation.
	MaxBytes int64
	// MaxBackups is the number of rotated files to keep (0 = keep all).
	MaxBackups int
}

// DefaultRotateOptions returns sensible defaults for log rotation.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		MaxBytes:   1 << 20, // 1 MiB
		MaxBackups: 5,
	}
}

// Rotate checks whether the log file at path exceeds opts.MaxBytes and, if so,
// renames it to a timestamped backup and removes old backups beyond MaxBackups.
// It returns true when a rotation was performed.
func Rotate(path string, opts RotateOptions) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("audit rotate: stat %s: %w", path, err)
	}

	if info.Size() < opts.MaxBytes {
		return false, nil
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]
	backup := fmt.Sprintf("%s.%s%s", base, timestamp, ext)

	if err := os.Rename(path, backup); err != nil {
		return false, fmt.Errorf("audit rotate: rename: %w", err)
	}

	if opts.MaxBackups > 0 {
		if err := pruneBackups(base, ext, opts.MaxBackups); err != nil {
			return true, err
		}
	}

	return true, nil
}

// pruneBackups removes the oldest backup files so that at most max remain.
func pruneBackups(base, ext string, max int) error {
	pattern := base + ".*" + ext
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("audit rotate: glob backups: %w", err)
	}

	// Glob returns names in lexicographic order; oldest timestamps sort first.
	for len(matches) > max {
		if err := os.Remove(matches[0]); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("audit rotate: remove old backup %s: %w", matches[0], err)
		}
		matches = matches[1:]
	}
	return nil
}
