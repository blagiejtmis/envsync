package env

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatchDetectsChange(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := WatchOptions{PollInterval: 50 * time.Millisecond}
	ch, err := Watch(ctx, path, opts)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(path, []byte("FOO=changed\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatchNoSpuriousEvents(t *testing.T) {
	path := writeTempEnv(t, "FOO=stable\n")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	opts := WatchOptions{PollInterval: 50 * time.Millisecond}
	ch, err := Watch(ctx, path, opts)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev := <-ch:
		t.Errorf("unexpected change event: %+v", ev)
	case <-ctx.Done():
		// expected — no changes were written
	}
}

func TestWatchChannelClosedOnCancel(t *testing.T) {
	path := writeTempEnv(t, "KEY=val\n")

	ctx, cancel := context.WithCancel(context.Background())
	opts := WatchOptions{PollInterval: 30 * time.Millisecond}
	ch, err := Watch(ctx, path, opts)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	cancel()

	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed after cancel")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel was not closed within timeout")
	}
}

func TestWatchNonExistentFile(t *testing.T) {
	_, err := Watch(context.Background(), "/nonexistent/.env", DefaultWatchOptions())
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}
