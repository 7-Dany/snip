// store_test.go
package storage

import (
	"testing"
)

func TestNewStore(t *testing.T) {
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
}

func TestStoreMetadataInitialization(t *testing.T) {
	store := newStore("test.json")

	if store.metadata.NextSnippetID != 1 {
		t.Errorf("Expected next_snippet_id=1, got %d", store.metadata.NextSnippetID)
	}

	if store.metadata.NextCategoryID != 1 {
		t.Errorf("Expected next_category_id=1, got %d", store.metadata.NextCategoryID)
	}

	if store.metadata.NextTagID != 1 {
		t.Errorf("Expected next_tag_id=1, got %d", store.metadata.NextTagID)
	}
}

func TestStoreMapsAreEmpty(t *testing.T) {
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
}

func TestSaveAndLoadStubsStore(t *testing.T) {
	store := newStore("test.json")

	// Currently stubs - should not error
	if err := store.save(); err != nil {
		t.Errorf("save() returned error: %v", err)
	}

	if err := store.load(); err != nil {
		t.Errorf("load() returned error: %v", err)
	}
}

// ID Generation Method Tests

func TestStore_NextSnippetID_StartsAt1(t *testing.T) {
	store := newStore("test.json")

	id := store.nextSnippetID()

	if id != 1 {
		t.Errorf("nextSnippetID() = %d, expected 1", id)
	}
}

func TestStore_NextSnippetID_Increments(t *testing.T) {
	store := newStore("test.json")

	id1 := store.nextSnippetID()
	id2 := store.nextSnippetID()
	id3 := store.nextSnippetID()

	if id1 != 1 || id2 != 2 || id3 != 3 {
		t.Errorf("nextSnippetID() sequence = %d, %d, %d, expected 1, 2, 3", id1, id2, id3)
	}
}

func TestStore_NextSnippetID_UpdatesMetadata(t *testing.T) {
	store := newStore("test.json")

	store.nextSnippetID()
	store.nextSnippetID()

	if store.metadata.NextSnippetID != 3 {
		t.Errorf("metadata.next_snippet_id = %d, expected 3", store.metadata.NextSnippetID)
	}
}

func TestStore_NextCategoryID_StartsAt1(t *testing.T) {
	store := newStore("test.json")

	id := store.nextCategoryID()

	if id != 1 {
		t.Errorf("nextCategoryID() = %d, expected 1", id)
	}
}

func TestStore_NextCategoryID_Increments(t *testing.T) {
	store := newStore("test.json")

	id1 := store.nextCategoryID()
	id2 := store.nextCategoryID()
	id3 := store.nextCategoryID()

	if id1 != 1 || id2 != 2 || id3 != 3 {
		t.Errorf("nextCategoryID() sequence = %d, %d, %d, expected 1, 2, 3", id1, id2, id3)
	}
}

func TestStore_NextCategoryID_UpdatesMetadata(t *testing.T) {
	store := newStore("test.json")

	store.nextCategoryID()
	store.nextCategoryID()

	if store.metadata.NextCategoryID != 3 {
		t.Errorf("metadata.next_category_id = %d, expected 3", store.metadata.NextCategoryID)
	}
}

func TestStore_NextTagID_StartsAt1(t *testing.T) {
	store := newStore("test.json")

	id := store.nextTagID()

	if id != 1 {
		t.Errorf("nextTagID() = %d, expected 1", id)
	}
}

func TestStore_NextTagID_Increments(t *testing.T) {
	store := newStore("test.json")

	id1 := store.nextTagID()
	id2 := store.nextTagID()
	id3 := store.nextTagID()

	if id1 != 1 || id2 != 2 || id3 != 3 {
		t.Errorf("nextTagID() sequence = %d, %d, %d, expected 1, 2, 3", id1, id2, id3)
	}
}

func TestStore_NextTagID_UpdatesMetadata(t *testing.T) {
	store := newStore("test.json")

	store.nextTagID()
	store.nextTagID()

	if store.metadata.NextTagID != 3 {
		t.Errorf("metadata.next_tag_id = %d, expected 3", store.metadata.NextTagID)
	}
}

func TestStore_IDGeneration_Independent(t *testing.T) {
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
		t.Errorf("next_snippet_id = %d, expected 2", store.metadata.NextSnippetID)
	}
	if store.metadata.NextCategoryID != 2 {
		t.Errorf("next_category_id = %d, expected 2", store.metadata.NextCategoryID)
	}
	if store.metadata.NextTagID != 2 {
		t.Errorf("next_tag_id = %d, expected 2", store.metadata.NextTagID)
	}
}
