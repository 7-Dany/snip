package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
// The configuration is stored in ~/.snip/config.json and contains
// the path to the snippets database file.
type Config struct {
	// StoragePath is the absolute path to the snippets database file.
	// Defaults to ~/.snip/snippets.json
	StoragePath string `json:"storage_path"`
}

// LoadConfig loads the application configuration from ~/.snip/config.json.
// If the configuration file does not exist, it creates a new one with default values.
// If the .snip directory does not exist, it will be created.
//
// Returns an error if:
//   - The home directory cannot be determined
//   - The .snip directory cannot be created
//   - The config file cannot be read or written
//   - The config file contains invalid JSON
func LoadConfig() (*Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	return loadConfigFromDir(dir)
}

// loadConfigFromDir loads configuration from a specific directory.
// This is an internal function used by LoadConfig and is also useful for testing.
// The config file will be located at <homeDir>/.snip/config.json
func loadConfigFromDir(homeDir string) (*Config, error) {
	snipPath := filepath.Join(homeDir, ".snip")

	// MkdirAll creates directory and all parents, does nothing if exists
	if err := os.MkdirAll(snipPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .snip directory: %w", err)
	}

	configPath := filepath.Join(snipPath, "config.json")

	// Check if config exists
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		// Create default config
		return createDefaultConfig(configPath, snipPath)
	} else if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	// Load existing config
	return loadExistingConfig(configPath)
}

// createDefaultConfig creates a new configuration file with default values.
// The default storage path is set to <snipPath>/snippets.json.
// The configuration is written as pretty-printed JSON for readability.
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

// loadExistingConfig reads and parses an existing configuration file.
// Returns an error if the file cannot be read or contains invalid JSON.
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
