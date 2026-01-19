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

type TagsTab struct {
	components.ViewportTab // Embed base viewport functionality
	repos                  *storage.Repositories
	mode                   viewMode

	// Reusable components
	tableView     components.SearchableTableView
	snippetsTable components.TableView
	menuView      components.MenuView
	helpView      components.HelpView
	formView      components.FormView
	confirmDialog components.ConfirmDialog

	// State
	selectedTag *domain.Tag
}

// Return pointer to match tea.Model interface
func NewTagsTab(repos *storage.Repositories) *TagsTab {
	tableView := components.NewSearchableTableView(
		[]table.Column{
			{Title: "ID", Width: 8},
			{Title: "Name", Width: 30},
			{Title: "Created", Width: 20},
			{Title: "Snippets", Width: 10},
		},
		"Tags",
		"No tags found.\n\nPress 'a' to add your first tag.\nPress '?' for help.",
		"Enter: menu | a: add | /: search | r: refresh | ?: help",
		10,
	)

	snippetsTable := components.NewTableView(
		[]table.Column{
			{Title: "ID", Width: 8},
			{Title: "Title", Width: 30},
			{Title: "Language", Width: 15},
			{Title: "Created", Width: 20},
		},
		"Snippets",
		"No snippets with this tag.",
		"Esc/q: back to tags",
		12,
	)

	helpView := components.NewHelpView(
		"Tags Help",
		[]components.HelpSection{
			{
				Title: "Available Actions",
				Items: []components.HelpItem{
					{Action: "Open tag menu", Key: "Enter"},
					{Action: "Add new tag", Key: "a"},
					{Action: "Search tags", Key: "/"},
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

	tab := &TagsTab{
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

func (t *TagsTab) Init() tea.Cmd {
	return nil
}

func (t *TagsTab) View() string {
	return t.ViewportTab.View()
}

func (t *TagsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case components.ViewportResizeMsg:
		t.ViewportTab, cmd = t.ViewportTab.HandleResize(msg)
		cmds = append(cmds, cmd)
		t.updateViewportContent()
		return t, tea.Batch(cmds...)

	case tea.KeyMsg:
		t.ClearMessages()
	}

	// Handle mode-specific updates
	switch t.mode {
	case viewModeList:
		cmd = t.updateList(msg)
		cmds = append(cmds, cmd)
	case viewModeMenu:
		cmd = t.updateMenu(msg)
		cmds = append(cmds, cmd)
	case viewModeAdd:
		cmd = t.updateForm(msg, true)
		cmds = append(cmds, cmd)
	case viewModeEdit:
		cmd = t.updateForm(msg, false)
		cmds = append(cmds, cmd)
	case viewModeDelete:
		cmd = t.updateDelete(msg)
		cmds = append(cmds, cmd)
	case viewModeSnippets:
		cmd = t.updateSnippets(msg)
		cmds = append(cmds, cmd)
	case viewModeHelp:
		cmd = t.updateHelp(msg)
		cmds = append(cmds, cmd)
	}

	// Update viewport for scrolling
	t.ViewportTab, cmd = t.ViewportTab.UpdateViewport(msg)
	cmds = append(cmds, cmd)

	// Update viewport content
	t.updateViewportContent()

	return t, tea.Batch(cmds...)
}

func (t *TagsTab) updateViewportContent() {
	if !t.IsReady() {
		return
	}

	var b strings.Builder

	// Render messages using base component
	b.WriteString(t.RenderMessages())

	switch t.mode {
	case viewModeList:
		b.WriteString(t.tableView.View())
	case viewModeMenu:
		b.WriteString(t.menuView.View())
	case viewModeAdd, viewModeEdit:
		b.WriteString(t.formView.View())
	case viewModeDelete:
		b.WriteString(t.confirmDialog.View())
	case viewModeSnippets:
		header := fmt.Sprintf("Tag: %s\n\n", t.selectedTag.Name())
		b.WriteString(header)
		b.WriteString(t.snippetsTable.View())
	case viewModeHelp:
		b.WriteString(t.helpView.View())
	}

	t.SetContent(b.String())
}

func (t *TagsTab) refreshTable() {
	tags, err := t.repos.Tags.List()
	if err != nil {
		t.SetError(fmt.Sprintf("Error loading tags: %v", err))
		return
	}

	var rows []table.Row
	for _, tag := range tags {
		snippets, _ := t.repos.Snippets.FindByTag(tag.ID())
		snippetCount := fmt.Sprintf("%d", len(snippets))

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", tag.ID()),
			tag.Name(),
			tag.CreatedAt().Format("2006-01-02 15:04"),
			snippetCount,
		})
	}

	t.tableView.SetRows(rows)
}

func (t *TagsTab) updateList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if t.tableView.IsSearchActive() {
			t.tableView, cmd = t.tableView.Update(msg)
			return cmd
		}

		switch msg.String() {
		case "/":
			t.tableView.ToggleSearch()
			return nil
		case "enter":
			selected := t.tableView.SelectedRow()
			if len(selected) > 0 {
				if tag := t.findTagByName(selected[1]); tag != nil {
					t.selectedTag = tag
					t.mode = viewModeMenu
					t.createTagMenu()
					t.GotoTop()
					return nil
				}
			}
		case "a":
			t.mode = viewModeAdd
			t.formView = components.NewFormView("Add New Tag", "", "Tag name")
			t.GotoTop()
			return t.formView.Focus()
		case "r":
			t.refreshTable()
			return nil
		case "?":
			t.mode = viewModeHelp
			t.GotoTop()
			return nil
		}
	}

	t.tableView, cmd = t.tableView.Update(msg)
	return cmd
}

func (t *TagsTab) updateHelp(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "?":
			t.mode = viewModeList
			t.GotoTop()
			return nil
		}
	}
	return nil
}

func (t *TagsTab) updateMenu(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			t.menuView.MoveUp()
			return nil
		case "down", "j":
			t.menuView.MoveDown()
			return nil
		case "enter":
			return t.executeMenuAction()
		case "esc":
			t.mode = viewModeList
			t.GotoTop()
			return nil
		case "v":
			t.loadSnippetsForTag(t.selectedTag.ID())
			t.mode = viewModeSnippets
			t.GotoTop()
			return nil
		case "e":
			t.mode = viewModeEdit
			t.formView = components.NewFormView(
				fmt.Sprintf("Edit Tag: %s", t.selectedTag.Name()),
				"",
				"Tag name",
			)
			t.formView.SetValue(t.selectedTag.Name())
			t.GotoTop()
			return t.formView.Focus()
		case "x":
			t.mode = viewModeDelete
			t.createDeleteDialog()
			t.GotoTop()
			return nil
		}
	}
	return nil
}

func (t *TagsTab) updateForm(msg tea.Msg, isAdd bool) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			name := strings.TrimSpace(t.formView.Value())
			if name == "" {
				t.SetError("Tag name cannot be empty")
				return nil
			}

			if isAdd {
				return t.handleAddTag(name)
			}
			return t.handleEditTag(name)

		case "esc":
			t.mode = viewModeList
			t.GotoTop()
			return nil
		}
	}

	t.formView, cmd = t.formView.Update(msg)
	return cmd
}

func (t *TagsTab) updateDelete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			t.confirmDialog.SelectNo()
			return nil
		case "right", "l":
			t.confirmDialog.SelectYes()
			return nil
		case "y":
			t.confirmDialog.SelectYes()
			return nil
		case "n":
			t.confirmDialog.SelectNo()
			return nil
		case "enter":
			if t.confirmDialog.IsYes() {
				err := t.repos.Tags.Delete(t.selectedTag.ID())
				if err != nil {
					t.SetError(fmt.Sprintf("Error deleting tag: %v", err))
				} else {
					t.SetSuccess("Tag deleted successfully")
					t.refreshTable()
				}
			}
			t.mode = viewModeList
			t.GotoTop()
			return nil
		case "esc":
			t.mode = viewModeList
			t.GotoTop()
			return nil
		}
	}
	return nil
}

func (t *TagsTab) updateSnippets(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			t.mode = viewModeList
			t.GotoTop()
			return nil
		}
	}

	t.snippetsTable, cmd = t.snippetsTable.Update(msg)
	return cmd
}

func (t *TagsTab) findTagByName(name string) *domain.Tag {
	tags, _ := t.repos.Tags.List()
	for _, tag := range tags {
		if tag.Name() == name {
			return tag
		}
	}
	return nil
}

func (t *TagsTab) createTagMenu() {
	snippets, _ := t.repos.Snippets.FindByTag(t.selectedTag.ID())
	subtitle := fmt.Sprintf("Snippets: %d", len(snippets))

	t.menuView = components.NewMenuView(
		fmt.Sprintf("Tag: %s", t.selectedTag.Name()),
		subtitle,
		[]components.MenuItem{
			{Label: "View Snippets", Shortcut: "v"},
			{Label: "Edit Tag", Shortcut: "e"},
			{Label: "Delete Tag", Shortcut: "x"},
		},
	)
}

func (t *TagsTab) createDeleteDialog() {
	snippets, _ := t.repos.Snippets.FindByTag(t.selectedTag.ID())
	snippetCount := len(snippets)

	message := fmt.Sprintf("Tag: %s\n\nAre you sure you want to delete this tag?", t.selectedTag.Name())
	warning := ""

	if snippetCount > 0 {
		warning = fmt.Sprintf(
			"⚠ Warning: This tag is used by %d snippet(s).\nDeleting will remove it from all snippets.",
			snippetCount,
		)
	}

	t.confirmDialog = components.NewConfirmDialog(
		"Delete Tag",
		message,
		warning,
	)
}

func (t *TagsTab) executeMenuAction() tea.Cmd {
	selection := t.menuView.GetSelection()

	if t.menuView.IsCancel() {
		t.mode = viewModeList
		t.GotoTop()
		return nil
	}

	switch selection {
	case 0: // View Snippets
		t.loadSnippetsForTag(t.selectedTag.ID())
		t.mode = viewModeSnippets
		t.GotoTop()
		return nil
	case 1: // Edit
		t.mode = viewModeEdit
		t.formView = components.NewFormView(
			fmt.Sprintf("Edit Tag: %s", t.selectedTag.Name()),
			"",
			"Tag name",
		)
		t.formView.SetValue(t.selectedTag.Name())
		t.GotoTop()
		return t.formView.Focus()
	case 2: // Delete
		t.mode = viewModeDelete
		t.createDeleteDialog()
		t.GotoTop()
		return nil
	}

	t.mode = viewModeList
	t.GotoTop()
	return nil
}

func (t *TagsTab) handleAddTag(name string) tea.Cmd {
	existing, _ := t.repos.Tags.FindByName(name)
	if existing != nil {
		t.SetError("Tag already exists")
		return nil
	}

	tag, err := domain.NewTag(name)
	if err != nil {
		t.SetError(fmt.Sprintf("Error: %v", err))
		return nil
	}

	err = t.repos.Tags.Create(tag)
	if err != nil {
		t.SetError(fmt.Sprintf("Error creating tag: %v", err))
		return nil
	}

	t.SetSuccess("Tag created successfully")
	t.mode = viewModeList
	t.GotoTop()
	t.refreshTable()
	return nil
}

func (t *TagsTab) handleEditTag(name string) tea.Cmd {
	if name == t.selectedTag.Name() {
		t.mode = viewModeList
		t.GotoTop()
		return nil
	}

	existing, _ := t.repos.Tags.FindByName(name)
	if existing != nil && existing.ID() != t.selectedTag.ID() {
		t.SetError("Tag name already exists")
		return nil
	}

	err := t.selectedTag.SetName(name)
	if err != nil {
		t.SetError(fmt.Sprintf("Error: %v", err))
		return nil
	}

	err = t.repos.Tags.Update(t.selectedTag)
	if err != nil {
		t.SetError(fmt.Sprintf("Error updating tag: %v", err))
		return nil
	}

	t.SetSuccess("Tag updated successfully")
	t.mode = viewModeList
	t.GotoTop()
	t.refreshTable()
	return nil
}

func (t *TagsTab) loadSnippetsForTag(tagID int) {
	snippets, err := t.repos.Snippets.FindByTag(tagID)
	if err != nil {
		t.SetError(fmt.Sprintf("Error loading snippets: %v", err))
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

	t.snippetsTable.SetRows(rows)
}
