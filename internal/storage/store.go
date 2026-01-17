// store.go
package storage

import (
	"github.com/7-Dany/snip/internal/domain"
)

// Metadata tracks auto-increment counters for ID generation.
// IDs are never reused, even after deletion.
//
// Design: Auto-increment strategy
// - NextSnippetID starts at 1 (0 means unassigned)
// - Incremented after each Create operation
// - Never decremented (even on Delete)
type Metadata struct {
	NextSnippetID  int // Next ID to assign to new snippet
	NextCategoryID int // Next ID to assign to new category
	NextTagID      int // Next ID to assign to new tag
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
			NextSnippetID:  1,
			NextCategoryID: 1,
			NextTagID:      1,
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

// Add these methods to store.go after the newStore function

// nextSnippetID returns the next available snippet ID and increments the counter.
// This encapsulates ID generation logic and prevents direct metadata access.
//
// Returns:
//
//	int - The next unique snippet ID (starting from 1)
//
// Side effects:
//
//	Increments next_snippet_id in metadata
func (s *Store) nextSnippetID() int {
	id := s.metadata.NextSnippetID
	s.metadata.NextSnippetID++
	return id
}

// nextCategoryID returns the next available category ID and increments the counter.
// This encapsulates ID generation logic and prevents direct metadata access.
//
// Returns:
//
//	int - The next unique category ID (starting from 1)
//
// Side effects:
//
//	Increments next_category_id in metadata
func (s *Store) nextCategoryID() int {
	id := s.metadata.NextCategoryID
	s.metadata.NextCategoryID++
	return id
}

// nextTagID returns the next available tag ID and increments the counter.
// This encapsulates ID generation logic and prevents direct metadata access.
//
// Returns:
//
//	int - The next unique tag ID (starting from 1)
//
// Side effects:
//
//	Increments next_tag_id in metadata
func (s *Store) nextTagID() int {
	id := s.metadata.NextTagID
	s.metadata.NextTagID++
	return id
}
