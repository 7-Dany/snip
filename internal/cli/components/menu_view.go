package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	Label    string
	Shortcut string
	Action   func() tea.Cmd
}

type MenuView struct {
	title       string
	subtitle    string
	items       []MenuItem
	selection   int
	cancelLabel string
}

func NewMenuView(title, subtitle string, items []MenuItem) MenuView {
	return MenuView{
		title:       title,
		subtitle:    subtitle,
		items:       items,
		selection:   0,
		cancelLabel: "Cancel",
	}
}

func (m *MenuView) MoveUp() {
	m.selection = max(0, m.selection-1)
}

func (m *MenuView) MoveDown() {
	m.selection = min(len(m.items), m.selection+1)
}

func (m MenuView) GetSelection() int {
	return m.selection
}

func (m MenuView) IsCancel() bool {
	return m.selection == len(m.items)
}

func (m MenuView) View() string {
	var b string

	// Title
	b += lipgloss.NewStyle().Bold(true).Render(m.title) + "\n\n"

	// Subtitle
	if m.subtitle != "" {
		b += m.subtitle + "\n\n"
	}

	b += "Select an action:\n\n"

	// Menu items
	menuStyle := lipgloss.NewStyle().Padding(0, 2)
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("13")).
		Padding(0, 2).
		Bold(true)

	// Render menu items
	for i, item := range m.items {
		label := fmt.Sprintf("%-20s", item.Label)
		if i == m.selection {
			b += selectedStyle.Render("▸ "+label) + fmt.Sprintf("  [%s]", item.Shortcut)
		} else {
			b += menuStyle.Render("  "+label) + fmt.Sprintf("  [%s]", item.Shortcut)
		}
		b += "\n"
	}

	// Cancel option
	cancelLabel := fmt.Sprintf("%-20s", m.cancelLabel)
	if m.selection == len(m.items) {
		b += selectedStyle.Render("▸ "+cancelLabel) + "  [esc]"
	} else {
		b += menuStyle.Render("  "+cancelLabel) + "  [esc]"
	}
	b += "\n"

	b += "\n"
	b += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).
		Render("↑↓/jk: navigate | Enter: select | Esc: cancel")
	b += "\n"

	// Build quick keys hint
	var shortcuts []string
	for _, item := range m.items {
		shortcuts = append(shortcuts, fmt.Sprintf("%s: %s", item.Shortcut, item.Label[:min(len(item.Label), 4)]))
	}
	b += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).
		Render("Quick keys: " + joinStrings(shortcuts, " | "))

	return b
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
