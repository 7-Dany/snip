package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type HomeTab struct{}

func NewHomeTab() HomeTab {
	return HomeTab{}
}

func (h HomeTab) Init() tea.Cmd {
	return nil
}

func (h HomeTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

func (h HomeTab) View() string {
	return `
  ████████
██        ██
██  ██  ██  ██
██  ██  ██  ██
██          ██
██  ████████ ██
████      ████

SNIP - Code Snippet Manager

Welcome! Use the tabs above to navigate:
- Home - You are here
- Categories - Manage snippet categories
- Tags - Manage snippet tags
- Snippets - Manage code snippets

Press Tab/Shift+Tab or ←→ to switch tabs
Press q or Ctrl+C to quit
`
}
