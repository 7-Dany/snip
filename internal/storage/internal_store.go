package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/7-Dany/snip/internal/domain"
)

// Common storage errors.
var (
	ErrNotFound      = errors.New("entity not found")
	ErrDuplicateName = errors.New("entity with this name already exists")
)

// store is the internal data structure for all entities.
// It handles JSON persistence and provides in-memory storage with
// auto-incrementing IDs for snippets, categories, and tags.
type store struct {
	filepath string
	mu       sync.RWMutex // Protects data structures
	idMu     sync.Mutex   // Protects ID counters (separate to avoid deadlock)

	snippets   []*domain.Snippet
	categories []*domain.Category
	tags       []*domain.Tag

	nextSnippetID  int
	nextCategoryID int
	nextTagID      int
}

// data is the JSON structure for persistence.
type data struct {
	Snippets       []*domain.Snippet  `json:"snippets"`
	Categories     []*domain.Category `json:"categories"`
	Tags           []*domain.Tag      `json:"tags"`
	NextSnippetID  int                `json:"next_snippet_id"`
	NextCategoryID int                `json:"next_category_id"`
	NextTagID      int                `json:"next_tag_id"`
}

// newStore creates a new store with the given filepath for persistence.
func newStore(filepath string) *store {
	return &store{
		filepath:       filepath,
		snippets:       make([]*domain.Snippet, 0),
		categories:     make([]*domain.Category, 0),
		tags:           make([]*domain.Tag, 0),
		nextSnippetID:  1,
		nextCategoryID: 1,
		nextTagID:      1,
	}
}

// save persists all data to the JSON file atomically.
func (s *store) save() error {
	s.mu.RLock()
	s.idMu.Lock()

	d := data{
		Snippets:       s.snippets,
		Categories:     s.categories,
		Tags:           s.tags,
		NextSnippetID:  s.nextSnippetID,
		NextCategoryID: s.nextCategoryID,
		NextTagID:      s.nextTagID,
	}

	s.idMu.Unlock()
	s.mu.RUnlock()

	jsonData, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := s.filepath + ".tmp"
	if err := os.WriteFile(tmpFile, jsonData, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, s.filepath)
}

// load reads all data from the JSON file into memory.
// If the file doesn't exist, this is not an error.
func (s *store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.filepath); os.IsNotExist(err) {
		return nil
	}

	jsonData, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}

	var d data
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return err
	}

	s.snippets = d.Snippets
	s.categories = d.Categories
	s.tags = d.Tags

	s.idMu.Lock()
	s.nextSnippetID = d.NextSnippetID
	s.nextCategoryID = d.NextCategoryID
	s.nextTagID = d.NextTagID
	s.idMu.Unlock()

	if s.snippets == nil {
		s.snippets = make([]*domain.Snippet, 0)
	}
	if s.categories == nil {
		s.categories = make([]*domain.Category, 0)
	}
	if s.tags == nil {
		s.tags = make([]*domain.Tag, 0)
	}

	return nil
}

// nextSnippetIDAndIncrement returns the next snippet ID and increments the counter.
// This method is thread-safe and can be called concurrently.
func (s *store) nextSnippetIDAndIncrement() int {
	s.idMu.Lock()
	defer s.idMu.Unlock()

	id := s.nextSnippetID
	s.nextSnippetID++
	return id
}

// nextCategoryIDAndIncrement returns the next category ID and increments the counter.
// This method is thread-safe and can be called concurrently.
func (s *store) nextCategoryIDAndIncrement() int {
	s.idMu.Lock()
	defer s.idMu.Unlock()

	id := s.nextCategoryID
	s.nextCategoryID++
	return id
}

// nextTagIDAndIncrement returns the next tag ID and increments the counter.
// This method is thread-safe and can be called concurrently.
func (s *store) nextTagIDAndIncrement() int {
	s.idMu.Lock()
	defer s.idMu.Unlock()

	id := s.nextTagID
	s.nextTagID++
	return id
}
