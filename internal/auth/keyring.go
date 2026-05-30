package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LucyHeres/xrxs-cli/pkg/config"
	"github.com/LucyHeres/xrxs-cli/pkg/encrypt"
)

// Keyring manages the encryption key for the session cookie jar.
// Priority: XRXS_KEY env > ~/.xrxs/.key file > generate and store.
type Keyring struct {
	key []byte
}

// NewKeyring loads or creates the encryption key.
func NewKeyring() (*Keyring, error) {
	// 1. Check env var
	if envKey := os.Getenv("XRXS_KEY"); envKey != "" {
		return &Keyring{key: []byte(envKey)}, nil
	}

	// 2. Try loading from key file
	keyPath := config.DefaultKeyFile()
	if key, err := loadKeyFile(keyPath); err == nil {
		return &Keyring{key: key}, nil
	}

	// 3. Generate and persist new key
	if err := os.MkdirAll(filepath.Dir(keyPath), config.DirPerm); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	// Write hex-encoded key to file
	if err := os.WriteFile(keyPath, []byte(hex.EncodeToString(key)), config.FilePerm); err != nil {
		return nil, fmt.Errorf("write key file: %w", err)
	}

	return &Keyring{key: key}, nil
}

// Key returns the raw encryption key bytes.
func (k *Keyring) Key() []byte {
	return k.key
}

// Encrypt encrypts data using the keyring key.
func (k *Keyring) Encrypt(plaintext []byte) ([]byte, error) {
	return encrypt.Encrypt(plaintext, k.key)
}

// Decrypt decrypts data using the keyring key.
func (k *Keyring) Decrypt(ciphertext []byte) ([]byte, error) {
	return encrypt.Decrypt(ciphertext, k.key)
}

func loadKeyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	raw := string(data)
	// Trim trailing newline if present
	if len(raw) > 0 && raw[len(raw)-1] == '\n' {
		raw = raw[:len(raw)-1]
	}
	return hex.DecodeString(raw)
}
