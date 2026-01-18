package tui

import (
	"fmt"
	"strings"

	"github.com/7-Dany/snip/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type CategoriesTab struct {
	repos *storage.Repositories
}

func NewCategoriesTab(repos *storage.Repositories) CategoriesTab {
	return CategoriesTab{repos: repos}
}

func (c CategoriesTab) Init() tea.Cmd {
	return nil
}

func (c CategoriesTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// For now, just handle quit - no interaction yet
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c CategoriesTab) View() string {
	// Get categories from repos
	categories, err := c.repos.Categories.List()
	if err != nil {
		return fmt.Sprintf("Error loading categories: %v", err)
	}

	if len(categories) == 0 {
		return "No categories found.\n\nPress 'a' to add a category."
	}

	// Build a simple table
	var b strings.Builder
	b.WriteString("CATEGORIES\n\n")
	b.WriteString(fmt.Sprintf("%-5s %-20s %-20s\n", "ID", "Name", "Created"))
	b.WriteString(strings.Repeat("─", 50) + "\n")

	for _, cat := range categories {
		b.WriteString(fmt.Sprintf("%-5d %-20s %-20s\n",
			cat.ID(),
			cat.Name(),
			cat.CreatedAt().Format("2006-01-02 15:04"),
		))
	}

	b.WriteString("\n")
	b.WriteString("↑↓ navigate | a:add | e:edit | d:delete | q:quit")

	return b.String()
}
