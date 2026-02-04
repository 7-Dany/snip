package storage

import "github.com/7-Dany/snip/internal/domain"

// categoryRepository implements domain.CategoryRepository using the internal store.
type categoryRepository struct {
	store *store
}

// newCategoryRepository creates a new category repository.
func newCategoryRepository(s *store) *categoryRepository {
	return &categoryRepository{store: s}
}

// List returns all categories.
func (r *categoryRepository) List() ([]*domain.Category, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	result := make([]*domain.Category, len(r.store.categories))
	copy(result, r.store.categories)
	return result, nil
}

// FindByID finds a category by its ID.
func (r *categoryRepository) FindByID(id int) (*domain.Category, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	for _, category := range r.store.categories {
		if category.ID() == id {
			return category, nil
		}
	}
	return nil, ErrNotFound
}

// FindByName finds a category by its name.
func (r *categoryRepository) FindByName(name string) (*domain.Category, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	for _, category := range r.store.categories {
		if category.Name() == name {
			return category, nil
		}
	}
	return nil, ErrNotFound
}

// Create adds a new category and assigns it an ID.
// Returns ErrDuplicateName if a category with the same name exists.
func (r *categoryRepository) Create(category *domain.Category) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	// Check for duplicate name
	for _, existing := range r.store.categories {
		if existing.Name() == category.Name() {
			return ErrDuplicateName
		}
	}

	id := r.store.nextCategoryIDAndIncrement()
	category.SetID(id)
	r.store.categories = append(r.store.categories, category)
	return nil
}

// Update replaces an existing category.
// Returns ErrDuplicateName if another category with the same name exists.
func (r *categoryRepository) Update(category *domain.Category) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	// Check for duplicate name (excluding current category)
	for _, existing := range r.store.categories {
		if existing.Name() == category.Name() && existing.ID() != category.ID() {
			return ErrDuplicateName
		}
	}

	for i, existing := range r.store.categories {
		if existing.ID() == category.ID() {
			r.store.categories[i] = category
			return nil
		}
	}
	return ErrNotFound
}

// Delete removes a category by ID.
func (r *categoryRepository) Delete(id int) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for i, category := range r.store.categories {
		if category.ID() == id {
			r.store.categories = append(r.store.categories[:i], r.store.categories[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
