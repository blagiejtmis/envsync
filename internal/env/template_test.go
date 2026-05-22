package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateTemplateRedactsValues(t *testing.T) {
	src := []Entry{
		{Key: "DB_URL", Value: "postgres://secret"},
		{Key: "API_KEY", Value: "abc123"},
	}
	opts := DefaultTemplateOptions()
	tpl := GenerateTemplate(src, opts)
	for _, e := range tpl {
		if e.Value != "" {
			t.Errorf("expected blank value for %s, got %q", e.Key, e.Value)
		}
	}
}

func TestGenerateTemplatePreservesKeys(t *testing.T) {
	src := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	tpl := GenerateTemplate(src, DefaultTemplateOptions())
	if len(tpl) != len(src) {
		t.Fatalf("expected %d entries, got %d", len(src), len(tpl))
	}
	for i, e := range tpl {
		if e.Key != src[i].Key {
			t.Errorf("key mismatch at %d: want %s got %s", i, src[i].Key, e.Key)
		}
	}
}

func TestApplyTemplateAddsNewKeys(t *testing.T) {
	dst := []Entry{{Key: "EXISTING", Value: "val"}}
	tmpl := []Entry{{Key: "EXISTING", Value: ""}, {Key: "NEW_KEY", Value: ""}}
	result, added := ApplyTemplate(dst, tmpl)
	if len(added) != 1 || added[0] != "NEW_KEY" {
		t.Fatalf("expected [NEW_KEY] added, got %v", added)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key == "EXISTING" && e.Value != "val" {
			t.Error("existing value should not be overwritten")
		}
	}
}

func TestCheckMissingKeys(t *testing.T) {
	env := []Entry{{Key: "A", Value: "filled"}, {Key: "B", Value: ""}}
	tmpl := []Entry{{Key: "A", Value: ""}, {Key: "B", Value: ""}, {Key: "C", Value: ""}}
	missing := CheckMissingKeys(env, tmpl)
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing, got %v", missing)
	}
}

func TestWriteTemplateFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.template")
	entries := []Entry{{Key: "SECRET", Value: "hidden"}, {Key: "PORT", Value: "8080"}}
	if err := WriteTemplateFile(path, entries, DefaultTemplateOptions()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("template file should not be empty")
	}
	// values must be blank
	parsed, err := Parse(string(data))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range parsed {
		if e.Value != "" {
			t.Errorf("expected blank value for %s in template, got %q", e.Key, e.Value)
		}
	}
}
