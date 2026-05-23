// Package sync provides functionality for pushing and pulling .env files
// to and from a shared secret store.
//
// # Filter Push
//
// The filter-push feature allows callers to restrict which environment
// entries are uploaded to the store during a push operation. This is
// useful when a repository contains entries that should remain local
// (e.g. developer-specific overrides or secrets that differ per machine).
//
// Filtering is applied after interpolation and before encryption, so the
// stored payload only contains the intended subset of keys.
//
// Usage:
//
//	opts := sync.FilterPushOptions{
//		Prefix:      "APP_",
//		ExcludeKeys: []string{"APP_LOCAL_OVERRIDE"},
//	}
package sync
