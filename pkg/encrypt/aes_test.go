package encrypt

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	password := []byte("my-secret-key")
	plaintext := []byte("sensitive cookie data")

	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatal("decrypted data does not match original")
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	password := []byte("correct-key")
	plaintext := []byte("secret data")

	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Decrypt(encrypted, []byte("wrong-key"))
	if err == nil {
		t.Fatal("expected error with wrong password")
	}
}

func TestDecryptCorruptedData(t *testing.T) {
	_, err := Decrypt([]byte("too short"), []byte("key"))
	if err == nil {
		t.Fatal("expected error with corrupted data")
	}
}

func TestEncryptEmpty(t *testing.T) {
	password := []byte("key")
	encrypted, err := Encrypt([]byte(""), password)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatal(err)
	}
	if len(decrypted) != 0 {
		t.Errorf("expected empty decrypted, got %q", decrypted)
	}
}
