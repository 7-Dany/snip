package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormView struct {
	title      string
	subtitle   string
	input      textinput.Model
	actionHint string
}

func NewFormView(title, subtitle, placeholder string) FormView {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 40

	return FormView{
		title:      title,
		subtitle:   subtitle,
		input:      ti,
		actionHint: "Enter: save | Esc: cancel",
	}
}

func (f *FormView) SetValue(value string) {
	f.input.SetValue(value)
}

func (f FormView) Value() string {
	return f.input.Value()
}

func (f *FormView) Update(msg tea.Msg) (FormView, tea.Cmd) {
	var cmd tea.Cmd
	f.input, cmd = f.input.Update(msg)
	return *f, cmd
}

func (f FormView) Focus() tea.Cmd {
	return f.input.Focus()
}

func (f FormView) View() string {
	var b string

	if f.title != "" {
		b += f.title + "\n\n"
	}

	if f.subtitle != "" {
		b += f.subtitle + "\n\n"
	}

	b += "Name: " + f.input.View() + "\n\n"

	b += lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(f.actionHint)

	return b
}
