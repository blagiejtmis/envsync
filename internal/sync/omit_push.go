package sync

import (
	"github.com/user/envsync/internal/env"
)

// OmitPushOptions configures which entries are stripped before pushing.
type OmitPushOptions struct {
	// Keys lists exact keys to omit from the push payload.
	Keys []string
	// Prefixes removes entries whose key starts with any of these prefixes.
	Prefixes []string
	// Blank removes entries with empty values before pushing.
	Blank bool
	// Sensitive removes sensitive entries (passwords, secrets, tokens) before
	// pushing. Use with care — this may cause remote state to lose values.
	Sensitive bool
}

// DefaultOmitPushOptions returns OmitPushOptions with no omissions enabled.
func DefaultOmitPushOptions() OmitPushOptions {
	return OmitPushOptions{}
}

// applyOmitPush strips entries from the push payload according to opts.
// It returns the filtered slice and the list of omitted keys for logging.
func applyOmitPush(entries []env.Entry, opts OmitPushOptions) ([]env.Entry, []string) {
	envOpts := env.OmitOptions{
		Keys:      opts.Keys,
		Prefixes:  opts.Prefixes,
		Blank:     opts.Blank,
		Sensitive: opts.Sensitive,
	}
	omitted := env.OmittedKeys(entries, envOpts)
	filtered := env.Omit(entries, envOpts)
	return filtered, omitted
}
