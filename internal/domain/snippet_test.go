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
		snippet, err := NewSnippet("Quick Sort", "", "code")

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
		snippet, err := NewSnippet("Quick Sort", "go", "")

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
		snippet, err := NewSnippet("Quick Sort", "go", "func quicksort() {}")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if snippet == nil {
			t.Fatal("expected snippet to be non-nil")
		}

		if snippet.Title() != "Quick Sort" {
			t.Errorf("expected title 'Quick Sort', got %q", snippet.Title())
		}

		if snippet.Language() != "go" {
			t.Errorf("expected language 'go', got %q", snippet.Language())
		}

		if snippet.Code() != "func quicksort() {}" {
			t.Errorf("expected code 'func quicksort() {}', got %q", snippet.Code())
		}

		if snippet.ID() != 0 {
			t.Errorf("expected ID 0, got %d", snippet.ID())
		}

		if snippet.CreatedAt().IsZero() {
			t.Error("expected CreatedAt to be set")
		}

		if snippet.UpdatedAt().IsZero() {
			t.Error("expected UpdatedAt to be set")
		}

		if !snippet.CreatedAt().Equal(snippet.UpdatedAt()) {
			t.Error("expected CreatedAt and UpdatedAt to be equal for new snippet")
		}

		if len(snippet.Tags()) != 0 {
			t.Errorf("expected empty tags, got %v", snippet.Tags())
		}
	})
}

func TestSnippet_SetID(t *testing.T) {
	t.Run("sets ID", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "lang", "code")

		snippet.SetID(42)

		if snippet.ID() != 42 {
			t.Errorf("expected ID 42, got %d", snippet.ID())
		}
	})
}

func TestSnippet_SetTitle(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "Binary Search", "go", "code")
		originalTitle := snippet.Title()
		originalUpdatedAt := snippet.UpdatedAt()

		err := snippet.SetTitle("")

		if err == nil {
			t.Fatal("expected error for empty title")
		}

		if !errors.Is(err, ErrEmptyTitle) {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}

		if snippet.Title() != originalTitle {
			t.Error("title should not change when SetTitle returns error")
		}

		if !snippet.UpdatedAt().Equal(originalUpdatedAt) {
			t.Error("UpdatedAt should not change when SetTitle returns error")
		}
	})

	t.Run("updates title", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "Binary Search", "go", "code")

		err := snippet.SetTitle("Linear Search")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if snippet.Title() != "Linear Search" {
			t.Errorf("expected title 'Linear Search', got %q", snippet.Title())
		}
	})

	t.Run("updates timestamp", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "Binary Search", "go", "code")
		originalTime := snippet.UpdatedAt()

		time.Sleep(2 * time.Millisecond)

		err := snippet.SetTitle("New Title")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !snippet.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to be updated after SetTitle")
		}
	})
}

func TestSnippet_SetLanguage(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		originalLanguage := snippet.Language()
		originalUpdatedAt := snippet.UpdatedAt()

		err := snippet.SetLanguage("")

		if err == nil {
			t.Fatal("expected error for empty language")
		}

		if !errors.Is(err, ErrEmptyLanguage) {
			t.Errorf("expected ErrEmptyLanguage, got %v", err)
		}

		if snippet.Language() != originalLanguage {
			t.Error("language should not change when SetLanguage returns error")
		}

		if !snippet.UpdatedAt().Equal(originalUpdatedAt) {
			t.Error("UpdatedAt should not change when SetLanguage returns error")
		}
	})

	t.Run("updates language", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")

		err := snippet.SetLanguage("python")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if snippet.Language() != "python" {
			t.Errorf("expected language 'python', got %q", snippet.Language())
		}
	})
}

func TestSnippet_SetCode(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		originalCode := snippet.Code()
		originalUpdatedAt := snippet.UpdatedAt()

		err := snippet.SetCode("")

		if err == nil {
			t.Fatal("expected error for empty code")
		}

		if !errors.Is(err, ErrEmptyCode) {
			t.Errorf("expected ErrEmptyCode, got %v", err)
		}

		if snippet.Code() != originalCode {
			t.Error("code should not change when SetCode returns error")
		}

		if !snippet.UpdatedAt().Equal(originalUpdatedAt) {
			t.Error("UpdatedAt should not change when SetCode returns error")
		}
	})

	t.Run("updates code", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")

		err := snippet.SetCode("new code")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if snippet.Code() != "new code" {
			t.Errorf("expected code 'new code', got %q", snippet.Code())
		}
	})
}

func TestSnippet_SetDescription(t *testing.T) {
	t.Run("updates description", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")

		snippet.SetDescription("A helpful description")

		if snippet.Description() != "A helpful description" {
			t.Errorf("expected description 'A helpful description', got %q", snippet.Description())
		}
	})

	t.Run("empty description is valid", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")

		snippet.SetDescription("")

		if snippet.Description() != "" {
			t.Errorf("expected empty description, got %q", snippet.Description())
		}
	})

	t.Run("updates timestamp", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		originalTime := snippet.UpdatedAt()

		time.Sleep(2 * time.Millisecond)

		snippet.SetDescription("description")

		if !snippet.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to be updated after SetDescription")
		}
	})
}

func TestSnippet_SetCategory(t *testing.T) {
	t.Run("updates category", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")

		snippet.SetCategory(5)

		if snippet.CategoryID() != 5 {
			t.Errorf("expected categoryID 5, got %d", snippet.CategoryID())
		}
	})

	t.Run("updates timestamp", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		originalTime := snippet.UpdatedAt()

		time.Sleep(2 * time.Millisecond)

		snippet.SetCategory(1)

		if !snippet.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to be updated after SetCategory")
		}
	})
}

func TestSnippet_Tags(t *testing.T) {
	t.Run("Tags returns a copy", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
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
		snippet := mustCreateSnippet(t, "title", "go", "code")

		snippet.AddTag(1)
		snippet.AddTag(1)
		snippet.AddTag(1)

		tags := snippet.Tags()
		if len(tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(tags))
		}
	})

	t.Run("AddTag updates timestamp", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		originalTime := snippet.UpdatedAt()

		time.Sleep(2 * time.Millisecond)

		snippet.AddTag(1)

		if !snippet.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to be updated after AddTag")
		}
	})

	t.Run("AddTag duplicate does not update timestamp", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		snippet.AddTag(1)

		time.Sleep(2 * time.Millisecond)
		lastUpdate := snippet.UpdatedAt()

		snippet.AddTag(1) // Duplicate

		if snippet.UpdatedAt().After(lastUpdate) {
			t.Error("UpdatedAt should not change when adding duplicate tag")
		}
	})

	t.Run("RemoveTag removes tag", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		snippet.AddTag(1)
		snippet.AddTag(2)
		snippet.AddTag(3)

		if len(snippet.Tags()) != 3 {
			t.Errorf("expected 3 tags after adding, got %d", len(snippet.Tags()))
		}

		snippet.RemoveTag(2)

		if snippet.HasTag(2) {
			t.Error("tag 2 should be removed")
		}

		if !snippet.HasTag(1) {
			t.Error("tag 1 should still exist")
		}

		if !snippet.HasTag(3) {
			t.Error("tag 3 should still exist")
		}

		if len(snippet.Tags()) != 2 {
			t.Errorf("expected 2 tags after removal, got %d", len(snippet.Tags()))
		}
	})

	t.Run("RemoveTag updates timestamp", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		snippet.AddTag(1)

		time.Sleep(2 * time.Millisecond)
		originalTime := snippet.UpdatedAt()

		snippet.RemoveTag(1)

		if !snippet.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to be updated after RemoveTag")
		}
	})

	t.Run("RemoveTag on non-existent tag does nothing", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		snippet.AddTag(1)

		time.Sleep(2 * time.Millisecond)
		lastUpdate := snippet.UpdatedAt()

		snippet.RemoveTag(999) // Remove tag that doesn't exist

		if !snippet.HasTag(1) {
			t.Error("RemoveTag on non-existent tag affected existing tags")
		}

		if len(snippet.Tags()) != 1 {
			t.Errorf("expected 1 tag, got %d", len(snippet.Tags()))
		}

		if snippet.UpdatedAt().After(lastUpdate) {
			t.Error("UpdatedAt should not change when removing non-existent tag")
		}
	})

	t.Run("HasTag returns correct value", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")
		snippet.AddTag(1)
		snippet.AddTag(3)

		if !snippet.HasTag(1) {
			t.Error("expected HasTag(1) to be true")
		}

		if snippet.HasTag(2) {
			t.Error("expected HasTag(2) to be false")
		}

		if !snippet.HasTag(3) {
			t.Error("expected HasTag(3) to be true")
		}
	})
}

func TestSnippet_String(t *testing.T) {
	t.Run("with ID, title, and language", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "Quick Sort", "go", "func quicksort() {}")
		snippet.SetID(42)

		got := snippet.String()
		want := `Snippet{id=42, title="Quick Sort", language="go"}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("with zero ID", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "Binary Search", "python", "def binary_search():")

		got := snippet.String()
		want := `Snippet{id=0, title="Binary Search", language="python"}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("with special characters", func(t *testing.T) {
		snippet := mustCreateSnippet(t, `"Hello World"`, "c++", "code")
		snippet.SetID(1)

		got := snippet.String()
		want := `Snippet{id=1, title="\"Hello World\"", language="c++"}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestSnippet_Equal(t *testing.T) {
	t.Run("identical snippets are equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{
			id: 1, title: "Quick Sort", language: "go", code: "code",
			description: "desc", categoryID: 5, tags: []int{1, 2},
			createdAt: now, updatedAt: now,
		}
		snippet2 := &Snippet{
			id: 1, title: "Quick Sort", language: "go", code: "code",
			description: "desc", categoryID: 5, tags: []int{1, 2},
			createdAt: now, updatedAt: now,
		}

		if !snippet1.Equal(snippet2) {
			t.Error("expected identical snippets to be equal")
		}
	})

	t.Run("same instance is equal to itself", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")

		if !snippet.Equal(snippet) {
			t.Error("expected snippet to be equal to itself")
		}
	})

	t.Run("different IDs are not equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{id: 1, title: "title", language: "go", code: "code", createdAt: now, updatedAt: now}
		snippet2 := &Snippet{id: 2, title: "title", language: "go", code: "code", createdAt: now, updatedAt: now}

		if snippet1.Equal(snippet2) {
			t.Error("expected snippets with different IDs to not be equal")
		}
	})

	t.Run("different titles are not equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{id: 1, title: "title1", language: "go", code: "code", createdAt: now, updatedAt: now}
		snippet2 := &Snippet{id: 1, title: "title2", language: "go", code: "code", createdAt: now, updatedAt: now}

		if snippet1.Equal(snippet2) {
			t.Error("expected snippets with different titles to not be equal")
		}
	})

	t.Run("different languages are not equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{id: 1, title: "title", language: "go", code: "code", createdAt: now, updatedAt: now}
		snippet2 := &Snippet{id: 1, title: "title", language: "python", code: "code", createdAt: now, updatedAt: now}

		if snippet1.Equal(snippet2) {
			t.Error("expected snippets with different languages to not be equal")
		}
	})

	t.Run("different code are not equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{id: 1, title: "title", language: "go", code: "code1", createdAt: now, updatedAt: now}
		snippet2 := &Snippet{id: 1, title: "title", language: "go", code: "code2", createdAt: now, updatedAt: now}

		if snippet1.Equal(snippet2) {
			t.Error("expected snippets with different code to not be equal")
		}
	})

	t.Run("different descriptions are not equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{id: 1, title: "title", language: "go", code: "code", description: "desc1", createdAt: now, updatedAt: now}
		snippet2 := &Snippet{id: 1, title: "title", language: "go", code: "code", description: "desc2", createdAt: now, updatedAt: now}

		if snippet1.Equal(snippet2) {
			t.Error("expected snippets with different descriptions to not be equal")
		}
	})

	t.Run("different categoryIDs are not equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{id: 1, title: "title", language: "go", code: "code", categoryID: 1, createdAt: now, updatedAt: now}
		snippet2 := &Snippet{id: 1, title: "title", language: "go", code: "code", categoryID: 2, createdAt: now, updatedAt: now}

		if snippet1.Equal(snippet2) {
			t.Error("expected snippets with different categoryIDs to not be equal")
		}
	})

	t.Run("different tags are not equal", func(t *testing.T) {
		now := time.Now()
		snippet1 := &Snippet{id: 1, title: "title", language: "go", code: "code", tags: []int{1, 2}, createdAt: now, updatedAt: now}
		snippet2 := &Snippet{id: 1, title: "title", language: "go", code: "code", tags: []int{1, 3}, createdAt: now, updatedAt: now}

		if snippet1.Equal(snippet2) {
			t.Error("expected snippets with different tags to not be equal")
		}
	})

	t.Run("nil other is not equal", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "title", "go", "code")

		if snippet.Equal(nil) {
			t.Error("expected snippet to not be equal to nil")
		}
	})
}

func TestSnippet_MarshalJSON(t *testing.T) {
	t.Run("marshals to valid JSON", func(t *testing.T) {
		snippet := mustCreateSnippet(t, "quicksort", "go", "func quicksort() {}")
		snippet.SetID(42)
		snippet.SetDescription("Fast sorting algorithm")
		snippet.SetCategory(5)
		snippet.AddTag(1)
		snippet.AddTag(2)

		data, err := json.Marshal(snippet)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("generated invalid JSON: %v", err)
		}

		if id, ok := result["id"].(float64); !ok || int(id) != 42 {
			t.Errorf("expected id=42, got %v", result["id"])
		}
		if title, ok := result["title"].(string); !ok || title != "quicksort" {
			t.Errorf("expected title='quicksort', got %v", result["title"])
		}
		if language, ok := result["language"].(string); !ok || language != "go" {
			t.Errorf("expected language='go', got %v", result["language"])
		}
		if code, ok := result["code"].(string); !ok || code != "func quicksort() {}" {
			t.Errorf("expected code='func quicksort() {}', got %v", result["code"])
		}
		if desc, ok := result["description"].(string); !ok || desc != "Fast sorting algorithm" {
			t.Errorf("expected description='Fast sorting algorithm', got %v", result["description"])
		}
		if categoryID, ok := result["category_id"].(float64); !ok || int(categoryID) != 5 {
			t.Errorf("expected category_id=5, got %v", result["category_id"])
		}
		if tags, ok := result["tags"].([]any); !ok || len(tags) != 2 {
			t.Errorf("expected 2 tags, got %v", result["tags"])
		}
		if _, ok := result["created_at"]; !ok {
			t.Error("missing created_at field")
		}
		if _, ok := result["updated_at"]; !ok {
			t.Error("missing updated_at field")
		}
	})
}

func TestSnippet_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals valid JSON", func(t *testing.T) {
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
			t.Errorf("expected ID=42, got %d", snippet.ID())
		}
		if snippet.Title() != "quicksort" {
			t.Errorf("expected title='quicksort', got %q", snippet.Title())
		}
		if snippet.Language() != "go" {
			t.Errorf("expected language='go', got %q", snippet.Language())
		}
		if snippet.Code() != "func quicksort() {}" {
			t.Errorf("expected code='func quicksort() {}', got %q", snippet.Code())
		}
		if snippet.Description() != "Fast sorting algorithm" {
			t.Errorf("expected description='Fast sorting algorithm', got %q", snippet.Description())
		}
		if snippet.CategoryID() != 5 {
			t.Errorf("expected category_id=5, got %d", snippet.CategoryID())
		}
		if len(snippet.Tags()) != 3 {
			t.Errorf("expected 3 tags, got %d", len(snippet.Tags()))
		}
		if snippet.CreatedAt().IsZero() {
			t.Error("CreatedAt should not be zero")
		}
		if snippet.UpdatedAt().IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
	})

	t.Run("empty title returns error", func(t *testing.T) {
		jsonData := []byte(`{
			"id": 1,
			"title": "",
			"language": "go",
			"code": "code",
			"created_at": "2024-01-15T10:30:00Z",
			"updated_at": "2024-01-15T10:30:00Z"
		}`)

		var snippet Snippet
		err := json.Unmarshal(jsonData, &snippet)

		if err == nil {
			t.Error("expected error for empty title")
		}

		if !errors.Is(err, ErrEmptyTitle) {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}
	})

	t.Run("empty language returns error", func(t *testing.T) {
		jsonData := []byte(`{
			"id": 1,
			"title": "title",
			"language": "",
			"code": "code",
			"created_at": "2024-01-15T10:30:00Z",
			"updated_at": "2024-01-15T10:30:00Z"
		}`)

		var snippet Snippet
		err := json.Unmarshal(jsonData, &snippet)

		if err == nil {
			t.Error("expected error for empty language")
		}

		if !errors.Is(err, ErrEmptyLanguage) {
			t.Errorf("expected ErrEmptyLanguage, got %v", err)
		}
	})

	t.Run("empty code returns error", func(t *testing.T) {
		jsonData := []byte(`{
			"id": 1,
			"title": "title",
			"language": "go",
			"code": "",
			"created_at": "2024-01-15T10:30:00Z",
			"updated_at": "2024-01-15T10:30:00Z"
		}`)

		var snippet Snippet
		err := json.Unmarshal(jsonData, &snippet)

		if err == nil {
			t.Error("expected error for empty code")
		}

		if !errors.Is(err, ErrEmptyCode) {
			t.Errorf("expected ErrEmptyCode, got %v", err)
		}
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		invalidJSON := []byte(`{"id": "not a number"}`)

		var snippet Snippet
		err := json.Unmarshal(invalidJSON, &snippet)

		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("malformed JSON returns error", func(t *testing.T) {
		malformedJSON := []byte(`{invalid}`)

		var snippet Snippet
		err := json.Unmarshal(malformedJSON, &snippet)

		if err == nil {
			t.Error("expected error for malformed JSON")
		}
	})

	t.Run("null tags become empty slice", func(t *testing.T) {
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

		if snippet.Tags() == nil {
			t.Error("Tags should not be nil after unmarshal")
		}

		if len(snippet.Tags()) != 0 {
			t.Errorf("expected empty tags, got %v", snippet.Tags())
		}
	})
}

func TestSnippet_JSONRoundTrip(t *testing.T) {
	t.Run("round trip preserves all fields", func(t *testing.T) {
		original := mustCreateSnippet(t, "binary search", "python", "def binary_search():")
		original.SetID(123)
		original.SetDescription("Efficient search")
		original.SetCategory(10)
		original.AddTag(5)
		original.AddTag(6)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var restored Snippet
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !restored.Equal(original) {
			t.Error("round trip failed: snippets are not equal")
			t.Logf("original:  %+v", original)
			t.Logf("restored:  %+v", restored)
		}
	})
}

// mustCreateSnippet creates a snippet or fails the test.
func mustCreateSnippet(t *testing.T, title, language, code string) *Snippet {
	t.Helper()
	snippet, err := NewSnippet(title, language, code)
	if err != nil {
		t.Fatalf("failed to create snippet: %v", err)
	}
	return snippet
}
