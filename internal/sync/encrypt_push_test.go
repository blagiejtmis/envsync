package sync

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/env"
)

var encPushEntries = []env.Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "s3cr3t"},
	{Key: "API_KEY", Value: "abc123"},
	{Key: "PORT", Value: "9000"},
}

func TestApplyEncryptPushDisabled(t *testing.T) {
	opts := DefaultEncryptPushOptions()
	out, err := applyEncryptPush(encPushEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range out {
		if e.Value != encPushEntries[i].Value {
			t.Errorf("key %s: value should be unchanged", e.Key)
		}
	}
}

func TestApplyEncryptPushEncryptsSensitive(t *testing.T) {
	opts := DefaultEncryptPushOptions()
	opts.Enabled = true
	opts.Passphrase = "hunter2"
	out, err := applyEncryptPush(encPushEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if env.IsSensitive(e.Key) {
			if !env.IsEncrypted(e) {
				t.Errorf("key %s should be encrypted", e.Key)
			}
		} else {
			if env.IsEncrypted(e) {
				t.Errorf("key %s should not be encrypted", e.Key)
			}
		}
	}
}

func TestApplyEncryptPushDoesNotMutateOriginal(t *testing.T) {
	opts := DefaultEncryptPushOptions()
	opts.Enabled = true
	opts.Passphrase = "hunter2"
	origVals := make([]string, len(encPushEntries))
	for i, e := range encPushEntries {
		origVals[i] = e.Value
	}
	_, _ = applyEncryptPush(encPushEntries, opts)
	for i, e := range encPushEntries {
		if e.Value != origVals[i] {
			t.Errorf("key %s: original mutated", e.Key)
		}
	}
}

func TestApplyDecryptPullRoundtrip(t *testing.T) {
	opts := DefaultEncryptPushOptions()
	opts.Enabled = true
	opts.Passphrase = "roundtrip-pass"
	encrypted, err := applyEncryptPush(encPushEntries, opts)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	decrypted, err := applyDecryptPull(encrypted, opts.Passphrase)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	for i, orig := range encPushEntries {
		if decrypted[i].Value != orig.Value {
			t.Errorf("key %s: want %q got %q", orig.Key, orig.Value, decrypted[i].Value)
		}
	}
}

func TestApplyDecryptPullEmptyPassphraseNoOp(t *testing.T) {
	out, err := applyDecryptPull(encPushEntries, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range out {
		if e.Value != encPushEntries[i].Value {
			t.Errorf("key %s: should be unchanged", e.Key)
		}
	}
	_ = strings.Contains // keep import used in other test files
}
