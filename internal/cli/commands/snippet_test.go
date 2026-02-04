// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestNewSnippetCommand(t *testing.T) {
	t.Run("creates snippet command with repos", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		if sc == nil {
			t.Fatal("NewSnippetCommand returned nil")
		}

		if sc.repos == nil {
			t.Error("repos is nil")
		}
	})
}

func TestSnippetCommand_list(t *testing.T) {
	t.Run("lists all snippets", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		snip1, _ := domain.NewSnippet("Binary Search", "go", "func binarySearch() {}")
		snip2, _ := domain.NewSnippet("Quicksort", "python", "def quicksort():")
		repos.Snippets.Create(snip1)
		repos.Snippets.Create(snip2)

		sc.list([]string{})

		snippets, _ := repos.Snippets.List()
		if len(snippets) != 2 {
			t.Errorf("Expected 2 snippets, got %d", len(snippets))
		}
	})

	t.Run("handles empty snippet list", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		sc.list([]string{})

		snippets, _ := repos.Snippets.List()
		if len(snippets) != 0 {
			t.Errorf("Expected 0 snippets, got %d", len(snippets))
		}
	})

	t.Run("filters by category", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		cat, _ := domain.NewCategory("algorithms")
		repos.Categories.Create(cat)

		snip1, _ := domain.NewSnippet("Binary Search", "go", "func binarySearch() {}")
		snip1.SetCategory(cat.ID())
		snip2, _ := domain.NewSnippet("Quicksort", "python", "def quicksort():")
		repos.Snippets.Create(snip1)
		repos.Snippets.Create(snip2)

		sc.list([]string{"--category", "1"})
	})

	t.Run("validates category ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		sc.list([]string{"--category", "not-a-number"})
	})

	t.Run("filters by tag", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		tag, _ := domain.NewTag("sorting")
		repos.Tags.Create(tag)

		snip1, _ := domain.NewSnippet("Binary Search", "go", "func binarySearch() {}")
		snip1.AddTag(tag.ID())
		repos.Snippets.Create(snip1)

		sc.list([]string{"--tag", "1"})
	})

	t.Run("validates tag ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		sc.list([]string{"--tag", "not-a-number"})
	})

	t.Run("filters by language", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		snip1, _ := domain.NewSnippet("Binary Search", "go", "func binarySearch() {}")
		snip2, _ := domain.NewSnippet("Quicksort", "python", "def quicksort():")
		repos.Snippets.Create(snip1)
		repos.Snippets.Create(snip2)

		sc.list([]string{"--language", "go"})
	})

	t.Run("validates flag has value", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		sc.list([]string{"--category"})
		sc.list([]string{"--tag"})
		sc.list([]string{"--language"})
	})
}

func TestSnippetCommand_show(t *testing.T) {
	t.Run("validates ID is required", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.show([]string{})
	})

	t.Run("validates ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.show([]string{"not-a-number"})
	})

	t.Run("shows error when snippet not found", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.show([]string{"999"})
	})

	t.Run("shows snippet details", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		snip, _ := domain.NewSnippet("Test Snippet", "go", "func test() {}")
		repos.Snippets.Create(snip)

		sc.show([]string{"1"})
	})

	t.Run("shows snippet with category and tags", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		cat, _ := domain.NewCategory("algorithms")
		repos.Categories.Create(cat)

		tag, _ := domain.NewTag("sorting")
		repos.Tags.Create(tag)

		snip, _ := domain.NewSnippet("Test Snippet", "go", "func test() {}")
		snip.SetCategory(cat.ID())
		snip.AddTag(tag.ID())
		repos.Snippets.Create(snip)

		sc.show([]string{"1"})
	})
}

func TestSnippetCommand_delete(t *testing.T) {
	t.Run("validates ID is required", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.delete([]string{})

		snippets, _ := repos.Snippets.List()
		if len(snippets) != 0 {
			t.Errorf("Expected 0 snippets, got %d", len(snippets))
		}
	})

	t.Run("validates ID is a number", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.delete([]string{"not-a-number"})
	})

	t.Run("shows error when snippet not found", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.delete([]string{"999"})
	})
}

func TestSnippetCommand_search(t *testing.T) {
	t.Run("validates keyword is required", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.search([]string{})
	})

	t.Run("searches snippets", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		snip, _ := domain.NewSnippet("Binary Search", "go", "func binarySearch() {}")
		repos.Snippets.Create(snip)

		sc.search([]string{"binary"})
	})

	t.Run("shows info when no results", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)
		sc.search([]string{"nonexistent"})
	})
}

func TestSnippetCommand_manage(t *testing.T) {
	repos := setupTestRepos(t)
	sc := NewSnippetCommand(repos)

	t.Run("shows error when no subcommand provided", func(t *testing.T) {
		sc.manage([]string{})
	})

	t.Run("handles unknown subcommand", func(t *testing.T) {
		sc.manage([]string{"unknown"})
	})

	t.Run("routes to list command", func(t *testing.T) {
		sc.manage([]string{"list"})
	})

	t.Run("routes to show command", func(t *testing.T) {
		sc.manage([]string{"show", "1"})
	})

	t.Run("routes to search command", func(t *testing.T) {
		sc.manage([]string{"search", "test"})
	})

	t.Run("handles case insensitive commands", func(t *testing.T) {
		sc.manage([]string{"LIST"})
		sc.manage([]string{"List"})
		sc.manage([]string{"SEARCH", "test"})
	})
}

func TestSnippetCommand_getTagNames(t *testing.T) {
	t.Run("returns N/A for empty tag list", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		names := sc.getTagNames([]int{})
		if len(names) != 1 || names[0] != "N/A" {
			t.Errorf("Expected [N/A], got %v", names)
		}
	})

	t.Run("returns tag names for valid IDs", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		tag1, _ := domain.NewTag("performance")
		tag2, _ := domain.NewTag("security")
		repos.Tags.Create(tag1)
		repos.Tags.Create(tag2)

		names := sc.getTagNames([]int{1, 2})
		if len(names) != 2 {
			t.Errorf("Expected 2 tag names, got %d", len(names))
		}
	})

	t.Run("returns N/A when no valid tags found", func(t *testing.T) {
		repos := setupTestRepos(t)
		sc := NewSnippetCommand(repos)

		names := sc.getTagNames([]int{999})
		if len(names) != 1 || names[0] != "N/A" {
			t.Errorf("Expected [N/A], got %v", names)
		}
	})
}
