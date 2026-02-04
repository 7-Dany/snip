// Package storage provides JSON file-based persistence for SNIP entities.
package storage

import "github.com/7-Dany/snip/internal/domain"

// Repositories bundles all repository implementations with shared state.
// All repositories share the same underlying store for data consistency.
type Repositories struct {
	Snippets   domain.SnippetRepository
	Categories domain.CategoryRepository
	Tags       domain.TagRepository
	store      *store
}

// New creates repositories with all implementations sharing the same store.
func New(filepath string) *Repositories {
	s := newStore(filepath)

	return &Repositories{
		Snippets:   newSnippetRepository(s),
		Categories: newCategoryRepository(s),
		Tags:       newTagRepository(s),
		store:      s,
	}
}

// Save persists all data to the JSON file atomically.
func (r *Repositories) Save() error {
	return r.store.save()
}

// Load reads all data from the JSON file into memory.
// If the file doesn't exist, this is not an error.
func (r *Repositories) Load() error {
	return r.store.load()
}
