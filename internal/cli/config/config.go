// Package config handles application configuration management.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type Config struct {
	StoragePath string `json:"storage_path"`
}

// LoadConfig loads configuration from ~/.snip/config.json.
// Creates the file with defaults if it doesn't exist.
func LoadConfig() (*Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	return loadConfigFromDir(dir)
}

// loadConfigFromDir loads configuration from a specific directory.
// Used by LoadConfig and for testing.
func loadConfigFromDir(homeDir string) (*Config, error) {
	snipPath := filepath.Join(homeDir, ".snip")

	// Note: MkdirAll error is not tested - would require simulating permission failures
	if err := os.MkdirAll(snipPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .snip directory: %w", err)
	}

	configPath := filepath.Join(snipPath, "config.json")

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return createDefaultConfig(configPath, snipPath)
	} else if err != nil {
		// Note: Stat errors other than NotExist are not tested - extremely rare
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	return loadExistingConfig(configPath)
}

// createDefaultConfig creates a new config file with default values.
func createDefaultConfig(configPath, snipPath string) (*Config, error) {
	config := &Config{
		StoragePath: filepath.Join(snipPath, "snippets.json"),
	}

	// Note: Marshal error not tested - won't fail with our simple struct
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Note: WriteFile error not tested - would require disk full or permission issues
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	return config, nil
}

// loadExistingConfig reads and parses an existing config file.
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
