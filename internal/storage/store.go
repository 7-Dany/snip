// store.go
package storage

import (
	"github.com/7-Dany/snip/internal/domain"
)

// Metadata tracks auto-increment counters for ID generation.
// IDs are never reused, even after deletion.
//
// Design: Auto-increment strategy
// - next_snippet_id starts at 1 (0 means unassigned)
// - Incremented after each Create operation
// - Never decremented (even on Delete)
type Metadata struct {
	next_snippet_id  int // Next ID to assign to new snippet
	next_category_id int // Next ID to assign to new category
	next_tag_id      int // Next ID to assign to new tag
}

// Store is the internal data structure managing all entities.
// It is unexported and should only be accessed through repository interfaces.
//
// Design: In-memory database with eventual persistence
// - All data loaded at startup (Load)
// - All mutations happen in-memory (fast!)
// - Persisted to disk on demand (Save)
//
// Performance characteristics:
// - O(1) lookups by ID (map-based)
// - O(1) full-text search (inverted index)
// - O(n) filtering operations (FindByCategory, FindByTag, etc.)
type Store struct {
	path        string                   // JSON file path for persistence
	snippets    map[int]*domain.Snippet  // ID → Snippet
	categories  map[int]*domain.Category // ID → Category
	tags        map[int]*domain.Tag      // ID → Tag
	searchIndex map[string][]int         // word → []snippetID (inverted index)
	metadata    Metadata                 // Auto-increment counters
}

// newStore creates a new Store with initialized maps and metadata.
// This function is unexported - users must go through storage.New().
//
// Parameters:
//
//	filepath - Path to JSON file for persistence
//
// Returns:
//
//	*Store with empty maps and metadata starting at ID 1
//
// Design note: IDs start at 1 because 0 means "unassigned"
func newStore(filepath string) *Store {
	return &Store{
		path:        filepath,
		snippets:    make(map[int]*domain.Snippet),
		categories:  make(map[int]*domain.Category),
		tags:        make(map[int]*domain.Tag),
		searchIndex: make(map[string][]int),
		metadata: Metadata{
			next_snippet_id:  1,
			next_category_id: 1,
			next_tag_id:      1,
		},
	}
}

// save persists all store data to the JSON file.
// TODO: Implement JSON marshaling
func (s *Store) save() error {
	return nil
}

// load reads all data from the JSON file into the store.
// TODO: Implement JSON unmarshaling
func (s *Store) load() error {
	return nil
}
