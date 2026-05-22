package env

import (
	"strings"
	"testing"
)

func schemaLines(s string) []string {
	return strings.Split(s, "\n")
}

func TestParseSchemaBasic(t *testing.T) {
	input := schemaLines("APP_NAME=\nAPP_PORT=8080")
	schema, err := ParseSchema(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(schema))
	}
	if schema[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", schema[0].Key)
	}
}

func TestParseSchemaAnnotations(t *testing.T) {
	input := schemaLines("# @required\n# Database URL\n# @default localhost\nDB_URL=")
	schema, err := ParseSchema(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(schema))
	}
	e := schema[0]
	if !e.Required {
		t.Error("expected Required=true")
	}
	if e.Default != "localhost" {
		t.Errorf("expected default localhost, got %q", e.Default)
	}
	if e.Comment != "Database URL" {
		t.Errorf("unexpected comment: %q", e.Comment)
	}
}

func TestParseSchemaInvalidLine(t *testing.T) {
	_, err := ParseSchema(schemaLines("NOTVALID"))
	if err == nil {
		t.Error("expected error for invalid schema line")
	}
}

func TestValidateAgainstSchemaAllPresent(t *testing.T) {
	schema := Schema{
		{Key: "APP_NAME", Required: true},
		{Key: "APP_PORT", Required: false},
	}
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_PORT", Value: "8080"},
	}
	violations := ValidateAgainstSchema(entries, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got: %v", violations)
	}
}

func TestValidateAgainstSchemaMissingRequired(t *testing.T) {
	schema := Schema{
		{Key: "SECRET_KEY", Required: true},
	}
	violations := ValidateAgainstSchema([]Entry{}, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0], "SECRET_KEY") {
		t.Errorf("violation should mention key name: %s", violations[0])
	}
}

func TestValidateAgainstSchemaBlankRequired(t *testing.T) {
	schema := Schema{
		{Key: "DB_PASS", Required: true},
	}
	entries := []Entry{{Key: "DB_PASS", Value: ""}}
	violations := ValidateAgainstSchema(entries, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateAgainstSchemaOptionalMissing(t *testing.T) {
	schema := Schema{
		{Key: "LOG_LEVEL", Required: false, Default: "info"},
	}
	violations := ValidateAgainstSchema([]Entry{}, schema)
	if len(violations) != 0 {
		t.Errorf("optional missing key should not be a violation")
	}
}
