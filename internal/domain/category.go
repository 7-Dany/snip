package domain

import (
	"encoding/json"
	"time"
)

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

// MarshalJSON implements the json.Marshaler interface for Category.
// This allows Categories with unexported fields to be serialized to JSON
// while maintaining encapsulation and preventing direct field access.
//
// Design: Creates anonymous struct with exported fields for marshaling
//
// Returns:
//
//	[]byte - JSON representation of the category
//	error  - Marshaling errors (rare, usually nil)
func (c *Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        c.id,
		Name:      c.name,
		CreatedAt: c.createdAt,
		UpdatedAt: c.updatedAt,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for Category.
// This allows JSON data to be deserialized into Categories with unexported fields
// while maintaining encapsulation.
//
// Design: Uses anonymous struct to unmarshal, then copies to unexported fields
//
// Parameters:
//
//	data - JSON bytes to unmarshal
//
// Returns:
//
//	error - Unmarshal errors (invalid JSON, type mismatches, etc.)
//
// Note: This bypasses validation in NewCategory() and SetName().
// The storage layer is trusted to only load valid data that was previously saved.
func (c *Category) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	c.id = aux.ID
	c.name = aux.Name
	c.createdAt = aux.CreatedAt
	c.updatedAt = aux.UpdatedAt
	return nil
}
