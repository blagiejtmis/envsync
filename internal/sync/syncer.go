package sync

import (
	"fmt"

	"github.com/yourorg/envsync/internal/crypto"
	"github.com/yourorg/envsync/internal/env"
	"github.com/yourorg/envsync/internal/store"
)

// Syncer handles pushing and pulling encrypted .env files to/from the store.
type Syncer struct {
	store      *store.SecretStore
	passphrase string
}

// New creates a new Syncer with the given store and passphrase.
func New(s *store.SecretStore, passphrase string) *Syncer {
	return &Syncer{store: s, passphrase: passphrase}
}

// Push reads the .env file at path, encrypts it, and stores it under the given key.
func (s *Syncer) Push(key, filePath string) error {
	vars, err := env.LoadFile(filePath)
	if err != nil {
		return fmt.Errorf("push: load file: %w", err)
	}

	plaintext := []byte(env.Serialize(env.FromMap(vars)))

	derived, err := crypto.DeriveKey(s.passphrase, nil)
	if err != nil {
		return fmt.Errorf("push: derive key: %w", err)
	}

	ciphertext, err := crypto.Encrypt(derived.Key, plaintext)
	if err != nil {
		return fmt.Errorf("push: encrypt: %w", err)
	}

	payload := append(derived.Salt, ciphertext...)
	if err := s.store.Put(key, payload); err != nil {
		return fmt.Errorf("push: store put: %w", err)
	}

	return nil
}

// Pull retrieves the encrypted payload for key, decrypts it, and writes it to filePath.
func (s *Syncer) Pull(key, filePath string) error {
	payload, err := s.store.Get(key)
	if err != nil {
		return fmt.Errorf("pull: store get: %w", err)
	}

	if len(payload) < crypto.SaltSize {
		return fmt.Errorf("pull: payload too short")
	}

	salt := payload[:crypto.SaltSize]
	ciphertext := payload[crypto.SaltSize:]

	derived, err := crypto.DeriveKey(s.passphrase, salt)
	if err != nil {
		return fmt.Errorf("pull: derive key: %w", err)
	}

	plaintext, err := crypto.Decrypt(derived.Key, ciphertext)
	if err != nil {
		return fmt.Errorf("pull: decrypt: %w", err)
	}

	vars, err := env.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("pull: parse env: %w", err)
	}

	if err := env.WriteFile(filePath, vars); err != nil {
		return fmt.Errorf("pull: write file: %w", err)
	}

	return nil
}
