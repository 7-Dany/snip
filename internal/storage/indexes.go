// Package storage provides JSON-based persistence with inverted index search.
// This file implements the search index functionality for efficient snippet searching.
package storage

import (
	"regexp"
	"slices"
	"strings"

	"github.com/7-Dany/snip/internal/domain"
)

// wordCleaner is a compiled regex for removing non-alphanumeric characters.
// Used by extractWords to normalize text for search indexing.
//
// Pattern: [^a-zA-Z0-9]+ (one or more non-alphanumeric chars)
var wordCleaner = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// extractWords takes text and returns unique lowercase words for indexing.
// This is the core text normalization function for the search engine.
//
// Performance: O(n) where n = text length
//
// Parameters:
//
//	text - Raw text to extract words from
//
// Returns:
//
//	Slice of unique lowercase words (order preserved)
//
// Normalization rules:
//  1. Removes dots between digits (1.21 → 121)
//  2. Removes hyphens (Quick-Sort → QuickSort)
//  3. Replaces non-alphanumeric with spaces
//  4. Converts to lowercase
//  5. Removes duplicates (case-insensitive)
//
// Examples:
//
//	"Quick-Sort 1.21" → ["quicksort", "121"]
//	"Hello World!" → ["hello", "world"]
//	"GO go Go" → ["go"] (deduplication)
func extractWords(text string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	// First: remove dots between digits (1.21 → 121)
	dotCleaner := regexp.MustCompile(`(\d)\.(\d)`)
	text = dotCleaner.ReplaceAllString(text, "$1$2")

	// Second: remove hyphens (Quick-Sort → QuickSort)
	text = strings.ReplaceAll(text, "-", "")

	// Third: replace remaining non-alphanumeric with spaces
	cleaned := wordCleaner.ReplaceAllString(text, " ")

	for w := range strings.FieldsSeq(cleaned) {
		word := strings.ToLower(w)
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

// buildSearchableText efficiently combines all searchable fields from a snippet.
// Pre-allocates string builder capacity to minimize allocations.
//
// Performance: O(n) where n = total text length, single allocation
//
// Parameters:
//
//	snippet - Snippet to build searchable text from
//
// Returns:
//
//	Space-separated concatenation of: title, description, code, language
//
// Implementation note: Uses strings.Builder with pre-calculated capacity
// to avoid multiple allocations during concatenation.
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

// containsID checks if a slice contains a specific ID.
// Helper function for preventing duplicate IDs in search index.
//
// Performance: O(n) where n = slice length (uses slices.Contains)
//
// Parameters:
//
//	ids - Slice of snippet IDs
//	target - ID to search for
//
// Returns:
//
//	true if target exists in ids, false otherwise
func containsID(ids []int, target int) bool {
	return slices.Contains(ids, target)
}

// indexSnippet adds a snippet to the inverted index.
// Extracts words from snippet and creates word→[snippetIDs] mappings.
//
// Performance: O(w) where w = unique words in snippet
//
// Parameters:
//
//	snippet - Snippet to index (must have valid ID assigned)
//
// Side effects:
//
//	Updates s.searchIndex with new word→ID mappings
//	Prevents duplicate snippet IDs for same word
//
// Implementation:
//  1. Extracts searchable text (title, desc, code, language)
//  2. Normalizes text into unique words
//  3. Adds snippet ID to each word's posting list
//  4. Skips if ID already in posting list (idempotent)
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

// removeFromIndex removes a snippet from the inverted index.
// Removes snippet ID from all word posting lists and cleans up empty entries.
//
// Performance: O(w*p) where w = words, p = avg posting list length
//
// Parameters:
//
//	snippetID - ID of snippet to remove from index
//
// Side effects:
//   - Removes snippetID from all posting lists
//   - Deletes word entries with empty posting lists
//   - Updates s.searchIndex
//
// Implementation note:
//
//	Iterates all word→[IDs] mappings because we don't store reverse index.
//	Cleans up empty word entries to prevent memory leaks.
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

// searchWithIndex searches using the inverted index.
// Performs multi-word OR search: returns snippets matching any query word.
//
// Performance: O(w*p) where w = query words, p = avg posting list length
//
// Parameters:
//
//	query - Search query (will be normalized into words)
//
// Returns:
//
//	Snippets matching any word in query (no ranking)
//	Empty slice if no matches or empty query
//	Never returns error (always succeeds)
//
// Search semantics:
//   - OR search: matches ANY word
//   - Case-insensitive
//   - Deduplicates results (each snippet once)
//   - No ranking/scoring (order undefined)
//
// Examples:
//
//	"quick sort" → snippets with "quick" OR "sort"
//	"" → empty result
//	"nonexistent" → empty result
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
