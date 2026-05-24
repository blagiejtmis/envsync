package env

import (
	"fmt"
	"strings"
)

// DefaultTruncateOptions returns sensible defaults for TruncateOptions.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxLength: 256,
		Suffix:    "...",
		OnlyValues: true,
	}
}

// TruncateOptions controls how Truncate behaves.
type TruncateOptions struct {
	// MaxLength is the maximum allowed length for a value (or key if OnlyValues is false).
	MaxLength int
	// Suffix is appended when a value is truncated. Counts toward MaxLength.
	Suffix string
	// OnlyValues, when true, truncates only values and leaves keys untouched.
	OnlyValues bool
	// Keys restricts truncation to the specified keys. Empty means all keys.
	Keys []string
}

// Truncate returns a new slice of entries where values (and optionally keys)
// exceeding MaxLength are clipped and the Suffix is appended.
func Truncate(entries []Entry, opts TruncateOptions) ([]Entry, error) {
	if opts.MaxLength <= 0 {
		return nil, fmt.Errorf("truncate: MaxLength must be greater than zero, got %d", opts.MaxLength)
	}
	if len(opts.Suffix) >= opts.MaxLength {
		return nil, fmt.Errorf("truncate: Suffix length %d must be less than MaxLength %d", len(opts.Suffix), opts.MaxLength)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}
	appliesTo := func(key string) bool {
		if len(keySet) == 0 {
			return true
		}
		_, ok := keySet[key]
		return ok
	}

	clip := func(s string) string {
		if len(s) <= opts.MaxLength {
			return s
		}
		return s[:opts.MaxLength-len(opts.Suffix)] + opts.Suffix
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		if appliesTo(e.Key) {
			e.Value = clip(e.Value)
			if !opts.OnlyValues {
				e.Key = strings.ToValidUTF8(clip(e.Key), "")
			}
		}
		out[i] = e
	}
	return out, nil
}
