package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func pinEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "API_KEY", Value: "secret"},
	}
}

func TestPinFoundKey(t *testing.T) {
	entries := pinEntries()
	p, err := Pin(entries, "DB_HOST", "alice", "keep stable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Key != "DB_HOST" || p.Value != "localhost" {
		t.Errorf("unexpected pin: %+v", p)
	}
	if p.PinnedBy != "alice" || p.Comment != "keep stable" {
		t.Errorf("metadata not set correctly: %+v", p)
	}
	if p.PinnedAt.IsZero() {
		t.Error("PinnedAt should not be zero")
	}
}

func TestPinMissingKey(t *testing.T) {
	_, err := Pin(pinEntries(), "MISSING", "bob", "")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestSaveAndLoadPins(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pins.json")
	pf := PinFile{
		Pins: []PinEntry{
			{Key: "DB_HOST", Value: "localhost", PinnedAt: time.Now().UTC(), PinnedBy: "alice"},
		},
	}
	if err := SavePins(path, pf); err != nil {
		t.Fatalf("SavePins: %v", err)
	}
	loaded, err := LoadPins(path)
	if err != nil {
		t.Fatalf("LoadPins: %v", err)
	}
	if len(loaded.Pins) != 1 || loaded.Pins[0].Key != "DB_HOST" {
		t.Errorf("unexpected loaded pins: %+v", loaded)
	}
	if loaded.Version != 1 {
		t.Errorf("expected version 1, got %d", loaded.Version)
	}
}

func TestLoadPinsNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "noexist.json")
	pf, err := LoadPins(path)
	if err != nil {
		t.Fatalf("expected empty PinFile, got error: %v", err)
	}
	if len(pf.Pins) != 0 {
		t.Errorf("expected empty pins, got %+v", pf)
	}
}

func TestUnpin(t *testing.T) {
	pf := PinFile{
		Pins: []PinEntry{
			{Key: "DB_HOST", Value: "localhost"},
			{Key: "API_KEY", Value: "secret"},
		},
	}
	if !Unpin(&pf, "DB_HOST") {
		t.Fatal("expected Unpin to return true")
	}
	if len(pf.Pins) != 1 || pf.Pins[0].Key != "API_KEY" {
		t.Errorf("unexpected pins after unpin: %+v", pf.Pins)
	}
	if Unpin(&pf, "MISSING") {
		t.Error("expected Unpin to return false for missing key")
	}
}

func TestPinnedKeys(t *testing.T) {
	pf := PinFile{
		Pins: []PinEntry{
			{Key: "X", Value: "1"},
			{Key: "Y", Value: "2"},
		},
	}
	m := PinnedKeys(pf)
	if len(m) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m))
	}
	if m["X"].Value != "1" || m["Y"].Value != "2" {
		t.Errorf("unexpected map contents: %+v", m)
	}
}

func TestSavePinsCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	if err := SavePins(path, PinFile{}); err != nil {
		t.Fatalf("SavePins: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
