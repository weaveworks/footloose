package cluster

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// KeyStore is a store for public keys.
type KeyStore struct {
	basePath string
}

// NewKeyStore creates a new KeyStore
func NewKeyStore(basePath string) *KeyStore {
	return &KeyStore{
		basePath: basePath,
	}
}

func fileExists(path string) bool {
	// XXX: There's a subtle bug: if stat fails for another reason that the file
	// not existing, we return the file exists.
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (s *KeyStore) keyPath(name string) string {
	return filepath.Join(s.basePath, name)
}

func (s *KeyStore) keyExists(name string) bool {
	return fileExists(s.keyPath(name))
}

// Store adds the key to the store.
func (s *KeyStore) Store(name, key string) error {
	if s.keyExists(name) {
		return errors.Errorf("key store: store: key '%s' already exists", name)
	}

	if err := ioutil.WriteFile(s.keyPath(name), []byte(key), 0644); err != nil {
		return errors.Wrap(err, "key store: write")
	}

	return nil
}

// Get retrieves a key from the store.
func (s *KeyStore) Get(name string) ([]byte, error) {
	if !s.keyExists(name) {
		return nil, errors.Errorf("key store: get: unknown key '%s'", name)
	}
	return ioutil.ReadFile(s.keyPath(name))
}

// Remove removes a key from the store.
func (s *KeyStore) Remove(name string) error {
	if !s.keyExists(name) {
		return errors.Errorf("key store: remove: unknown key '%s'", name)
	}
	if err := os.Remove(s.keyPath(name)); err != nil {
		return errors.Wrap(err, "key store: remove")
	}
	return nil
}