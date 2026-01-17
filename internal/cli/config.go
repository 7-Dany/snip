package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	StoragePath string `json:"storage_path"`
}

func LoadConfig() (*Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	snipPath := filepath.Join(dir, ".snip")

	// MkdirAll creates directory and all parents, does nothing if exists
	if err := os.MkdirAll(snipPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .snip directory: %w", err)
	}

	configPath := filepath.Join(snipPath, "config.json")

	// Check if config exists
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		// Create default config
		return createDefaultConfig(configPath, snipPath)
	} else if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	// Load existing config
	return loadExistingConfig(configPath)
}

func createDefaultConfig(configPath, snipPath string) (*Config, error) {
	config := &Config{
		StoragePath: filepath.Join(snipPath, "snippets.json"),
	}

	data, err := json.MarshalIndent(config, "", "  ") // Pretty JSON!
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	return config, nil
}

func loadExistingConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
