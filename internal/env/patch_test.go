package env

import (
	"testing"
)

func patchBase() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
}

func TestPatchSetExistingKey(t *testing.T) {
	out, err := Patch(patchBase(), []PatchOp{{Key: "DB_HOST", Value: "prod-db"}}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := findValue(out, "DB_HOST"); got != "prod-db" {
		t.Errorf("expected prod-db, got %s", got)
	}
}

func TestPatchAddNewKey(t *testing.T) {
	out, err := Patch(patchBase(), []PatchOp{{Key: "REDIS_URL", Value: "redis://localhost"}}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := findValue(out, "REDIS_URL"); got != "redis://localhost" {
		t.Errorf("expected redis://localhost, got %s", got)
	}
}

func TestPatchDisallowNewKey(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.AllowNewKeys = false
	_, err := Patch(patchBase(), []PatchOp{{Key: "NEW_KEY", Value: "v"}}, opts)
	if err == nil {
		t.Fatal("expected error for new key when AllowNewKeys=false")
	}
}

func TestPatchDeleteKey(t *testing.T) {
	out, err := Patch(patchBase(), []PatchOp{{Key: "DB_PORT", Delete: true}}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if indexOfKey(out, "DB_PORT") != -1 {
		t.Error("expected DB_PORT to be deleted")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestPatchDeleteMissingKeyError(t *testing.T) {
	_, err := Patch(patchBase(), []PatchOp{{Key: "MISSING", Delete: true}}, DefaultPatchOptions())
	if err == nil {
		t.Fatal("expected error when deleting missing key")
	}
}

func TestPatchDeleteMissingKeyIgnored(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.IgnoreMissingDeletes = true
	out, err := Patch(patchBase(), []PatchOp{{Key: "MISSING", Delete: true}}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries unchanged, got %d", len(out))
	}
}

func TestPatchDoesNotMutateOriginal(t *testing.T) {
	base := patchBase()
	_, err := Patch(base, []PatchOp{{Key: "APP_ENV", Value: "production"}}, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if base[0].Value != "development" {
		t.Error("original slice was mutated")
	}
}

func TestPatchEmptyKeyError(t *testing.T) {
	_, err := Patch(patchBase(), []PatchOp{{Key: "", Value: "v"}}, DefaultPatchOptions())
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

// findValue is a small helper for tests.
func findValue(entries []Entry, key string) string {
	for _, e := range entries {
		if e.Key == key {
			return e.Value
		}
	}
	return ""
}
