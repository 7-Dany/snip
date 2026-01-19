package components

import "github.com/charmbracelet/lipgloss"

type ConfirmDialog struct {
	title    string
	message  string
	warning  string
	choice   int // 0=No, 1=Yes
	yesLabel string
	noLabel  string
	yesColor string
	noColor  string
}

func NewConfirmDialog(title, message, warning string) ConfirmDialog {
	return ConfirmDialog{
		title:    title,
		message:  message,
		warning:  warning,
		choice:   0,
		yesLabel: "Yes",
		noLabel:  "No",
		yesColor: "9",
		noColor:  "10",
	}
}

func (c *ConfirmDialog) SelectYes() {
	c.choice = 1
}

func (c *ConfirmDialog) SelectNo() {
	c.choice = 0
}

func (c ConfirmDialog) IsYes() bool {
	return c.choice == 1
}

func (c ConfirmDialog) View() string {
	var b string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(c.yesColor))

	b += titleStyle.Render(c.title) + "\n\n"
	b += c.message + "\n\n"

	if c.warning != "" {
		warningBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("11")).
			Padding(0, 1).
			Render(c.warning)
		b += warningBox + "\n\n"
	}

	// Buttons
	noStyle := lipgloss.NewStyle().
		Padding(0, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	yesStyle := lipgloss.NewStyle().
		Padding(0, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	selectedNoStyle := lipgloss.NewStyle().
		Padding(0, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.noColor)).
		Background(lipgloss.Color(c.noColor)).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	selectedYesStyle := lipgloss.NewStyle().
		Padding(0, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.yesColor)).
		Background(lipgloss.Color(c.yesColor)).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	var noButton, yesButton string
	if c.choice == 0 {
		noButton = selectedNoStyle.Render(c.noLabel)
		yesButton = yesStyle.Render(c.yesLabel)
	} else {
		noButton = noStyle.Render(c.noLabel)
		yesButton = selectedYesStyle.Render(c.yesLabel)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, noButton, "  ", yesButton)
	b += buttons + "\n\n"

	b += lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("←→/hl: select | Enter: confirm | n/y: quick select | Esc: cancel")

	return b
}
