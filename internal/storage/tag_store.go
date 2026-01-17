package storage

import (
	"strings"

	"github.com/7-Dany/snip/internal/domain"
)

// tagStore implements domain.TagRepository using the shared Store.
// This type is unexported - users interact through the TagRepository interface.
//
// Design: Repository pattern with shared state
// - Wraps *Store to access shared data
// - Implements all 6 TagRepository methods
// - Uses simple map storage (no search index needed)
type tagStore struct {
	store *Store
}

// newTagStore creates a new tagStore that wraps the given Store.
// This is an unexported constructor - tags are accessed through the
// public Repositories.Tags field returned by storage.New().
func newTagStore(store *Store) *tagStore {
	return &tagStore{store: store}
}

// List returns all tags in the store.
//
// Performance: O(n) where n = total tags
//
// Returns:
//
//	Empty slice (not nil) if no tags exist
func (ts *tagStore) List() ([]*domain.Tag, error) {
	tags := make([]*domain.Tag, 0, len(ts.store.tags))
	for _, v := range ts.store.tags {
		tags = append(tags, v)
	}
	return tags, nil
}

// FindByID retrieves a tag by its unique ID.
//
// Performance: O(1) map lookup
//
// Parameters:
//
//	id - Tag ID to find
//
// Returns:
//
//	domain.ErrNotFound if tag ID doesn't exist
func (ts *tagStore) FindByID(id int) (*domain.Tag, error) {
	tag, ok := ts.store.tags[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return tag, nil
}

// FindByName retrieves a tag by its name (case-insensitive).
//
// Performance: O(n) where n = total tags
//
// Parameters:
//
//	name - Tag name to search for (case-insensitive match)
//
// Returns:
//
//	First matching tag (if multiple exist with same name)
//	domain.ErrNotFound if no tag matches
//
// Note: Uses case-insensitive comparison (strings.EqualFold) for better UX.
// Storage layer allows duplicate names - CLI layer should enforce uniqueness.
func (ts *tagStore) FindByName(name string) (*domain.Tag, error) {
	for _, v := range ts.store.tags {
		if strings.EqualFold(v.Name(), name) {
			return v, nil
		}
	}
	return nil, domain.ErrNotFound
}

// Create adds a new tag to the store and assigns it a unique ID.
//
// Performance: O(1)
//
// Parameters:
//
//	tag - Must be valid Tag from domain.NewTag()
//
// Side effects:
//   - Assigns auto-incremented ID to tag
//   - Updates store metadata (next_tag_id++)
//   - Stores tag in map
//
// Note: Ignores any pre-set ID on the tag - always assigns new ID.
func (ts *tagStore) Create(tag *domain.Tag) error {
	tag.SetID(ts.store.metadata.next_tag_id)
	ts.store.metadata.next_tag_id++

	ts.store.tags[tag.ID()] = tag

	return nil
}

// Update replaces an existing tag.
//
// Performance: O(1)
//
// Parameters:
//
//	tag - Must have valid ID from previous Create/FindByID
//
// Returns:
//
//	domain.ErrNotFound if tag ID doesn't exist
//
// Side effects:
//
//	Replaces tag in map (including updated timestamp)
func (ts *tagStore) Update(tag *domain.Tag) error {
	if _, ok := ts.store.tags[tag.ID()]; !ok {
		return domain.ErrNotFound
	}

	ts.store.tags[tag.ID()] = tag
	return nil
}

// Delete removes a tag from the store.
//
// Performance: O(1)
//
// Parameters:
//
//	id - Tag ID to delete
//
// Returns:
//
//	domain.ErrNotFound if tag ID doesn't exist
//
// Side effects:
//
//	Removes tag from map
//
// Note: Does not remove tag references from snippets - CLI layer should
// handle cleanup before deletion or filter orphaned tags when displaying.
func (ts *tagStore) Delete(id int) error {
	if _, ok := ts.store.tags[id]; !ok {
		return domain.ErrNotFound
	}

	delete(ts.store.tags, id)
	return nil
}
