package encrypt

import (
	"crypto/rc4"
	"encoding/hex"
)

// RC4Encrypt encrypts plaintext with RC4 using the given key, returning the hex-encoded result.
// This matches the JavaScript XrxsEncryptDecryptUtils.encrypt(str, key) implementation.
func RC4Encrypt(plaintext, key string) (string, error) {
	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	src := []byte(plaintext)
	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)
	return hex.EncodeToString(dst), nil
}

// DefaultEncryptKey is the default RC4 key used for password encryption,
// matching the JS side default key 'qjydxone'.
const DefaultEncryptKey = "qjydxone"
