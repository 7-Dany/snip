package storage

import (
	"strings"

	"github.com/7-Dany/snip/internal/domain"
)

// searchIndex provides efficient search across snippets.
// It maintains in-memory indexes for common search operations.
type searchIndex struct {
	store *store
}

// newSearchIndex creates a new search index for the given store.
func newSearchIndex(s *store) *searchIndex {
	return &searchIndex{store: s}
}

// search finds snippets matching the given query string.
// It searches across title, language, code, and description fields.
// The search is case-insensitive and matches partial strings.
func (idx *searchIndex) search(query string) []*domain.Snippet {
	if query == "" {
		return nil
	}

	idx.store.mu.Lock()
	defer idx.store.mu.Unlock()

	query = strings.ToLower(query)
	results := make([]*domain.Snippet, 0)

	for _, snippet := range idx.store.snippets {
		if idx.matches(snippet, query) {
			results = append(results, snippet)
		}
	}

	return results
}

// matches checks if a snippet matches the search query.
func (idx *searchIndex) matches(snippet *domain.Snippet, query string) bool {
	return strings.Contains(strings.ToLower(snippet.Title()), query) ||
		strings.Contains(strings.ToLower(snippet.Language()), query) ||
		strings.Contains(strings.ToLower(snippet.Code()), query) ||
		strings.Contains(strings.ToLower(snippet.Description()), query)
}

// findByLanguage finds all snippets with the given language.
// The search is case-insensitive.
func (idx *searchIndex) findByLanguage(language string) []*domain.Snippet {
	if language == "" {
		return nil
	}

	idx.store.mu.Lock()
	defer idx.store.mu.Unlock()

	language = strings.ToLower(language)
	results := make([]*domain.Snippet, 0)

	for _, snippet := range idx.store.snippets {
		if strings.EqualFold(snippet.Language(), language) {
			results = append(results, snippet)
		}
	}

	return results
}

// findByCategory finds all snippets in the given category.
func (idx *searchIndex) findByCategory(categoryID int) []*domain.Snippet {
	idx.store.mu.Lock()
	defer idx.store.mu.Unlock()

	results := make([]*domain.Snippet, 0)

	for _, snippet := range idx.store.snippets {
		if snippet.CategoryID() == categoryID {
			results = append(results, snippet)
		}
	}

	return results
}

// findByTag finds all snippets with the given tag.
func (idx *searchIndex) findByTag(tagID int) []*domain.Snippet {
	idx.store.mu.Lock()
	defer idx.store.mu.Unlock()

	results := make([]*domain.Snippet, 0)

	for _, snippet := range idx.store.snippets {
		if snippet.HasTag(tagID) {
			results = append(results, snippet)
		}
	}

	return results
}
