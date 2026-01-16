package domain

import "time"

// Category represents a code snippet category.
// Categories are used to organize snippets into logical groups.
// All fields are unexported to maintain encapsulation and enforce validation.
type Category struct {
	id        int
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewCategory creates a new Category with the given name.
// Returns ErrEmptyName if the name is empty.
func NewCategory(name string) (*Category, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	return &Category{
		name:      name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ID returns the category's unique identifier.
// The ID is 0 until set by the storage layer.
func (c *Category) ID() int { return c.id }

// Name returns the category's name.
func (c *Category) Name() string { return c.name }

// CreatedAt returns when the category was created.
func (c *Category) CreatedAt() time.Time { return c.createdAt }

// UpdatedAt returns when the category was last modified.
func (c *Category) UpdatedAt() time.Time { return c.updatedAt }

// SetName updates the category's name.
// Returns ErrEmptyName if the name is empty.
// Updates the category's modification timestamp.
func (c *Category) SetName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	c.name = name
	c.updatedAt = time.Now()
	return nil
}

// SetID sets the category's unique identifier.
// This method is intended for use by the storage layer when persisting categories.
func (c *Category) SetID(id int) {
	c.id = id
}
