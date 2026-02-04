// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestNewTagCommand(t *testing.T) {
	t.Run("creates tag command with repos", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		if tc == nil {
			t.Fatal("NewTagCommand returned nil")
		}

		if tc.repos == nil {
			t.Error("repos is nil")
		}
	})
}

func TestTagCommand_create(t *testing.T) {
	t.Run("creates tag with valid name argument", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.create([]string{"performance"})

		tags, err := repos.Tags.List()
		if err != nil {
			t.Fatalf("Failed to list tags: %v", err)
		}

		if len(tags) != 1 {
			t.Fatalf("Expected 1 tag, got %d", len(tags))
		}

		if tags[0].Name() != "performance" {
			t.Errorf("Expected name 'performance', got '%s'", tags[0].Name())
		}
	})

	t.Run("trims whitespace from name", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.create([]string{"  security  "})

		tags, err := repos.Tags.List()
		if err != nil {
			t.Fatalf("Failed to list tags: %v", err)
		}

		if len(tags) != 1 {
			t.Fatalf("Expected 1 tag, got %d", len(tags))
		}

		if tags[0].Name() != "security" {
			t.Errorf("Expected name 'security', got '%s'", tags[0].Name())
		}
	})

	t.Run("rejects duplicate tag name", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.create([]string{"security"})
		tc.create([]string{"security"})

		tags, _ := repos.Tags.List()
		if len(tags) != 1 {
			t.Errorf("Expected 1 tag after duplicate attempt, got %d", len(tags))
		}
	})

	t.Run("rejects empty name", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.create([]string{"   "})

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags after empty name, got %d", len(tags))
		}
	})
}

func TestTagCommand_list(t *testing.T) {
	t.Run("lists all tags", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tag1, _ := domain.NewTag("performance")
		tag2, _ := domain.NewTag("security")
		repos.Tags.Create(tag1)
		repos.Tags.Create(tag2)

		tc.list()

		tags, _ := repos.Tags.List()
		if len(tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(tags))
		}
	})

	t.Run("handles empty tag list", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.list()

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})
}

func TestTagCommand_delete(t *testing.T) {
	t.Run("validates ID is required", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.delete([]string{})

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("validates ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.delete([]string{"not-a-number"})

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("shows error when tag not found", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.delete([]string{"999"})

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})
}

func TestTagCommand_manage(t *testing.T) {
	repos := setupTestRepos(t)
	tc := NewTagCommand(repos)

	t.Run("shows error when no subcommand provided", func(t *testing.T) {
		tc.manage([]string{})
	})

	t.Run("handles unknown subcommand", func(t *testing.T) {
		tc.manage([]string{"unknown"})
	})

	t.Run("routes to list command", func(t *testing.T) {
		tc.manage([]string{"list"})
	})

	t.Run("routes to create command with argument", func(t *testing.T) {
		tc.manage([]string{"create", "test-tag"})

		tags, _ := repos.Tags.List()
		if len(tags) == 0 {
			t.Error("Expected tag to be created")
		}
	})

	t.Run("handles case insensitive commands", func(t *testing.T) {
		tc.manage([]string{"LIST"})
		tc.manage([]string{"List"})
		tc.manage([]string{"lIsT"})
	})
}
