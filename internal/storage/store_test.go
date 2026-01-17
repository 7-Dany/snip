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

	if store.metadata.next_snippet_id != 1 {
		t.Errorf("Expected next_snippet_id=1, got %d", store.metadata.next_snippet_id)
	}

	if store.metadata.next_category_id != 1 {
		t.Errorf("Expected next_category_id=1, got %d", store.metadata.next_category_id)
	}

	if store.metadata.next_tag_id != 1 {
		t.Errorf("Expected next_tag_id=1, got %d", store.metadata.next_tag_id)
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
