package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrNotFound is returned when a secret key does not exist in the store.
var ErrNotFound = errors.New("secret not found")

// Entry represents a single encrypted secret entry.
type Entry struct {
	Key       string    `json:"key"`
	Ciphertext []byte   `json:"ciphertext"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SecretStore manages encrypted secret entries persisted to a JSON file.
type SecretStore struct {
	path    string
	entries map[string]Entry
}

// New opens or creates a SecretStore backed by the file at path.
func New(path string) (*SecretStore, error) {
	s := &SecretStore{
		path:    path,
		entries: make(map[string]Entry),
	}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Put stores an encrypted ciphertext under key.
func (s *SecretStore) Put(key string, ciphertext []byte) error {
	s.entries[key] = Entry{
		Key:        key,
		Ciphertext: ciphertext,
		UpdatedAt:  time.Now().UTC(),
	}
	return s.save()
}

// Get retrieves the ciphertext for key.
func (s *SecretStore) Get(key string) ([]byte, error) {
	e, ok := s.entries[key]
	if !ok {
		return nil, ErrNotFound
	}
	return e.Ciphertext, nil
}

// Delete removes key from the store.
func (s *SecretStore) Delete(key string) error {
	if _, ok := s.entries[key]; !ok {
		return ErrNotFound
	}
	delete(s.entries, key)
	return s.save()
}

// Keys returns all stored keys.
func (s *SecretStore) Keys() []string {
	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	return keys
}

// Has reports whether key exists in the store.
func (s *SecretStore) Has(key string) bool {
	_, ok := s.entries[key]
	return ok
}

// GetEntry returns the full Entry for key, including metadata.
func (s *SecretStore) GetEntry(key string) (Entry, error) {
	e, ok := s.entries[key]
	if !ok {
		return Entry{}, ErrNotFound
	}
	return e, nil
}

func (s *SecretStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	for _, e := range entries {
		s.entries[e.Key] = e
	}
	return nil
}

func (s *SecretStore) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	entries := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		entries = append(entries, e)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
