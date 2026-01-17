package storage

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func setupTestStore(t *testing.T) (*Store, func()) {
	t.Helper()
	// Replace slashes with underscores to create valid filename
	name := t.Name()
	name = strings.ReplaceAll(name, "/", "_")
	tmpfile := "test_store_" + name + ".json"
	store := newStore(tmpfile)
	cleanup := func() {
		os.Remove(tmpfile)
	}
	return store, cleanup
}

func TestNewStore(t *testing.T) {
	t.Run("initializes all fields", func(t *testing.T) {
		store := newStore("test.json")

		if store == nil {
			t.Fatal("newStore() returned nil")
		}

		if store.path != "test.json" {
			t.Errorf("Expected path 'test.json', got '%s'", store.path)
		}

		if store.snippets == nil {
			t.Error("snippets map is nil")
		}

		if store.categories == nil {
			t.Error("categories map is nil")
		}

		if store.tags == nil {
			t.Error("tags map is nil")
		}

		if store.searchIndex == nil {
			t.Error("searchIndex map is nil")
		}
	})

	t.Run("initializes metadata", func(t *testing.T) {
		store := newStore("test.json")

		if store.metadata.NextSnippetID != 1 {
			t.Errorf("Expected NextSnippetID=1, got %d", store.metadata.NextSnippetID)
		}

		if store.metadata.NextCategoryID != 1 {
			t.Errorf("Expected NextCategoryID=1, got %d", store.metadata.NextCategoryID)
		}

		if store.metadata.NextTagID != 1 {
			t.Errorf("Expected NextTagID=1, got %d", store.metadata.NextTagID)
		}
	})

	t.Run("creates empty maps", func(t *testing.T) {
		store := newStore("test.json")

		if len(store.snippets) != 0 {
			t.Errorf("Expected empty snippets map, got %d entries", len(store.snippets))
		}

		if len(store.categories) != 0 {
			t.Errorf("Expected empty categories map, got %d entries", len(store.categories))
		}

		if len(store.tags) != 0 {
			t.Errorf("Expected empty tags map, got %d entries", len(store.tags))
		}

		if len(store.searchIndex) != 0 {
			t.Errorf("Expected empty searchIndex map, got %d entries", len(store.searchIndex))
		}
	})
}

func TestStoreSaveAndLoad(t *testing.T) {
	t.Run("empty store", func(t *testing.T) {
		store, cleanup := setupTestStore(t)
		defer cleanup()

		// Save empty store
		if err := store.save(); err != nil {
			t.Fatalf("save() failed: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(store.path); os.IsNotExist(err) {
			t.Fatal("save() did not create file")
		}

		// Load into new store
		newStore := newStore(store.path)
		if err := newStore.load(); err != nil {
			t.Fatalf("load() failed: %v", err)
		}

		// Verify empty data
		if len(newStore.snippets) != 0 {
			t.Errorf("Expected 0 snippets, got %d", len(newStore.snippets))
		}
		if len(newStore.categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(newStore.categories))
		}
		if len(newStore.tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(newStore.tags))
		}

		// Verify metadata preserved
		if newStore.metadata.NextSnippetID != 1 {
			t.Errorf("Expected NextSnippetID=1, got %d", newStore.metadata.NextSnippetID)
		}
	})

	t.Run("round trip with data", func(t *testing.T) {
		store, cleanup := setupTestStore(t)
		defer cleanup()

		// Create test data
		category, _ := domain.NewCategory("algorithms")
		category.SetID(store.nextCategoryID())
		store.categories[category.ID()] = category

		tag1, _ := domain.NewTag("sorting")
		tag1.SetID(store.nextTagID())
		store.tags[tag1.ID()] = tag1

		tag2, _ := domain.NewTag("recursion")
		tag2.SetID(store.nextTagID())
		store.tags[tag2.ID()] = tag2

		snippet, _ := domain.NewSnippet("quicksort", "go", "func quicksort() {}")
		snippet.SetID(store.nextSnippetID())
		snippet.SetCategory(category.ID())
		snippet.AddTag(tag1.ID())
		snippet.AddTag(tag2.ID())
		snippet.SetDescription("Quick sort implementation")
		store.snippets[snippet.ID()] = snippet
		store.indexSnippet(snippet)

		// Save
		if err := store.save(); err != nil {
			t.Fatalf("save() failed: %v", err)
		}

		// Load into new store
		newStore := newStore(store.path)
		if err := newStore.load(); err != nil {
			t.Fatalf("load() failed: %v", err)
		}

		// Verify snippets
		if len(newStore.snippets) != 1 {
			t.Fatalf("Expected 1 snippet, got %d", len(newStore.snippets))
		}

		loadedSnippet := newStore.snippets[snippet.ID()]
		if loadedSnippet == nil {
			t.Fatal("Snippet not found after load")
		}
		if loadedSnippet.Title() != snippet.Title() {
			t.Errorf("Title mismatch: expected %s, got %s", snippet.Title(), loadedSnippet.Title())
		}
		if loadedSnippet.Language() != snippet.Language() {
			t.Errorf("Language mismatch: expected %s, got %s", snippet.Language(), loadedSnippet.Language())
		}
		if loadedSnippet.Code() != snippet.Code() {
			t.Errorf("Code mismatch")
		}
		if loadedSnippet.Description() != snippet.Description() {
			t.Errorf("Description mismatch")
		}
		if loadedSnippet.CategoryID() != snippet.CategoryID() {
			t.Errorf("CategoryID mismatch: expected %d, got %d", snippet.CategoryID(), loadedSnippet.CategoryID())
		}
		if len(loadedSnippet.Tags()) != len(snippet.Tags()) {
			t.Errorf("Tags count mismatch: expected %d, got %d", len(snippet.Tags()), len(loadedSnippet.Tags()))
		}

		// Verify categories
		if len(newStore.categories) != 1 {
			t.Fatalf("Expected 1 category, got %d", len(newStore.categories))
		}
		loadedCategory := newStore.categories[category.ID()]
		if loadedCategory == nil {
			t.Fatal("Category not found after load")
		}
		if loadedCategory.Name() != category.Name() {
			t.Errorf("Category name mismatch: expected %s, got %s", category.Name(), loadedCategory.Name())
		}

		// Verify tags
		if len(newStore.tags) != 2 {
			t.Fatalf("Expected 2 tags, got %d", len(newStore.tags))
		}
		loadedTag1 := newStore.tags[tag1.ID()]
		if loadedTag1 == nil || loadedTag1.Name() != tag1.Name() {
			t.Error("Tag1 not loaded correctly")
		}
		loadedTag2 := newStore.tags[tag2.ID()]
		if loadedTag2 == nil || loadedTag2.Name() != tag2.Name() {
			t.Error("Tag2 not loaded correctly")
		}

		// Verify metadata (ID counters should be preserved)
		if newStore.metadata.NextSnippetID != store.metadata.NextSnippetID {
			t.Errorf("NextSnippetID mismatch: expected %d, got %d",
				store.metadata.NextSnippetID, newStore.metadata.NextSnippetID)
		}
		if newStore.metadata.NextCategoryID != store.metadata.NextCategoryID {
			t.Errorf("NextCategoryID mismatch: expected %d, got %d",
				store.metadata.NextCategoryID, newStore.metadata.NextCategoryID)
		}
		if newStore.metadata.NextTagID != store.metadata.NextTagID {
			t.Errorf("NextTagID mismatch: expected %d, got %d",
				store.metadata.NextTagID, newStore.metadata.NextTagID)
		}
	})

	t.Run("multiple entities", func(t *testing.T) {
		store, cleanup := setupTestStore(t)
		defer cleanup()

		// Create multiple categories
		for i := range 3 {
			cat, _ := domain.NewCategory("category" + string(rune('A'+i)))
			cat.SetID(store.nextCategoryID())
			store.categories[cat.ID()] = cat
		}

		// Create multiple tags
		for i := range 5 {
			tag, _ := domain.NewTag("tag" + string(rune('A'+i)))
			tag.SetID(store.nextTagID())
			store.tags[tag.ID()] = tag
		}

		// Create multiple snippets
		for i := range 10 {
			snip, _ := domain.NewSnippet("snippet"+string(rune('A'+i)), "go", "code")
			snip.SetID(store.nextSnippetID())
			store.snippets[snip.ID()] = snip
			store.indexSnippet(snip)
		}

		// Save
		if err := store.save(); err != nil {
			t.Fatalf("save() failed: %v", err)
		}

		data, _ := os.ReadFile(store.path)
		t.Logf("JSON content:\n%s", string(data))

		// Load into new store
		newStore := newStore(store.path)
		if err := newStore.load(); err != nil {
			t.Fatalf("load() failed: %v", err)
		}

		// Verify counts
		if len(newStore.snippets) != 10 {
			t.Errorf("Expected 10 snippets, got %d", len(newStore.snippets))
		}
		if len(newStore.categories) != 3 {
			t.Errorf("Expected 3 categories, got %d", len(newStore.categories))
		}
		if len(newStore.tags) != 5 {
			t.Errorf("Expected 5 tags, got %d", len(newStore.tags))
		}
	})
}

func TestStoreLoad(t *testing.T) {
	t.Run("file not exist", func(t *testing.T) {
		store := newStore("nonexistent_file.json")

		err := store.load()
		if err != nil {
			t.Errorf("load() should return nil when file doesn't exist, got: %v", err)
		}

		// Verify store is empty (initialized state)
		if len(store.snippets) != 0 {
			t.Error("Snippets should be empty after load from nonexistent file")
		}
		if store.metadata.NextSnippetID != 1 {
			t.Error("Metadata should remain at initial state")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		tmpfile := "invalid_test.json"
		defer os.Remove(tmpfile)

		// Write invalid JSON
		invalidJSON := []byte(`{"snippets": [invalid json}`)
		if err := os.WriteFile(tmpfile, invalidJSON, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		store := newStore(tmpfile)
		err := store.load()

		if err == nil {
			t.Error("load() should return error for invalid JSON")
		}
	})

	t.Run("rebuilds search index", func(t *testing.T) {
		store, cleanup := setupTestStore(t)
		defer cleanup()

		// Create snippet with searchable content
		snippet, _ := domain.NewSnippet("binary search", "python", "def binary_search():")
		snippet.SetID(store.nextSnippetID())
		snippet.SetDescription("Fast search algorithm")
		store.snippets[snippet.ID()] = snippet
		store.indexSnippet(snippet)

		// Verify index exists
		if len(store.searchIndex) == 0 {
			t.Fatal("Search index should not be empty")
		}

		// Save
		if err := store.save(); err != nil {
			t.Fatalf("save() failed: %v", err)
		}

		// Load into new store
		newStore := newStore(store.path)
		if err := newStore.load(); err != nil {
			t.Fatalf("load() failed: %v", err)
		}

		// Verify search index was rebuilt
		if len(newStore.searchIndex) == 0 {
			t.Fatal("Search index should be rebuilt on load")
		}

		// Verify search works
		results, _ := newStore.searchWithIndex("binary")
		if len(results) != 1 {
			t.Errorf("Expected 1 search result, got %d", len(results))
		}
		if len(results) > 0 && results[0].ID() != snippet.ID() {
			t.Error("Search index not rebuilt correctly")
		}
	})
}

func TestStoreJSONFormat(t *testing.T) {
	t.Run("has correct structure", func(t *testing.T) {
		store, cleanup := setupTestStore(t)
		defer cleanup()

		// Create minimal test data
		snippet, _ := domain.NewSnippet("test", "go", "code")
		snippet.SetID(store.nextSnippetID())
		store.snippets[snippet.ID()] = snippet

		// Save
		if err := store.save(); err != nil {
			t.Fatalf("save() failed: %v", err)
		}

		// Read raw JSON
		data, err := os.ReadFile(store.path)
		if err != nil {
			t.Fatalf("Failed to read saved file: %v", err)
		}

		// Parse as generic map to verify structure
		var jsonData map[string]any
		if err := json.Unmarshal(data, &jsonData); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Verify top-level keys exist
		if _, ok := jsonData["snippets"]; !ok {
			t.Error("JSON missing 'snippets' key")
		}
		if _, ok := jsonData["categories"]; !ok {
			t.Error("JSON missing 'categories' key")
		}
		if _, ok := jsonData["tags"]; !ok {
			t.Error("JSON missing 'tags' key")
		}
		if _, ok := jsonData["next_snippet_id"]; !ok {
			t.Error("JSON missing 'next_snippet_id' key")
		}
		if _, ok := jsonData["next_category_id"]; !ok {
			t.Error("JSON missing 'next_category_id' key")
		}
		if _, ok := jsonData["next_tag_id"]; !ok {
			t.Error("JSON missing 'next_tag_id' key")
		}

		// Verify arrays (not objects)
		snippets, ok := jsonData["snippets"].([]any)
		if !ok {
			t.Error("'snippets' should be an array")
		}
		if len(snippets) != 1 {
			t.Errorf("Expected 1 snippet in JSON, got %d", len(snippets))
		}
	})
}

func TestStoreNextSnippetID(t *testing.T) {
	t.Run("starts at 1", func(t *testing.T) {
		store := newStore("test.json")

		id := store.nextSnippetID()

		if id != 1 {
			t.Errorf("nextSnippetID() = %d, expected 1", id)
		}
	})

	t.Run("increments", func(t *testing.T) {
		store := newStore("test.json")

		id1 := store.nextSnippetID()
		id2 := store.nextSnippetID()
		id3 := store.nextSnippetID()

		if id1 != 1 || id2 != 2 || id3 != 3 {
			t.Errorf("nextSnippetID() sequence = %d, %d, %d, expected 1, 2, 3", id1, id2, id3)
		}
	})

	t.Run("updates metadata", func(t *testing.T) {
		store := newStore("test.json")

		store.nextSnippetID()
		store.nextSnippetID()

		if store.metadata.NextSnippetID != 3 {
			t.Errorf("metadata.NextSnippetID = %d, expected 3", store.metadata.NextSnippetID)
		}
	})
}

func TestStoreNextCategoryID(t *testing.T) {
	t.Run("starts at 1", func(t *testing.T) {
		store := newStore("test.json")

		id := store.nextCategoryID()

		if id != 1 {
			t.Errorf("nextCategoryID() = %d, expected 1", id)
		}
	})

	t.Run("increments", func(t *testing.T) {
		store := newStore("test.json")

		id1 := store.nextCategoryID()
		id2 := store.nextCategoryID()
		id3 := store.nextCategoryID()

		if id1 != 1 || id2 != 2 || id3 != 3 {
			t.Errorf("nextCategoryID() sequence = %d, %d, %d, expected 1, 2, 3", id1, id2, id3)
		}
	})

	t.Run("updates metadata", func(t *testing.T) {
		store := newStore("test.json")

		store.nextCategoryID()
		store.nextCategoryID()

		if store.metadata.NextCategoryID != 3 {
			t.Errorf("metadata.NextCategoryID = %d, expected 3", store.metadata.NextCategoryID)
		}
	})
}

func TestStoreNextTagID(t *testing.T) {
	t.Run("starts at 1", func(t *testing.T) {
		store := newStore("test.json")

		id := store.nextTagID()

		if id != 1 {
			t.Errorf("nextTagID() = %d, expected 1", id)
		}
	})

	t.Run("increments", func(t *testing.T) {
		store := newStore("test.json")

		id1 := store.nextTagID()
		id2 := store.nextTagID()
		id3 := store.nextTagID()

		if id1 != 1 || id2 != 2 || id3 != 3 {
			t.Errorf("nextTagID() sequence = %d, %d, %d, expected 1, 2, 3", id1, id2, id3)
		}
	})

	t.Run("updates metadata", func(t *testing.T) {
		store := newStore("test.json")

		store.nextTagID()
		store.nextTagID()

		if store.metadata.NextTagID != 3 {
			t.Errorf("metadata.NextTagID = %d, expected 3", store.metadata.NextTagID)
		}
	})
}

func TestStoreIDGeneration(t *testing.T) {
	t.Run("independent counters", func(t *testing.T) {
		store := newStore("test.json")

		snippetID := store.nextSnippetID()   // 1
		categoryID := store.nextCategoryID() // 1
		tagID := store.nextTagID()           // 1

		if snippetID != 1 || categoryID != 1 || tagID != 1 {
			t.Errorf("ID generation not independent: snippet=%d, category=%d, tag=%d",
				snippetID, categoryID, tagID)
		}

		// Each counter should be at 2 now
		if store.metadata.NextSnippetID != 2 {
			t.Errorf("NextSnippetID = %d, expected 2", store.metadata.NextSnippetID)
		}
		if store.metadata.NextCategoryID != 2 {
			t.Errorf("NextCategoryID = %d, expected 2", store.metadata.NextCategoryID)
		}
		if store.metadata.NextTagID != 2 {
			t.Errorf("NextTagID = %d, expected 2", store.metadata.NextTagID)
		}
	})
}
