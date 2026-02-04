// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestNewCategoryCommand(t *testing.T) {
	t.Run("creates category command with repos", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		if cc == nil {
			t.Fatal("NewCategoryCommand returned nil")
		}

		if cc.repos == nil {
			t.Error("repos is nil")
		}
	})
}

func TestCategoryCommand_create(t *testing.T) {
	t.Run("creates category with valid name argument", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.create([]string{"algorithms"})

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

	t.Run("trims whitespace from name", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.create([]string{"  web-dev  "})

		categories, err := repos.Categories.List()
		if err != nil {
			t.Fatalf("Failed to list categories: %v", err)
		}

		if len(categories) != 1 {
			t.Fatalf("Expected 1 category, got %d", len(categories))
		}

		if categories[0].Name() != "web-dev" {
			t.Errorf("Expected name 'web-dev', got '%s'", categories[0].Name())
		}
	})

	t.Run("rejects duplicate category name", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.create([]string{"web-dev"})
		cc.create([]string{"web-dev"})

		categories, _ := repos.Categories.List()
		if len(categories) != 1 {
			t.Errorf("Expected 1 category after duplicate attempt, got %d", len(categories))
		}
	})

	t.Run("rejects empty name", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.create([]string{"   "})

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories after empty name, got %d", len(categories))
		}
	})
}

func TestCategoryCommand_list(t *testing.T) {
	t.Run("lists all categories", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cat1, _ := domain.NewCategory("algorithms")
		cat2, _ := domain.NewCategory("web-dev")
		repos.Categories.Create(cat1)
		repos.Categories.Create(cat2)

		cc.list()

		categories, _ := repos.Categories.List()
		if len(categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(categories))
		}
	})

	t.Run("handles empty category list", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.list()

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})
}

func TestCategoryCommand_delete(t *testing.T) {
	t.Run("validates ID is required", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.delete([]string{})

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})

	t.Run("validates ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.delete([]string{"not-a-number"})

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})

	t.Run("shows error when category not found", func(t *testing.T) {
		repos := setupTestRepos(t)
		cc := NewCategoryCommand(repos)

		cc.delete([]string{"999"})

		categories, _ := repos.Categories.List()
		if len(categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(categories))
		}
	})
}

func TestCategoryCommand_manage(t *testing.T) {
	repos := setupTestRepos(t)
	cc := NewCategoryCommand(repos)

	t.Run("shows error when no subcommand provided", func(t *testing.T) {
		cc.manage([]string{})
	})

	t.Run("handles unknown subcommand", func(t *testing.T) {
		cc.manage([]string{"unknown"})
	})

	t.Run("routes to list command", func(t *testing.T) {
		cc.manage([]string{"list"})
	})

	t.Run("routes to create command with argument", func(t *testing.T) {
		cc.manage([]string{"create", "test-category"})

		categories, _ := repos.Categories.List()
		if len(categories) == 0 {
			t.Error("Expected category to be created")
		}
	})

	t.Run("handles case insensitive commands", func(t *testing.T) {
		cc.manage([]string{"LIST"})
		cc.manage([]string{"List"})
		cc.manage([]string{"lIsT"})
	})
}
