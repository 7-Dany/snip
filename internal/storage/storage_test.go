package storage

import (
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestNew(t *testing.T) {
	t.Run("creates all repositories", func(t *testing.T) {
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
			t.Error("Internal store is nil")
		}
	})
}

func TestRepositories(t *testing.T) {
	t.Run("share same store", func(t *testing.T) {
		repos := New("test.json")

		// Create a category
		category, _ := domain.NewCategory("algorithms")
		if err := repos.Categories.Create(category); err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}

		// Create a snippet referencing that category
		snippet, _ := domain.NewSnippet("quicksort", "go", "func quicksort() {}")
		snippet.SetCategory(category.ID())
		if err := repos.Snippets.Create(snippet); err != nil {
			t.Fatalf("Failed to create snippet: %v", err)
		}

		// Verify snippet can find snippets by the category we just created
		snippets, err := repos.Snippets.FindByCategory(category.ID())
		if err != nil {
			t.Fatalf("Failed to find by category: %v", err)
		}

		if len(snippets) != 1 {
			t.Errorf("Expected 1 snippet, got %d", len(snippets))
		}

		if snippets[0].ID() != snippet.ID() {
			t.Error("Retrieved snippet doesn't match created snippet")
		}
	})
}
