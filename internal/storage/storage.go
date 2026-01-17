package storage

import "github.com/7-Dany/snip/internal/domain"

// Repositories bundles all repository implementations with shared state.
// This is the main entry point for the storage layer.
//
// Design: Repository Factory Pattern
// - All repositories share the same underlying Store
// - Provides clean separation between public API and internal implementation
// - Allows easy mocking for tests (replace individual repositories)
type Repositories struct {
	Snippets   domain.SnippetRepository  // Snippet CRUD + search operations
	Categories domain.CategoryRepository // Category CRUD operations
	Tags       domain.TagRepository      // Tag CRUD operations
	store      *Store                    // Unexported - internal state
}

// New creates a new Repositories instance with all repository implementations.
// All repositories share the same underlying Store for consistent state.
//
// Parameters:
//
//	filepath - Path to JSON file for persistence (e.g., "snippets.json")
//
// Returns:
//
//	*Repositories with all three repository interfaces initialized
//
// Example:
//
//	repos := storage.New("data.json")
//	repos.Snippets.Create(snippet)
//	repos.Categories.List()
func New(filepath string) *Repositories {
	store := newStore(filepath)

	return &Repositories{
		Snippets:   newSnippetStore(store),
		Categories: newCategoryStore(store),
		Tags:       newTagStore(store),
		store:      store,
	}
}

// Save persists all data (snippets, categories, tags) to the JSON file.
//
// Returns:
//
//	error - File I/O errors, JSON marshaling errors, etc.
func (r *Repositories) Save() error {
	return r.store.save()
}

// Load reads all data from the JSON file into memory.
//
// Returns:
//
//	error - File not found, JSON parsing errors, etc.
func (r *Repositories) Load() error {
	return r.store.load()
}
