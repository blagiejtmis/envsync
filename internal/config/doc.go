// Package config provides loading and validation of the envsync
// project configuration file (.envsync.toml).
//
// A minimal configuration file looks like:
//
//	[store]
//	path = ".envsync-store"
//
//	[env]
//	file = ".env"
//
// The store.path field is required; env.file defaults to ".env" when
// omitted. Use config.Load to read a file from disk, or config.Default
// to obtain a Config value pre-populated with sensible defaults.
package config
