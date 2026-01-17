package domain

import (
	"encoding/json"
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

func TestCategorySetID(t *testing.T) {
	t.Run("updates ID", func(t *testing.T) {
		category, err := NewCategory("cat")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		category.SetID(42)

		if category.ID() != 42 {
			t.Errorf("expected ID 42, got %d", category.ID())
		}
	})
}

func TestCategorySetName(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		category, _ := NewCategory("algorithms")

		err := category.SetName("")

		if err == nil {
			t.Error("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("updates name", func(t *testing.T) {
		category, _ := NewCategory("algorithms")

		err := category.SetName("algo")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if category.Name() != "algo" {
			t.Errorf("expected name 'algo', got %s", category.Name())
		}
	})

	t.Run("updates timestamp", func(t *testing.T) {
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

func TestCategoryJSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		category, _ := NewCategory("algorithms")
		category.SetID(42)

		data, err := json.Marshal(category)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		// Verify it's valid JSON
		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Generated invalid JSON: %v", err)
		}

		// Verify fields are present
		if id, ok := result["id"].(float64); !ok || int(id) != 42 {
			t.Errorf("Expected id=42, got %v", result["id"])
		}
		if name, ok := result["name"].(string); !ok || name != "algorithms" {
			t.Errorf("Expected name='algorithms', got %v", result["name"])
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
		"name": "algorithms",
		"created_at": "2024-01-15T10:30:00Z",
		"updated_at": "2024-01-16T14:20:00Z"
	}`)

		var category Category
		err := json.Unmarshal(jsonData, &category)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if category.ID() != 42 {
			t.Errorf("Expected ID=42, got %d", category.ID())
		}
		if category.Name() != "algorithms" {
			t.Errorf("Expected name='algorithms', got '%s'", category.Name())
		}
		if category.CreatedAt().IsZero() {
			t.Error("CreatedAt should not be zero")
		}
		if category.UpdatedAt().IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
	})

	t.Run("round trip", func(t *testing.T) {
		original, _ := NewCategory("databases")
		original.SetID(123)

		// Marshal
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// Unmarshal
		var restored Category
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// Verify all fields match
		if restored.ID() != original.ID() {
			t.Errorf("ID mismatch: expected %d, got %d", original.ID(), restored.ID())
		}
		if restored.Name() != original.Name() {
			t.Errorf("Name mismatch: expected %s, got %s", original.Name(), restored.Name())
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

		var category Category
		err := json.Unmarshal(invalidJSON, &category)

		if err == nil {
			t.Error("UnmarshalJSON should return error for invalid JSON")
		}
	})
}
