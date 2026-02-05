package tui

import (
	"fmt"
	"strings"

	"github.com/7-Dany/snip/internal/cli/components"
	"github.com/7-Dany/snip/internal/domain"
	"github.com/7-Dany/snip/internal/storage"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type viewMode int

const (
	viewModeList viewMode = iota
	viewModeMenu
	viewModeAdd
	viewModeEdit
	viewModeDelete
	viewModeSnippets
	viewModeHelp
)

type CategoriesTab struct {
	components.ViewportTab // Embed base viewport functionality
	repos                  *storage.Repositories
	mode                   viewMode

	// Components
	tableView     components.SearchableTableView
	snippetsTable components.SearchableTableView
	menuView      components.MenuView
	helpView      components.HelpView
	formView      components.FormView
	confirmDialog components.ConfirmDialog

	// State
	selectedCat *domain.Category
}

// Return pointer to match tea.Model interface
func NewCategoriesTab(repos *storage.Repositories) *CategoriesTab {
	tableView := components.NewSearchableTableView(
		[]table.Column{
			{Title: "ID", Width: 8},
			{Title: "Name", Width: 30},
			{Title: "Created", Width: 20},
			{Title: "Snippets", Width: 10},
		},
		"Categories",
		"No categories found.\n\nPress 'a' to add your first category.\nPress '?' for help.",
		"Enter: menu | a: add | /: search | r: refresh | ?: help",
		10,
	)

	snippetsTable := components.NewSearchableTableView(
		[]table.Column{
			{Title: "ID", Width: 8},
			{Title: "Title", Width: 30},
			{Title: "Language", Width: 15},
			{Title: "Created", Width: 20},
		},
		"Snippets",
		"No snippets in this category.",
		"Esc/q: back to categories",
		12,
	)

	helpView := components.NewHelpView(
		"Categories Help",
		[]components.HelpSection{
			{
				Title: "Available Actions",
				Items: []components.HelpItem{
					{Action: "Open category menu", Key: "Enter"},
					{Action: "Add new category", Key: "a"},
					{Action: "Search categories", Key: "/"},
					{Action: "Refresh list", Key: "r"},
					{Action: "Show this help", Key: "?"},
				},
			},
			{
				Title: "Navigation",
				Items: []components.HelpItem{
					{Action: "Move selection up/down", Key: "↑↓ / j/k"},
					{Action: "Scroll viewport", Key: "PgUp/PgDn"},
					{Action: "Toggle view mode", Key: "Ctrl+F"},
				},
			},
			{
				Title: "Search Mode",
				Items: []components.HelpItem{
					{Action: "Type to filter results", Key: "letters/numbers"},
					{Action: "Apply filter and exit", Key: "Enter"},
					{Action: "Cancel search", Key: "Esc"},
				},
			},
		},
	)

	tab := &CategoriesTab{
		ViewportTab:   components.NewViewportTab(),
		repos:         repos,
		mode:          viewModeList,
		tableView:     tableView,
		snippetsTable: snippetsTable,
		helpView:      helpView,
	}

	tab.refreshTable()
	return tab
}

func (c *CategoriesTab) Init() tea.Cmd {
	return nil
}

func (c *CategoriesTab) View() string {
	return c.ViewportTab.View()
}

func (c *CategoriesTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case components.ViewportResizeMsg:
		c.ViewportTab, cmd = c.ViewportTab.HandleResize(msg)
		cmds = append(cmds, cmd)
		c.updateViewportContent()
		return c, tea.Batch(cmds...)

	case tea.KeyMsg:
		c.ClearMessages()
	}

	// Handle mode-specific updates
	switch c.mode {
	case viewModeList:
		cmd = c.updateList(msg)
		cmds = append(cmds, cmd)
	case viewModeMenu:
		cmd = c.updateMenu(msg)
		cmds = append(cmds, cmd)
	case viewModeAdd:
		cmd = c.updateForm(msg, true)
		cmds = append(cmds, cmd)
	case viewModeEdit:
		cmd = c.updateForm(msg, false)
		cmds = append(cmds, cmd)
	case viewModeDelete:
		cmd = c.updateDelete(msg)
		cmds = append(cmds, cmd)
	case viewModeSnippets:
		cmd = c.updateSnippets(msg)
		cmds = append(cmds, cmd)
	case viewModeHelp:
		cmd = c.updateHelp(msg)
		cmds = append(cmds, cmd)
	}

	// Update viewport for scrolling
	c.ViewportTab, cmd = c.ViewportTab.UpdateViewport(msg)
	cmds = append(cmds, cmd)

	// Update viewport content
	c.updateViewportContent()

	return c, tea.Batch(cmds...)
}

func (c *CategoriesTab) updateViewportContent() {
	if !c.IsReady() {
		return
	}

	var b strings.Builder

	// Render messages using base component
	b.WriteString(c.RenderMessages())

	// Render mode-specific content
	switch c.mode {
	case viewModeList:
		b.WriteString(c.tableView.View())
	case viewModeMenu:
		b.WriteString(c.menuView.View())
	case viewModeAdd, viewModeEdit:
		b.WriteString(c.formView.View())
	case viewModeDelete:
		b.WriteString(c.confirmDialog.View())
	case viewModeSnippets:
		header := fmt.Sprintf("Category: %s\n\n", c.selectedCat.Name())
		b.WriteString(header)
		b.WriteString(c.snippetsTable.View())
	case viewModeHelp:
		b.WriteString(c.helpView.View())
	}

	c.SetContent(b.String())
}

func (c *CategoriesTab) refreshTable() {
	categories, err := c.repos.Categories.List()
	if err != nil {
		c.SetError(fmt.Sprintf("Error loading categories: %v", err))
		return
	}

	var rows []table.Row
	for _, cat := range categories {
		snippets, _ := c.repos.Snippets.FindByCategory(cat.ID())
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", cat.ID()),
			cat.Name(),
			cat.CreatedAt().Format("2006-01-02 15:04"),
			fmt.Sprintf("%d", len(snippets)),
		})
	}

	c.tableView.SetRows(rows)
}

func (c *CategoriesTab) updateList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.tableView.IsSearchActive() {
			c.tableView, cmd = c.tableView.Update(msg)
			return cmd
		}

		switch msg.String() {
		case "/":
			c.tableView.ToggleSearch()
			return nil
		case "enter":
			selected := c.tableView.SelectedRow()
			if len(selected) > 0 {
				if cat := c.findCategoryByName(selected[1]); cat != nil {
					c.selectedCat = cat
					c.mode = viewModeMenu
					c.createCategoryMenu()
					c.GotoTop()
					return nil
				}
			}
		case "a":
			c.mode = viewModeAdd
			c.formView = components.NewFormView("Add New Category", "", "Category name")
			c.GotoTop()
			return c.formView.Focus()
		case "r":
			c.refreshTable()
			return nil
		case "?":
			c.mode = viewModeHelp
			c.GotoTop()
			return nil
		}
	}

	c.tableView, cmd = c.tableView.Update(msg)
	return cmd
}

func (c *CategoriesTab) updateHelp(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "?":
			c.mode = viewModeList
			c.GotoTop()
			return nil
		}
	}
	return nil
}

func (c *CategoriesTab) updateMenu(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			c.menuView.MoveUp()
			return nil
		case "down", "j":
			c.menuView.MoveDown()
			return nil
		case "enter":
			return c.executeMenuAction()
		case "esc":
			c.mode = viewModeList
			c.GotoTop()
			return nil
		case "v":
			c.loadSnippetsForCategory(c.selectedCat.ID())
			c.mode = viewModeSnippets
			c.GotoTop()
			return nil
		case "e":
			c.mode = viewModeEdit
			c.formView = components.NewFormView(
				fmt.Sprintf("Edit Category: %s", c.selectedCat.Name()),
				"",
				"Category name",
			)
			c.formView.SetValue(c.selectedCat.Name())
			c.GotoTop()
			return c.formView.Focus()
		case "x":
			c.mode = viewModeDelete
			c.createDeleteDialog()
			c.GotoTop()
			return nil
		}
	}
	return nil
}

func (c *CategoriesTab) updateForm(msg tea.Msg, isAdd bool) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			name := strings.TrimSpace(c.formView.Value())
			if name == "" {
				c.SetError("Category name cannot be empty")
				return nil
			}

			if isAdd {
				return c.handleAddCategory(name)
			}
			return c.handleEditCategory(name)

		case "esc":
			c.mode = viewModeList
			c.GotoTop()
			return nil
		}
	}

	c.formView, cmd = c.formView.Update(msg)
	return cmd
}

func (c *CategoriesTab) updateDelete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			c.confirmDialog.SelectNo()
			return nil
		case "right", "l":
			c.confirmDialog.SelectYes()
			return nil
		case "y":
			c.confirmDialog.SelectYes()
			return nil
		case "n":
			c.confirmDialog.SelectNo()
			return nil
		case "enter":
			if c.confirmDialog.IsYes() {
				err := c.repos.Categories.Delete(c.selectedCat.ID())
				if err != nil {
					c.SetError(fmt.Sprintf("Error deleting category: %v", err))
				} else {
					c.SetSuccess("Category deleted successfully")
					c.refreshTable()
				}
			}
			c.mode = viewModeList
			c.GotoTop()
			return nil
		case "esc":
			c.mode = viewModeList
			c.GotoTop()
			return nil
		}
	}
	return nil
}

func (c *CategoriesTab) updateSnippets(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			c.mode = viewModeList
			c.GotoTop()
			return nil
		}
	}

	c.snippetsTable, cmd = c.snippetsTable.Update(msg)
	return cmd
}

func (c *CategoriesTab) findCategoryByName(name string) *domain.Category {
	categories, _ := c.repos.Categories.List()
	for _, cat := range categories {
		if cat.Name() == name {
			return cat
		}
	}
	return nil
}

func (c *CategoriesTab) createCategoryMenu() {
	snippets, _ := c.repos.Snippets.FindByCategory(c.selectedCat.ID())
	subtitle := fmt.Sprintf("Snippets: %d", len(snippets))

	c.menuView = components.NewMenuView(
		fmt.Sprintf("Category: %s", c.selectedCat.Name()),
		subtitle,
		[]components.MenuItem{
			{Label: "View Snippets", Shortcut: "v"},
			{Label: "Edit Category", Shortcut: "e"},
			{Label: "Delete Category", Shortcut: "x"},
		},
	)
}

func (c *CategoriesTab) createDeleteDialog() {
	snippets, _ := c.repos.Snippets.FindByCategory(c.selectedCat.ID())
	snippetCount := len(snippets)

	message := fmt.Sprintf("Category: %s\n\nAre you sure you want to delete this category?", c.selectedCat.Name())
	warning := ""

	if snippetCount > 0 {
		warning = fmt.Sprintf(
			"⚠ Warning: This category has %d snippet(s).\nDeleting will unassign them from this category.",
			snippetCount,
		)
	}

	c.confirmDialog = components.NewConfirmDialog(
		"Delete Category",
		message,
		warning,
	)
}

func (c *CategoriesTab) executeMenuAction() tea.Cmd {
	selection := c.menuView.GetSelection()

	if c.menuView.IsCancel() {
		c.mode = viewModeList
		c.GotoTop()
		return nil
	}

	switch selection {
	case 0: // View Snippets
		c.loadSnippetsForCategory(c.selectedCat.ID())
		c.mode = viewModeSnippets
		c.GotoTop()
		return nil
	case 1: // Edit
		c.mode = viewModeEdit
		c.formView = components.NewFormView(
			fmt.Sprintf("Edit Category: %s", c.selectedCat.Name()),
			"",
			"Category name",
		)
		c.formView.SetValue(c.selectedCat.Name())
		c.GotoTop()
		return c.formView.Focus()
	case 2: // Delete
		c.mode = viewModeDelete
		c.createDeleteDialog()
		c.GotoTop()
		return nil
	}

	c.mode = viewModeList
	c.GotoTop()
	return nil
}

func (c *CategoriesTab) handleAddCategory(name string) tea.Cmd {
	existing, _ := c.repos.Categories.FindByName(name)
	if existing != nil {
		c.SetError("Category already exists")
		return nil
	}

	cat, err := domain.NewCategory(name)
	if err != nil {
		c.SetError(fmt.Sprintf("Error: %v", err))
		return nil
	}

	err = c.repos.Categories.Create(cat)
	if err != nil {
		c.SetError(fmt.Sprintf("Error creating category: %v", err))
		return nil
	}

	c.SetSuccess("Category created successfully")
	c.mode = viewModeList
	c.GotoTop()
	c.refreshTable()
	return nil
}

func (c *CategoriesTab) handleEditCategory(name string) tea.Cmd {
	if name == c.selectedCat.Name() {
		c.mode = viewModeList
		c.GotoTop()
		return nil
	}

	existing, _ := c.repos.Categories.FindByName(name)
	if existing != nil && existing.ID() != c.selectedCat.ID() {
		c.SetError("Category name already exists")
		return nil
	}

	err := c.selectedCat.SetName(name)
	if err != nil {
		c.SetError(fmt.Sprintf("Error: %v", err))
		return nil
	}

	err = c.repos.Categories.Update(c.selectedCat)
	if err != nil {
		c.SetError(fmt.Sprintf("Error updating category: %v", err))
		return nil
	}

	c.SetSuccess("Category updated successfully")
	c.mode = viewModeList
	c.GotoTop()
	c.refreshTable()
	return nil
}

func (c *CategoriesTab) loadSnippetsForCategory(categoryID int) {
	snippets, err := c.repos.Snippets.FindByCategory(categoryID)
	if err != nil {
		c.SetError(fmt.Sprintf("Error loading snippets: %v", err))
		return
	}

	var rows []table.Row
	for _, snip := range snippets {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", snip.ID()),
			snip.Title(),
			snip.Language(),
			snip.CreatedAt().Format("2006-01-02 15:04"),
		})
	}

	c.snippetsTable.SetRows(rows)
}
