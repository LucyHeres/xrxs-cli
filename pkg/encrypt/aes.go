package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	pbkdf2Iterations = 600_000
	saltSize         = 16
	keySize          = 32
)

// Encrypt encrypts plaintext with AES-256-GCM using a password-derived key.
// Output format: salt (16) + nonce (12) + ciphertext.
func Encrypt(plaintext, password []byte) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	key := pbkdf2.Key(password, salt, pbkdf2Iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	// Output: salt || nonce || ciphertext
	out := make([]byte, 0, len(salt)+len(nonce)+len(plaintext)+aesgcm.Overhead())
	out = append(out, salt...)
	out = append(out, nonce...)
	out = aesgcm.Seal(out, nonce, plaintext, nil)
	return out, nil
}

// Decrypt decrypts data encrypted with Encrypt.
func Decrypt(data, password []byte) ([]byte, error) {
	if len(data) < saltSize+aesGCMNonceSize() {
		return nil, fmt.Errorf("data too short: %d bytes", len(data))
	}

	salt := data[:saltSize]
	key := pbkdf2.Key(password, salt, pbkdf2Iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	return plaintext, nil
}

// EncryptString is a convenience wrapper for Encrypt with string I/O.
func EncryptString(plaintext, password string) ([]byte, error) {
	return Encrypt([]byte(plaintext), []byte(password))
}

// DecryptString is a convenience wrapper for Decrypt with string I/O.
func DecryptString(data []byte, password string) ([]byte, error) {
	return Decrypt(data, []byte(password))
}

func aesGCMNonceSize() int {
	block, err := aes.NewCipher(make([]byte, keySize))
	if err != nil {
		return 12 // fallback to AES-256-GCM default nonce size
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return 12
	}
	return aesgcm.NonceSize()
}
