package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize   = 16
	keySize    = 32
	pbkdf2Iter = 100_000
)

// DerivedKey holds a derived AES key and the salt used to derive it.
type DerivedKey struct {
	Key  []byte
	Salt []byte
}

// DeriveKey derives a 32-byte AES key from passphrase using PBKDF2.
// If salt is nil, a new random salt is generated.
func DeriveKey(passphrase string, salt []byte) (*DerivedKey, error) {
	if passphrase == "" {
		return nil, errors.New("passphrase must not be empty")
	}
	if salt == nil {
		salt = make([]byte, SaltSize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, err
		}
	}
	key := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iter, keySize, sha256.New)
	return &DerivedKey{Key: key, Salt: salt}, nil
}

// Encrypt encrypts plaintext using AES-GCM with the provided key.
// Returns nonce+ciphertext.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts AES-GCM ciphertext (nonce+ciphertext) with the provided key.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, data := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, errors.New("decryption failed: invalid key or corrupted data")
	}
	return plaintext, nil
}
