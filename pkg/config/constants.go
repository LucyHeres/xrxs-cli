package config

import (
	"os"
	"path/filepath"
	"time"
)

// Permission constants.
const (
	DirPerm  = 0o700
	FilePerm = 0o600
)

// Timeout constants.
const (
	HTTPTimeout       = 30 * time.Second
	LoginTimeout      = 15 * time.Second
	MaxResponseBytes  = 10 * 1024 * 1024
)

// File names.
const (
	ConfigFileName  = "config.json"
	CookiesFileName = "cookies.enc"
	KeyFileName     = ".key"
)

// DefaultConfigDir returns the xrxs config directory (default ~/.xrxs/).
func DefaultConfigDir() string {
	if dir := os.Getenv("XRXS_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".xrxs")
	}
	return filepath.Join(home, ".xrxs")
}

// DefaultConfigFile returns the path to config.json.
func DefaultConfigFile() string {
	return filepath.Join(DefaultConfigDir(), ConfigFileName)
}

// DefaultCookiesFile returns the path to the encrypted cookie jar.
func DefaultCookiesFile() string {
	return filepath.Join(DefaultConfigDir(), CookiesFileName)
}

// DefaultKeyFile returns the path to the encryption key file.
func DefaultKeyFile() string {
	return filepath.Join(DefaultConfigDir(), KeyFileName)
}
