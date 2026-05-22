package env

import (
	"fmt"
	"strings"
)

// ValidationError describes a single validation issue on a key.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) OK() bool { return len(r.Errors) == 0 }

func (r *ValidationResult) Error() string {
	if r.OK() {
		return ""
	}
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// ValidateKeys checks that every key in entries conforms to the
// conventional env-var naming rules: non-empty, starts with a letter
// or underscore, and contains only letters, digits, and underscores.
func ValidateKeys(entries []Entry) *ValidationResult {
	res := &ValidationResult{}
	for _, e := range entries {
		if err := validateKey(e.Key); err != nil {
			res.Errors = append(res.Errors, ValidationError{Key: e.Key, Message: err.Error()})
		}
	}
	return res
}

// ValidateNoBlanks returns errors for any entry whose value is the
// empty string, useful as an optional lint step before pushing.
func ValidateNoBlanks(entries []Entry) *ValidationResult {
	res := &ValidationResult{}
	for _, e := range entries {
		if strings.TrimSpace(e.Value) == "" {
			res.Errors = append(res.Errors, ValidationError{
				Key:     e.Key,
				Message: "value is blank",
			})
		}
	}
	return res
}

func validateKey(k string) error {
	if k == "" {
		return fmt.Errorf("key must not be empty")
	}
	for i, ch := range k {
		switch {
		case ch == '_':
			// always ok
		case ch >= 'A' && ch <= 'Z':
			// ok
		case ch >= 'a' && ch <= 'z':
			// ok
		case ch >= '0' && ch <= '9':
			if i == 0 {
				return fmt.Errorf("must start with a letter or underscore, got %q", ch)
			}
		default:
			return fmt.Errorf("invalid character %q at position %d", ch, i)
		}
	}
	return nil
}
