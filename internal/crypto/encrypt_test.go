package crypto_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envsync/internal/crypto"
)

func TestDeriveKey(t *testing.T) {
	key := crypto.DeriveKey("my-shared-secret")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	key2 := crypto.DeriveKey("my-shared-secret")
	if !bytes.Equal(key, key2) {
		t.Fatal("same secret should produce same key")
	}

	key3 := crypto.DeriveKey("different-secret")
	if bytes.Equal(key, key3) {
		t.Fatal("different secrets should produce different keys")
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := crypto.DeriveKey("test-secret")
	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := crypto.DeriveKey("correct-secret")
	wrongKey := crypto.DeriveKey("wrong-secret")
	plaintext := []byte("SECRET=value")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptShortCiphertext(t *testing.T) {
	key := crypto.DeriveKey("secret")
	_, err := crypto.Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}
