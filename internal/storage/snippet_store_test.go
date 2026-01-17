// snippet_store_test.go
package storage

import (
	"errors"
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

// Helper function to create a test store with sample data
func setupTestStore(t *testing.T) *Store {
	store := newStore("test.json")

	// Create sample category
	cat, _ := domain.NewCategory("algorithms")
	cat.SetID(1)
	store.categories[1] = cat
	store.metadata.NextCategoryID = 2

	// Create sample tags
	tag1, _ := domain.NewTag("sorting")
	tag1.SetID(1)
	store.tags[1] = tag1

	tag2, _ := domain.NewTag("recursion")
	tag2.SetID(2)
	store.tags[2] = tag2

	store.metadata.NextTagID = 3

	return store
}

// === List Tests ===

func TestSnippetStore_List_Empty(t *testing.T) {
	store := newStore("test.json")
	ss := newSnippetStore(store)

	snippets, err := ss.List()

	if err != nil {
		t.Errorf("List() returned error: %v", err)
	}

	if snippets == nil {
		t.Fatal("List() returned nil slice")
	}

	if len(snippets) != 0 {
		t.Errorf("Expected empty slice, got %d snippets", len(snippets))
	}
}

func TestSnippetStore_List_MultipleSnippets(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	// Create 3 snippets
	for i := range 3 {
		snippet, _ := domain.NewSnippet("snippet"+string(rune(i+'0')), "go", "code")
		ss.Create(snippet)
	}

	snippets, err := ss.List()

	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(snippets) != 3 {
		t.Errorf("Expected 3 snippets, got %d", len(snippets))
	}
}

// === FindByID Tests ===

func TestSnippetStore_FindByID_Success(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	original, _ := domain.NewSnippet("quicksort", "go", "func quicksort() {}")
	ss.Create(original)

	found, err := ss.FindByID(original.ID())

	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}

	if found.ID() != original.ID() {
		t.Errorf("Expected ID %d, got %d", original.ID(), found.ID())
	}

	if found.Title() != "quicksort" {
		t.Errorf("Expected title 'quicksort', got '%s'", found.Title())
	}
}

func TestSnippetStore_FindByID_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	_, err := ss.FindByID(999)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// === FindByCategory Tests ===

func TestSnippetStore_FindByCategory_Success(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet1, _ := domain.NewSnippet("quicksort", "go", "code1")
	snippet1.SetCategory(1)
	ss.Create(snippet1)

	snippet2, _ := domain.NewSnippet("mergesort", "go", "code2")
	snippet2.SetCategory(1)
	ss.Create(snippet2)

	snippet3, _ := domain.NewSnippet("hello", "go", "code3")
	snippet3.SetCategory(2) // Different category
	ss.Create(snippet3)

	results, err := ss.FindByCategory(1)

	if err != nil {
		t.Fatalf("FindByCategory() error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 snippets, got %d", len(results))
	}
}

func TestSnippetStore_FindByCategory_EmptyResult(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	results, err := ss.FindByCategory(999)

	if err != nil {
		t.Errorf("FindByCategory() returned error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected empty result, got %d snippets", len(results))
	}
}

// === FindByTag Tests ===

func TestSnippetStore_FindByTag_Success(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet1, _ := domain.NewSnippet("quicksort", "go", "code1")
	snippet1.AddTag(1) // sorting tag
	ss.Create(snippet1)

	snippet2, _ := domain.NewSnippet("mergesort", "go", "code2")
	snippet2.AddTag(1) // sorting tag
	snippet2.AddTag(2) // recursion tag
	ss.Create(snippet2)

	snippet3, _ := domain.NewSnippet("hello", "go", "code3")
	snippet3.AddTag(2) // only recursion tag
	ss.Create(snippet3)

	results, err := ss.FindByTag(1)

	if err != nil {
		t.Fatalf("FindByTag() error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 snippets with tag 1, got %d", len(results))
	}
}

func TestSnippetStore_FindByTag_NonexistentTag(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	results, err := ss.FindByTag(999)

	if err != nil {
		t.Errorf("FindByTag() returned error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected empty result, got %d snippets", len(results))
	}
}

// === FindByLanguage Tests ===

func TestSnippetStore_FindByLanguage_Success(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	go1, _ := domain.NewSnippet("snippet1", "go", "code1")
	ss.Create(go1)

	py1, _ := domain.NewSnippet("snippet2", "python", "code2")
	ss.Create(py1)

	go2, _ := domain.NewSnippet("snippet3", "go", "code3")
	ss.Create(go2)

	results, err := ss.FindByLanguage("go")

	if err != nil {
		t.Fatalf("FindByLanguage() error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 Go snippets, got %d", len(results))
	}
}

func TestSnippetStore_FindByLanguage_NoMatches(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	results, err := ss.FindByLanguage("rust")

	if err != nil {
		t.Errorf("FindByLanguage() returned error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected empty result, got %d snippets", len(results))
	}
}

// === Search Tests ===

func TestSnippetStore_Search_Success(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("quicksort algorithm", "go", "func quicksort() { sort() }")
	ss.Create(snippet)

	results, err := ss.Search("quicksort")

	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Title() != "quicksort algorithm" {
		t.Errorf("Wrong snippet returned")
	}
}

func TestSnippetStore_Search_NoMatches(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("hello world", "go", "fmt.Println")
	ss.Create(snippet)

	results, err := ss.Search("nonexistent")

	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected no results, got %d", len(results))
	}
}

// === Create Tests ===

func TestSnippetStore_Create_AssignsID(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("test", "go", "code")

	if snippet.ID() != 0 {
		t.Errorf("New snippet should have ID 0, got %d", snippet.ID())
	}

	err := ss.Create(snippet)

	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	if snippet.ID() != 1 {
		t.Errorf("Expected ID 1, got %d", snippet.ID())
	}
}

func TestSnippetStore_Create_IncrementsMetadata(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet1, _ := domain.NewSnippet("test1", "go", "code1")
	ss.Create(snippet1)
	snippet2, _ := domain.NewSnippet("test2", "go", "code2")
	ss.Create(snippet2)

	if snippet1.ID() != 1 {
		t.Errorf("First snippet should have ID 1, got %d", snippet1.ID())
	}
	if snippet2.ID() != 2 {
		t.Errorf("Second snippet should have ID 2, got %d", snippet2.ID())
	}
	if store.metadata.NextSnippetID != 3 {
		t.Errorf("Expected next ID 3, got %d", store.metadata.NextSnippetID)
	}
}

func TestSnippetStore_Create_UpdatesSearchIndex(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("quicksort", "go", "func sort() {}")
	ss.Create(snippet) // Search should find it

	results, _ := ss.Search("quicksort")
	if len(results) != 1 {
		t.Errorf("Search didn't find created snippet")
	}
}

// === Update Tests ===
func TestSnippetStore_Update_Success(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("original", "go", "code")
	ss.Create(snippet)
	snippet.SetTitle("updated")

	err := ss.Update(snippet)
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	found, _ := ss.FindByID(snippet.ID())
	if found.Title() != "updated" {
		t.Errorf("Title not updated")
	}
}

func TestSnippetStore_Update_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("test", "go", "code")
	snippet.SetID(999) // Non-existent ID

	err := ss.Update(snippet)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestSnippetStore_Update_ReindexesSearch(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("quicksort", "go", "func sort() {}")
	ss.Create(snippet) // Change title
	snippet.SetTitle("mergesort")
	ss.Update(snippet) // Old title should not be found

	results, _ := ss.Search("quicksort")
	if len(results) != 0 {
		t.Error("Old title still in search index")
	}

	// New title should be found
	results, _ = ss.Search("mergesort")
	if len(results) != 1 {
		t.Error("New title not in search index")
	}
}

// === Delete Tests ===
func TestSnippetStore_Delete_Success(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("test", "go", "code")
	ss.Create(snippet)
	id := snippet.ID()

	err := ss.Delete(id)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	// Should not be found
	_, err = ss.FindByID(id)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Error("Deleted snippet still exists")
	}
}

func TestSnippetStore_Delete_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	err := ss.Delete(999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestSnippetStore_Delete_RemovesFromSearchIndex(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet, _ := domain.NewSnippet("quicksort", "go", "code")
	ss.Create(snippet)
	id := snippet.ID()
	ss.Delete(id)

	// Should not be found in search
	results, _ := ss.Search("quicksort")
	if len(results) != 0 {
		t.Error("Deleted snippet still in search index")
	}
}

func TestSnippetStore_Delete_DoesNotReuseID(t *testing.T) {
	store := setupTestStore(t)
	ss := newSnippetStore(store)

	snippet1, _ := domain.NewSnippet("test1", "go", "code1")
	ss.Create(snippet1)
	ss.Delete(snippet1.ID())

	snippet2, _ := domain.NewSnippet("test2", "go", "code2")
	ss.Create(snippet2) // ID should increment, not reuse deleted ID

	if snippet2.ID() == snippet1.ID() {
		t.Error("ID was reused after deletion")
	}
	if snippet2.ID() != 2 {
		t.Errorf("Expected ID 2, got %d", snippet2.ID())
	}
}
