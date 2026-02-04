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
		category, err := NewCategory("algorithms")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if category == nil {
			t.Fatal("expected category to be non-nil")
		}

		if category.Name() != "algorithms" {
			t.Errorf("expected name 'algorithms', got %q", category.Name())
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

		if !category.CreatedAt().Equal(category.UpdatedAt()) {
			t.Error("expected CreatedAt and UpdatedAt to be equal for new category")
		}
	})
}

func TestCategory_SetID(t *testing.T) {
	t.Run("sets positive ID", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")

		category.SetID(42)

		if category.ID() != 42 {
			t.Errorf("expected ID 42, got %d", category.ID())
		}
	})

	t.Run("sets zero ID", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")

		category.SetID(0)

		if category.ID() != 0 {
			t.Errorf("expected ID 0, got %d", category.ID())
		}
	})

	t.Run("sets negative ID", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")

		category.SetID(-1)

		if category.ID() != -1 {
			t.Errorf("expected ID -1, got %d", category.ID())
		}
	})
}

func TestCategory_SetName(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")
		originalName := category.Name()
		originalUpdatedAt := category.UpdatedAt()

		err := category.SetName("")

		if err == nil {
			t.Fatal("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}

		if category.Name() != originalName {
			t.Error("name should not change when SetName returns error")
		}

		if !category.UpdatedAt().Equal(originalUpdatedAt) {
			t.Error("UpdatedAt should not change when SetName returns error")
		}
	})

	t.Run("updates name", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")

		err := category.SetName("data structures")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if category.Name() != "data structures" {
			t.Errorf("expected name 'data structures', got %q", category.Name())
		}
	})

	t.Run("updates timestamp", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")
		originalTime := category.UpdatedAt()

		time.Sleep(2 * time.Millisecond)

		err := category.SetName("new name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !category.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to be updated after SetName")
		}
	})

	t.Run("same name is valid", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")

		err := category.SetName("algorithms")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if category.Name() != "algorithms" {
			t.Errorf("expected name 'algorithms', got %q", category.Name())
		}
	})
}

func TestCategory_String(t *testing.T) {
	t.Run("with ID and name", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")
		category.SetID(42)

		got := category.String()
		want := `Category{id=42, name="algorithms"}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("with zero ID", func(t *testing.T) {
		category := mustCreateCategory(t, "databases")

		got := category.String()
		want := `Category{id=0, name="databases"}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("with special characters in name", func(t *testing.T) {
		category := mustCreateCategory(t, `"quoted"`)
		category.SetID(1)

		got := category.String()
		want := `Category{id=1, name="\"quoted\""}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestCategory_Equal(t *testing.T) {
	t.Run("identical categories are equal", func(t *testing.T) {
		now := time.Now()
		cat1 := &Category{id: 1, name: "algorithms", createdAt: now, updatedAt: now}
		cat2 := &Category{id: 1, name: "algorithms", createdAt: now, updatedAt: now}

		if !cat1.Equal(cat2) {
			t.Error("expected identical categories to be equal")
		}
	})

	t.Run("same instance is equal to itself", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")

		if !category.Equal(category) {
			t.Error("expected category to be equal to itself")
		}
	})

	t.Run("different IDs are not equal", func(t *testing.T) {
		now := time.Now()
		cat1 := &Category{id: 1, name: "algorithms", createdAt: now, updatedAt: now}
		cat2 := &Category{id: 2, name: "algorithms", createdAt: now, updatedAt: now}

		if cat1.Equal(cat2) {
			t.Error("expected categories with different IDs to not be equal")
		}
	})

	t.Run("different names are not equal", func(t *testing.T) {
		now := time.Now()
		cat1 := &Category{id: 1, name: "algorithms", createdAt: now, updatedAt: now}
		cat2 := &Category{id: 1, name: "databases", createdAt: now, updatedAt: now}

		if cat1.Equal(cat2) {
			t.Error("expected categories with different names to not be equal")
		}
	})

	t.Run("different CreatedAt are not equal", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)
		cat1 := &Category{id: 1, name: "algorithms", createdAt: now, updatedAt: now}
		cat2 := &Category{id: 1, name: "algorithms", createdAt: later, updatedAt: now}

		if cat1.Equal(cat2) {
			t.Error("expected categories with different CreatedAt to not be equal")
		}
	})

	t.Run("different UpdatedAt are not equal", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)
		cat1 := &Category{id: 1, name: "algorithms", createdAt: now, updatedAt: now}
		cat2 := &Category{id: 1, name: "algorithms", createdAt: now, updatedAt: later}

		if cat1.Equal(cat2) {
			t.Error("expected categories with different UpdatedAt to not be equal")
		}
	})

	t.Run("nil other is not equal", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")

		if category.Equal(nil) {
			t.Error("expected category to not be equal to nil")
		}
	})
}

func TestCategory_MarshalJSON(t *testing.T) {
	t.Run("marshals to valid JSON", func(t *testing.T) {
		category := mustCreateCategory(t, "algorithms")
		category.SetID(42)

		data, err := json.Marshal(category)
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
		if name, ok := result["name"].(string); !ok || name != "algorithms" {
			t.Errorf("expected name='algorithms', got %v", result["name"])
		}
		if _, ok := result["created_at"]; !ok {
			t.Error("missing created_at field")
		}
		if _, ok := result["updated_at"]; !ok {
			t.Error("missing updated_at field")
		}
	})
}

func TestCategory_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals valid JSON", func(t *testing.T) {
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
			t.Errorf("expected ID=42, got %d", category.ID())
		}
		if category.Name() != "algorithms" {
			t.Errorf("expected name='algorithms', got %q", category.Name())
		}
		if category.CreatedAt().IsZero() {
			t.Error("CreatedAt should not be zero")
		}
		if category.UpdatedAt().IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
	})

	t.Run("empty name returns error", func(t *testing.T) {
		jsonData := []byte(`{
			"id": 1,
			"name": "",
			"created_at": "2024-01-15T10:30:00Z",
			"updated_at": "2024-01-15T10:30:00Z"
		}`)

		var category Category
		err := json.Unmarshal(jsonData, &category)

		if err == nil {
			t.Error("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		invalidJSON := []byte(`{"id": "not a number"}`)

		var category Category
		err := json.Unmarshal(invalidJSON, &category)

		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("malformed JSON returns error", func(t *testing.T) {
		malformedJSON := []byte(`{invalid}`)

		var category Category
		err := json.Unmarshal(malformedJSON, &category)

		if err == nil {
			t.Error("expected error for malformed JSON")
		}
	})
}

func TestCategory_JSONRoundTrip(t *testing.T) {
	t.Run("round trip preserves all fields", func(t *testing.T) {
		original := mustCreateCategory(t, "databases")
		original.SetID(123)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var restored Category
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !restored.Equal(original) {
			t.Error("round trip failed: categories are not equal")
			t.Logf("original:  ID=%d Name=%q CreatedAt=%v UpdatedAt=%v",
				original.ID(), original.Name(), original.CreatedAt(), original.UpdatedAt())
			t.Logf("restored:  ID=%d Name=%q CreatedAt=%v UpdatedAt=%v",
				restored.ID(), restored.Name(), restored.CreatedAt(), restored.UpdatedAt())
		}
	})
}

// mustCreateCategory creates a category or fails the test.
func mustCreateCategory(t *testing.T, name string) *Category {
	t.Helper()
	category, err := NewCategory(name)
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}
	return category
}
