// Package domain contains the core business logic and entities for the snippet manager.
package domain

import (
	"slices"
	"time"
)

// Snippet represents a code snippet with metadata including title, language,
// code content, optional description, category, and tags.
// All fields are unexported to maintain encapsulation and enforce validation.
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

// NewSnippet creates a new Snippet with the given title, language, and code.
// It validates that all required fields are non-empty and initializes timestamps.
// Returns ErrEmptyTitle, ErrEmptyLanguage, or ErrEmptyCode if validation fails.
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
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ID returns the snippet's unique identifier.
// The ID is 0 until set by the storage layer.
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

// CreatedAt returns when the snippet was created.
func (s *Snippet) CreatedAt() time.Time { return s.createdAt }

// UpdatedAt returns when the snippet was last modified.
func (s *Snippet) UpdatedAt() time.Time { return s.updatedAt }

// Tags returns a copy of the snippet's tag IDs.
// Modifying the returned slice does not affect the snippet's internal tags.
func (s *Snippet) Tags() []int {
	tags := make([]int, len(s.tags))
	copy(tags, s.tags)
	return tags
}

// SetTitle updates the snippet's title and modification timestamp, if empty it returns ErrEmptyTitle.
func (s *Snippet) SetTitle(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}

	s.title = title
	s.updatedAt = time.Now()
	return nil
}

// SetLanguage updates the snippet's langauge and modification timestamp, if empty it returns ErrEmptyLanguage.
func (s *Snippet) SetLanguage(language string) error {
	if language == "" {
		return ErrEmptyLanguage
	}

	s.language = language
	s.updatedAt = time.Now()
	return nil
}

// SetCode updates the snippet's code and modification timestamp, if empty it returns ErrEmptyCode.
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
// Updates the snippet's modification timestamp.
func (s *Snippet) AddTag(tagID int) {
	if slices.Contains(s.tags, tagID) {
		return
	}
	s.tags = append(s.tags, tagID)
	s.updatedAt = time.Now()
}

// RemoveTag removes a tag ID from the snippet.
// If the tag ID is not present, this is a no-op.
// Updates the snippet's modification timestamp.
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
// This method is intended for use by the storage layer when persisting snippets.
func (s *Snippet) SetID(id int) {
	s.id = id
}
