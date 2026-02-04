package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestNew(t *testing.T) {
	t.Run("creates all repositories successfully", func(t *testing.T) {
		repos := New("test.json")

		if repos == nil {
			t.Fatal("New() returned nil")
		}

		if repos.Snippets == nil {
			t.Error("Snippets repository is nil")
		}

		if repos.Categories == nil {
			t.Error("Categories repository is nil")
		}

		if repos.Tags == nil {
			t.Error("Tags repository is nil")
		}

		if repos.store == nil {
			t.Error("internal store is nil")
		}
	})

	t.Run("repositories share same underlying store", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test.json")
		repos := New(tempFile)

		// Create a category
		category := mustCreateCategory(t, "algorithms")
		if err := repos.Categories.Create(category); err != nil {
			t.Fatalf("failed to create category: %v", err)
		}

		// Create a snippet referencing that category
		snippet := mustCreateSnippet(t, "quicksort", "go", "func quicksort() {}")
		snippet.SetCategory(category.ID())
		if err := repos.Snippets.Create(snippet); err != nil {
			t.Fatalf("failed to create snippet: %v", err)
		}

		// Verify snippet repository can access category data
		snippets, err := repos.Snippets.FindByCategory(category.ID())
		if err != nil {
			t.Fatalf("failed to find by category: %v", err)
		}

		if len(snippets) != 1 {
			t.Errorf("expected 1 snippet, got %d", len(snippets))
		}

		if snippets[0].ID() != snippet.ID() {
			t.Error("retrieved snippet ID doesn't match created snippet ID")
		}
	})
}

func TestRepositories_SaveAndLoad(t *testing.T) {
	t.Run("persists and restores all data", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test.json")
		repos := New(tempFile)

		// Create test data
		category := mustCreateCategory(t, "algorithms")
		if err := repos.Categories.Create(category); err != nil {
			t.Fatalf("failed to create category: %v", err)
		}

		tag := mustCreateTag(t, "sorting")
		if err := repos.Tags.Create(tag); err != nil {
			t.Fatalf("failed to create tag: %v", err)
		}

		snippet := mustCreateSnippet(t, "quicksort", "go", "func quicksort() {}")
		snippet.SetCategory(category.ID())
		snippet.AddTag(tag.ID())
		if err := repos.Snippets.Create(snippet); err != nil {
			t.Fatalf("failed to create snippet: %v", err)
		}

		// Save data
		if err := repos.Save(); err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		// Create new repositories instance and load
		repos2 := New(tempFile)
		if err := repos2.Load(); err != nil {
			t.Fatalf("failed to load: %v", err)
		}

		// Verify categories
		categories, err := repos2.Categories.List()
		if err != nil {
			t.Fatalf("failed to list categories: %v", err)
		}
		if len(categories) != 1 {
			t.Errorf("expected 1 category, got %d", len(categories))
		}
		if !categories[0].Equal(category) {
			t.Error("loaded category doesn't match saved category")
		}

		// Verify tags
		tags, err := repos2.Tags.List()
		if err != nil {
			t.Fatalf("failed to list tags: %v", err)
		}
		if len(tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(tags))
		}
		if !tags[0].Equal(tag) {
			t.Error("loaded tag doesn't match saved tag")
		}

		// Verify snippets
		snippets, err := repos2.Snippets.List()
		if err != nil {
			t.Fatalf("failed to list snippets: %v", err)
		}
		if len(snippets) != 1 {
			t.Errorf("expected 1 snippet, got %d", len(snippets))
		}
		if !snippets[0].Equal(snippet) {
			t.Error("loaded snippet doesn't match saved snippet")
		}
	})

	t.Run("handles missing file on load", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "nonexistent.json")
		repos := New(tempFile)

		// Loading non-existent file should not error
		if err := repos.Load(); err != nil {
			t.Errorf("expected no error on missing file, got %v", err)
		}

		// Should have empty data
		snippets, err := repos.Snippets.List()
		if err != nil {
			t.Fatalf("failed to list snippets: %v", err)
		}
		if len(snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(snippets))
		}
	})

	t.Run("creates file if it doesn't exist on save", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "new.json")
		repos := New(tempFile)

		// Add some data
		category := mustCreateCategory(t, "test")
		if err := repos.Categories.Create(category); err != nil {
			t.Fatalf("failed to create category: %v", err)
		}

		// Save should create the file
		if err := repos.Save(); err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(tempFile); os.IsNotExist(err) {
			t.Error("save did not create file")
		}
	})
}

// mustCreateCategory creates a category or fails the test.
func mustCreateCategory(t *testing.T, name string) *domain.Category {
	t.Helper()
	category, err := domain.NewCategory(name)
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}
	return category
}

// mustCreateTag creates a tag or fails the test.
func mustCreateTag(t *testing.T, name string) *domain.Tag {
	t.Helper()
	tag, err := domain.NewTag(name)
	if err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}
	return tag
}

// mustCreateSnippet creates a snippet or fails the test.
func mustCreateSnippet(t *testing.T, title, language, code string) *domain.Snippet {
	t.Helper()
	snippet, err := domain.NewSnippet(title, language, code)
	if err != nil {
		t.Fatalf("failed to create snippet: %v", err)
	}
	return snippet
}
