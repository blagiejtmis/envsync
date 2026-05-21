package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/store"
	"github.com/yourorg/envsync/internal/sync"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envsync-sync-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func newTestSyncer(t *testing.T, passphrase string) (*sync.Syncer, string) {
	t.Helper()
	dir := tempDir(t)
	storePath := filepath.Join(dir, "store.json")
	s, err := store.New(storePath)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return sync.New(s, passphrase), dir
}

func TestPushAndPull(t *testing.T) {
	syncer, dir := newTestSyncer(t, "supersecret")

	srcFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(srcFile, []byte("APP_ENV=production\nDB_HOST=localhost\n"), 0600); err != nil {
		t.Fatalf("write src file: %v", err)
	}

	if err := syncer.Push("myapp", srcFile); err != nil {
		t.Fatalf("Push: %v", err)
	}

	dstFile := filepath.Join(dir, ".env.out")
	if err := syncer.Pull("myapp", dstFile); err != nil {
		t.Fatalf("Pull: %v", err)
	}

	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("read dst file: %v", err)
	}

	if !containsLine(string(got), "APP_ENV=production") || !containsLine(string(got), "DB_HOST=localhost") {
		t.Errorf("pulled content missing expected vars, got:\n%s", got)
	}
}

func TestPullWrongPassphrase(t *testing.T) {
	syncer, dir := newTestSyncer(t, "correctpassphrase")

	srcFile := filepath.Join(dir, ".env")
	_ = os.WriteFile(srcFile, []byte("SECRET=value\n"), 0600)

	if err := syncer.Push("myapp", srcFile); err != nil {
		t.Fatalf("Push: %v", err)
	}

	s2, err := store.New(filepath.Join(dir, "store.json"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	wrongSyncer := sync.New(s2, "wrongpassphrase")

	dstFile := filepath.Join(dir, ".env.out")
	if err := wrongSyncer.Pull("myapp", dstFile); err == nil {
		t.Error("expected error pulling with wrong passphrase, got nil")
	}
}

func containsLine(content, line string) bool {
	for _, l := range splitLines(content) {
		if l == line {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
