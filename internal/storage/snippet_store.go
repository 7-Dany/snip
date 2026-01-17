// snippet_store.go
package storage

import (
	"slices"

	"github.com/7-Dany/snip/internal/domain"
)

// snippetStore implements domain.SnippetRepository using the shared Store.
// This type is unexported - users interact through the SnippetRepository interface.
//
// Design: Repository pattern with shared state
// - Wraps *Store to access shared data
// - Implements all 9 SnippetRepository methods
// - Maintains search index consistency on mutations
type snippetStore struct {
	store *Store
}

// newSnippetStore creates a new snippetStore wrapping the given Store.
// This function is unexported - called only by storage.New().
func newSnippetStore(store *Store) *snippetStore {
	return &snippetStore{store: store}
}

// List returns all snippets in the store.
//
// Performance: O(n) where n = number of snippets
// Returns: Empty slice if no snippets (not an error)
//
// Implementation note:
// - Pre-allocates slice with capacity for efficiency
// - Uses capacity (not length) to avoid nil entries
func (ss *snippetStore) List() ([]*domain.Snippet, error) {
	snippets := make([]*domain.Snippet, 0, len(ss.store.snippets))
	for _, v := range ss.store.snippets {
		snippets = append(snippets, v)
	}
	return snippets, nil
}

// FindByID retrieves a snippet by its ID.
//
// Performance: O(1) map lookup
//
// Parameters:
//
//	id - Snippet ID to find
//
// Returns:
//
//	*Snippet if found
//	domain.ErrNotFound if snippet doesn't exist
func (ss *snippetStore) FindByID(id int) (*domain.Snippet, error) {
	snippet, ok := ss.store.snippets[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return snippet, nil
}

// FindByCategory returns all snippets in the given category.
//
// Performance: O(n) - iterates all snippets
//
// Parameters:
//
//	categoryID - Category ID to filter by
//
// Returns:
//
//	Empty slice if no matches (not an error)
//	Does NOT validate if category exists (caller's responsibility)
//
// Design note: Returning empty slice (not ErrNotFound) because:
// - User asked "what snippets are in category X?"
// - Answer: "none" (empty slice)
// - NOT: "category doesn't exist" (ErrNotFound)
func (ss *snippetStore) FindByCategory(categoryID int) ([]*domain.Snippet, error) {
	snippets := []*domain.Snippet{}
	for _, v := range ss.store.snippets {
		if v.CategoryID() == categoryID {
			snippets = append(snippets, v)
		}
	}
	return snippets, nil
}

// FindByTag returns all snippets with the given tag.
//
// Performance: O(n*m) where n=snippets, m=avg tags per snippet
//
// Parameters:
//
//	tagID - Tag ID to filter by
//
// Returns:
//
//	Empty slice if tag doesn't exist or no snippets have it
//
// Implementation: Uses slices.Contains for clean tag membership check
func (ss *snippetStore) FindByTag(tagID int) ([]*domain.Snippet, error) {
	if _, ok := ss.store.tags[tagID]; !ok {
		return []*domain.Snippet{}, nil
	}

	snippets := []*domain.Snippet{}
	for _, v := range ss.store.snippets {
		if slices.Contains(v.Tags(), tagID) {
			snippets = append(snippets, v)
		}
	}
	return snippets, nil
}

// FindByLanguage returns all snippets in the given programming language.
//
// Performance: O(n)
//
// Parameters:
//
//	language - Programming language (e.g., "go", "python", "javascript")
//
// Returns:
//
//	Empty slice if no matches
func (ss *snippetStore) FindByLanguage(language string) ([]*domain.Snippet, error) {
	snippets := []*domain.Snippet{}
	for _, v := range ss.store.snippets {
		if v.Language() == language {
			snippets = append(snippets, v)
		}
	}
	return snippets, nil
}

// Search performs full-text search across title, description, code, and language.
//
// Performance: O(k) where k = number of matching snippets (not total snippets!)
// This is achieved through the inverted index.
//
// Parameters:
//
//	value - Search query (e.g., "quicksort algorithm")
//
// Returns:
//
//	All snippets containing ANY word from the query
//	Empty slice if no matches
//
// Example:
//
//	Search("quick sort") → snippets containing "quick" OR "sort"
//
// Implementation: Delegates to Store.searchWithIndex (see indexes.go)
func (ss *snippetStore) Search(value string) ([]*domain.Snippet, error) {
	return ss.store.searchWithIndex(value)
}

// Create stores a new snippet and assigns it an ID.
//
// Performance: O(k) where k = words in snippet (for indexing)
//
// Parameters:
//
//	snippet - Must be created via domain.NewSnippet (validated)
//
// Side effects:
//   - Assigns next available ID via SetID()
//   - Increments metadata.next_snippet_id
//   - Adds to snippets map
//   - Updates search index
//
// Design note: Always assigns new ID (ignores any pre-set ID)
// Rationale: Repository controls ID generation for consistency
func (ss *snippetStore) Create(snippet *domain.Snippet) error {
	snippet.SetID(ss.store.metadata.next_snippet_id)
	ss.store.metadata.next_snippet_id++

	ss.store.snippets[snippet.ID()] = snippet
	ss.store.indexSnippet(snippet)

	return nil
}

// Update replaces an existing snippet and re-indexes it.
//
// Performance: O(k) where k = words in snippet (for re-indexing)
//
// Parameters:
//
//	snippet - Must have valid ID from previous Create/FindByID
//
// Returns:
//
//	domain.ErrNotFound if snippet ID doesn't exist
//
// Side effects:
//   - Replaces snippet in map
//   - Removes old snippet from search index
//   - Adds updated snippet to search index
//
// Critical: Must re-index because snippet content may have changed!
// Example: Title changed from "quicksort" to "mergesort"
//   - Old index: {"quicksort": [1]}
//   - Must remove snippet 1 from "quicksort"
//   - Must add snippet 1 to "mergesort"
func (ss *snippetStore) Update(snippet *domain.Snippet) error {
	if _, ok := ss.store.snippets[snippet.ID()]; !ok {
		return domain.ErrNotFound
	}

	ss.store.snippets[snippet.ID()] = snippet
	ss.store.removeFromIndex(snippet.ID()) // Remove old content from index
	ss.store.indexSnippet(snippet)         // Add new content to index

	return nil
}

// Delete removes a snippet and cleans up the search index.
//
// Performance: O(w) where w = unique words in all snippets
// (must scan index to remove snippet ID)
//
// Parameters:
//
//	id - Snippet ID to delete
//
// Returns:
//
//	domain.ErrNotFound if snippet doesn't exist
//
// Side effects:
//   - Removes from snippets map
//   - Removes from search index
//
// Design note: ID is never reused after deletion
func (ss *snippetStore) Delete(id int) error {
	if _, ok := ss.store.snippets[id]; !ok {
		return domain.ErrNotFound
	}

	delete(ss.store.snippets, id)
	ss.store.removeFromIndex(id)

	return nil
}
