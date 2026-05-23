package env

import (
	"errors"
	"strings"
	"testing"
)

func encFields() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_SECRET", Value: "topsecret"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestEncryptFieldsAllSensitive(t *testing.T) {
	encFn := func(p string) (string, error) { return "ENC(" + p + ")", nil }
	out, err := EncryptFields(encFields(), DefaultEncryptFieldOptions(), encFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if IsSensitive(e.Key) {
			if !strings.HasPrefix(e.Value, encryptedPrefix) {
				t.Errorf("key %s should be encrypted, got %q", e.Key, e.Value)
			}
		} else {
			if strings.HasPrefix(e.Value, encryptedPrefix) {
				t.Errorf("key %s should not be encrypted", e.Key)
			}
		}
	}
}

func TestEncryptFieldsExplicitKeys(t *testing.T) {
	encFn := func(p string) (string, error) { return "ENC(" + p + ")", nil }
	opts := EncryptFieldOptions{Keys: []string{"PORT"}, AllSensitive: false}
	out, err := EncryptFields(encFields(), opts, encFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "PORT" {
			if !strings.HasPrefix(e.Value, encryptedPrefix) {
				t.Errorf("PORT should be encrypted")
			}
		} else if strings.HasPrefix(e.Value, encryptedPrefix) {
			t.Errorf("key %s should not be encrypted", e.Key)
		}
	}
}

func TestEncryptFieldsSkipsAlreadyEncrypted(t *testing.T) {
	called := 0
	encFn := func(p string) (string, error) { called++; return "ENC(" + p + ")", nil }
	entries := []Entry{{Key: "DB_PASSWORD", Value: encryptedPrefix + "already"}}
	_, err := EncryptFields(entries, DefaultEncryptFieldOptions(), encFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 0 {
		t.Errorf("encryptFn should not be called for already-encrypted values")
	}
}

func TestEncryptFieldsNilFnError(t *testing.T) {
	_, err := EncryptFields(encFields(), DefaultEncryptFieldOptions(), nil)
	if err == nil {
		t.Fatal("expected error for nil encryptFn")
	}
}

func TestDecryptFieldsRoundtrip(t *testing.T) {
	encFn := func(p string) (string, error) { return "ENC(" + p + ")", nil }
	decFn := func(c string) (string, error) {
		if !strings.HasPrefix(c, "ENC(") {
			return "", errors.New("bad cipher")
		}
		return strings.TrimSuffix(strings.TrimPrefix(c, "ENC("), ")"), nil
	}
	original := encFields()
	encrypted, err := EncryptFields(original, DefaultEncryptFieldOptions(), encFn)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	decrypted, err := DecryptFields(encrypted, decFn)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	for i, orig := range original {
		if decrypted[i].Value != orig.Value {
			t.Errorf("key %s: want %q got %q", orig.Key, orig.Value, decrypted[i].Value)
		}
	}
}

func TestIsEncrypted(t *testing.T) {
	if IsEncrypted(Entry{Key: "X", Value: "plain"}) {
		t.Error("plain value should not be encrypted")
	}
	if !IsEncrypted(Entry{Key: "X", Value: encryptedPrefix + "abc"}) {
		t.Error("prefixed value should be encrypted")
	}
}
