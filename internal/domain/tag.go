package domain

import "time"

// Tag represents a label that can be applied to code snippets.
// Tags are used to categorize and search snippets by topics or characteristics.
// All fields are unexported to maintain encapsulation and enforce validation.
type Tag struct {
	id        int
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewTag creates a new Tag with the given name.
// Returns ErrEmptyName if the name is empty.
func NewTag(name string) (*Tag, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	return &Tag{
		name:      name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ID returns the tag's unique identifier.
// The ID is 0 until set by the storage layer.
func (t *Tag) ID() int { return t.id }

// Name returns the tag's name.
func (t *Tag) Name() string { return t.name }

// CreatedAt returns when the tag was created.
func (t *Tag) CreatedAt() time.Time { return t.createdAt }

// UpdatedAt returns when the tag was last modified.
func (t *Tag) UpdatedAt() time.Time { return t.updatedAt }

// SetName updates the tag's name.
// Returns ErrEmptyName if the name is empty.
// Updates the tag's modification timestamp.
func (t *Tag) SetName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	t.name = name
	t.updatedAt = time.Now()
	return nil
}

// SetID sets the tag's unique identifier.
// This method is intended for use by the storage layer when persisting tags.
func (t *Tag) SetID(id int) {
	t.id = id
}
