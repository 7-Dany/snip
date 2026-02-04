package domain

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestNewTag(t *testing.T) {
	t.Run("empty name returns error", func(t *testing.T) {
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
		tag, err := NewTag("performance")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tag == nil {
			t.Fatal("expected tag to be non-nil")
		}

		if tag.Name() != "performance" {
			t.Errorf("expected name 'performance', got %q", tag.Name())
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

		if !tag.CreatedAt().Equal(tag.UpdatedAt()) {
			t.Error("expected CreatedAt and UpdatedAt to be equal for new tag")
		}
	})
}

func TestTag_SetID(t *testing.T) {
	t.Run("sets positive ID", func(t *testing.T) {
		tag := mustCreateTag(t, "sorting")

		tag.SetID(42)

		if tag.ID() != 42 {
			t.Errorf("expected ID 42, got %d", tag.ID())
		}
	})

	t.Run("sets zero ID", func(t *testing.T) {
		tag := mustCreateTag(t, "sorting")

		tag.SetID(0)

		if tag.ID() != 0 {
			t.Errorf("expected ID 0, got %d", tag.ID())
		}
	})

	t.Run("sets negative ID", func(t *testing.T) {
		tag := mustCreateTag(t, "sorting")

		tag.SetID(-1)

		if tag.ID() != -1 {
			t.Errorf("expected ID -1, got %d", tag.ID())
		}
	})
}

func TestTag_SetName(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		tag := mustCreateTag(t, "performance")
		originalName := tag.Name()
		originalUpdatedAt := tag.UpdatedAt()

		err := tag.SetName("")

		if err == nil {
			t.Fatal("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}

		if tag.Name() != originalName {
			t.Error("name should not change when SetName returns error")
		}

		if !tag.UpdatedAt().Equal(originalUpdatedAt) {
			t.Error("UpdatedAt should not change when SetName returns error")
		}
	})

	t.Run("updates name", func(t *testing.T) {
		tag := mustCreateTag(t, "performance")

		err := tag.SetName("optimization")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tag.Name() != "optimization" {
			t.Errorf("expected name 'optimization', got %q", tag.Name())
		}
	})

	t.Run("updates timestamp", func(t *testing.T) {
		tag := mustCreateTag(t, "performance")
		originalTime := tag.UpdatedAt()

		time.Sleep(2 * time.Millisecond)

		err := tag.SetName("new name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !tag.UpdatedAt().After(originalTime) {
			t.Error("expected UpdatedAt to be updated after SetName")
		}
	})

	t.Run("same name is valid", func(t *testing.T) {
		tag := mustCreateTag(t, "performance")

		err := tag.SetName("performance")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tag.Name() != "performance" {
			t.Errorf("expected name 'performance', got %q", tag.Name())
		}
	})
}

func TestTag_String(t *testing.T) {
	t.Run("with ID and name", func(t *testing.T) {
		tag := mustCreateTag(t, "sorting")
		tag.SetID(7)

		got := tag.String()
		want := `Tag{id=7, name="sorting"}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("with zero ID", func(t *testing.T) {
		tag := mustCreateTag(t, "performance")

		got := tag.String()
		want := `Tag{id=0, name="performance"}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("with special characters in name", func(t *testing.T) {
		tag := mustCreateTag(t, `"important"`)
		tag.SetID(1)

		got := tag.String()
		want := `Tag{id=1, name="\"important\""}`

		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestTag_Equal(t *testing.T) {
	t.Run("identical tags are equal", func(t *testing.T) {
		now := time.Now()
		tag1 := &Tag{id: 1, name: "sorting", createdAt: now, updatedAt: now}
		tag2 := &Tag{id: 1, name: "sorting", createdAt: now, updatedAt: now}

		if !tag1.Equal(tag2) {
			t.Error("expected identical tags to be equal")
		}
	})

	t.Run("same instance is equal to itself", func(t *testing.T) {
		tag := mustCreateTag(t, "sorting")

		if !tag.Equal(tag) {
			t.Error("expected tag to be equal to itself")
		}
	})

	t.Run("different IDs are not equal", func(t *testing.T) {
		now := time.Now()
		tag1 := &Tag{id: 1, name: "sorting", createdAt: now, updatedAt: now}
		tag2 := &Tag{id: 2, name: "sorting", createdAt: now, updatedAt: now}

		if tag1.Equal(tag2) {
			t.Error("expected tags with different IDs to not be equal")
		}
	})

	t.Run("different names are not equal", func(t *testing.T) {
		now := time.Now()
		tag1 := &Tag{id: 1, name: "sorting", createdAt: now, updatedAt: now}
		tag2 := &Tag{id: 1, name: "searching", createdAt: now, updatedAt: now}

		if tag1.Equal(tag2) {
			t.Error("expected tags with different names to not be equal")
		}
	})

	t.Run("different CreatedAt are not equal", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)
		tag1 := &Tag{id: 1, name: "sorting", createdAt: now, updatedAt: now}
		tag2 := &Tag{id: 1, name: "sorting", createdAt: later, updatedAt: now}

		if tag1.Equal(tag2) {
			t.Error("expected tags with different CreatedAt to not be equal")
		}
	})

	t.Run("different UpdatedAt are not equal", func(t *testing.T) {
		now := time.Now()
		later := now.Add(1 * time.Hour)
		tag1 := &Tag{id: 1, name: "sorting", createdAt: now, updatedAt: now}
		tag2 := &Tag{id: 1, name: "sorting", createdAt: now, updatedAt: later}

		if tag1.Equal(tag2) {
			t.Error("expected tags with different UpdatedAt to not be equal")
		}
	})

	t.Run("nil other is not equal", func(t *testing.T) {
		tag := mustCreateTag(t, "sorting")

		if tag.Equal(nil) {
			t.Error("expected tag to not be equal to nil")
		}
	})
}

func TestTag_MarshalJSON(t *testing.T) {
	t.Run("marshals to valid JSON", func(t *testing.T) {
		tag := mustCreateTag(t, "sorting")
		tag.SetID(7)

		data, err := json.Marshal(tag)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("generated invalid JSON: %v", err)
		}

		if id, ok := result["id"].(float64); !ok || int(id) != 7 {
			t.Errorf("expected id=7, got %v", result["id"])
		}
		if name, ok := result["name"].(string); !ok || name != "sorting" {
			t.Errorf("expected name='sorting', got %v", result["name"])
		}
		if _, ok := result["created_at"]; !ok {
			t.Error("missing created_at field")
		}
		if _, ok := result["updated_at"]; !ok {
			t.Error("missing updated_at field")
		}
	})
}

func TestTag_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals valid JSON", func(t *testing.T) {
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
			t.Errorf("expected ID=7, got %d", tag.ID())
		}
		if tag.Name() != "sorting" {
			t.Errorf("expected name='sorting', got %q", tag.Name())
		}
		if tag.CreatedAt().IsZero() {
			t.Error("CreatedAt should not be zero")
		}
		if tag.UpdatedAt().IsZero() {
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

		var tag Tag
		err := json.Unmarshal(jsonData, &tag)

		if err == nil {
			t.Error("expected error for empty name")
		}

		if !errors.Is(err, ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		invalidJSON := []byte(`{"id": "not a number"}`)

		var tag Tag
		err := json.Unmarshal(invalidJSON, &tag)

		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("malformed JSON returns error", func(t *testing.T) {
		malformedJSON := []byte(`{invalid}`)

		var tag Tag
		err := json.Unmarshal(malformedJSON, &tag)

		if err == nil {
			t.Error("expected error for malformed JSON")
		}
	})
}

func TestTag_JSONRoundTrip(t *testing.T) {
	t.Run("round trip preserves all fields", func(t *testing.T) {
		original := mustCreateTag(t, "performance")
		original.SetID(99)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var restored Tag
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !restored.Equal(original) {
			t.Error("round trip failed: tags are not equal")
			t.Logf("original:  ID=%d Name=%q CreatedAt=%v UpdatedAt=%v",
				original.ID(), original.Name(), original.CreatedAt(), original.UpdatedAt())
			t.Logf("restored:  ID=%d Name=%q CreatedAt=%v UpdatedAt=%v",
				restored.ID(), restored.Name(), restored.CreatedAt(), restored.UpdatedAt())
		}
	})
}

// mustCreateTag creates a tag or fails the test.
func mustCreateTag(t *testing.T, name string) *Tag {
	t.Helper()
	tag, err := NewTag(name)
	if err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}
	return tag
}
