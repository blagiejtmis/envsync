package env

import (
	"testing"
)

var castEntries = []Entry{
	{Key: "STRING_VAL", Value: "hello"},
	{Key: "BOOL_TRUE", Value: "true"},
	{Key: "BOOL_YES", Value: "yes"},
	{Key: "BOOL_ONE", Value: "1"},
	{Key: "BOOL_FALSE", Value: "false"},
	{Key: "BOOL_NO", Value: "no"},
	{Key: "INT_VAL", Value: "42"},
	{Key: "NEG_INT", Value: "-7"},
	{Key: "FLOAT_VAL", Value: "3.14"},
	{Key: "LIST_VAL", Value: "a,b,c"},
	{Key: "BAD_INT", Value: "notanint"},
	{Key: "BAD_BOOL", Value: "maybe"},
}

func TestAsString(t *testing.T) {
	v, err := AsString(castEntries, "STRING_VAL")
	if err != nil || v != "hello" {
		t.Fatalf("expected hello, got %q err %v", v, err)
	}
}

func TestAsStringNotFound(t *testing.T) {
	_, err := AsString(castEntries, "MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestAsBoolTrue(t *testing.T) {
	for _, key := range []string{"BOOL_TRUE", "BOOL_YES", "BOOL_ONE"} {
		v, err := AsBool(castEntries, key)
		if err != nil || !v {
			t.Errorf("key %s: expected true, got %v err %v", key, v, err)
		}
	}
}

func TestAsBoolFalse(t *testing.T) {
	for _, key := range []string{"BOOL_FALSE", "BOOL_NO"} {
		v, err := AsBool(castEntries, key)
		if err != nil || v {
			t.Errorf("key %s: expected false, got %v err %v", key, v, err)
		}
	}
}

func TestAsBoolInvalid(t *testing.T) {
	_, err := AsBool(castEntries, "BAD_BOOL")
	if err == nil {
		t.Fatal("expected error for invalid bool")
	}
	if ce, ok := err.(*CastError); !ok || ce.TargetType != "bool" {
		t.Fatalf("expected CastError with type bool, got %T", err)
	}
}

func TestAsInt(t *testing.T) {
	v, err := AsInt(castEntries, "INT_VAL")
	if err != nil || v != 42 {
		t.Fatalf("expected 42, got %d err %v", v, err)
	}
}

func TestAsIntNegative(t *testing.T) {
	v, err := AsInt(castEntries, "NEG_INT")
	if err != nil || v != -7 {
		t.Fatalf("expected -7, got %d err %v", v, err)
	}
}

func TestAsIntInvalid(t *testing.T) {
	_, err := AsInt(castEntries, "BAD_INT")
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestAsFloat(t *testing.T) {
	v, err := AsFloat(castEntries, "FLOAT_VAL")
	if err != nil || v < 3.13 || v > 3.15 {
		t.Fatalf("expected ~3.14, got %f err %v", v, err)
	}
}

func TestAsStringSlice(t *testing.T) {
	v, err := AsStringSlice(castEntries, "LIST_VAL", ",")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 3 || v[0] != "a" || v[1] != "b" || v[2] != "c" {
		t.Fatalf("unexpected slice: %v", v)
	}
}

func TestAsStringSliceNotFound(t *testing.T) {
	_, err := AsStringSlice(castEntries, "MISSING", ",")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
