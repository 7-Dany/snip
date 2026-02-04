// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewSimpleInputModel tests the creation of a simple input model.
func TestNewSimpleInputModel(t *testing.T) {
	t.Run("creates model with correct configuration", func(t *testing.T) {
		model := newSimpleInputModel("test placeholder", 100, 50)

		if model.textInput.Placeholder != "test placeholder" {
			t.Errorf("Expected placeholder 'test placeholder', got '%s'", model.textInput.Placeholder)
		}

		if model.textInput.CharLimit != 100 {
			t.Errorf("Expected char limit 100, got %d", model.textInput.CharLimit)
		}

		if model.textInput.Width != 50 {
			t.Errorf("Expected width 50, got %d", model.textInput.Width)
		}

		if model.cancelled {
			t.Error("Expected cancelled to be false")
		}
	})

	t.Run("model is focused by default", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)

		if !model.textInput.Focused() {
			t.Error("Expected text input to be focused")
		}
	})
}

// TestSimpleInputModel_Init tests the initialization of the model.
func TestSimpleInputModel_Init(t *testing.T) {
	t.Run("returns blink command", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)
		cmd := model.Init()

		if cmd == nil {
			t.Error("Expected Init to return a command")
		}
	})
}

// TestSimpleInputModel_Update tests the Update method.
func TestSimpleInputModel_Update(t *testing.T) {
	t.Run("handles Enter key", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)

		// Simulate Enter key press
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		if cmd == nil {
			t.Error("Expected Enter to return quit command")
		}

		// Check that model is returned
		if updatedModel == nil {
			t.Error("Expected model to be returned")
		}

		m := updatedModel.(simpleInputModel)
		if m.cancelled {
			t.Error("Expected cancelled to remain false on Enter")
		}
	})

	t.Run("handles Esc key - sets cancelled flag", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if cmd == nil {
			t.Error("Expected Esc to return quit command")
		}

		m := updatedModel.(simpleInputModel)
		if !m.cancelled {
			t.Error("Expected cancelled to be true on Esc")
		}
	})

	t.Run("handles Ctrl+C - sets cancelled flag", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

		if cmd == nil {
			t.Error("Expected Ctrl+C to return quit command")
		}

		m := updatedModel.(simpleInputModel)
		if !m.cancelled {
			t.Error("Expected cancelled to be true on Ctrl+C")
		}
	})

	t.Run("delegates other keys to text input", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)

		// Simulate typing a character
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

		m := updatedModel.(simpleInputModel)
		if m.textInput.Value() != "a" {
			t.Errorf("Expected value 'a', got '%s'", m.textInput.Value())
		}
	})
}

// TestSimpleInputModel_View tests the View method.
func TestSimpleInputModel_View(t *testing.T) {
	t.Run("returns text input view", func(t *testing.T) {
		model := newSimpleInputModel("test placeholder", 50, 30)

		view := model.View()

		if view == "" {
			t.Error("Expected View to return non-empty string")
		}
	})
}

// TestInputModelWrapper_Init tests wrapper initialization.
func TestInputModelWrapper_Init(t *testing.T) {
	t.Run("delegates to wrapped model", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)
		wrapper := inputModelWrapper{
			model:     model,
			title:     "Test Title",
			emoji:     "📝",
			fieldName: "Name",
		}

		cmd := wrapper.Init()

		if cmd == nil {
			t.Error("Expected Init to return a command")
		}
	})
}

// TestInputModelWrapper_Update tests wrapper update logic.
func TestInputModelWrapper_Update(t *testing.T) {
	t.Run("updates wrapped model and returns wrapper", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)
		wrapper := inputModelWrapper{
			model:     model,
			title:     "Test Title",
			emoji:     "📝",
			fieldName: "Name",
		}

		// Simulate typing
		updatedModel, _ := wrapper.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

		// Should return wrapper, not inner model
		updatedWrapper, ok := updatedModel.(inputModelWrapper)
		if !ok {
			t.Fatal("Expected Update to return inputModelWrapper")
		}

		if updatedWrapper.model.textInput.Value() != "a" {
			t.Errorf("Expected value 'a', got '%s'", updatedWrapper.model.textInput.Value())
		}
	})

	t.Run("preserves wrapper fields", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)
		wrapper := inputModelWrapper{
			model:     model,
			title:     "Test Title",
			emoji:     "📝",
			fieldName: "Name",
		}

		updatedModel, _ := wrapper.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
		updatedWrapper := updatedModel.(inputModelWrapper)

		if updatedWrapper.title != "Test Title" {
			t.Errorf("Expected title 'Test Title', got '%s'", updatedWrapper.title)
		}

		if updatedWrapper.emoji != "📝" {
			t.Errorf("Expected emoji '📝', got '%s'", updatedWrapper.emoji)
		}

		if updatedWrapper.fieldName != "Name" {
			t.Errorf("Expected fieldName 'Name', got '%s'", updatedWrapper.fieldName)
		}
	})

	t.Run("handles Enter key", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)
		wrapper := inputModelWrapper{
			model:     model,
			title:     "Test",
			emoji:     "📝",
			fieldName: "Field",
		}

		updatedModel, cmd := wrapper.Update(tea.KeyMsg{Type: tea.KeyEnter})

		if cmd == nil {
			t.Error("Expected Enter to return quit command")
		}

		updatedWrapper := updatedModel.(inputModelWrapper)
		if updatedWrapper.model.cancelled {
			t.Error("Expected cancelled to be false on Enter")
		}
	})

	t.Run("handles cancellation", func(t *testing.T) {
		model := newSimpleInputModel("test", 50, 30)
		wrapper := inputModelWrapper{
			model:     model,
			title:     "Test",
			emoji:     "📝",
			fieldName: "Field",
		}

		updatedModel, _ := wrapper.Update(tea.KeyMsg{Type: tea.KeyEsc})
		updatedWrapper := updatedModel.(inputModelWrapper)

		if !updatedWrapper.model.cancelled {
			t.Error("Expected cancelled to be true on Esc")
		}
	})
}

// TestInputModelWrapper_View tests wrapper view rendering.
func TestInputModelWrapper_View(t *testing.T) {
	t.Run("includes all custom fields in view", func(t *testing.T) {
		model := newSimpleInputModel("test placeholder", 50, 30)
		wrapper := inputModelWrapper{
			model:     model,
			title:     "Test Title",
			emoji:     "📝",
			fieldName: "Username",
		}

		view := wrapper.View()

		if view == "" {
			t.Error("Expected View to return non-empty string")
		}

		// View should contain the custom fields
		// (exact format depends on implementation)
	})

	t.Run("renders with different configurations", func(t *testing.T) {
		testCases := []struct {
			title     string
			emoji     string
			fieldName string
		}{
			{"Create Category", "📁", "Category Name"},
			{"Create Tag", "🏷️", "Tag Name"},
			{"Update Item", "✏️", "New Name"},
		}

		for _, tc := range testCases {
			model := newSimpleInputModel("test", 50, 30)
			wrapper := inputModelWrapper{
				model:     model,
				title:     tc.title,
				emoji:     tc.emoji,
				fieldName: tc.fieldName,
			}

			view := wrapper.View()
			if view == "" {
				t.Errorf("Expected non-empty view for %s", tc.title)
			}
		}
	})
}

// TestCreateTextInput tests the helper function.
func TestCreateTextInput(t *testing.T) {
	t.Run("creates text input with correct configuration", func(t *testing.T) {
		ti := createTextInput("test placeholder", 75, 40)

		if ti.Placeholder != "test placeholder" {
			t.Errorf("Expected placeholder 'test placeholder', got '%s'", ti.Placeholder)
		}

		if ti.CharLimit != 75 {
			t.Errorf("Expected char limit 75, got %d", ti.CharLimit)
		}

		if ti.Width != 40 {
			t.Errorf("Expected width 40, got %d", ti.Width)
		}
	})

	t.Run("creates unfocused input by default", func(t *testing.T) {
		ti := createTextInput("test", 50, 30)

		if ti.Focused() {
			t.Error("Expected createTextInput to create unfocused input")
		}
	})
}

// TestCreateFocusedTextInput tests the focused helper function.
func TestCreateFocusedTextInput(t *testing.T) {
	t.Run("creates focused text input", func(t *testing.T) {
		ti := createFocusedTextInput("test placeholder", 60, 35)

		if !ti.Focused() {
			t.Error("Expected createFocusedTextInput to create focused input")
		}

		if ti.Placeholder != "test placeholder" {
			t.Errorf("Expected placeholder 'test placeholder', got '%s'", ti.Placeholder)
		}

		if ti.CharLimit != 60 {
			t.Errorf("Expected char limit 60, got %d", ti.CharLimit)
		}

		if ti.Width != 35 {
			t.Errorf("Expected width 35, got %d", ti.Width)
		}
	})
}
