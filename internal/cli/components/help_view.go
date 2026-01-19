package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type HelpSection struct {
	Title string
	Items []HelpItem
}

type HelpItem struct {
	Action string
	Key    string
}

type HelpView struct {
	title    string
	sections []HelpSection
}

func NewHelpView(title string, sections []HelpSection) HelpView {
	return HelpView{
		title:    title,
		sections: sections,
	}
}

func (h HelpView) View() string {
	var b string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13"))

	b += titleStyle.Render(h.title) + "\n\n"

	menuStyle := lipgloss.NewStyle().Padding(0, 2)

	for _, section := range h.sections {
		b += section.Title + ":\n\n"

		for _, item := range section.Items {
			line := fmt.Sprintf("  %-25s [%s]", item.Action, item.Key)
			b += menuStyle.Render(line) + "\n"
		}

		b += "\n"
	}

	b += lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Esc/?: close help")

	return b
}
