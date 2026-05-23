package env

import "fmt"

// RenameOptions controls the behaviour of Rename.
type RenameOptions struct {
	// FailOnMissing returns an error when the source key does not exist.
	// When false the entry slice is returned unchanged.
	FailOnMissing bool

	// FailOnCollision returns an error when the destination key already exists.
	// When false the existing destination entry is overwritten.
	FailOnCollision bool
}

// DefaultRenameOptions returns sensible defaults for Rename.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		FailOnMissing:   true,
		FailOnCollision: true,
	}
}

// Rename returns a new slice where the entry with key src is renamed to dst.
// The relative order of entries is preserved; the renamed entry keeps its
// original position in the slice.
//
// Rename does not mutate the input slice.
func Rename(entries []Entry, src, dst string, opts RenameOptions) ([]Entry, error) {
	srcIdx := -1
	dstIdx := -1

	for i, e := range entries {
		switch e.Key {
		case src:
			srcIdx = i
		case dst:
			dstIdx = i
		}
	}

	if srcIdx == -1 {
		if opts.FailOnMissing {
			return nil, fmt.Errorf("env rename: key %q not found", src)
		}
		// Return a shallow copy so the caller always owns the slice.
		out := make([]Entry, len(entries))
		copy(out, entries)
		return out, nil
	}

	if dstIdx != -1 && opts.FailOnCollision {
		return nil, fmt.Errorf("env rename: destination key %q already exists", dst)
	}

	out := make([]Entry, 0, len(entries))
	for i, e := range entries {
		switch {
		case i == srcIdx:
			out = append(out, Entry{Key: dst, Value: e.Value})
		case i == dstIdx:
			// Drop the old destination entry when collision is allowed.
			continue
		default:
			out = append(out, e)
		}
	}
	return out, nil
}
