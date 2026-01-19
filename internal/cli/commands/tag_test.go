package commands

import (
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestTagCommandCreateWithArgument(t *testing.T) {
	t.Run("creates tag with valid name", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.create([]string{"performance"})

		// Verify creation
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

	t.Run("rejects duplicate tag name", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		// Create first tag
		tc.create([]string{"security"})

		// Try to create duplicate
		tc.create([]string{"security"})

		// Should still have only 1 tag
		tags, _ := repos.Tags.List()
		if len(tags) != 1 {
			t.Errorf("Expected 1 tag after duplicate attempt, got %d", len(tags))
		}
	})
}

func TestTagCommandList(t *testing.T) {
	t.Run("lists all tags", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		// Create test tags
		tag1, _ := domain.NewTag("performance")
		tag2, _ := domain.NewTag("security")
		repos.Tags.Create(tag1)
		repos.Tags.Create(tag2)

		// List should not panic and should work
		tc.list()

		// Verify data exists
		tags, _ := repos.Tags.List()
		if len(tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(tags))
		}
	})

	t.Run("handles empty tag list", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		// Should show info message, not panic
		tc.list()

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})
}

func TestTagCommandDelete(t *testing.T) {
	t.Run("validates ID is required", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		// Should show error for missing ID (not hang)
		tc.delete([]string{})

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("validates ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		// Should show error for invalid ID
		tc.delete([]string{"not-a-number"})

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("shows error when tag not found", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		// Should show error for non-existent ID
		tc.delete([]string{"999"})

		tags, _ := repos.Tags.List()
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(tags))
		}
	})

	// Note: We can't test actual deletion without mocking the confirmation prompt
	// because it reads from stdin and would block waiting for y/n input.
}

func TestTagCommandManage(t *testing.T) {
	repos := setupTestRepos(t)
	tc := NewTagCommand(repos)

	t.Run("shows error when no subcommand provided", func(t *testing.T) {
		// Should show error, not panic
		tc.manage([]string{})
	})

	t.Run("handles unknown subcommand", func(t *testing.T) {
		// Should show error, not panic
		tc.manage([]string{"unknown"})
	})

	t.Run("routes to list command", func(t *testing.T) {
		// Should not panic
		tc.manage([]string{"list"})
	})

	t.Run("handles case insensitive commands", func(t *testing.T) {
		tc.manage([]string{"LIST"})
		tc.manage([]string{"List"})
		tc.manage([]string{"lIsT"})
	})
}

func TestTagCommandManageWithCreate(t *testing.T) {
	t.Run("routes to create with argument", func(t *testing.T) {
		repos := setupTestRepos(t)
		tc := NewTagCommand(repos)

		tc.manage([]string{"create", "test-tag"})

		tags, _ := repos.Tags.List()
		if len(tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(tags))
		}
	})

	// Note: We intentionally skip testing "create" with no args
	// because it launches an interactive Bubbletea prompt.
}
