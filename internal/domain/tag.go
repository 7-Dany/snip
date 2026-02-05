// Package domain contains the core business entities and logic for the snippet manager.
package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// Tag represents a label that can be applied to code snippets for categorization.
type Tag struct {
	id        int
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewTag creates and returns a new Tag with the given name.
// It returns ErrEmptyName if name is empty.
func NewTag(name string) (*Tag, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	now := time.Now()
	return &Tag{
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the tag's unique identifier.
func (t *Tag) ID() int { return t.id }

// Name returns the tag's name.
func (t *Tag) Name() string { return t.name }

// CreatedAt returns the tag's creation time.
func (t *Tag) CreatedAt() time.Time { return t.createdAt }

// UpdatedAt returns the tag's last modification time.
func (t *Tag) UpdatedAt() time.Time { return t.updatedAt }

// SetName updates the tag's name and modification timestamp.
// It returns ErrEmptyName if name is empty.
func (t *Tag) SetName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	t.name = name
	t.updatedAt = time.Now()
	return nil
}

// SetID sets the tag's unique identifier.
// This should only be called by the storage layer.
func (t *Tag) SetID(id int) {
	t.id = id
}

// String returns a string representation of the tag.
func (t *Tag) String() string {
	return fmt.Sprintf("Tag{id=%d, name=%q}", t.id, t.name)
}

// Equal returns true if this tag has the same data as other.
// Two tags are considered equal if all their fields match.
func (t *Tag) Equal(other *Tag) bool {
	if other == nil {
		return false
	}
	return t.id == other.id &&
		t.name == other.name &&
		t.createdAt.Equal(other.createdAt) &&
		t.updatedAt.Equal(other.updatedAt)
}

// MarshalJSON implements json.Marshaler.
func (t *Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        t.id,
		Name:      t.name,
		CreatedAt: t.createdAt,
		UpdatedAt: t.updatedAt,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Tag) UnmarshalJSON(data []byte) error {
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

	t.id = aux.ID
	t.name = aux.Name
	t.createdAt = aux.CreatedAt
	t.updatedAt = aux.UpdatedAt
	return nil
}
