package env

import (
	"testing"
)

func TestValidateKeysAllValid(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "_PRIVATE", Value: "secret"},
		{Key: "Port8080", Value: "8080"},
	}
	res := ValidateKeys(entries)
	if !res.OK() {
		t.Fatalf("expected no errors, got: %s", res.Error())
	}
}

func TestValidateKeysStartsWithDigit(t *testing.T) {
	entries := []Entry{{Key: "1INVALID", Value: "x"}}
	res := ValidateKeys(entries)
	if res.OK() {
		t.Fatal("expected error for key starting with digit")
	}
	if len(res.Errors) != 1 || res.Errors[0].Key != "1INVALID" {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

func TestValidateKeysInvalidCharacter(t *testing.T) {
	entries := []Entry{
		{Key: "BAD-KEY", Value: "v"},
		{Key: "ALSO BAD", Value: "v"},
	}
	res := ValidateKeys(entries)
	if res.OK() {
		t.Fatal("expected errors for invalid characters")
	}
	if len(res.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(res.Errors))
	}
}

func TestValidateKeysEmptyKey(t *testing.T) {
	entries := []Entry{{Key: "", Value: "v"}}
	res := ValidateKeys(entries)
	if res.OK() {
		t.Fatal("expected error for empty key")
	}
}

func TestValidateNoBlanksAllFilled(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "hello"},
		{Key: "B", Value: "world"},
	}
	res := ValidateNoBlanks(entries)
	if !res.OK() {
		t.Fatalf("expected no errors, got: %s", res.Error())
	}
}

func TestValidateNoBlanksDetectsBlanks(t *testing.T) {
	entries := []Entry{
		{Key: "FILLED", Value: "ok"},
		{Key: "EMPTY", Value: ""},
		{Key: "SPACES", Value: "   "},
	}
	res := ValidateNoBlanks(entries)
	if res.OK() {
		t.Fatal("expected errors for blank values")
	}
	if len(res.Errors) != 2 {
		t.Fatalf("expected 2 blank errors, got %d: %s", len(res.Errors), res.Error())
	}
}

func TestValidationResultErrorString(t *testing.T) {
	res := &ValidationResult{
		Errors: []ValidationError{
			{Key: "BAD", Message: "invalid character"},
		},
	}
	s := res.Error()
	if s == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestValidationResultOKWhenEmpty(t *testing.T) {
	res := &ValidationResult{}
	if !res.OK() {
		t.Fatal("empty ValidationResult should be OK")
	}
	if res.Error() != "" {
		t.Fatalf("expected empty error string, got %q", res.Error())
	}
}
