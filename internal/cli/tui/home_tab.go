package tui

import (
	"github.com/7-Dany/snip/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
)

// HomeTab displays welcome message and navigation instructions.
type HomeTab struct {
	components.ViewportTab
}

// NewHomeTab creates a new home tab instance.
func NewHomeTab() HomeTab {
	return HomeTab{
		ViewportTab: components.NewViewportTab(),
	}
}

func (h HomeTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle viewport resize
	switch msg := msg.(type) {
	case components.ViewportResizeMsg:
		h.ViewportTab, cmd = h.ViewportTab.HandleResize(msg)
		if h.IsReady() {
			h.SetContent(h.getContent())
		}
		return h, cmd
	}

	// Update viewport for scrolling
	h.ViewportTab, cmd = h.ViewportTab.UpdateViewport(msg)
	return h, cmd
}

func (h HomeTab) getContent() string {
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

Navigation:
  Ctrl+F      - Toggle between Interactive and Navigation mode
  Tab         - Next tab (in Navigation mode)
  Shift+Tab   - Previous tab (in Navigation mode)
  ←/→         - Switch tabs (in Navigation mode)
  PgUp/PgDn   - Scroll content (in Interactive mode)
  Ctrl+←/→    - Scroll horizontally
  Ctrl+C      - Quit

Mode Indicators:
  ⌨ (green)   - Interactive mode: Components receive keyboard input
  ⇄ (blue)    - Navigation mode: Tab switching and scrolling enabled

The mode indicator appears in the active tab header.

---

Additional Information:
This is a terminal-based code snippet manager that helps you organize
and quickly access your code snippets. You can categorize snippets,
add tags for better organization, and search through your collection.

Features:
- Category management
- Tag-based organization
- Full-text search
- Syntax highlighting
- Export/Import functionality
- Quick snippet insertion

Getting Started:
1. Create categories to organize your snippets
2. Add tags for flexible organization
3. Start adding your code snippets
4. Use search to quickly find what you need

Tips:
- Use descriptive names for categories
- Tag snippets with multiple relevant tags
- Include context in snippet descriptions
- Regularly review and update your snippets

This content is scrollable when it exceeds the viewport height.
Try resizing your terminal to see the viewport adjust automatically!
`
}
