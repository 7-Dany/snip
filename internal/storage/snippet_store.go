package storage

import (
	"slices"

	"github.com/7-Dany/snip/internal/domain"
)

type snippetStore struct {
	store *Store
}

func newSnippetStore(store *Store) *snippetStore {
	return &snippetStore{store: store}
}

func (ss *snippetStore) List() ([]*domain.Snippet, error) {
	snippets := make([]*domain.Snippet, 0, len(ss.store.snippets))
	for _, v := range ss.store.snippets {
		snippets = append(snippets, v)
	}

	return snippets, nil
}

func (ss *snippetStore) FindByID(id int) (*domain.Snippet, error) {
	snippet, ok := ss.store.snippets[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return snippet, nil
}

func (ss *snippetStore) FindByCategory(categoryID int) ([]*domain.Snippet, error) {
	snippets := []*domain.Snippet{}
	for _, v := range ss.store.snippets {
		if v.CategoryID() == categoryID {
			snippets = append(snippets, v)
		}
	}
	return snippets, nil
}

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

func (ss *snippetStore) FindByLanguage(language string) ([]*domain.Snippet, error) {
	snippets := []*domain.Snippet{}
	for _, v := range ss.store.snippets {
		if v.Language() == language {
			snippets = append(snippets, v)
		}
	}
	return snippets, nil
}

func (ss *snippetStore) Search(value string) ([]*domain.Snippet, error) {
	return ss.store.searchWithIndex(value)
}

func (ss *snippetStore) Create(snippet *domain.Snippet) error {
	snippet.SetID(ss.store.metadata.next_snippet_id)
	ss.store.metadata.next_snippet_id++

	ss.store.snippets[snippet.ID()] = snippet
	ss.store.indexSnippet(snippet)

	return nil
}

func (ss *snippetStore) Update(snippet *domain.Snippet) error {
	if _, ok := ss.store.snippets[snippet.ID()]; !ok {
		return domain.ErrNotFound
	}

	ss.store.snippets[snippet.ID()] = snippet
	ss.store.removeFromIndex(snippet.ID())
	ss.store.indexSnippet(snippet)

	return nil
}

func (ss *snippetStore) Delete(id int) error {
	if _, ok := ss.store.snippets[id]; !ok {
		return domain.ErrNotFound
	}

	delete(ss.store.snippets, id)
	ss.store.removeFromIndex(id)

	return nil
}
