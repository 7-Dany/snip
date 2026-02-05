package storage

import "github.com/7-Dany/snip/internal/domain"

// snippetRepository implements domain.SnippetRepository using the internal store.
type snippetRepository struct {
	store *store
	index *searchIndex
}

// newSnippetRepository creates a new snippet repository.
func newSnippetRepository(s *store) *snippetRepository {
	return &snippetRepository{
		store: s,
		index: newSearchIndex(s),
	}
}

// List returns all snippets.
func (r *snippetRepository) List() ([]*domain.Snippet, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	result := make([]*domain.Snippet, len(r.store.snippets))
	copy(result, r.store.snippets)
	return result, nil
}

// FindByID finds a snippet by its ID.
func (r *snippetRepository) FindByID(id int) (*domain.Snippet, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for _, snippet := range r.store.snippets {
		if snippet.ID() == id {
			return snippet, nil
		}
	}
	return nil, ErrNotFound
}

// FindByCategory finds all snippets in a category.
func (r *snippetRepository) FindByCategory(categoryID int) ([]*domain.Snippet, error) {
	return r.index.findByCategory(categoryID), nil
}

// FindByTag finds all snippets with a tag.
func (r *snippetRepository) FindByTag(tagID int) ([]*domain.Snippet, error) {
	return r.index.findByTag(tagID), nil
}

// FindByLanguage finds all snippets with a language.
func (r *snippetRepository) FindByLanguage(language string) ([]*domain.Snippet, error) {
	return r.index.findByLanguage(language), nil
}

// Search finds snippets matching the query.
func (r *snippetRepository) Search(query string) ([]*domain.Snippet, error) {
	return r.index.search(query), nil
}

// Create adds a new snippet and assigns it an ID.
func (r *snippetRepository) Create(snippet *domain.Snippet) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	id := r.store.nextSnippetIDAndIncrement()
	snippet.SetID(id)
	r.store.snippets = append(r.store.snippets, snippet)
	return nil
}

// Update replaces an existing snippet.
func (r *snippetRepository) Update(snippet *domain.Snippet) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for i, existing := range r.store.snippets {
		if existing.ID() == snippet.ID() {
			r.store.snippets[i] = snippet
			return nil
		}
	}
	return ErrNotFound
}

// Delete removes a snippet by ID.
func (r *snippetRepository) Delete(id int) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for i, snippet := range r.store.snippets {
		if snippet.ID() == id {
			r.store.snippets = append(r.store.snippets[:i], r.store.snippets[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
