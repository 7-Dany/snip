package storage

import (
	"strings"

	"github.com/7-Dany/snip/internal/domain"
)

// categoryStore implements domain.CategoryRepository using the shared Store.
// This type is unexported - users interact through the CategoryRepository interface.
//
// Design: Repository pattern with shared state
// - Wraps *Store to access shared data
// - Implements all 6 CategoryRepository methods
// - Uses simple map storage (no search index needed)
type categoryStore struct {
	store *Store
}

// newCategoryStore creates a new categoryStore that wraps the given Store.
// This is an unexported constructor - categories are accessed through the
// public Repositories.Categories field returned by storage.New().
func newCategoryStore(store *Store) *categoryStore {
	return &categoryStore{store: store}
}

// List returns all categories in the store.
//
// Performance: O(n) where n = total categories
//
// Returns:
//
//	Empty slice (not nil) if no categories exist
func (cs *categoryStore) List() ([]*domain.Category, error) {
	categories := make([]*domain.Category, 0, len(cs.store.categories))
	for _, v := range cs.store.categories {
		categories = append(categories, v)
	}
	return categories, nil
}

// FindByID retrieves a category by its unique ID.
//
// Performance: O(1) map lookup
//
// Parameters:
//
//	id - Category ID to find
//
// Returns:
//
//	domain.ErrNotFound if category ID doesn't exist
func (cs *categoryStore) FindByID(id int) (*domain.Category, error) {
	category, ok := cs.store.categories[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return category, nil
}

// FindByName retrieves a category by its name (case-insensitive).
//
// Performance: O(n) where n = total categories
//
// Parameters:
//
//	name - Category name to search for (case-insensitive match)
//
// Returns:
//
//	First matching category (if multiple exist with same name)
//	domain.ErrNotFound if no category matches
//
// Note: Uses case-insensitive comparison (strings.EqualFold) for better UX.
// Storage layer allows duplicate names - CLI layer should enforce uniqueness.
func (cs *categoryStore) FindByName(name string) (*domain.Category, error) {
	for _, v := range cs.store.categories {
		if strings.EqualFold(v.Name(), name) {
			return v, nil
		}
	}

	return nil, domain.ErrNotFound
}

// Create adds a new category to the store and assigns it a unique ID.
//
// Performance: O(1)
//
// Parameters:
//
//	category - Must be valid Category from domain.NewCategory()
//
// Side effects:
//   - Assigns auto-incremented ID to category
//   - Updates store metadata (next_category_id++)
//   - Stores category in map
//
// Note: Ignores any pre-set ID on the category - always assigns new ID.
func (cs *categoryStore) Create(category *domain.Category) error {
	category.SetID(cs.store.metadata.next_category_id)
	cs.store.metadata.next_category_id++

	cs.store.categories[category.ID()] = category

	return nil
}

// Update replaces an existing category.
//
// Performance: O(1)
//
// Parameters:
//
//	category - Must have valid ID from previous Create/FindByID
//
// Returns:
//
//	domain.ErrNotFound if category ID doesn't exist
//
// Side effects:
//
//	Replaces category in map (including updated timestamp)
func (cs *categoryStore) Update(category *domain.Category) error {
	if _, ok := cs.store.categories[category.ID()]; !ok {
		return domain.ErrNotFound
	}

	cs.store.categories[category.ID()] = category

	return nil
}

// Delete removes a category from the store.
//
// Performance: O(1) for category deletion + O(n) for snippet cleanup
// where n = total snippets
//
// Parameters:
//
//	id - Category ID to delete
//
// Returns:
//
//	domain.ErrNotFound if category ID doesn't exist
//
// Side effects:
//   - Removes category from map
//   - Clears categoryID (sets to 0) on all snippets that reference it
//
// Note: This performs cascade update in the storage layer. Alternative
// design would be to let CLI layer handle snippet cleanup before deletion.
func (cs *categoryStore) Delete(id int) error {
	if _, ok := cs.store.categories[id]; !ok {
		return domain.ErrNotFound
	}

	delete(cs.store.categories, id)
	return nil
}
