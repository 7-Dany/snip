package commands

import (
	"path/filepath"
	"testing"

	"github.com/7-Dany/snip/internal/domain"
	"github.com/7-Dany/snip/internal/storage"
)

// setupTestRepos creates a temporary storage repository for testing
func setupTestRepos(t *testing.T) *storage.Repositories {
	t.Helper()
	tmpfile := filepath.Join(t.TempDir(), "test.json")
	repos := storage.New(tmpfile)
	if err := repos.Load(); err != nil {
		t.Fatalf("Failed to load test repos: %v", err)
	}
	return repos
}

func TestCategoryCommandCreateWithArgument(t *testing.T) {
	t.Run("creates category with valid name", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.create([]string{"algorithms"})

		// Verify creation
		categories, err := repos.Categories.List()
		if err != nil {
			t.Fatalf("Failed to list categories: %v", err)
		}

		if len(categories) != 1 {
			t.Fatalf("Expected 1 category, got %d", len(categories))
		}

		if categories[0].Name() != "algorithms" {
			t.Errorf("Expected name 'algorithms', got '%s'", categories[0].Name())
		}
	})

	t.Run("rejects duplicate category name", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		// Create first category
		cc.create([]string{"web-dev"})

		// Try to create duplicate
		cc.create([]string{"web-dev"})

		// Should still have only 1 category
		categories, _ := repos.Categories.List()
		if len(categories) != 1 {
			t.Errorf("Expected 1 category after duplicate attempt, got %d", len(categories))
		}
	})
}

func TestCategoryCommandList(t *testing.T) {
	t.Run("lists all categories", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		// Create test categories
		cat1, _ := domain.NewCategory("algorithms")
		cat2, _ := domain.NewCategory("web-dev")
		repos.Categories.Create(cat1)
		repos.Categories.Create(cat2)

		// List should not panic and should work
		cc.list()

		// Verify data exists
		categories, _ := repos.Categories.List()
		if len(categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(categories))
		}
	})

	t.Run("handles empty category list", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		// Should show info message, not panic
		cc.list()

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})
}

func TestCategoryCommandDelete(t *testing.T) {
	t.Run("validates ID is required", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		// Should show error for missing ID (not hang)
		cc.delete([]string{})

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})

	t.Run("validates ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		// Should show error for invalid ID
		cc.delete([]string{"not-a-number"})

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})

	t.Run("shows error when category not found", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		// Should show error for non-existent ID
		cc.delete([]string{"999"})

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})

	// Note: We can't test actual deletion without mocking the confirmation prompt
	// because it reads from stdin. The delete function would block waiting for y/n input.
}

func TestCategoryCommandManage(t *testing.T) {
	repos := setupTestRepos(t)
	cc := NewCategoryCommand(repos)

	t.Run("shows error when no subcommand provided", func(t *testing.T) {
		// Should show error, not panic
		cc.manage([]string{})
	})

	t.Run("handles unknown subcommand", func(t *testing.T) {
		// Should show error, not panic
		cc.manage([]string{"unknown"})
	})

	t.Run("routes to list command", func(t *testing.T) {
		// Should not panic
		cc.manage([]string{"list"})
	})

	t.Run("handles case insensitive commands", func(t *testing.T) {
		// Should work with uppercase
		cc.manage([]string{"LIST"})
		cc.manage([]string{"List"})
		cc.manage([]string{"lIsT"})
	})
}

func TestCategoryCommandManageWithCreate(t *testing.T) {
	t.Run("routes to create with argument", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		// This should work - provides name as argument
		cc.manage([]string{"create", "test-category"})

		categories, _ := repos.Categories.List()
		if len(categories) != 1 {
			t.Errorf("Expected 1 category, got %d", len(categories))
		}
	})

	// Note: We SKIP testing "create" with no args because it would hang
	// waiting for interactive Bubbletea input
}
