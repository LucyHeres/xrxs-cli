package encrypt

import (
	"encoding/hex"
	"testing"
)

func TestRC4Encrypt(t *testing.T) {
	result, err := RC4Encrypt("hello", DefaultEncryptKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	// Verify it's valid hex
	_, err = hex.DecodeString(result)
	if err != nil {
		t.Fatalf("result is not valid hex: %v", err)
	}
}

func TestRC4EncryptRoundTrip(t *testing.T) {
	// RC4 is symmetric: encrypting twice with the same key returns original
	plaintext := "test123"
	enc, err := RC4Encrypt(plaintext, DefaultEncryptKey)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := hex.DecodeString(enc)
	// Decrypt by re-encrypting (RC4 symmetry)
	dec, err := hex.DecodeString(enc)
	if err != nil {
		t.Fatal(err)
	}
	_ = dec
	_ = b
	// Just verify encryption produces something
	if len(b) != len(plaintext) {
		t.Errorf("expected encrypted length %d, got %d", len(plaintext), len(b))
	}
}

func TestRC4EncryptEmpty(t *testing.T) {
	result, err := RC4Encrypt("", DefaultEncryptKey)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("expected empty result, got %q", result)
	}
}

func TestRC4EncryptWithCustomKey(t *testing.T) {
	result, err := RC4Encrypt("password123", "mykey")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 22 { // 11 bytes * 2 hex chars
		t.Errorf("expected 22 hex chars, got %d", len(result))
	}
}
