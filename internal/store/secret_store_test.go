package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/store"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "secrets.json")
}

func TestPutAndGet(t *testing.T) {
	s, err := store.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ciphertext := []byte("encrypted-payload")
	if err := s.Put("DB_PASSWORD", ciphertext); err != nil {
		t.Fatalf("Put: %v", err)
	}

	got, err := s.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(got) != string(ciphertext) {
		t.Errorf("Get = %q; want %q", got, ciphertext)
	}
}

func TestGetNotFound(t *testing.T) {
	s, err := store.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = s.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestDelete(t *testing.T) {
	s, err := store.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_ = s.Put("API_KEY", []byte("secret"))
	if err := s.Delete("API_KEY"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get("API_KEY"); err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestPersistence(t *testing.T) {
	path := tempStorePath(t)

	s1, _ := store.New(path)
	_ = s1.Put("TOKEN", []byte("abc123"))

	s2, err := store.New(path)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	got, err := s2.Get("TOKEN")
	if err != nil {
		t.Fatalf("Get after reload: %v", err)
	}
	if string(got) != "abc123" {
		t.Errorf("persisted value = %q; want %q", got, "abc123")
	}
}

func TestFilePermissions(t *testing.T) {
	path := tempStorePath(t)
	s, _ := store.New(path)
	_ = s.Put("X", []byte("y"))

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o; want 0600", perm)
	}
}
