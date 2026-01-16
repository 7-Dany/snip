package domain

import (
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

func TestSnippetMutations(t *testing.T) {
	t.Run("SetTitle with empty string returns error", func(t *testing.T) {
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

	t.Run("SetTitle updates title", func(t *testing.T) {
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

	t.Run("SetLanguage with empty string returns error", func(t *testing.T) {
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

	t.Run("SetLanguage updates language", func(t *testing.T) {
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

	t.Run("SetCode with empty string returns error", func(t *testing.T) {
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

	t.Run("SetCode updates code", func(t *testing.T) {
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

	t.Run("SetDescription updates description", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.SetDescription("description")

		if snippet.Description() != "description" {
			t.Errorf("expected description 'description', got %s", snippet.Description())
		}
	})

	t.Run("SetCategory updates category", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.SetCategory(1)

		if snippet.CategoryID() != 1 {
			t.Errorf("expected categoryID 1, got %d", snippet.CategoryID())
		}
	})

	t.Run("SetID updates ID", func(t *testing.T) {
		snippet, err := NewSnippet("title", "lang", "code")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		snippet.SetID(42)

		if snippet.ID() != 42 {
			t.Errorf("expected ID 42, got %d", snippet.ID())
		}
	})

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
