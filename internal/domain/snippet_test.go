package domain

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestNewSnippet(t *testing.T) {
	t.Run("empty title returns error", func(t *testing.T) {
		snippet, err := NewSnippet("", "go", "code")

		if err == nil {
			t.Fatal("expected error for empty title, got nil")
		}

		if !errors.Is(err, ErrEmptyTitle) {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}

		if snippet != nil {
			t.Error("expected nil snippet on error")
		}
	})

	t.Run("empty language returns error", func(t *testing.T) {
		snippet, err := NewSnippet("Quick", "", "code")

		if err == nil {
			t.Fatal("expected error for empty language, got nil")
		}

		if !errors.Is(err, ErrEmptyLanguage) {
			t.Errorf("expected ErrEmptyLanguage, got %v", err)
		}

		if snippet != nil {
			t.Error("expected nil snippet on error")
		}
	})

	t.Run("empty code returns error", func(t *testing.T) {
		snippet, err := NewSnippet("Quick", "lang", "")

		if err == nil {
			t.Fatal("expected error for empty code, got nil")
		}

		if !errors.Is(err, ErrEmptyCode) {
			t.Errorf("expected ErrEmptyCode, got %v", err)
		}

		if snippet != nil {
			t.Error("expected nil snippet on error")
		}
	})

	t.Run("valid snippet is created", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")

		if err != nil {
			t.Fatalf("expected valid snippet, got error: %v", err)
		}

		if snippet == nil {
			t.Fatal("expected valid snippet, got nil")
		}

		if snippet.Title() != "title" {
			t.Errorf("expected title 'title', got %s", snippet.Title())
		}

		if snippet.Language() != "lang" {
			t.Errorf("expected language 'lang', got %s", snippet.Language())
		}

		if snippet.Code() != "code" {
			t.Errorf("expected code 'code', got %s", snippet.Code())
		}

		if snippet.ID() != 0 {
			t.Errorf("expected id to be 0, got %v", snippet.ID())
		}

		if snippet.CreatedAt().IsZero() {
			t.Error("expected createdAt to be set")
		}

		if snippet.UpdatedAt().IsZero() {
			t.Error("expected updatedAt to be set")
		}
	})
}

func TestSnippetTags(t *testing.T) {
	t.Run("Tags returns a copy", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.AddTag(1)
		snippet.AddTag(2)

		tags := snippet.Tags()
		tags[0] = 999

		if snippet.HasTag(999) {
			t.Error("modifying returned tags affected original snippet")
		}

		if !snippet.HasTag(1) {
			t.Error("original tag 1 should still exist")
		}
	})

	t.Run("AddTag prevents duplicates", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.AddTag(1)
		snippet.AddTag(1)

		tags := snippet.Tags()
		if len(tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(tags))
		}
	})

	t.Run("RemoveTag removes tag", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.AddTag(1)
		snippet.AddTag(2)

		if len(snippet.Tags()) != 2 {
			t.Errorf("expected 2 tags after adding, got %d", len(snippet.Tags()))
		}

		snippet.RemoveTag(1)

		if snippet.HasTag(1) {
			t.Error("tag 1 should be removed")
		}

		if !snippet.HasTag(2) {
			t.Error("tag 2 should still exist")
		}
	})

	t.Run("RemoveTag on non-existent tag does nothing", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.AddTag(1)
		snippet.RemoveTag(999) // Remove tag that doesn't exist

		if !snippet.HasTag(1) {
			t.Error("RemoveTag on non-existent tag affected existing tags")
		}

		if len(snippet.Tags()) != 1 {
			t.Errorf("expected 1 tag, got %d", len(snippet.Tags()))
		}
	})
}

func TestSnippetSetTitle(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = snippet.SetTitle("")

		if err == nil {
			t.Error("expected error for empty title")
		}

		if !errors.Is(err, ErrEmptyTitle) {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}
	})

	t.Run("updates title", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = snippet.SetTitle("hello")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if snippet.Title() != "hello" {
			t.Errorf("expected title 'hello', got %s", snippet.Title())
		}
	})
}

func TestSnippetSetLanguage(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = snippet.SetLanguage("")

		if err == nil {
			t.Error("expected error for empty language")
		}

		if !errors.Is(err, ErrEmptyLanguage) {
			t.Errorf("expected ErrEmptyLanguage, got %v", err)
		}
	})

	t.Run("updates language", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = snippet.SetLanguage("hello")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if snippet.Language() != "hello" {
			t.Errorf("expected language 'hello', got %s", snippet.Language())
		}
	})
}

func TestSnippetSetCode(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = snippet.SetCode("")

		if err == nil {
			t.Error("expected error for empty code")
		}

		if !errors.Is(err, ErrEmptyCode) {
			t.Errorf("expected ErrEmptyCode, got %v", err)
		}
	})

	t.Run("updates code", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = snippet.SetCode("hello")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if snippet.Code() != "hello" {
			t.Errorf("expected code 'hello', got %s", snippet.Code())
		}
	})
}

func TestSnippetSetDescription(t *testing.T) {
	t.Run("updates description", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.SetDescription("description")

		if snippet.Description() != "description" {
			t.Errorf("expected description 'description', got %s", snippet.Description())
		}
	})
}

func TestSnippetSetCategory(t *testing.T) {
	t.Run("updates category", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.SetCategory(1)

		if snippet.CategoryID() != 1 {
			t.Errorf("expected categoryID 1, got %d", snippet.CategoryID())
		}
	})
}

func TestSnippetSetID(t *testing.T) {
	t.Run("updates ID", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.SetID(42)

		if snippet.ID() != 42 {
			t.Errorf("expected ID 42, got %d", snippet.ID())
		}
	})
}

func TestSnippetMutations(t *testing.T) {
	t.Run("mutations update timestamp", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		originalTime := snippet.UpdatedAt()

		// Sleep to ensure time difference
		time.Sleep(2 * time.Millisecond)

		snippet.SetDescription("new desc")

		if !snippet.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to change after mutation")
		}
	})
}

func TestSnippetJSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		snippet, _ := NewSnippet("quicksort", "go", "func quicksort() {}")
		snippet.SetID(42)
		snippet.SetDescription("Fast sorting algorithm")
		snippet.SetCategory(5)
		snippet.AddTag(1)
		snippet.AddTag(2)

		data, err := json.Marshal(snippet)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		// Verify it's valid JSON
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Generated invalid JSON: %v", err)
		}

		// Verify fields are present
		if id, ok := result["id"].(float64); !ok || int(id) != 42 {
			t.Errorf("Expected id=42, got %v", result["id"])
		}
		if title, ok := result["title"].(string); !ok || title != "quicksort" {
			t.Errorf("Expected title='quicksort', got %v", result["title"])
		}
		if language, ok := result["language"].(string); !ok || language != "go" {
			t.Errorf("Expected language='go', got %v", result["language"])
		}
		if code, ok := result["code"].(string); !ok || code != "func quicksort() {}" {
			t.Errorf("Expected code='func quicksort() {}', got %v", result["code"])
		}
		if desc, ok := result["description"].(string); !ok || desc != "Fast sorting algorithm" {
			t.Errorf("Expected description='Fast sorting algorithm', got %v", result["description"])
		}
		if categoryID, ok := result["category_id"].(float64); !ok || int(categoryID) != 5 {
			t.Errorf("Expected category_id=5, got %v", result["category_id"])
		}
		if tags, ok := result["tags"].([]interface{}); !ok || len(tags) != 2 {
			t.Errorf("Expected 2 tags, got %v", result["tags"])
		}
		if _, ok := result["created_at"]; !ok {
			t.Error("Missing created_at field")
		}
		if _, ok := result["updated_at"]; !ok {
			t.Error("Missing updated_at field")
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		jsonData := []byte(`{
		"id": 42,
		"title": "quicksort",
		"language": "go",
		"code": "func quicksort() {}",
		"description": "Fast sorting algorithm",
		"category_id": 5,
		"tags": [1, 2, 3],
		"created_at": "2024-01-15T10:30:00Z",
		"updated_at": "2024-01-16T14:20:00Z"
	}`)

		var snippet Snippet
		err := json.Unmarshal(jsonData, &snippet)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if snippet.ID() != 42 {
			t.Errorf("Expected ID=42, got %d", snippet.ID())
		}
		if snippet.Title() != "quicksort" {
			t.Errorf("Expected title='quicksort', got '%s'", snippet.Title())
		}
		if snippet.Language() != "go" {
			t.Errorf("Expected language='go', got '%s'", snippet.Language())
		}
		if snippet.Code() != "func quicksort() {}" {
			t.Errorf("Expected code='func quicksort() {}', got '%s'", snippet.Code())
		}
		if snippet.Description() != "Fast sorting algorithm" {
			t.Errorf("Expected description='Fast sorting algorithm', got '%s'", snippet.Description())
		}
		if snippet.CategoryID() != 5 {
			t.Errorf("Expected category_id=5, got %d", snippet.CategoryID())
		}
		if len(snippet.Tags()) != 3 {
			t.Errorf("Expected 3 tags, got %d", len(snippet.Tags()))
		}
		if snippet.CreatedAt().IsZero() {
			t.Error("CreatedAt should not be zero")
		}
		if snippet.UpdatedAt().IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
	})

	t.Run("round trip", func(t *testing.T) {
		original, _ := NewSnippet("binary search", "python", "def binary_search():")
		original.SetID(123)
		original.SetDescription("Efficient search")
		original.SetCategory(10)
		original.AddTag(5)
		original.AddTag(6)

		// Marshal
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// Unmarshal
		var restored Snippet
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// Verify all fields match
		if restored.ID() != original.ID() {
			t.Errorf("ID mismatch: expected %d, got %d", original.ID(), restored.ID())
		}
		if restored.Title() != original.Title() {
			t.Errorf("Title mismatch: expected %s, got %s", original.Title(), restored.Title())
		}
		if restored.Language() != original.Language() {
			t.Errorf("Language mismatch: expected %s, got %s", original.Language(), restored.Language())
		}
		if restored.Code() != original.Code() {
			t.Errorf("Code mismatch")
		}
		if restored.Description() != original.Description() {
			t.Errorf("Description mismatch")
		}
		if restored.CategoryID() != original.CategoryID() {
			t.Errorf("CategoryID mismatch: expected %d, got %d", original.CategoryID(), restored.CategoryID())
		}
		if len(restored.Tags()) != len(original.Tags()) {
			t.Errorf("Tags count mismatch: expected %d, got %d", len(original.Tags()), len(restored.Tags()))
		}
		if !restored.CreatedAt().Equal(original.CreatedAt()) {
			t.Error("CreatedAt mismatch")
		}
		if !restored.UpdatedAt().Equal(original.UpdatedAt()) {
			t.Error("UpdatedAt mismatch")
		}
	})

	t.Run("unmarshal invalid JSON", func(t *testing.T) {
		invalidJSON := []byte(`{"id": "not a number"}`)

		var snippet Snippet
		err := json.Unmarshal(invalidJSON, &snippet)

		if err == nil {
			t.Error("UnmarshalJSON should return error for invalid JSON")
		}
	})

	t.Run("unmarshal empty tags", func(t *testing.T) {
		jsonData := []byte(`{
		"id": 1,
		"title": "test",
		"language": "go",
		"code": "code",
		"description": "",
		"category_id": 0,
		"tags": null,
		"created_at": "2024-01-15T10:30:00Z",
		"updated_at": "2024-01-15T10:30:00Z"
	}`)

		var snippet Snippet
		if err := json.Unmarshal(jsonData, &snippet); err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		// nil tags should work fine
		if snippet.Tags() == nil {
			t.Error("Tags should not be nil after unmarshal")
		}
	})
}
