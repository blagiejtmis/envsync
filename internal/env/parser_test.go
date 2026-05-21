package env

import (
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	input := `
DB_HOST=localhost
DB_PORT=5432
APP_SECRET="mysecret"
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "DB_HOST" || entries[0].Value != "localhost" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
	if entries[2].Value != "mysecret" {
		t.Errorf("expected quotes stripped, got %q", entries[2].Value)
	}
}

func TestParseSkipsCommentsAndBlanks(t *testing.T) {
	input := `
# This is a comment

FOO=bar
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestParseInvalidLine(t *testing.T) {
	input := "INVALID_LINE_NO_EQUALS\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseEmptyKey(t *testing.T) {
	input := "=value\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSerialize(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
	out := Serialize(entries)
	expected := "HOST=localhost\nPORT=8080\n"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestParseSerializeRoundtrip(t *testing.T) {
	original := "KEY1=value1\nKEY2=value2\n"
	entries, err := Parse(strings.NewReader(original))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	result := Serialize(entries)
	if result != original {
		t.Errorf("roundtrip mismatch: got %q", result)
	}
}
