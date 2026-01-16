package storage

import (
	"regexp"
	"strings"

	"github.com/7-Dany/snip/internal/domain"
)

var wordCleaner = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// extractWords takes text and returns unique lowercase words
func extractWords(text string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, w := range strings.Fields(text) {
		word := strings.ToLower(wordCleaner.ReplaceAllString(w, ""))
		if word == "" {
			continue
		}
		if _, exists := seen[word]; !exists {
			seen[word] = struct{}{}
			result = append(result, word)
		}
	}
	return result
}

// buildSearchableText efficiently combines all searchable fields
func buildSearchableText(snippet *domain.Snippet) string {
	estimatedSize := len(snippet.Title()) +
		len(snippet.Description()) +
		len(snippet.Code()) +
		len(snippet.Language()) +
		3
	var sb strings.Builder
	sb.Grow(estimatedSize)
	sb.WriteString(snippet.Title())
	sb.WriteString(" ")
	sb.WriteString(snippet.Description())
	sb.WriteString(" ")
	sb.WriteString(snippet.Code())
	sb.WriteString(" ")
	sb.WriteString(snippet.Language())
	return sb.String()
}

// containsID checks if a slice contains a specific ID
func containsID(ids []int, target int) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

// indexSnippet adds a snippet to the inverted index
func (s *Store) indexSnippet(snippet *domain.Snippet) {
	text := buildSearchableText(snippet)
	words := extractWords(text)

	for _, word := range words {
		// Prevent duplicate IDs
		if !containsID(s.searchIndex[word], snippet.ID()) {
			s.searchIndex[word] = append(s.searchIndex[word], snippet.ID())
		}
	}
}

// removeFromIndex removes a snippet from the inverted index
func (s *Store) removeFromIndex(snippetID int) {
	for k, v := range s.searchIndex {
		for i, id := range v {
			if id == snippetID {
				s.searchIndex[k] = append(s.searchIndex[k][:i], s.searchIndex[k][i+1:]...)

				// Clean up empty word entries
				if len(s.searchIndex[k]) == 0 {
					delete(s.searchIndex, k)
				}
				break
			}
		}
	}
}

// searchWithIndex searches using the inverted index
func (s *Store) searchWithIndex(query string) ([]*domain.Snippet, error) {
	words := extractWords(query)

	if len(words) == 0 {
		return []*domain.Snippet{}, nil
	}

	set := make(map[int]struct{})
	snippets := []*domain.Snippet{}

	for _, word := range words {
		ids := s.searchIndex[word]
		for _, id := range ids {
			if _, ok := set[id]; !ok {
				snippets = append(snippets, s.snippets[id])
				set[id] = struct{}{}
			}
		}
	}

	return snippets, nil
}
