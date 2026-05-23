// Package env provides utilities for working with .env files.
package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CastError is returned when a value cannot be cast to the requested type.
type CastError struct {
	Key      string
	Value    string
	TargetType string
	Err      error
}

func (e *CastError) Error() string {
	return fmt.Sprintf("cast: key %q value %q cannot be cast to %s: %v", e.Key, e.Value, e.TargetType, e.Err)
}

// AsString returns the raw string value for the given key.
// Returns an error if the key is not found.
func AsString(entries []Entry, key string) (string, error) {
	for _, e := range entries {
		if e.Key == key {
			return e.Value, nil
		}
	}
	return "", fmt.Errorf("cast: key %q not found", key)
}

// AsBool returns the value for key interpreted as a boolean.
// Accepts: true/false, 1/0, yes/no (case-insensitive).
func AsBool(entries []Entry, key string) (bool, error) {
	v, err := AsString(entries, key)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	}
	return false, &CastError{Key: key, Value: v, TargetType: "bool", Err: fmt.Errorf("unrecognised value")}
}

// AsInt returns the value for key interpreted as an int64.
func AsInt(entries []Entry, key string) (int64, error) {
	v, err := AsString(entries, key)
	if err != nil {
		return 0, err
	}
	n, parseErr := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if parseErr != nil {
		return 0, &CastError{Key: key, Value: v, TargetType: "int64", Err: parseErr}
	}
	return n, nil
}

// AsFloat returns the value for key interpreted as a float64.
func AsFloat(entries []Entry, key string) (float64, error) {
	v, err := AsString(entries, key)
	if err != nil {
		return 0, err
	}
	f, parseErr := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if parseErr != nil {
		return 0, &CastError{Key: key, Value: v, TargetType: "float64", Err: parseErr}
	}
	return f, nil
}

// AsStringSlice splits the value for key by sep and returns the parts.
func AsStringSlice(entries []Entry, key, sep string) ([]string, error) {
	v, err := AsString(entries, key)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(v, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result, nil
}
