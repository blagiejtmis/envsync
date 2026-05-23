// Package sync provides push/pull synchronisation between local .env files
// and the shared secret store.
package sync

import "github.com/yourorg/envsync/internal/env"

// TransformPushOptions controls entry transformation applied before pushing
// entries to the secret store.
type TransformPushOptions struct {
	// UppercaseKeys normalises all keys to uppercase before pushing.
	UppercaseKeys bool
	// TrimValues strips leading/trailing whitespace from values before pushing.
	TrimValues bool
	// StripQuotes removes surrounding quotes from values before pushing.
	StripQuotes bool
}

// DefaultTransformPushOptions returns a TransformPushOptions with safe
// defaults that match common CI/CD expectations.
func DefaultTransformPushOptions() TransformPushOptions {
	return TransformPushOptions{
		UppercaseKeys: false,
		TrimValues:    true,
		StripQuotes:   false,
	}
}

// applyTransformPush transforms entries according to opts before they are
// handed off to the encryptor and stored. It returns a new slice and never
// mutates the input.
func applyTransformPush(entries []env.Entry, opts TransformPushOptions) []env.Entry {
	topts := env.DefaultTransformOptions()
	topts.TrimKeys = true
	topts.TrimValues = opts.TrimValues
	topts.UppercaseKeys = opts.UppercaseKeys
	topts.StripQuotes = opts.StripQuotes
	return env.Transform(entries, topts)
}
