package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	t.Run("creates empty store with correct initial state", func(t *testing.T) {
		s := newStore("test.json")

		if s == nil {
			t.Fatal("newStore returned nil")
		}

		if s.filepath != "test.json" {
			t.Errorf("expected filepath %q, got %q", "test.json", s.filepath)
		}

		if len(s.snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(s.snippets))
		}

		if len(s.categories) != 0 {
			t.Errorf("expected 0 categories, got %d", len(s.categories))
		}

		if len(s.tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(s.tags))
		}

		if s.nextSnippetID != 1 {
			t.Errorf("expected nextSnippetID 1, got %d", s.nextSnippetID)
		}

		if s.nextCategoryID != 1 {
			t.Errorf("expected nextCategoryID 1, got %d", s.nextCategoryID)
		}

		if s.nextTagID != 1 {
			t.Errorf("expected nextTagID 1, got %d", s.nextTagID)
		}
	})
}

func TestStore_SaveAndLoad(t *testing.T) {
	t.Run("saves and loads all data correctly", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test.json")
		s := newStore(tempFile)

		// Add test data
		category := mustCreateCategory(t, "algorithms")
		category.SetID(1)
		s.categories = append(s.categories, category)
		s.nextCategoryID = 2

		tag := mustCreateTag(t, "sorting")
		tag.SetID(1)
		s.tags = append(s.tags, tag)
		s.nextTagID = 2

		snippet := mustCreateSnippet(t, "quicksort", "go", "func quicksort() {}")
		snippet.SetID(1)
		snippet.SetCategory(category.ID())
		snippet.AddTag(tag.ID())
		s.snippets = append(s.snippets, snippet)
		s.nextSnippetID = 2

		// Save
		if err := s.save(); err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		// Load into new store
		s2 := newStore(tempFile)
		if err := s2.load(); err != nil {
			t.Fatalf("failed to load: %v", err)
		}

		// Verify categories
		if len(s2.categories) != 1 {
			t.Fatalf("expected 1 category, got %d", len(s2.categories))
		}
		if !s2.categories[0].Equal(category) {
			t.Error("loaded category doesn't match saved category")
		}
		if s2.nextCategoryID != 2 {
			t.Errorf("expected nextCategoryID 2, got %d", s2.nextCategoryID)
		}

		// Verify tags
		if len(s2.tags) != 1 {
			t.Fatalf("expected 1 tag, got %d", len(s2.tags))
		}
		if !s2.tags[0].Equal(tag) {
			t.Error("loaded tag doesn't match saved tag")
		}
		if s2.nextTagID != 2 {
			t.Errorf("expected nextTagID 2, got %d", s2.nextTagID)
		}

		// Verify snippets
		if len(s2.snippets) != 1 {
			t.Fatalf("expected 1 snippet, got %d", len(s2.snippets))
		}
		if !s2.snippets[0].Equal(snippet) {
			t.Error("loaded snippet doesn't match saved snippet")
		}
		if s2.nextSnippetID != 2 {
			t.Errorf("expected nextSnippetID 2, got %d", s2.nextSnippetID)
		}
	})

	t.Run("handles nonexistent file gracefully", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "nonexistent.json")
		s := newStore(tempFile)

		if err := s.load(); err != nil {
			t.Errorf("expected no error on missing file, got %v", err)
		}

		if len(s.snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(s.snippets))
		}

		if len(s.categories) != 0 {
			t.Errorf("expected 0 categories, got %d", len(s.categories))
		}

		if len(s.tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(s.tags))
		}
	})

	t.Run("creates file if it doesn't exist on save", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "new.json")
		s := newStore(tempFile)

		category := mustCreateCategory(t, "test")
		category.SetID(1)
		s.categories = append(s.categories, category)

		if err := s.save(); err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		if _, err := os.Stat(tempFile); os.IsNotExist(err) {
			t.Error("save did not create file")
		}
	})

	t.Run("converts nil slices to empty slices on load", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test.json")

		// Write JSON with null arrays
		jsonData := []byte(`{
			"snippets": null,
			"categories": null,
			"tags": null,
			"next_snippet_id": 1,
			"next_category_id": 1,
			"next_tag_id": 1
		}`)
		if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		s := newStore(tempFile)
		if err := s.load(); err != nil {
			t.Fatalf("failed to load: %v", err)
		}

		if s.snippets == nil {
			t.Error("snippets should not be nil")
		}
		if s.categories == nil {
			t.Error("categories should not be nil")
		}
		if s.tags == nil {
			t.Error("tags should not be nil")
		}
	})

	t.Run("save is atomic using temp file", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test.json")
		s := newStore(tempFile)

		category := mustCreateCategory(t, "test")
		category.SetID(1)
		s.categories = append(s.categories, category)

		if err := s.save(); err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		// Verify temp file doesn't exist after save
		tmpFile := tempFile + ".tmp"
		if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
			t.Error("temp file should not exist after save")
		}

		// Verify actual file exists
		if _, err := os.Stat(tempFile); os.IsNotExist(err) {
			t.Error("actual file should exist after save")
		}
	})
}

func TestStore_IDIncrements(t *testing.T) {
	t.Run("snippet ID increments correctly", func(t *testing.T) {
		s := newStore("test.json")

		id1 := s.nextSnippetIDAndIncrement()
		id2 := s.nextSnippetIDAndIncrement()
		id3 := s.nextSnippetIDAndIncrement()

		if id1 != 1 {
			t.Errorf("expected first ID to be 1, got %d", id1)
		}
		if id2 != 2 {
			t.Errorf("expected second ID to be 2, got %d", id2)
		}
		if id3 != 3 {
			t.Errorf("expected third ID to be 3, got %d", id3)
		}
		if s.nextSnippetID != 4 {
			t.Errorf("expected next ID to be 4, got %d", s.nextSnippetID)
		}
	})

	t.Run("category ID increments correctly", func(t *testing.T) {
		s := newStore("test.json")

		id1 := s.nextCategoryIDAndIncrement()
		id2 := s.nextCategoryIDAndIncrement()

		if id1 != 1 {
			t.Errorf("expected first ID to be 1, got %d", id1)
		}
		if id2 != 2 {
			t.Errorf("expected second ID to be 2, got %d", id2)
		}
		if s.nextCategoryID != 3 {
			t.Errorf("expected next ID to be 3, got %d", s.nextCategoryID)
		}
	})

	t.Run("tag ID increments correctly", func(t *testing.T) {
		s := newStore("test.json")

		id1 := s.nextTagIDAndIncrement()
		id2 := s.nextTagIDAndIncrement()

		if id1 != 1 {
			t.Errorf("expected first ID to be 1, got %d", id1)
		}
		if id2 != 2 {
			t.Errorf("expected second ID to be 2, got %d", id2)
		}
		if s.nextTagID != 3 {
			t.Errorf("expected next ID to be 3, got %d", s.nextTagID)
		}
	})
}
