package env

import (
	"testing"
)

func normEntries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestNormalizeUppercaseKeys(t *testing.T) {
	input := normEntries("db_host", "localhost", "api_key", "secret")
	got := Normalize(input, NormalizeOptions{UppercaseKeys: true})
	if got[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", got[0].Key)
	}
	if got[1].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", got[1].Key)
	}
}

func TestNormalizeTrimValues(t *testing.T) {
	input := normEntries("HOST", "  localhost  ", "PORT", "\t8080\n")
	got := Normalize(input, NormalizeOptions{TrimValues: true})
	if got[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", got[0].Value)
	}
	if got[1].Value != "8080" {
		t.Errorf("expected '8080', got %q", got[1].Value)
	}
}

func TestNormalizeRemovesDuplicatesKeepsLast(t *testing.T) {
	input := normEntries("HOST", "first", "PORT", "3000", "HOST", "last")
	got := Normalize(input, NormalizeOptions{RemoveDuplicates: true})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Value != "last" {
		t.Errorf("expected last value 'last', got %q", got[0].Value)
	}
}

func TestNormalizeDoesNotMutateOriginal(t *testing.T) {
	input := normEntries("db_pass", "  secret  ")
	orig := input[0].Value
	Normalize(input, DefaultNormalizeOptions())
	if input[0].Value != orig {
		t.Error("original slice was mutated")
	}
}

func TestNormalizeDefaultOptions(t *testing.T) {
	input := normEntries("api_url", "  http://example.com  ", "api_url", "  http://other.com  ")
	got := Normalize(input, DefaultNormalizeOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 entry after dedup, got %d", len(got))
	}
	if got[0].Key != "API_URL" {
		t.Errorf("expected API_URL, got %s", got[0].Key)
	}
	if got[0].Value != "http://other.com" {
		t.Errorf("expected trimmed value, got %q", got[0].Value)
	}
}

func TestNormalizeEmptyInput(t *testing.T) {
	got := Normalize(nil, DefaultNormalizeOptions())
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d entries", len(got))
	}
}
