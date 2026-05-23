// Package env — cast.go
//
// # Type Casting for Env Entries
//
// The cast helpers allow callers to read typed values directly from a
// []Entry slice without manually converting strings.
//
// Supported target types:
//
//   - string  — raw value via AsString
//   - bool    — true/false, yes/no, 1/0 (case-insensitive) via AsBool
//   - int64   — decimal integer via AsInt
//   - float64 — decimal float via AsFloat
//   - []string — comma (or custom) separated list via AsStringSlice
//
// All functions return a descriptive CastError when the value cannot be
// converted, making it easy to surface actionable messages to the user.
//
// Example:
//
//	entries, _ := LoadFile(".env")
//	port, err := AsInt(entries, "PORT")
//	if err != nil {
//		log.Fatal(err)
//	}
package env
