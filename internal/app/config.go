package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LucyHeres/xrxs-cli/pkg/config"
)

// Config holds the CLI configuration.
type Config struct {
	BaseURL string `json:"base_url"`
}

// DefaultConfig returns a config with default values.
func DefaultConfig() *Config {
	return &Config{
		BaseURL: "",
	}
}

// LoadConfig reads the config file from the given path.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// SaveConfig writes the config to the given path.
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, config.FilePerm); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
