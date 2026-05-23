package env

import (
	"errors"
	"strings"
)

const encryptedPrefix = "enc:"

// EncryptFieldOptions controls field-level encryption behaviour.
type EncryptFieldOptions struct {
	// Keys lists the entry keys that should be encrypted.
	Keys []string
	// AllSensitive automatically encrypts any key deemed sensitive by IsSensitive.
	AllSensitive bool
}

// DefaultEncryptFieldOptions returns options with AllSensitive enabled.
func DefaultEncryptFieldOptions() EncryptFieldOptions {
	return EncryptFieldOptions{AllSensitive: true}
}

// EncryptFields returns a new slice where values for selected keys are
// replaced with ciphertext produced by encryptFn. Values that are already
// prefixed with encryptedPrefix are left unchanged.
func EncryptFields(entries []Entry, opts EncryptFieldOptions, encryptFn func(plain string) (string, error)) ([]Entry, error) {
	if encryptFn == nil {
		return nil, errors.New("encryptFn must not be nil")
	}
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		copy := e
		if shouldEncrypt(e, keySet, opts.AllSensitive) && !strings.HasPrefix(e.Value, encryptedPrefix) {
			cipher, err := encryptFn(e.Value)
			if err != nil {
				return nil, err
			}
			copy.Value = encryptedPrefix + cipher
		}
		out[i] = copy
	}
	return out, nil
}

// DecryptFields returns a new slice where encrypted values are decrypted using
// decryptFn. Non-encrypted values are passed through unchanged.
func DecryptFields(entries []Entry, decryptFn func(cipher string) (string, error)) ([]Entry, error) {
	if decryptFn == nil {
		return nil, errors.New("decryptFn must not be nil")
	}
	out := make([]Entry, len(entries))
	for i, e := range entries {
		copy := e
		if strings.HasPrefix(e.Value, encryptedPrefix) {
			plain, err := decryptFn(strings.TrimPrefix(e.Value, encryptedPrefix))
			if err != nil {
				return nil, err
			}
			copy.Value = plain
		}
		out[i] = copy
	}
	return out, nil
}

// IsEncrypted reports whether the entry value is field-encrypted.
func IsEncrypted(e Entry) bool {
	return strings.HasPrefix(e.Value, encryptedPrefix)
}

func shouldEncrypt(e Entry, keySet map[string]struct{}, allSensitive bool) bool {
	if _, ok := keySet[e.Key]; ok {
		return true
	}
	return allSensitive && IsSensitive(e.Key)
}
