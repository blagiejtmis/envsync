package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestFprintSchema(t *testing.T) {
	schema := Schema{
		{Key: "APP_NAME", Required: true, Comment: "Application name"},
		{Key: "LOG_LEVEL", Required: false, Default: "info"},
	}
	var buf bytes.Buffer
	FprintSchema(&buf, schema)
	out := buf.String()
	if !strings.Contains(out, "APP_NAME") {
		t.Error("output should contain APP_NAME")
	}
	if !strings.Contains(out, "required") {
		t.Error("output should contain 'required'")
	}
	if !strings.Contains(out, "Application name") {
		t.Error("output should contain comment")
	}
	if !strings.Contains(out, "default=info") {
		t.Error("output should contain default value")
	}
}

func TestFprintSchemaViolations(t *testing.T) {
	var buf bytes.Buffer
	n := FprintSchemaViolations(&buf, []string{"missing required key: FOO", "required key has blank value: BAR"})
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Error("output should mention FOO")
	}
	if !strings.Contains(out, "BAR") {
		t.Error("output should mention BAR")
	}
}

func TestFprintSchemaViolationsEmpty(t *testing.T) {
	var buf bytes.Buffer
	n := FprintSchemaViolations(&buf, nil)
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
	if buf.Len() != 0 {
		t.Error("expected empty output for no violations")
	}
}

func TestSchemaDefaults(t *testing.T) {
	schema := Schema{
		{Key: "A", Default: "alpha"},
		{Key: "B", Default: ""},
		{Key: "C", Default: "gamma"},
	}
	defaults := SchemaDefaults(schema)
	if defaults["A"] != "alpha" {
		t.Errorf("expected alpha, got %s", defaults["A"])
	}
	if _, ok := defaults["B"]; ok {
		t.Error("B has no default and should not appear")
	}
	if defaults["C"] != "gamma" {
		t.Errorf("expected gamma, got %s", defaults["C"])
	}
}

func TestApplySchemaDefaults(t *testing.T) {
	schema := Schema{
		{Key: "HOST", Default: "localhost"},
		{Key: "PORT", Default: "5432"},
	}
	entries := []Entry{{Key: "HOST", Value: "db.internal"}}
	result := ApplySchemaDefaults(entries, schema)
	keys := make(map[string]string)
	for _, e := range result {
		keys[e.Key] = e.Value
	}
	if keys["HOST"] != "db.internal" {
		t.Error("existing value should not be overwritten")
	}
	if keys["PORT"] != "5432" {
		t.Errorf("expected default PORT=5432, got %s", keys["PORT"])
	}
}
