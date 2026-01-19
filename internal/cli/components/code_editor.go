package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FocusChangedMsg is sent when focus changes to request viewport scroll.
type FocusChangedMsg struct {
	FocusedField int
	FieldLine    int // Approximate line number of the field
}

// CodeEditor provides a multi-input form for creating/editing code snippets.
type CodeEditor struct {
	titleInput       textinput.Model
	languageInput    textinput.Model
	descriptionInput textinput.Model
	codeArea         textarea.Model
	focusedField     int // 0=title, 1=language, 2=description, 3=code, 4=cancel button, 5=save button
	width            int
	height           int

	// Category/Tag management
	categoryID   int
	categoryName string
	tagIDs       []int
	tagNames     []string
}

// NewCodeEditor creates a new code editor form.
func NewCodeEditor(width, height int) CodeEditor {
	// Calculate input width based on available space
	inputWidth := width - 8
	if inputWidth < 40 {
		inputWidth = 40
	}

	// Title input
	ti := textinput.New()
	ti.Placeholder = "Snippet title (e.g., 'Binary Search Algorithm')"
	ti.CharLimit = 100
	ti.Width = inputWidth
	ti.Focus()

	// Language input
	li := textinput.New()
	li.Placeholder = "Language (e.g., 'go', 'python', 'javascript')"
	li.CharLimit = 30
	li.Width = inputWidth

	// Description input
	di := textinput.New()
	di.Placeholder = "Description (optional)"
	di.CharLimit = 200
	di.Width = inputWidth

	// Code textarea - make it much larger to show more code
	ta := textarea.New()
	ta.Placeholder = "Paste or type your code here..."
	ta.SetWidth(inputWidth)
	// Increase height significantly - subtract less for more code space
	codeHeight := height - 12
	if codeHeight < 10 {
		codeHeight = 10
	}
	ta.SetHeight(codeHeight)
	ta.ShowLineNumbers = true
	ta.CharLimit = 50000

	return CodeEditor{
		titleInput:       ti,
		languageInput:    li,
		descriptionInput: di,
		codeArea:         ta,
		focusedField:     0,
		width:            width,
		height:           height,
		categoryID:       0,
		categoryName:     "None",
		tagIDs:           []int{},
		tagNames:         []string{},
	}
}

// SetValues populates the editor with existing snippet data.
func (ce *CodeEditor) SetValues(title, language, description, code string) {
	ce.titleInput.SetValue(title)
	ce.languageInput.SetValue(language)
	ce.descriptionInput.SetValue(description)
	ce.codeArea.SetValue(code)
}

// SetCategory sets the category for the snippet.
func (ce *CodeEditor) SetCategory(id int, name string) {
	ce.categoryID = id
	ce.categoryName = name
}

// SetTags sets the tags for the snippet.
func (ce *CodeEditor) SetTags(ids []int, names []string) {
	ce.tagIDs = ids
	ce.tagNames = names
}

// GetCategory returns the category ID.
func (ce CodeEditor) GetCategory() int {
	return ce.categoryID
}

// GetTags returns the tag IDs.
func (ce CodeEditor) GetTags() []int {
	return ce.tagIDs
}

// GetValues returns all field values.
func (ce CodeEditor) GetValues() (title, language, description, code string) {
	return strings.TrimSpace(ce.titleInput.Value()),
		strings.TrimSpace(ce.languageInput.Value()),
		strings.TrimSpace(ce.descriptionInput.Value()),
		strings.TrimSpace(ce.codeArea.Value())
}

// GetFocusedField returns the currently focused field index.
func (ce CodeEditor) GetFocusedField() int {
	return ce.focusedField
}

// IsOnTextInput returns true if focus is on a text input field (not buttons).
func (ce CodeEditor) IsOnTextInput() bool {
	return ce.focusedField < 4
}

// NextField moves focus to the next input field.
func (ce *CodeEditor) NextField() tea.Cmd {
	ce.focusedField = (ce.focusedField + 1) % 6
	ce.updateFocus()
	return func() tea.Msg {
		return FocusChangedMsg{
			FocusedField: ce.focusedField,
			FieldLine:    ce.getFieldLine(ce.focusedField),
		}
	}
}

// PrevField moves focus to the previous input field.
func (ce *CodeEditor) PrevField() tea.Cmd {
	ce.focusedField = (ce.focusedField - 1 + 6) % 6
	ce.updateFocus()
	return func() tea.Msg {
		return FocusChangedMsg{
			FocusedField: ce.focusedField,
			FieldLine:    ce.getFieldLine(ce.focusedField),
		}
	}
}

// getFieldLine returns the approximate line number for a field.
func (ce CodeEditor) getFieldLine(field int) int {
	switch field {
	case 0: // Title
		return 2
	case 1: // Language
		return 6
	case 2: // Description
		return 10
	case 3: // Code (start)
		return 18
	case 4, 5: // Buttons at bottom
		// Calculate button position based on code height
		return 18 + ce.codeArea.Height() + 4
	default:
		return 0
	}
}

// IsCancelFocused returns true if Cancel button is focused.
func (ce CodeEditor) IsCancelFocused() bool {
	return ce.focusedField == 4
}

// IsSaveFocused returns true if Save button is focused.
func (ce CodeEditor) IsSaveFocused() bool {
	return ce.focusedField == 5
}

// updateFocus sets focus to the current field.
func (ce *CodeEditor) updateFocus() {
	ce.titleInput.Blur()
	ce.languageInput.Blur()
	ce.descriptionInput.Blur()
	ce.codeArea.Blur()

	switch ce.focusedField {
	case 0:
		ce.titleInput.Focus()
	case 1:
		ce.languageInput.Focus()
	case 2:
		ce.descriptionInput.Focus()
	case 3:
		ce.codeArea.Focus()
	}
}

// Update handles input for the code editor.
func (ce *CodeEditor) Update(msg tea.Msg) (CodeEditor, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Always allow tab navigation
		switch msg.String() {
		case "tab":
			cmd = ce.NextField()
			return *ce, cmd
		case "shift+tab":
			cmd = ce.PrevField()
			return *ce, cmd
		}

		// If on buttons, handle button-specific keys
		if ce.focusedField >= 4 {
			switch msg.String() {
			case "enter", "esc", "ctrl+s", "alt+s", "alt+c", "alt+t":
				// Let parent handle these when on buttons
				return *ce, nil
			}
		}

		// If on text input fields, only let parent handle specific control keys
		if ce.focusedField < 4 {
			switch msg.String() {
			case "ctrl+s", "alt+s", "esc", "alt+c", "alt+t":
				// These are always handled by parent
				return *ce, nil
			case "enter":
				// In title/language/description, enter moves to next field
				// In code area, enter creates new line
				if ce.focusedField < 3 {
					cmd = ce.NextField()
					return *ce, cmd
				}
			}
		}
	}

	// Update focused field
	if ce.focusedField < 4 {
		switch ce.focusedField {
		case 0:
			ce.titleInput, cmd = ce.titleInput.Update(msg)
			cmds = append(cmds, cmd)
		case 1:
			ce.languageInput, cmd = ce.languageInput.Update(msg)
			cmds = append(cmds, cmd)
		case 2:
			ce.descriptionInput, cmd = ce.descriptionInput.Update(msg)
			cmds = append(cmds, cmd)
		case 3:
			ce.codeArea, cmd = ce.codeArea.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return *ce, tea.Batch(cmds...)
}

// View renders the code editor.
func (ce CodeEditor) View() string {
	var b strings.Builder

	// Styles
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	focusedLabelStyle := labelStyle.Copy().
		Foreground(lipgloss.Color("13"))

	fieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	focusedFieldStyle := fieldStyle.Copy().
		BorderForeground(lipgloss.Color("13"))

	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Title field
	if ce.focusedField == 0 {
		b.WriteString(focusedLabelStyle.Render("▸ Title:"))
	} else {
		b.WriteString(labelStyle.Render("  Title:"))
	}
	b.WriteString("\n")
	if ce.focusedField == 0 {
		b.WriteString(focusedFieldStyle.Render(ce.titleInput.View()))
	} else {
		b.WriteString(fieldStyle.Render(ce.titleInput.View()))
	}
	b.WriteString("\n\n")

	// Language field
	if ce.focusedField == 1 {
		b.WriteString(focusedLabelStyle.Render("▸ Language:"))
	} else {
		b.WriteString(labelStyle.Render("  Language:"))
	}
	b.WriteString("\n")
	if ce.focusedField == 1 {
		b.WriteString(focusedFieldStyle.Render(ce.languageInput.View()))
	} else {
		b.WriteString(fieldStyle.Render(ce.languageInput.View()))
	}
	b.WriteString("\n\n")

	// Description field
	if ce.focusedField == 2 {
		b.WriteString(focusedLabelStyle.Render("▸ Description:"))
	} else {
		b.WriteString(labelStyle.Render("  Description:"))
	}
	b.WriteString("\n")
	if ce.focusedField == 2 {
		b.WriteString(focusedFieldStyle.Render(ce.descriptionInput.View()))
	} else {
		b.WriteString(fieldStyle.Render(ce.descriptionInput.View()))
	}
	b.WriteString("\n\n")

	// Category field (read-only display)
	b.WriteString(labelStyle.Render("  Category:"))
	b.WriteString(" ")
	b.WriteString(metaStyle.Render(ce.categoryName))
	b.WriteString(" ")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("(Alt+C to change)"))
	b.WriteString("\n\n")

	// Tags field (read-only display)
	b.WriteString(labelStyle.Render("  Tags:"))
	b.WriteString(" ")
	tagsDisplay := "None"
	if len(ce.tagNames) > 0 {
		tagsDisplay = strings.Join(ce.tagNames, ", ")
	}
	b.WriteString(metaStyle.Render(tagsDisplay))
	b.WriteString(" ")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("(Alt+T to manage)"))
	b.WriteString("\n\n")

	// Code area
	if ce.focusedField == 3 {
		b.WriteString(focusedLabelStyle.Render("▸ Code:"))
	} else {
		b.WriteString(labelStyle.Render("  Code:"))
	}
	b.WriteString("\n")

	codeBox := fieldStyle
	if ce.focusedField == 3 {
		codeBox = focusedFieldStyle
	}
	b.WriteString(codeBox.Render(ce.codeArea.View()))
	b.WriteString("\n\n")

	// Buttons
	cancelButton := lipgloss.NewStyle().
		Padding(0, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render("Cancel")

	saveButton := lipgloss.NewStyle().
		Padding(0, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render("Save")

	if ce.focusedField == 4 {
		cancelButton = lipgloss.NewStyle().
			Padding(0, 3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Background(lipgloss.Color("9")).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Render("Cancel")
	}

	if ce.focusedField == 5 {
		saveButton = lipgloss.NewStyle().
			Padding(0, 3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("10")).
			Background(lipgloss.Color("10")).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Render("Save")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, cancelButton, "  ", saveButton)
	b.WriteString(buttons)

	return b.String()
}
