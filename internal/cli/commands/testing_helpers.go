// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"

	"github.com/7-Dany/snip/internal/storage"
)

// setupTestRepos creates a temporary storage repository for testing.
// This helper is used across all command test files.
func setupTestRepos(t *testing.T) *storage.Repositories {
	t.Helper()
	repos := storage.New(t.TempDir() + "/test.json")
	if err := repos.Load(); err != nil {
		t.Fatalf("Failed to load test repos: %v", err)
	}
	return repos
}
