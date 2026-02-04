// Package domain contains the core business entities and logic for the snippet manager.
package domain

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"
)

// Snippet represents a code snippet with metadata including title, language,
// code content, optional description, category, and tags.
type Snippet struct {
	id          int
	title       string
	language    string
	categoryID  int
	tags        []int
	description string
	code        string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewSnippet creates and returns a new Snippet with the given title, language, and code.
// It returns ErrEmptyTitle, ErrEmptyLanguage, or ErrEmptyCode if validation fails.
func NewSnippet(title, language, code string) (*Snippet, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if language == "" {
		return nil, ErrEmptyLanguage
	}
	if code == "" {
		return nil, ErrEmptyCode
	}
	return &Snippet{
		title:     title,
		language:  language,
		code:      code,
		tags:      []int{},
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ID returns the snippet's unique identifier.
func (s *Snippet) ID() int { return s.id }

// Title returns the snippet's title.
func (s *Snippet) Title() string { return s.title }

// Language returns the programming language of the snippet.
func (s *Snippet) Language() string { return s.language }

// Code returns the snippet's code content.
func (s *Snippet) Code() string { return s.code }

// Description returns the snippet's optional description.
func (s *Snippet) Description() string { return s.description }

// CategoryID returns the ID of the snippet's category, or 0 if uncategorized.
func (s *Snippet) CategoryID() int { return s.categoryID }

// CreatedAt returns the snippet's creation time.
func (s *Snippet) CreatedAt() time.Time { return s.createdAt }

// UpdatedAt returns the snippet's last modification time.
func (s *Snippet) UpdatedAt() time.Time { return s.updatedAt }

// Tags returns a copy of the snippet's tag IDs.
// Modifying the returned slice does not affect the snippet's internal tags.
func (s *Snippet) Tags() []int {
	if s.tags == nil {
		return []int{}
	}
	tags := make([]int, len(s.tags))
	copy(tags, s.tags)
	return tags
}

// SetTitle updates the snippet's title and modification timestamp.
// It returns ErrEmptyTitle if title is empty.
func (s *Snippet) SetTitle(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	s.title = title
	s.updatedAt = time.Now()
	return nil
}

// SetLanguage updates the snippet's language and modification timestamp.
// It returns ErrEmptyLanguage if language is empty.
func (s *Snippet) SetLanguage(language string) error {
	if language == "" {
		return ErrEmptyLanguage
	}
	s.language = language
	s.updatedAt = time.Now()
	return nil
}

// SetCode updates the snippet's code and modification timestamp.
// It returns ErrEmptyCode if code is empty.
func (s *Snippet) SetCode(code string) error {
	if code == "" {
		return ErrEmptyCode
	}
	s.code = code
	s.updatedAt = time.Now()
	return nil
}

// SetDescription updates the snippet's description and modification timestamp.
func (s *Snippet) SetDescription(description string) {
	s.description = description
	s.updatedAt = time.Now()
}

// SetCategory updates the snippet's category ID and modification timestamp.
func (s *Snippet) SetCategory(catID int) {
	s.categoryID = catID
	s.updatedAt = time.Now()
}

// AddTag adds a tag ID to the snippet if not already present.
// Duplicate tag IDs are silently ignored.
func (s *Snippet) AddTag(tagID int) {
	if slices.Contains(s.tags, tagID) {
		return
	}
	s.tags = append(s.tags, tagID)
	s.updatedAt = time.Now()
}

// RemoveTag removes a tag ID from the snippet.
// If the tag ID is not present, this is a no-op.
func (s *Snippet) RemoveTag(tagID int) {
	for i, t := range s.tags {
		if t == tagID {
			s.tags = append(s.tags[:i], s.tags[i+1:]...)
			s.updatedAt = time.Now()
			return
		}
	}
}

// HasTag reports whether the snippet has the given tag ID.
func (s *Snippet) HasTag(tagID int) bool {
	return slices.Contains(s.tags, tagID)
}

// SetID sets the snippet's unique identifier.
// This should only be called by the storage layer.
func (s *Snippet) SetID(id int) {
	s.id = id
}

// String returns a string representation of the snippet.
func (s *Snippet) String() string {
	return fmt.Sprintf("Snippet{id=%d, title=%q, language=%q}", s.id, s.title, s.language)
}

// Equal returns true if this snippet has the same data as other.
// Two snippets are considered equal if all their fields match.
func (s *Snippet) Equal(other *Snippet) bool {
	if other == nil {
		return false
	}
	return s.id == other.id &&
		s.title == other.title &&
		s.language == other.language &&
		s.code == other.code &&
		s.description == other.description &&
		s.categoryID == other.categoryID &&
		slices.Equal(s.tags, other.tags) &&
		s.createdAt.Equal(other.createdAt) &&
		s.updatedAt.Equal(other.updatedAt)
}

// MarshalJSON implements json.Marshaler.
func (s *Snippet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          int       `json:"id"`
		Title       string    `json:"title"`
		Language    string    `json:"language"`
		Code        string    `json:"code"`
		Description string    `json:"description"`
		CategoryID  int       `json:"category_id"`
		Tags        []int     `json:"tags"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}{
		ID:          s.id,
		Title:       s.title,
		Language:    s.language,
		Code:        s.code,
		Description: s.description,
		CategoryID:  s.categoryID,
		Tags:        s.tags,
		CreatedAt:   s.createdAt,
		UpdatedAt:   s.updatedAt,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Snippet) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID          int       `json:"id"`
		Title       string    `json:"title"`
		Language    string    `json:"language"`
		Code        string    `json:"code"`
		Description string    `json:"description"`
		CategoryID  int       `json:"category_id"`
		Tags        []int     `json:"tags"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Validate loaded data
	if aux.Title == "" {
		return ErrEmptyTitle
	}
	if aux.Language == "" {
		return ErrEmptyLanguage
	}
	if aux.Code == "" {
		return ErrEmptyCode
	}

	s.id = aux.ID
	s.title = aux.Title
	s.language = aux.Language
	s.code = aux.Code
	s.description = aux.Description
	s.categoryID = aux.CategoryID
	s.tags = aux.Tags
	if s.tags == nil {
		s.tags = []int{}
	}
	s.createdAt = aux.CreatedAt
	s.updatedAt = aux.UpdatedAt
	return nil
}
