package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/env"
)

func TestMergePullNoConflicts(t *testing.T) {
	s, dir := newTestSyncer(t, "passphrase123")

	// Push a remote state.
	if err := s.Push(); err != nil {
		t.Fatalf("push: %v", err)
	}

	// Write a different local key.
	local := filepath.Join(dir, ".env")
	if err := os.WriteFile(local, []byte("LOCAL_ONLY=yes\n"), 0o600); err != nil {
		t.Fatalf("write local: %v", err)
	}

	res, err := s.MergePull(MergePullOptions{Strategy: env.PreferLocal})
	if err != nil {
		t.Fatalf("merge pull: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if !res.Written {
		t.Error("expected file to be written")
	}
}

func TestMergePullConflictPreferLocal(t *testing.T) {
	s, dir := newTestSyncer(t, "passphrase123")

	// Push initial remote.
	if err := s.Push(); err != nil {
		t.Fatalf("push: %v", err)
	}

	// Overwrite local with conflicting value for an existing key.
	local := filepath.Join(dir, ".env")
	data, _ := os.ReadFile(local)
	data = append(data, []byte("\nCONFLICT_KEY=local_value\n")...)
	if err := os.WriteFile(local, data, 0o600); err != nil {
		t.Fatalf("write local: %v", err)
	}

	// Push the conflict key with a different value via a second syncer.
	s2, _ := newTestSyncerAt(t, "passphrase123", dir, "CONFLICT_KEY=remote_value")
	if err := s2.Push(); err != nil {
		t.Fatalf("push s2: %v", err)
	}

	res, err := s.MergePull(MergePullOptions{Strategy: env.PreferLocal})
	if err != nil {
		t.Fatalf("merge pull: %v", err)
	}
	if len(res.Conflicts) == 0 {
		t.Error("expected at least one conflict")
	}
	for _, e := range res.Merged {
		if e.Key == "CONFLICT_KEY" && e.Value != "local_value" {
			t.Errorf("PreferLocal: expected local_value, got %q", e.Value)
		}
	}
}

func TestMergePullDryRun(t *testing.T) {
	s, dir := newTestSyncer(t, "passphrase123")
	if err := s.Push(); err != nil {
		t.Fatalf("push: %v", err)
	}

	local := filepath.Join(dir, ".env")
	original, _ := os.ReadFile(local)

	_, err := s.MergePull(MergePullOptions{Strategy: env.PreferRemote, DryRun: true})
	if err != nil {
		t.Fatalf("dry-run merge pull: %v", err)
	}

	after, _ := os.ReadFile(local)
	if string(original) != string(after) {
		t.Error("dry-run should not modify the file")
	}
}

// newTestSyncerAt creates a syncer using an existing store dir but different env content.
func newTestSyncerAt(t *testing.T, passphrase, dir, content string) (*Syncer, string) {
	t.Helper()
	envPath := filepath.Join(dir, ".env.alt")
	if err := os.WriteFile(envPath, []byte(content+"\n"), 0o600); err != nil {
		t.Fatalf("write alt env: %v", err)
	}
	s, err := New(filepath.Join(dir, "store"), envPath, passphrase)
	if err != nil {
		t.Fatalf("new syncer: %v", err)
	}
	return s, dir
}
