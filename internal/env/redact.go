// Package env provides utilities for parsing, serializing, and manipulating
// .env files.
package env

import "strings"

// SensitivePatterns is a list of substrings that, when found (case-insensitively)
// in a key name, indicate the value should be redacted.
var SensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

// RedactedPlaceholder is the value substituted for sensitive entries.
const RedactedPlaceholder = "***REDACTED***"

// IsSensitive reports whether the given key matches any sensitive pattern.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range SensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// Redact returns a copy of entries where values for sensitive keys are replaced
// with RedactedPlaceholder. The original slice is not modified.
func Redact(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if IsSensitive(e.Key) {
			out[i] = Entry{Key: e.Key, Value: RedactedPlaceholder}
		} else {
			out[i] = e
		}
	}
	return out
}

// RedactMap returns a copy of m where values for sensitive keys are replaced
// with RedactedPlaceholder. The original map is not modified.
func RedactMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if IsSensitive(k) {
			out[k] = RedactedPlaceholder
		} else {
			out[k] = v
		}
	}
	return out
}
