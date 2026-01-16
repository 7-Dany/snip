package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewTag(t *testing.T) {
	t.Run("empty name should return error", func(t *testing.T) {
		tag, err := NewTag("")
		if err == nil {
			t.Fatal("expected error for empty name, got nil")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}

		if tag != nil {
			t.Error("expected tag to be nil")
		}
	})

	t.Run("valid tag", func(t *testing.T) {
		tag, err := NewTag("cat")
		if err != nil {
			t.Fatal("expected err to be nil")
		}

		if tag == nil {
			t.Fatal("expected tag not be nil")
		}

		if tag.Name() != "cat" {
			t.Errorf("name to be 'cat', got %v", tag.Name())
		}

		if tag.ID() != 0 {
			t.Errorf("expected ID 0, got %d", tag.ID())
		}

		if tag.CreatedAt().IsZero() {
			t.Error("expected CreatedAt to be set")
		}

		if tag.UpdatedAt().IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})
}

func TestTagMutations(t *testing.T) {
	t.Run("SetID updates ID", func(t *testing.T) {
		tag, err := NewTag("tag")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		tag.SetID(42)

		if tag.ID() != 42 {
			t.Errorf("expected ID 42, got %d", tag.ID())
		}
	})

	t.Run("SetName with empty string returns error", func(t *testing.T) {
		tag, _ := NewTag("red")

		err := tag.SetName("")

		if err == nil {
			t.Error("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("SetName updates name", func(t *testing.T) {
		tag, _ := NewTag("red")

		err := tag.SetName("blue")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if tag.Name() != "blue" {
			t.Errorf("expected name 'blue', got %s", tag.Name())
		}
	})

	t.Run("mutations update timestamp", func(t *testing.T) {
		tag, err := NewTag("red")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		originalTime := tag.UpdatedAt()

		// Sleep to ensure time difference
		time.Sleep(2 * time.Millisecond)

		tag.SetName("new name")

		if !tag.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to change after mutation")
		}
	})
}
