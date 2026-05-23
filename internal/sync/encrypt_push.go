package sync

import (
	"github.com/yourorg/envsync/internal/crypto"
	"github.com/yourorg/envsync/internal/env"
)

// EncryptPushOptions configures field-level encryption applied before pushing.
type EncryptPushOptions struct {
	// Passphrase is used to derive the encryption key.
	Passphrase string
	// FieldOpts controls which fields are encrypted.
	FieldOpts env.EncryptFieldOptions
	// Enabled turns field encryption on or off.
	Enabled bool
}

// DefaultEncryptPushOptions returns options with AllSensitive enabled but
// encryption disabled until a passphrase is provided.
func DefaultEncryptPushOptions() EncryptPushOptions {
	return EncryptPushOptions{
		FieldOpts: env.DefaultEncryptFieldOptions(),
		Enabled:   false,
	}
}

// applyEncryptPush encrypts selected fields in entries using the passphrase
// derived key before they are stored. If opts.Enabled is false the original
// slice is returned unchanged.
func applyEncryptPush(entries []env.Entry, opts EncryptPushOptions) ([]env.Entry, error) {
	if !opts.Enabled || opts.Passphrase == "" {
		return entries, nil
	}
	key, err := crypto.DeriveKey(opts.Passphrase, nil)
	if err != nil {
		return nil, err
	}
	encFn := func(plain string) (string, error) {
		return crypto.Encrypt(key, []byte(plain))
	}
	return env.EncryptFields(entries, opts.FieldOpts, encFn)
}

// applyDecryptPull decrypts field-encrypted values after pulling from the
// store. If passphrase is empty the entries are returned unchanged.
func applyDecryptPull(entries []env.Entry, passphrase string) ([]env.Entry, error) {
	if passphrase == "" {
		return entries, nil
	}
	key, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return nil, err
	}
	decFn := func(cipher string) (string, error) {
		b, err := crypto.Decrypt(key, cipher)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return env.DecryptFields(entries, decFn)
}
