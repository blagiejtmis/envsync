package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/config"
)

func writeTOML(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTOML: %v", err)
	}
	return p
}

func TestLoadValid(t *testing.T) {
	dir := t.TempDir()
	p := writeTOML(t, dir, ".envsync.toml", `
[store]
path = "/tmp/mystore"

[env]
file = ".env.local"
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Store.Path != "/tmp/mystore" {
		t.Errorf("store.path = %q, want /tmp/mystore", cfg.Store.Path)
	}
	if cfg.Env.File != ".env.local" {
		t.Errorf("env.file = %q, want .env.local", cfg.Env.File)
	}
}

func TestLoadDefaultsEnvFile(t *testing.T) {
	dir := t.TempDir()
	p := writeTOML(t, dir, ".envsync.toml", `
[store]
path = "/tmp/store"
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Env.File != ".env" {
		t.Errorf("env.file = %q, want .env", cfg.Env.File)
	}
}

func TestLoadMissingStorePath(t *testing.T) {
	dir := t.TempDir()
	p := writeTOML(t, dir, ".envsync.toml", `
[env]
file = ".env"
`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/.envsync.toml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.Store.Path == "" {
		t.Error("default store.path should not be empty")
	}
	if cfg.Env.File == "" {
		t.Error("default env.file should not be empty")
	}
}
