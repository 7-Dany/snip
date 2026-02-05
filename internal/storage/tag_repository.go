package storage

import "github.com/7-Dany/snip/internal/domain"

// tagRepository implements domain.TagRepository using the internal store.
type tagRepository struct {
	store *store
}

// newTagRepository creates a new tag repository.
func newTagRepository(s *store) *tagRepository {
	return &tagRepository{store: s}
}

// List returns all tags.
func (r *tagRepository) List() ([]*domain.Tag, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	result := make([]*domain.Tag, len(r.store.tags))
	copy(result, r.store.tags)
	return result, nil
}

// FindByID finds a tag by its ID.
func (r *tagRepository) FindByID(id int) (*domain.Tag, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for _, tag := range r.store.tags {
		if tag.ID() == id {
			return tag, nil
		}
	}
	return nil, ErrNotFound
}

// FindByName finds a tag by its name.
func (r *tagRepository) FindByName(name string) (*domain.Tag, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for _, tag := range r.store.tags {
		if tag.Name() == name {
			return tag, nil
		}
	}
	return nil, ErrNotFound
}

// Create adds a new tag and assigns it an ID.
// Returns ErrDuplicateName if a tag with the same name exists.
func (r *tagRepository) Create(tag *domain.Tag) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	// Check for duplicate name
	for _, existing := range r.store.tags {
		if existing.Name() == tag.Name() {
			return ErrDuplicateName
		}
	}

	id := r.store.nextTagIDAndIncrement()
	tag.SetID(id)
	r.store.tags = append(r.store.tags, tag)
	return nil
}

// Update replaces an existing tag.
// Returns ErrDuplicateName if another tag with the same name exists.
func (r *tagRepository) Update(tag *domain.Tag) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	// Check for duplicate name (excluding current tag)
	for _, existing := range r.store.tags {
		if existing.Name() == tag.Name() && existing.ID() != tag.ID() {
			return ErrDuplicateName
		}
	}

	for i, existing := range r.store.tags {
		if existing.ID() == tag.ID() {
			r.store.tags[i] = tag
			return nil
		}
	}
	return ErrNotFound
}

// Delete removes a tag by ID.
func (r *tagRepository) Delete(id int) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for i, tag := range r.store.tags {
		if tag.ID() == id {
			r.store.tags = append(r.store.tags[:i], r.store.tags[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
