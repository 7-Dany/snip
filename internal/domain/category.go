// Package domain contains the core business entities and logic for the snippet manager.
package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// Category represents a code snippet category for organizing snippets into logical groups.
type Category struct {
	id        int
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewCategory creates and returns a new Category with the given name.
// It returns ErrEmptyName if name is empty.
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
func (c *Category) ID() int { return c.id }

// Name returns the category's name.
func (c *Category) Name() string { return c.name }

// CreatedAt returns the category's creation time.
func (c *Category) CreatedAt() time.Time { return c.createdAt }

// UpdatedAt returns the category's last modification time.
func (c *Category) UpdatedAt() time.Time { return c.updatedAt }

// SetName updates the category's name and modification timestamp.
// It returns ErrEmptyName if name is empty.
func (c *Category) SetName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	c.name = name
	c.updatedAt = time.Now()
	return nil
}

// SetID sets the category's unique identifier.
// This should only be called by the storage layer.
func (c *Category) SetID(id int) {
	c.id = id
}

// String returns a string representation of the category.
func (c *Category) String() string {
	return fmt.Sprintf("Category{id=%d, name=%q}", c.id, c.name)
}

// Equal returns true if this category has the same data as other.
// Two categories are considered equal if all their fields match.
func (c *Category) Equal(other *Category) bool {
	if other == nil {
		return false
	}
	return c.id == other.id &&
		c.name == other.name &&
		c.createdAt.Equal(other.createdAt) &&
		c.updatedAt.Equal(other.updatedAt)
}

// MarshalJSON implements json.Marshaler.
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

// UnmarshalJSON implements json.Unmarshaler.
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

	// Validate loaded data
	if aux.Name == "" {
		return ErrEmptyName
	}

	c.id = aux.ID
	c.name = aux.Name
	c.createdAt = aux.CreatedAt
	c.updatedAt = aux.UpdatedAt
	return nil
}
