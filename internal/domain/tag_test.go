package domain

import (
	"encoding/json"
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

func TestTagSetID(t *testing.T) {
	t.Run("updates ID", func(t *testing.T) {
		tag, err := NewTag("tag")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		tag.SetID(42)

		if tag.ID() != 42 {
			t.Errorf("expected ID 42, got %d", tag.ID())
		}
	})
}

func TestTagSetName(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		tag, _ := NewTag("red")

		err := tag.SetName("")

		if err == nil {
			t.Error("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("updates name", func(t *testing.T) {
		tag, _ := NewTag("red")

		err := tag.SetName("blue")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if tag.Name() != "blue" {
			t.Errorf("expected name 'blue', got %s", tag.Name())
		}
	})

	t.Run("updates timestamp", func(t *testing.T) {
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

func TestTagJSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		tag, _ := NewTag("sorting")
		tag.SetID(7)

		data, err := json.Marshal(tag)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		// Verify it's valid JSON
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Generated invalid JSON: %v", err)
		}

		// Verify fields are present
		if id, ok := result["id"].(float64); !ok || int(id) != 7 {
			t.Errorf("Expected id=7, got %v", result["id"])
		}
		if name, ok := result["name"].(string); !ok || name != "sorting" {
			t.Errorf("Expected name='sorting', got %v", result["name"])
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
		"id": 7,
		"name": "sorting",
		"created_at": "2024-01-15T10:30:00Z",
		"updated_at": "2024-01-16T14:20:00Z"
	}`)

		var tag Tag
		err := json.Unmarshal(jsonData, &tag)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if tag.ID() != 7 {
			t.Errorf("Expected ID=7, got %d", tag.ID())
		}
		if tag.Name() != "sorting" {
			t.Errorf("Expected name='sorting', got '%s'", tag.Name())
		}
		if tag.CreatedAt().IsZero() {
			t.Error("CreatedAt should not be zero")
		}
		if tag.UpdatedAt().IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
	})

	t.Run("round trip", func(t *testing.T) {
		original, _ := NewTag("performance")
		original.SetID(99)

		// Marshal
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// Unmarshal
		var restored Tag
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

		var tag Tag
		err := json.Unmarshal(invalidJSON, &tag)

		if err == nil {
			t.Error("UnmarshalJSON should return error for invalid JSON")
		}
	})
}
