package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

// Session holds the authentication state for API calls.
type Session struct {
	Cookies   []*http.Cookie `json:"cookies"`
	CSRFToken string         `json:"csrf_token"`
	BaseURL   string         `json:"base_url"`
	CreatedAt time.Time      `json:"created_at"`
}

// Save serializes and encrypts the session to a file.
func (s *Session) Save(path string, keyring *Keyring) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	encrypted, err := keyring.Encrypt(data)
	if err != nil {
		return fmt.Errorf("encrypt session: %w", err)
	}

	if err := os.WriteFile(path, encrypted, 0o600); err != nil {
		return fmt.Errorf("write session: %w", err)
	}
	return nil
}

// LoadSession decrypts and deserializes a session from a file.
func LoadSession(path string, keyring *Keyring) (*Session, error) {
	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read session file: %w", err)
	}

	data, err := keyring.Decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt session: %w", err)
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}
	return &s, nil
}

// IsExpired returns true if the session cookies have expired.
func (s *Session) IsExpired() bool {
	now := time.Now()
	for _, c := range s.Cookies {
		if !c.Expires.IsZero() && now.After(c.Expires) {
			return true
		}
	}
	return false
}

// CookieJar returns an http.CookieJar pre-populated with the session cookies.
func (s *Session) CookieJar() http.CookieJar {
	jar, _ := cookiejar.New(nil)
	u, err := url.Parse(s.BaseURL)
	if err != nil {
		return jar
	}
	jar.SetCookies(u, s.Cookies)
	return jar
}
