package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewCategory(t *testing.T) {
	t.Run("empty name returns error", func(t *testing.T) {
		category, err := NewCategory("")
		if err == nil {
			t.Fatal("expected error for empty name, got nil")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}

		if category != nil {
			t.Error("expected category to be nil")
		}
	})

	t.Run("valid category", func(t *testing.T) {
		category, err := NewCategory("cat")
		if err != nil {
			t.Fatal("expected err to be nil")
		}

		if category == nil {
			t.Fatal("expected category not be nil")
		}

		if category.Name() != "cat" {
			t.Errorf("name to be 'cat', got %v", category.Name())
		}

		if category.ID() != 0 {
			t.Errorf("expected ID 0, got %d", category.ID())
		}

		if category.CreatedAt().IsZero() {
			t.Error("expected CreatedAt to be set")
		}

		if category.UpdatedAt().IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})
}

func TestCategoryMutations(t *testing.T) {
	t.Run("SetID updates ID", func(t *testing.T) {
		category, err := NewCategory("cat")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		category.SetID(42)

		if category.ID() != 42 {
			t.Errorf("expected ID 42, got %d", category.ID())
		}
	})

	t.Run("SetName with empty string returns error", func(t *testing.T) {
		category, _ := NewCategory("algorithms")

		err := category.SetName("")

		if err == nil {
			t.Error("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("SetName updates name", func(t *testing.T) {
		category, _ := NewCategory("algorithms")

		err := category.SetName("algo")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if category.Name() != "algo" {
			t.Errorf("expected name 'algo', got %s", category.Name())
		}
	})

	t.Run("mutations update timestamp", func(t *testing.T) {
		category, err := NewCategory("cat")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		originalTime := category.UpdatedAt()

		// Sleep to ensure time difference
		time.Sleep(2 * time.Millisecond)

		category.SetName("new name")

		if !category.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to change after mutation")
		}
	})
}
