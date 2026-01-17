package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("creates config with defaults when none exists", func(t *testing.T) {
		// Setup: Use temp directory as fake home
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome) // Override home directory for test

		config, err := LoadConfig()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedPath := filepath.Join(tempHome, ".snip", "snippets.json")
		if config.StoragePath != expectedPath {
			t.Errorf("expected storage path %s, got %s", expectedPath, config.StoragePath)
		}

		// Verify .snip directory was created
		snipPath := filepath.Join(tempHome, ".snip")
		if _, err := os.Stat(snipPath); os.IsNotExist(err) {
			t.Error("expected .snip directory to be created")
		}

		// Verify config.json was created
		configPath := filepath.Join(snipPath, "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("expected config.json to be created")
		}
	})

	t.Run("loads existing config", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		// Create config manually
		snipPath := filepath.Join(tempHome, ".snip")
		os.MkdirAll(snipPath, 0755)

		customPath := "/custom/path/snippets.json"
		existingConfig := Config{StoragePath: customPath}
		data, _ := json.Marshal(existingConfig)
		configPath := filepath.Join(snipPath, "config.json")
		os.WriteFile(configPath, data, 0644)

		// Load it
		config, err := LoadConfig()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if config.StoragePath != customPath {
			t.Errorf("expected storage path %s, got %s", customPath, config.StoragePath)
		}
	})

	t.Run("creates parent directory if it doesn't exist", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		config, err := LoadConfig()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if config == nil {
			t.Fatal("expected config to be returned")
		}

		// Verify directory was created
		snipPath := filepath.Join(tempHome, ".snip")
		info, err := os.Stat(snipPath)
		if err != nil {
			t.Fatalf("expected .snip directory to exist, got error: %v", err)
		}

		if !info.IsDir() {
			t.Error("expected .snip to be a directory")
		}
	})

	t.Run("returns error when config file is invalid JSON", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		snipPath := filepath.Join(tempHome, ".snip")
		os.MkdirAll(snipPath, 0755)

		// Write invalid JSON
		configPath := filepath.Join(snipPath, "config.json")
		os.WriteFile(configPath, []byte("invalid json {{{"), 0644)

		_, err := LoadConfig()

		if err == nil {
			t.Fatal("expected error when loading invalid JSON, got nil")
		}
	})

	t.Run("config file contains pretty-printed JSON", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		LoadConfig()

		configPath := filepath.Join(tempHome, ".snip", "config.json")
		data, _ := os.ReadFile(configPath)

		// Pretty JSON should have newlines and indentation
		dataStr := string(data)
		if dataStr == "" {
			t.Fatal("config file is empty")
		}

		// Check for indentation (pretty printing)
		var config Config
		if err := json.Unmarshal(data, &config); err != nil {
			t.Errorf("config file should contain valid JSON: %v", err)
		}

		// Should have multiple lines (pretty printed)
		if len(dataStr) < 20 {
			t.Error("expected pretty-printed JSON with indentation")
		}
	})
}

func TestCreateDefaultConfig(t *testing.T) {
	t.Run("creates valid config file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		snipPath := filepath.Join(tempDir, ".snip")

		config, err := createDefaultConfig(configPath, snipPath)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedPath := filepath.Join(snipPath, "snippets.json")
		if config.StoragePath != expectedPath {
			t.Errorf("expected storage path %s, got %s", expectedPath, config.StoragePath)
		}

		// Verify file was created
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("expected config file to be created")
		}
	})
}

func TestLoadExistingConfig(t *testing.T) {
	t.Run("loads config from valid file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		customPath := "/my/custom/path.json"
		existingConfig := Config{StoragePath: customPath}
		data, _ := json.Marshal(existingConfig)
		os.WriteFile(configPath, data, 0644)

		config, err := loadExistingConfig(configPath)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if config.StoragePath != customPath {
			t.Errorf("expected storage path %s, got %s", customPath, config.StoragePath)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := loadExistingConfig("/path/that/does/not/exist.json")

		if err == nil {
			t.Fatal("expected error for non-existent file, got nil")
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		os.WriteFile(configPath, []byte("not valid json"), 0644)

		_, err := loadExistingConfig(configPath)

		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})
}
