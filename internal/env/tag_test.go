package env

import (
	"testing"
)

func tagEntries() []Entry {
	return []Entry{
		{Key: "# @tag:db DB_HOST,DB_PORT,DB_NAME"},
		{Key: "# @tag:cache REDIS_HOST,REDIS_PORT"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "DB_NAME", Value: "mydb"},
		{Key: "REDIS_HOST", Value: "localhost"},
		{Key: "REDIS_PORT", Value: "6379"},
		{Key: "APP_SECRET", Value: "secret"},
	}
}

func TestParseTagsFindsAnnotations(t *testing.T) {
	tags := ParseTags(tagEntries())
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if len(tags["db"]) != 3 {
		t.Errorf("expected 3 keys for db tag, got %d", len(tags["db"]))
	}
	if len(tags["cache"]) != 2 {
		t.Errorf("expected 2 keys for cache tag, got %d", len(tags["cache"]))
	}
}

func TestParseTagsNoAnnotations(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	tags := ParseTags(entries)
	if len(tags) != 0 {
		t.Errorf("expected empty TagMap, got %v", tags)
	}
}

func TestFilterByTagReturnsMatchingKeys(t *testing.T) {
	tags := ParseTags(tagEntries())
	result, err := FilterByTag(tagEntries(), tags, "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != "DB_HOST" && e.Key != "DB_PORT" && e.Key != "DB_NAME" {
			t.Errorf("unexpected key %q in result", e.Key)
		}
	}
}

func TestFilterByTagUnknownReturnsError(t *testing.T) {
	tags := ParseTags(tagEntries())
	_, err := FilterByTag(tagEntries(), tags, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown tag, got nil")
	}
}

func TestTagNamesSorted(t *testing.T) {
	tags := TagMap{"zebra": nil, "alpha": nil, "middle": nil}
	names := TagNames(tags)
	expected := []string{"alpha", "middle", "zebra"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("pos %d: expected %q got %q", i, expected[i], n)
		}
	}
}
