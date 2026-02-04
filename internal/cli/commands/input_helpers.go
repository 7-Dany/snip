// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// simpleInputModel is a reusable Bubble Tea model for single text input prompts.
type simpleInputModel struct {
	textInput textinput.Model
	cancelled bool
}

// newSimpleInputModel creates a new simple input model with the given configuration.
func newSimpleInputModel(placeholder string, charLimit, width int) simpleInputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = charLimit
	ti.Width = width
	return simpleInputModel{
		textInput: ti,
	}
}

// Init initializes the simple input model.
func (m simpleInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input events for the simple input model.
func (m simpleInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the simple input model (basic implementation).
func (m simpleInputModel) View() string {
	return m.textInput.View()
}

// promptForInput displays an interactive prompt and returns the trimmed user input.
// Returns empty string if cancelled or if input is empty after trimming.
func promptForInput(title, emoji, fieldName, placeholder string, charLimit, width int) string {
	model := newSimpleInputModel(placeholder, charLimit, width)
	p := tea.NewProgram(inputModelWrapper{model: model, title: title, emoji: emoji, fieldName: fieldName})
	finalModel, err := p.Run()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to run input prompt: %v", err))
		return ""
	}

	wrapper := finalModel.(inputModelWrapper)
	if wrapper.model.cancelled {
		return ""
	}

	value := strings.TrimSpace(wrapper.model.textInput.Value())
	if value == "" {
		PrintError(fmt.Sprintf("%s cannot be empty", fieldName))
		return ""
	}

	return value
}

// inputModelWrapper wraps simpleInputModel to provide view customization.
type inputModelWrapper struct {
	model     simpleInputModel
	title     string
	emoji     string
	fieldName string
}

// Init initializes the wrapper.
func (w inputModelWrapper) Init() tea.Cmd {
	return w.model.Init()
}

// Update delegates to the wrapped model.
func (w inputModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := w.model.Update(msg)
	w.model = updatedModel.(simpleInputModel)
	return w, cmd
}

// View renders the wrapped model with custom formatting.
func (w inputModelWrapper) View() string {
	return fmt.Sprintf(
		"\n%s %s\n\n"+
			"%s: %s\n\n"+
			"(press Enter to confirm, Esc to cancel)\n",
		w.emoji,
		w.title,
		w.fieldName,
		w.model.textInput.View(),
	)
}

// createTextInput is a helper to create and configure a textinput.Model.
// This eliminates repetitive textinput creation code.
func createTextInput(placeholder string, charLimit, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = charLimit
	ti.Width = width
	return ti
}

// createFocusedTextInput creates a text input that is focused by default.
func createFocusedTextInput(placeholder string, charLimit, width int) textinput.Model {
	ti := createTextInput(placeholder, charLimit, width)
	ti.Focus()
	return ti
}
