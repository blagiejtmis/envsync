package env

import (
	"bytes"
	"strings"
	"testing"
)

func compareEntries() ([]Entry, []Entry) {
	a := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "ONLY_A", Value: "aval"},
	}
	b := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "remotehost"},
		{Key: "DB_PASS", Value: "newsecret"},
		{Key: "ONLY_B", Value: "bval"},
	}
	return a, b
}

func TestCompareIdentical(t *testing.T) {
	a := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	r := Compare(a, a)
	if r.HasDifferences() {
		t.Fatal("expected no differences for identical inputs")
	}
	if len(r.Identical) != 2 {
		t.Fatalf("expected 2 identical, got %d", len(r.Identical))
	}
}

func TestCompareDetectsChanges(t *testing.T) {
	a, b := compareEntries()
	r := Compare(a, b)

	if !r.HasDifferences() {
		t.Fatal("expected differences")
	}
	if len(r.OnlyInA) != 1 || r.OnlyInA[0].Key != "ONLY_A" {
		t.Errorf("expected ONLY_A in OnlyInA, got %v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0].Key != "ONLY_B" {
		t.Errorf("expected ONLY_B in OnlyInB, got %v", r.OnlyInB)
	}
	if len(r.Changed) != 2 {
		t.Errorf("expected 2 changed, got %d", len(r.Changed))
	}
	if len(r.Identical) != 1 || r.Identical[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME identical, got %v", r.Identical)
	}
}

func TestCompareSortedOutput(t *testing.T) {
	a := []Entry{{Key: "Z", Value: "1"}, {Key: "A", Value: "2"}, {Key: "M", Value: "3"}}
	b := []Entry{{Key: "Z", Value: "1"}, {Key: "A", Value: "9"}, {Key: "M", Value: "3"}}
	r := Compare(a, b)
	if len(r.Changed) != 1 || r.Changed[0].A.Key != "A" {
		t.Errorf("unexpected changed: %v", r.Changed)
	}
}

func TestFprintCompareNoChanges(t *testing.T) {
	a := []Entry{{Key: "X", Value: "1"}}
	r := Compare(a, a)
	var buf bytes.Buffer
	FprintCompare(&buf, r)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no differences message, got: %s", buf.String())
	}
}

func TestFprintCompareVerboseRedactsSensitive(t *testing.T) {
	a := []Entry{{Key: "DB_PASSWORD", Value: "old"}}
	b := []Entry{{Key: "DB_PASSWORD", Value: "new"}}
	r := Compare(a, b)
	var buf bytes.Buffer
	FprintCompareVerbose(&buf, r)
	if strings.Contains(buf.String(), "old") || strings.Contains(buf.String(), "new") {
		t.Errorf("sensitive value should be redacted, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "[redacted]") {
		t.Errorf("expected [redacted] in output, got: %s", buf.String())
	}
}

func TestCompareSummary(t *testing.T) {
	a, b := compareEntries()
	r := Compare(a, b)
	s := CompareSummary(r)
	if !strings.Contains(s, "added") || !strings.Contains(s, "removed") {
		t.Errorf("unexpected summary: %s", s)
	}
}
