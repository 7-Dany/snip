package storage

import (
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestSearchIndex_Search(t *testing.T) {
	t.Run("finds snippets by title", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet1 := mustCreateSnippet(t, "quicksort algorithm", "go", "func quicksort() {}")
		snippet1.SetID(1)
		snippet2 := mustCreateSnippet(t, "binary search", "go", "func binarySearch() {}")
		snippet2.SetID(2)
		s.snippets = []*domain.Snippet{snippet1, snippet2}

		results := idx.search("quicksort")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].ID() != snippet1.ID() {
			t.Error("wrong snippet returned")
		}
	})

	t.Run("finds snippets by language", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "python", "print('hello')")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.search("python")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("finds snippets by code content", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "func fibonacci() int { return 42 }")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.search("fibonacci")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("finds snippets by description", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "func test() {}")
		snippet.SetDescription("implements bubble sort algorithm")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.search("bubble sort")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("search is case insensitive", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "QuickSort", "Go", "func QuickSort() {}")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.search("QUICKSORT")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("returns empty for no matches", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "func test() {}")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.search("nonexistent")

		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("returns nil for empty query", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		results := idx.search("")

		if results != nil {
			t.Error("expected nil for empty query")
		}
	})

	t.Run("finds multiple matching snippets", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet1 := mustCreateSnippet(t, "quicksort", "go", "func quicksort() {}")
		snippet1.SetID(1)
		snippet2 := mustCreateSnippet(t, "mergesort", "go", "func mergesort() {}")
		snippet2.SetID(2)
		snippet3 := mustCreateSnippet(t, "bubblesort", "go", "func bubblesort() {}")
		snippet3.SetID(3)
		s.snippets = []*domain.Snippet{snippet1, snippet2, snippet3}

		results := idx.search("sort")

		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
	})
}

func TestSearchIndex_FindByLanguage(t *testing.T) {
	t.Run("finds snippets with exact language match", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		snippet1.SetID(1)
		snippet2 := mustCreateSnippet(t, "test2", "python", "code2")
		snippet2.SetID(2)
		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")
		snippet3.SetID(3)
		s.snippets = []*domain.Snippet{snippet1, snippet2, snippet3}

		results := idx.findByLanguage("go")

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("language search is case insensitive", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "Go", "code")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.findByLanguage("GO")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("returns empty for no matches", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.findByLanguage("rust")

		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("returns nil for empty language", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		results := idx.findByLanguage("")

		if results != nil {
			t.Error("expected nil for empty language")
		}
	})
}

func TestSearchIndex_FindByCategory(t *testing.T) {
	t.Run("finds snippets in category", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		snippet1.SetID(1)
		snippet1.SetCategory(5)
		snippet2 := mustCreateSnippet(t, "test2", "go", "code2")
		snippet2.SetID(2)
		snippet2.SetCategory(10)
		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")
		snippet3.SetID(3)
		snippet3.SetCategory(5)
		s.snippets = []*domain.Snippet{snippet1, snippet2, snippet3}

		results := idx.findByCategory(5)

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("returns empty for category with no snippets", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		snippet.SetID(1)
		snippet.SetCategory(5)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.findByCategory(99)

		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("finds uncategorized snippets", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		snippet.SetID(1)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.findByCategory(0)

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})
}

func TestSearchIndex_FindByTag(t *testing.T) {
	t.Run("finds snippets with tag", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		snippet1.SetID(1)
		snippet1.AddTag(5)
		snippet2 := mustCreateSnippet(t, "test2", "go", "code2")
		snippet2.SetID(2)
		snippet2.AddTag(10)
		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")
		snippet3.SetID(3)
		snippet3.AddTag(5)
		snippet3.AddTag(10)
		s.snippets = []*domain.Snippet{snippet1, snippet2, snippet3}

		results := idx.findByTag(5)

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("returns empty for tag with no snippets", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		snippet.SetID(1)
		snippet.AddTag(5)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.findByTag(99)

		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("finds snippets with multiple tags", func(t *testing.T) {
		s := newStore("test.json")
		idx := newSearchIndex(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		snippet.SetID(1)
		snippet.AddTag(5)
		snippet.AddTag(10)
		snippet.AddTag(15)
		s.snippets = []*domain.Snippet{snippet}

		results := idx.findByTag(10)

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})
}
