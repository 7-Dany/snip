// store.go
package storage

import (
	"encoding/json"
	"fmt"
	"os"

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
	NextSnippetID  int `json:"next_snippet_id"`  // Next ID to assign to new snippet
	NextCategoryID int `json:"next_category_id"` // Next ID to assign to new category
	NextTagID      int `json:"next_tag_id"`      // Next ID to assign to new tag
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

type storeJSON struct {
	Snippets       []*domain.Snippet  `json:"snippets"`
	Categories     []*domain.Category `json:"categories"`
	Tags           []*domain.Tag      `json:"tags"`
	NextSnippetID  int                `json:"next_snippet_id"`
	NextCategoryID int                `json:"next_category_id"`
	NextTagID      int                `json:"next_tag_id"`
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

// Add this documentation to your existing store.go methods:

// toJSON converts Store's maps to slices for JSON serialization.
// This separates internal representation (maps for O(1) lookup) from
// persistence format (slices for human-readable JSON).
//
// Design: Map → Slice conversion
// - Iterates through entity maps
// - Appends to pre-allocated slices
// - Includes metadata for ID continuity
//
// Returns:
//
//	*storeJSON - Serializable representation with slices
//
// Performance: O(n) where n = total entities
func (s *Store) toJSON() *storeJSON {
	jsonStore := &storeJSON{
		Snippets:       make([]*domain.Snippet, 0, len(s.snippets)),
		Tags:           make([]*domain.Tag, 0, len(s.tags)),
		Categories:     make([]*domain.Category, 0, len(s.categories)),
		NextSnippetID:  s.metadata.NextSnippetID,
		NextCategoryID: s.metadata.NextCategoryID,
		NextTagID:      s.metadata.NextTagID,
	}

	for _, v := range s.snippets {
		jsonStore.Snippets = append(jsonStore.Snippets, v)
	}

	for _, v := range s.categories {
		jsonStore.Categories = append(jsonStore.Categories, v)
	}

	for _, v := range s.tags {
		jsonStore.Tags = append(jsonStore.Tags, v)
	}

	return jsonStore
}

// fromJSON converts JSON slices back to Store's map representation.
// This reconstructs internal state from persisted data.
//
// Design: Slice → Map conversion
// - Populates entity maps using entity IDs as keys
// - Rebuilds search index for snippets
// - Restores metadata counters
//
// Parameters:
//
//	data - Deserialized JSON data with entity slices
//
// Side effects:
//   - Populates snippets, categories, tags maps
//   - Rebuilds searchIndex via indexSnippet()
//   - Updates metadata counters
//
// Performance: O(n*k) where:
//
//	n = total snippets
//	k = average words per snippet (for indexing)
func (s *Store) fromJSON(data *storeJSON) {
	for _, v := range data.Snippets {
		s.snippets[v.ID()] = v
		s.indexSnippet(v)
	}

	for _, v := range data.Categories {
		s.categories[v.ID()] = v
	}

	for _, v := range data.Tags {
		s.tags[v.ID()] = v
	}

	s.metadata.NextSnippetID = data.NextSnippetID
	s.metadata.NextCategoryID = data.NextCategoryID
	s.metadata.NextTagID = data.NextTagID
}

// save persists all store data to the JSON file.
// This enables data persistence between application runs.
//
// Design: In-memory → Disk persistence
// - Converts maps to slices via toJSON()
// - Marshals to JSON bytes
// - Writes to file with 0644 permissions
//
// File format: See storeJSON struct for JSON schema
// File permissions: 0644 (owner: rw, others: r)
//
// Returns:
//
//	error - Marshal or file write errors
//
// Performance: O(n) where n = total entities
func (s *Store) save() error {
	data, err := json.Marshal(s.toJSON())
	if err != nil {
		return fmt.Errorf("failed to marshal data, %w", err)
	}

	err = os.WriteFile(s.path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save json data into a file %s: %w", s.path, err)
	}

	return nil
}

// load reads all data from the JSON file into the store.
// This restores application state from previous session.
//
// Design: Disk → In-memory restoration
// - Reads JSON file bytes
// - Unmarshals to storeJSON
// - Converts slices to maps via fromJSON()
// - Rebuilds search index
//
// Behavior:
//   - Returns nil if file doesn't exist (first run)
//   - Returns error if file exists but JSON is invalid
//   - Rebuilds search index automatically
//
// Returns:
//
//	error - File read or unmarshal errors (nil if file doesn't exist)
//
// Performance: O(n*k) where:
//
//	n = total snippets
//	k = average words per snippet (for re-indexing)
func (s *Store) load() error {
	// read data from a file
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Start with empty store
		}
		return fmt.Errorf("failed to load file %s: %w", s.path, err)
	}

	// marshall data
	var jsonData *storeJSON = &storeJSON{}
	err = json.Unmarshal(data, jsonData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	s.fromJSON(jsonData)
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
