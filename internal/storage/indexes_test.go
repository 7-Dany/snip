package storage

import (
	"reflect"
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func TestExtractWords(t *testing.T) {
	t.Run("normal text with multiple words", func(t *testing.T) {
		result := extractWords("Quick Sort Algorithm")
		expected := []string{"quick", "sort", "algorithm"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("text with punctuation", func(t *testing.T) {
		result := extractWords("Hello, World!")
		expected := []string{"hello", "world"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("text with numbers", func(t *testing.T) {
		result := extractWords("Go 1.21 is great")
		expected := []string{"go", "121", "is", "great"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("text with multiple spaces", func(t *testing.T) {
		result := extractWords("Quick  Sort   Algorithm")
		expected := []string{"quick", "sort", "algorithm"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		result := extractWords("")
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("only whitespace", func(t *testing.T) {
		result := extractWords("   ")
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("code with special characters", func(t *testing.T) {
		result := extractWords("func quicksort(arr []int)")
		expected := []string{"func", "quicksort", "arr", "int"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("duplicate words are removed", func(t *testing.T) {
		result := extractWords("quick quick sort")
		expected := []string{"quick", "sort"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("case insensitive - all lowercase", func(t *testing.T) {
		result := extractWords("Quick QUICK quick")
		expected := []string{"quick"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("mixed punctuation and case", func(t *testing.T) {
		result := extractWords("Quick-Sort: A Fast Algorithm!")
		expected := []string{"quicksort", "a", "fast", "algorithm"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestBuildSearchableText(t *testing.T) {
	t.Run("combines all fields with spaces", func(t *testing.T) {
		snippet, err := domain.NewSnippet("Quick Sort", "go", "func quicksort() {}")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		snippet.SetDescription("Fast sorting")

		result := buildSearchableText(snippet)
		expected := "Quick Sort Fast sorting func quicksort() {} go"

		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("handles empty description", func(t *testing.T) {
		snippet, err := domain.NewSnippet("Title", "go", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		result := buildSearchableText(snippet)
		expected := "Title  code go"

		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("handles large code snippet", func(t *testing.T) {
		largeCode := ""
		for i := 0; i < 1000; i++ {
			largeCode += "func test() { return nil }\n"
		}

		snippet, err := domain.NewSnippet("Large File", "go", largeCode)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		result := buildSearchableText(snippet)

		// Should contain all parts
		if len(result) < len(largeCode) {
			t.Error("expected result to contain full code")
		}
	})
}

func TestContainsID(t *testing.T) {
	t.Run("finds ID in slice", func(t *testing.T) {
		ids := []int{1, 2, 3, 4, 5}
		if !containsID(ids, 3) {
			t.Error("expected to find ID 3")
		}
	})

	t.Run("returns false for missing ID", func(t *testing.T) {
		ids := []int{1, 2, 3}
		if containsID(ids, 99) {
			t.Error("expected not to find ID 99")
		}
	})

	t.Run("returns false for empty slice", func(t *testing.T) {
		ids := []int{}
		if containsID(ids, 1) {
			t.Error("expected not to find ID in empty slice")
		}
	})

	t.Run("finds first ID", func(t *testing.T) {
		ids := []int{42, 2, 3}
		if !containsID(ids, 42) {
			t.Error("expected to find first ID")
		}
	})

	t.Run("finds last ID", func(t *testing.T) {
		ids := []int{1, 2, 99}
		if !containsID(ids, 99) {
			t.Error("expected to find last ID")
		}
	})
}

func TestIndexSnippet(t *testing.T) {
	t.Run("indexes all words from snippet", func(t *testing.T) {
		store := NewStore("test.json")

		snippet, err := domain.NewSnippet("Quick Sort", "go", "func quicksort() {}")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		snippet.SetID(1)

		store.indexSnippet(snippet)

		// Check "quick" is indexed
		if ids, ok := store.searchIndex["quick"]; !ok {
			t.Error("expected 'quick' to be in index")
		} else if len(ids) != 1 || ids[0] != 1 {
			t.Errorf("searchIndex[\"quick\"] = %v, want [1]", ids)
		}

		// Check "sort" is indexed
		if _, ok := store.searchIndex["sort"]; !ok {
			t.Error("expected 'sort' to be in index")
		}

		// Check "func" from code is indexed
		if _, ok := store.searchIndex["func"]; !ok {
			t.Error("expected 'func' to be in index")
		}

		// Check "go" from language is indexed
		if _, ok := store.searchIndex["go"]; !ok {
			t.Error("expected 'go' to be in index")
		}
	})

	t.Run("does not add duplicate IDs", func(t *testing.T) {
		store := NewStore("test.json")

		snippet, err := domain.NewSnippet("Quick Sort", "go", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		snippet.SetID(1)

		// Index twice
		store.indexSnippet(snippet)
		store.indexSnippet(snippet)

		// Should only have one entry
		if ids := store.searchIndex["quick"]; len(ids) != 1 {
			t.Errorf("expected 1 ID for 'quick', got %d: %v", len(ids), ids)
		}
	})

	t.Run("indexes multiple snippets with overlapping words", func(t *testing.T) {
		store := NewStore("test.json")

		s1, _ := domain.NewSnippet("Quick Sort", "go", "code1")
		s1.SetID(1)
		store.indexSnippet(s1)

		s2, _ := domain.NewSnippet("Quick Select", "go", "code2")
		s2.SetID(2)
		store.indexSnippet(s2)

		// "quick" should have both IDs
		if ids := store.searchIndex["quick"]; len(ids) != 2 {
			t.Errorf("expected 2 IDs for 'quick', got %d: %v", len(ids), ids)
		}

		// "sort" should only have ID 1
		if ids := store.searchIndex["sort"]; len(ids) != 1 || ids[0] != 1 {
			t.Errorf("expected [1] for 'sort', got %v", ids)
		}

		// "select" should only have ID 2
		if ids := store.searchIndex["select"]; len(ids) != 1 || ids[0] != 2 {
			t.Errorf("expected [2] for 'select', got %v", ids)
		}
	})
}

func TestRemoveFromIndex(t *testing.T) {
	t.Run("removes snippet from all words", func(t *testing.T) {
		store := NewStore("test.json")

		snippet, err := domain.NewSnippet("Quick Sort", "go", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		snippet.SetID(1)
		store.indexSnippet(snippet)

		// Verify it's indexed
		if _, ok := store.searchIndex["quick"]; !ok {
			t.Fatal("expected 'quick' to be in index before removal")
		}

		// Remove it
		store.removeFromIndex(1)

		// Verify it's gone
		if _, ok := store.searchIndex["quick"]; ok {
			t.Error("expected 'quick' to be removed from index")
		}
		if _, ok := store.searchIndex["sort"]; ok {
			t.Error("expected 'sort' to be removed from index")
		}
	})

	t.Run("removes only specified snippet from shared words", func(t *testing.T) {
		store := NewStore("test.json")

		s1, _ := domain.NewSnippet("Quick Sort", "go", "code1")
		s1.SetID(1)
		store.indexSnippet(s1)

		s2, _ := domain.NewSnippet("Quick Select", "go", "code2")
		s2.SetID(2)
		store.indexSnippet(s2)

		// Both should be under "quick"
		if len(store.searchIndex["quick"]) != 2 {
			t.Fatalf("expected 2 snippets with 'quick', got %d",
				len(store.searchIndex["quick"]))
		}

		// Remove snippet 1
		store.removeFromIndex(1)

		// "quick" should still have snippet 2
		if ids := store.searchIndex["quick"]; len(ids) != 1 || ids[0] != 2 {
			t.Errorf("expected [2] for 'quick', got %v", ids)
		}

		// "sort" should be gone (only in snippet 1)
		if _, ok := store.searchIndex["sort"]; ok {
			t.Error("expected 'sort' to be removed")
		}

		// "select" should still be there (only in snippet 2)
		if ids := store.searchIndex["select"]; len(ids) != 1 || ids[0] != 2 {
			t.Errorf("expected [2] for 'select', got %v", ids)
		}
	})

	t.Run("removing non-existent ID does not crash", func(t *testing.T) {
		store := NewStore("test.json")

		// Remove from empty index - should not panic
		store.removeFromIndex(999)

		// Add a snippet
		snippet, _ := domain.NewSnippet("Quick Sort", "go", "code")
		snippet.SetID(1)
		store.indexSnippet(snippet)

		// Remove different ID - should not panic
		store.removeFromIndex(999)

		// Original should still be there
		if _, ok := store.searchIndex["quick"]; !ok {
			t.Error("expected 'quick' to still be in index")
		}
	})

	t.Run("cleans up empty word entries", func(t *testing.T) {
		store := NewStore("test.json")

		snippet, _ := domain.NewSnippet("Unique", "go", "code")
		snippet.SetID(1)
		store.indexSnippet(snippet)

		// "unique" should exist
		if _, ok := store.searchIndex["unique"]; !ok {
			t.Fatal("expected 'unique' to be in index")
		}

		// Remove snippet
		store.removeFromIndex(1)

		// "unique" key should be deleted (not just empty slice)
		if _, ok := store.searchIndex["unique"]; ok {
			t.Error("expected 'unique' key to be deleted from index")
		}
	})
}

func TestSearchWithIndex(t *testing.T) {
	setup := func() *Store {
		store := NewStore("test.json")
		store.snippets = make(map[int]*domain.Snippet)

		s1, _ := domain.NewSnippet("Quick Sort", "go", "fast sorting")
		s1.SetID(1)
		store.snippets[1] = s1
		store.indexSnippet(s1)

		s2, _ := domain.NewSnippet("Merge Sort", "go", "stable sorting")
		s2.SetID(2)
		store.snippets[2] = s2
		store.indexSnippet(s2)

		s3, _ := domain.NewSnippet("Binary Search", "go", "fast search")
		s3.SetID(3)
		store.snippets[3] = s3
		store.indexSnippet(s3)

		return store
	}

	t.Run("finds single match", func(t *testing.T) {
		store := setup()

		results, err := store.searchWithIndex("quick")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}

		if results[0].ID() != 1 {
			t.Errorf("expected ID 1, got %d", results[0].ID())
		}
	})

	t.Run("finds multiple matches", func(t *testing.T) {
		store := setup()

		results, err := store.searchWithIndex("sorting")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}

		// Check both snippets 1 and 2 are present
		ids := []int{results[0].ID(), results[1].ID()}
		foundOne := ids[0] == 1 || ids[1] == 1
		foundTwo := ids[0] == 2 || ids[1] == 2

		if !foundOne || !foundTwo {
			t.Errorf("expected IDs 1 and 2, got %v", ids)
		}
	})

	t.Run("OR logic - any word matches", func(t *testing.T) {
		store := setup()

		results, err := store.searchWithIndex("quick merge")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("no matches returns empty slice", func(t *testing.T) {
		store := setup()

		results, err := store.searchWithIndex("nonexistent")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("empty query returns empty slice", func(t *testing.T) {
		store := setup()

		results, err := store.searchWithIndex("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("case insensitive search", func(t *testing.T) {
		store := setup()

		results, err := store.searchWithIndex("QUICK")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("expected 1 result for uppercase query, got %d", len(results))
		}
	})

	t.Run("no duplicate results", func(t *testing.T) {
		store := setup()

		// "fast" appears in both snippet 1 and 3
		// "sorting" appears in snippet 1 and 2
		// Snippet 1 matches both words, but should only appear once
		results, err := store.searchWithIndex("fast sorting")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should get 3 unique snippets (1, 2, 3)
		if len(results) != 3 {
			t.Errorf("expected 3 unique results, got %d", len(results))
		}

		// Verify no duplicates
		seen := make(map[int]bool)
		for _, r := range results {
			if seen[r.ID()] {
				t.Errorf("found duplicate snippet ID %d", r.ID())
			}
			seen[r.ID()] = true
		}
	})
}
