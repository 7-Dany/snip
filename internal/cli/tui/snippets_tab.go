package tui

import (
	"fmt"
	"strings"

	"github.com/7-Dany/snip/internal/cli/components"
	"github.com/7-Dany/snip/internal/domain"
	"github.com/7-Dany/snip/internal/storage"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type snippetViewMode int

const (
	snippetViewList snippetViewMode = iota
	snippetViewMenu
	snippetViewAdd
	snippetViewEdit
	snippetViewDelete
	snippetViewCode
	snippetViewHelp
	snippetViewSelectCategory
	snippetViewSelectTags
)

type SnippetsTab struct {
	components.ViewportTab // Embed base viewport functionality
	repos                  *storage.Repositories
	mode                   snippetViewMode

	// Components
	tableView        components.SearchableTableView
	menuView         components.MenuView
	helpView         components.HelpView
	codeEditor       components.CodeEditor
	codeViewer       components.CodeViewer
	confirmDialog    components.ConfirmDialog
	categorySelector components.SelectorView
	tagSelector      components.SelectorView

	// State
	selectedSnippet *domain.Snippet
	width           int
	height          int
}

func NewSnippetsTab(repos *storage.Repositories) *SnippetsTab {
	tableView := components.NewSearchableTableView(
		[]table.Column{
			{Title: "ID", Width: 6},
			{Title: "Title", Width: 25},
			{Title: "Language", Width: 12},
			{Title: "Category", Width: 12},
			{Title: "Tags", Width: 20},
			{Title: "Created", Width: 18},
		},
		"Snippets",
		"No snippets found.\n\nPress 'a' to add your first snippet.\nPress '?' for help.",
		"Enter: menu | a: add | /: search | r: refresh | ?: help",
		10,
	)

	helpView := components.NewHelpView(
		"Snippets Help",
		[]components.HelpSection{
			{
				Title: "Available Actions",
				Items: []components.HelpItem{
					{Action: "Open snippet menu", Key: "Enter"},
					{Action: "Add new snippet", Key: "a"},
					{Action: "Search snippets", Key: "/"},
					{Action: "Refresh list", Key: "r"},
					{Action: "Show this help", Key: "?"},
				},
			},
			{
				Title: "Editor Navigation",
				Items: []components.HelpItem{
					{Action: "Next field", Key: "Tab"},
					{Action: "Previous field", Key: "Shift+Tab"},
					{Action: "Select category", Key: "Alt+C"},
					{Action: "Manage tags", Key: "Alt+T"},
					{Action: "Save snippet", Key: "Ctrl+S"},
					{Action: "Cancel editing", Key: "Esc"},
				},
			},
			{
				Title: "Code Viewer",
				Items: []components.HelpItem{
					{Action: "Scroll code", Key: "↑↓ / PgUp/PgDn"},
					{Action: "Back to list", Key: "Esc / q"},
				},
			},
		},
	)

	tab := &SnippetsTab{
		ViewportTab: components.NewViewportTab(),
		repos:       repos,
		mode:        snippetViewList,
		tableView:   tableView,
		helpView:    helpView,
		width:       80,
		height:      24,
	}

	tab.refreshTable()
	return tab
}

func (s *SnippetsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case components.ViewportResizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.ViewportTab, cmd = s.ViewportTab.HandleResize(msg)
		cmds = append(cmds, cmd)

		// Update code editor size if in editor mode
		if s.mode == snippetViewAdd || s.mode == snippetViewEdit {
			s.codeEditor = components.NewCodeEditor(s.width, s.height)
			// Restore values if in edit mode
			if s.mode == snippetViewEdit && s.selectedSnippet != nil {
				s.restoreEditorValues()
			}
		}

		s.updateViewportContent()
		return s, tea.Batch(cmds...)

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case components.FocusChangedMsg:
		// Handle focus changes from CodeEditor to scroll viewport
		s.scrollToFocusedField(msg.FieldLine)
		return s, nil

	case tea.KeyMsg:
		s.ClearMessages()
	}

	// Handle mode-specific updates
	switch s.mode {
	case snippetViewList:
		cmd = s.updateList(msg)
		cmds = append(cmds, cmd)
	case snippetViewMenu:
		cmd = s.updateMenu(msg)
		cmds = append(cmds, cmd)
	case snippetViewAdd:
		cmd = s.updateEditor(msg, true)
		cmds = append(cmds, cmd)
	case snippetViewEdit:
		cmd = s.updateEditor(msg, false)
		cmds = append(cmds, cmd)
	case snippetViewDelete:
		cmd = s.updateDelete(msg)
		cmds = append(cmds, cmd)
	case snippetViewCode:
		cmd = s.updateCodeView(msg)
		cmds = append(cmds, cmd)
	case snippetViewHelp:
		cmd = s.updateHelp(msg)
		cmds = append(cmds, cmd)
	case snippetViewSelectCategory:
		cmd = s.updateCategorySelector(msg)
		cmds = append(cmds, cmd)
	case snippetViewSelectTags:
		cmd = s.updateTagSelector(msg)
		cmds = append(cmds, cmd)
	}

	// Only update viewport for scrolling if NOT in editor mode or if on buttons
	// This prevents arrow keys from scrolling the viewport when editing code
	shouldScroll := true
	if s.mode == snippetViewAdd || s.mode == snippetViewEdit {
		// Don't scroll if in text input fields
		if s.codeEditor.IsOnTextInput() {
			shouldScroll = false
		}
	}

	if shouldScroll {
		s.ViewportTab, cmd = s.ViewportTab.UpdateViewport(msg)
		cmds = append(cmds, cmd)
	}

	// Update viewport content
	s.updateViewportContent()

	return s, tea.Batch(cmds...)
}

func (s *SnippetsTab) Init() tea.Cmd {
	return nil
}

func (s *SnippetsTab) View() string {
	return s.ViewportTab.View()
}

func (s *SnippetsTab) updateViewportContent() {
	if !s.IsReady() {
		return
	}

	var b strings.Builder

	// Render messages
	b.WriteString(s.RenderMessages())

	switch s.mode {
	case snippetViewList:
		b.WriteString(s.tableView.View())
	case snippetViewMenu:
		b.WriteString(s.menuView.View())
	case snippetViewAdd, snippetViewEdit:
		header := "Add New Snippet"
		if s.mode == snippetViewEdit {
			header = "Edit Snippet: " + s.selectedSnippet.Title()
		}
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).Render(header))
		b.WriteString("\n\n")
		b.WriteString(s.codeEditor.View())
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).
			Render("Tab/Shift+Tab: navigate | Alt+C: category | Alt+T: tags | Ctrl+S: save | Esc: cancel"))
	case snippetViewDelete:
		b.WriteString(s.confirmDialog.View())
	case snippetViewCode:
		b.WriteString(s.codeViewer.View())
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).
			Render("Esc/q: back to list"))
	case snippetViewHelp:
		b.WriteString(s.helpView.View())
	case snippetViewSelectCategory:
		b.WriteString(s.categorySelector.View())
	case snippetViewSelectTags:
		b.WriteString(s.tagSelector.View())
	}

	s.SetContent(b.String())
}

func (s *SnippetsTab) scrollToFocusedField(fieldLine int) {
	if !s.IsReady() {
		return
	}

	// Calculate viewport middle
	viewportMiddle := s.Height() / 2

	// Calculate desired offset to center the field
	desiredYOffset := fieldLine - viewportMiddle

	// Clamp to valid range
	if desiredYOffset < 0 {
		desiredYOffset = 0
	}

	// Set viewport Y offset
	s.SetYOffset(desiredYOffset)
}

func (s *SnippetsTab) restoreEditorValues() {
	s.codeEditor.SetValues(
		s.selectedSnippet.Title(),
		s.selectedSnippet.Language(),
		s.selectedSnippet.Description(),
		s.selectedSnippet.Code(),
	)
	if s.selectedSnippet.CategoryID() != 0 {
		cat, err := s.repos.Categories.FindByID(s.selectedSnippet.CategoryID())
		if err == nil {
			s.codeEditor.SetCategory(cat.ID(), cat.Name())
		}
	}
	tagIDs := s.selectedSnippet.Tags()
	tagNames := []string{}
	for _, tagID := range tagIDs {
		tag, err := s.repos.Tags.FindByID(tagID)
		if err == nil {
			tagNames = append(tagNames, tag.Name())
		}
	}
	s.codeEditor.SetTags(tagIDs, tagNames)
}

func (s *SnippetsTab) refreshTable() {
	snippets, err := s.repos.Snippets.List()
	if err != nil {
		s.SetError(fmt.Sprintf("Error loading snippets: %v", err))
		return
	}

	var rows []table.Row
	for _, snip := range snippets {
		categoryName := "None"
		if snip.CategoryID() != 0 {
			cat, err := s.repos.Categories.FindByID(snip.CategoryID())
			if err == nil {
				categoryName = cat.Name()
			}
		}

		tagNames := []string{}
		for _, tagID := range snip.Tags() {
			tag, err := s.repos.Tags.FindByID(tagID)
			if err == nil {
				tagNames = append(tagNames, tag.Name())
			}
		}
		tagsStr := "None"
		if len(tagNames) > 0 {
			tagsStr = strings.Join(tagNames, ", ")
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", snip.ID()),
			truncate(snip.Title(), 25),
			snip.Language(),
			truncate(categoryName, 12),
			truncate(tagsStr, 20),
			snip.CreatedAt().Format("2006-01-02 15:04"),
		})
	}

	s.tableView.SetRows(rows)
}

func truncate(str string, max int) string {
	if len(str) <= max {
		return str
	}
	return str[:max-3] + "..."
}

func (s *SnippetsTab) loadCategoriesIntoSelector() {
	categories, _ := s.repos.Categories.List()

	var rows []table.Row
	for _, cat := range categories {
		snippets, _ := s.repos.Snippets.FindByCategory(cat.ID())
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", cat.ID()),
			cat.Name(),
			fmt.Sprintf("%d", len(snippets)),
		})
	}

	s.categorySelector.SetRows(rows)
}

func (s *SnippetsTab) loadTagsIntoSelector() {
	tags, _ := s.repos.Tags.List()

	currentTags := s.codeEditor.GetTags()
	s.tagSelector.SetSelected(currentTags)

	var rows []table.Row
	for _, tag := range tags {
		snippets, _ := s.repos.Snippets.FindByTag(tag.ID())

		checkbox := " "
		for _, id := range currentTags {
			if id == tag.ID() {
				checkbox = "✓"
				break
			}
		}

		rows = append(rows, table.Row{
			checkbox,
			fmt.Sprintf("%d", tag.ID()),
			tag.Name(),
			fmt.Sprintf("%d", len(snippets)),
		})
	}

	s.tagSelector.SetRows(rows)
	s.tagSelector.UpdateCheckboxDisplay()
}

func (s *SnippetsTab) updateList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.tableView.IsSearchActive() {
			s.tableView, cmd = s.tableView.Update(msg)
			return cmd
		}

		switch msg.String() {
		case "/":
			s.tableView.ToggleSearch()
			return nil
		case "enter":
			selected := s.tableView.SelectedRow()
			if len(selected) > 0 {
				if snip := s.findSnippetByTitle(selected[1]); snip != nil {
					s.selectedSnippet = snip
					s.mode = snippetViewMenu
					s.createSnippetMenu()
					s.GotoTop()
					return nil
				}
			}
		case "a":
			s.mode = snippetViewAdd
			s.codeEditor = components.NewCodeEditor(s.width, s.height)
			s.GotoTop()
			return nil
		case "r":
			s.refreshTable()
			return nil
		case "?":
			s.mode = snippetViewHelp
			s.GotoTop()
			return nil
		}
	}

	s.tableView, cmd = s.tableView.Update(msg)
	return cmd
}

func (s *SnippetsTab) updateHelp(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "?":
			s.mode = snippetViewList
			s.GotoTop()
			return nil
		}
	}
	return nil
}

func (s *SnippetsTab) updateMenu(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.menuView.MoveUp()
			return nil
		case "down", "j":
			s.menuView.MoveDown()
			return nil
		case "enter":
			return s.executeMenuAction()
		case "esc":
			s.mode = snippetViewList
			s.GotoTop()
			return nil
		case "v":
			s.mode = snippetViewCode
			s.codeViewer = components.NewCodeViewer(
				s.selectedSnippet.Title(),
				s.selectedSnippet.Language(),
				s.selectedSnippet.Description(),
				s.selectedSnippet.Code(),
				s.width,
			)
			s.GotoTop()
			return nil
		case "e":
			s.mode = snippetViewEdit
			s.codeEditor = components.NewCodeEditor(s.width, s.height)
			s.restoreEditorValues()
			s.GotoTop()
			return nil
		case "x":
			s.mode = snippetViewDelete
			s.createDeleteDialog()
			s.GotoTop()
			return nil
		}
	}
	return nil
}

func (s *SnippetsTab) updateEditor(msg tea.Msg, isAdd bool) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keys that should work regardless of focus
		switch msg.String() {
		case "ctrl+s", "alt+s":
			return s.handleSaveSnippet(isAdd)
		case "esc":
			s.mode = snippetViewList
			s.GotoTop()
			return nil
		case "alt+c":
			// Open category selector
			s.mode = snippetViewSelectCategory
			s.categorySelector = components.NewSelectorView(
				"Select Category",
				"Choose a category for this snippet",
				false,
				s.width,
				s.height,
			)
			s.loadCategoriesIntoSelector()
			s.GotoTop()
			return nil
		case "alt+t":
			// Open tag selector
			s.mode = snippetViewSelectTags
			s.tagSelector = components.NewSelectorView(
				"Select Tags",
				"Choose tags for this snippet (Space to toggle)",
				true,
				s.width,
				s.height,
			)
			s.loadTagsIntoSelector()
			s.GotoTop()
			return nil
		case "enter":
			// Handle enter on buttons
			if s.codeEditor.IsCancelFocused() {
				s.mode = snippetViewList
				s.GotoTop()
				return nil
			}
			if s.codeEditor.IsSaveFocused() {
				return s.handleSaveSnippet(isAdd)
			}
		}
	}

	// Let code editor handle the message
	s.codeEditor, cmd = s.codeEditor.Update(msg)
	return cmd
}

func (s *SnippetsTab) handleSaveSnippet(isAdd bool) tea.Cmd {
	title, language, description, code := s.codeEditor.GetValues()

	if title == "" {
		s.SetError("Title cannot be empty")
		return nil
	}
	if language == "" {
		s.SetError("Language cannot be empty")
		return nil
	}
	if code == "" {
		s.SetError("Code cannot be empty")
		return nil
	}

	if isAdd {
		return s.handleAddSnippet(title, language, description, code)
	}
	return s.handleEditSnippet(title, language, description, code)
}

func (s *SnippetsTab) updateDelete(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "n":
			s.confirmDialog.SelectNo()
			return nil
		case "right", "l", "y":
			s.confirmDialog.SelectYes()
			return nil
		case "enter":
			if s.confirmDialog.IsYes() {
				err := s.repos.Snippets.Delete(s.selectedSnippet.ID())
				if err != nil {
					s.SetError(fmt.Sprintf("Error deleting snippet: %v", err))
				} else {
					s.SetSuccess("Snippet deleted successfully")
					s.refreshTable()
				}
			}
			s.mode = snippetViewList
			s.GotoTop()
			return nil
		case "esc":
			s.mode = snippetViewList
			s.GotoTop()
			return nil
		}
	}
	return nil
}

func (s *SnippetsTab) updateCodeView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			s.mode = snippetViewList
			s.GotoTop()
			return nil
		}
	}
	return nil
}

func (s *SnippetsTab) updateCategorySelector(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			categoryID := s.categorySelector.GetSelectedID()
			if categoryID != 0 {
				cat, err := s.repos.Categories.FindByID(categoryID)
				if err == nil {
					s.codeEditor.SetCategory(cat.ID(), cat.Name())
				}
			}
			if s.selectedSnippet != nil {
				s.mode = snippetViewEdit
			} else {
				s.mode = snippetViewAdd
			}
			s.GotoTop()
			return nil

		case "n":
			s.codeEditor.SetCategory(0, "None")
			if s.selectedSnippet != nil {
				s.mode = snippetViewEdit
			} else {
				s.mode = snippetViewAdd
			}
			s.GotoTop()
			return nil

		case "esc":
			if s.selectedSnippet != nil {
				s.mode = snippetViewEdit
			} else {
				s.mode = snippetViewAdd
			}
			s.GotoTop()
			return nil
		}
	}

	updatedSelector, cmd := s.categorySelector.Update(msg)
	s.categorySelector = updatedSelector
	return cmd
}

func (s *SnippetsTab) updateTagSelector(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			tagIDs := s.tagSelector.GetSelectedIDs()
			tagNames := []string{}
			for _, id := range tagIDs {
				tag, err := s.repos.Tags.FindByID(id)
				if err == nil {
					tagNames = append(tagNames, tag.Name())
				}
			}

			s.codeEditor.SetTags(tagIDs, tagNames)

			if s.selectedSnippet != nil {
				s.mode = snippetViewEdit
			} else {
				s.mode = snippetViewAdd
			}
			s.GotoTop()
			return nil

		case "n":
			s.codeEditor.SetTags([]int{}, []string{})
			if s.selectedSnippet != nil {
				s.mode = snippetViewEdit
			} else {
				s.mode = snippetViewAdd
			}
			s.GotoTop()
			return nil

		case "esc":
			if s.selectedSnippet != nil {
				s.mode = snippetViewEdit
			} else {
				s.mode = snippetViewAdd
			}
			s.GotoTop()
			return nil
		}
	}

	updatedSelector, cmd := s.tagSelector.Update(msg)
	s.tagSelector = updatedSelector
	return cmd
}

func (s *SnippetsTab) findSnippetByTitle(title string) *domain.Snippet {
	title = strings.TrimSuffix(title, "...")

	snippets, _ := s.repos.Snippets.List()
	for _, snip := range snippets {
		if strings.HasPrefix(snip.Title(), title) {
			return snip
		}
	}
	return nil
}

func (s *SnippetsTab) createSnippetMenu() {
	categoryName := "None"
	if s.selectedSnippet.CategoryID() != 0 {
		cat, err := s.repos.Categories.FindByID(s.selectedSnippet.CategoryID())
		if err == nil {
			categoryName = cat.Name()
		}
	}

	tagNames := []string{}
	for _, tagID := range s.selectedSnippet.Tags() {
		tag, err := s.repos.Tags.FindByID(tagID)
		if err == nil {
			tagNames = append(tagNames, tag.Name())
		}
	}
	tagsStr := "None"
	if len(tagNames) > 0 {
		tagsStr = strings.Join(tagNames, ", ")
	}

	subtitle := fmt.Sprintf("Language: %s | Category: %s\nTags: %s",
		s.selectedSnippet.Language(),
		categoryName,
		tagsStr,
	)

	s.menuView = components.NewMenuView(
		s.selectedSnippet.Title(),
		subtitle,
		[]components.MenuItem{
			{Label: "View Code", Shortcut: "v"},
			{Label: "Edit Snippet", Shortcut: "e"},
			{Label: "Delete Snippet", Shortcut: "x"},
		},
	)
}

func (s *SnippetsTab) createDeleteDialog() {
	message := fmt.Sprintf("Title: %s\nLanguage: %s\n\nAre you sure you want to delete this snippet?",
		s.selectedSnippet.Title(),
		s.selectedSnippet.Language(),
	)

	warning := "⚠ Warning: This action cannot be undone.\nThe snippet will be permanently deleted."

	s.confirmDialog = components.NewConfirmDialog(
		"Delete Snippet",
		message,
		warning,
	)
}

func (s *SnippetsTab) executeMenuAction() tea.Cmd {
	selection := s.menuView.GetSelection()

	if s.menuView.IsCancel() {
		s.mode = snippetViewList
		s.GotoTop()
		return nil
	}

	switch selection {
	case 0: // View Code
		s.mode = snippetViewCode
		s.codeViewer = components.NewCodeViewer(
			s.selectedSnippet.Title(),
			s.selectedSnippet.Language(),
			s.selectedSnippet.Description(),
			s.selectedSnippet.Code(),
			s.width,
		)
		s.GotoTop()
		return nil
	case 1: // Edit
		s.mode = snippetViewEdit
		s.codeEditor = components.NewCodeEditor(s.width, s.height)
		s.restoreEditorValues()
		s.GotoTop()
		return nil
	case 2: // Delete
		s.mode = snippetViewDelete
		s.createDeleteDialog()
		s.GotoTop()
		return nil
	}

	s.mode = snippetViewList
	s.GotoTop()
	return nil
}

func (s *SnippetsTab) handleAddSnippet(title, language, description, code string) tea.Cmd {
	snippet, err := domain.NewSnippet(title, language, code)
	if err != nil {
		s.SetError(fmt.Sprintf("Error: %v", err))
		return nil
	}

	if description != "" {
		snippet.SetDescription(description)
	}

	categoryID := s.codeEditor.GetCategory()
	if categoryID != 0 {
		snippet.SetCategory(categoryID)
	}

	tagIDs := s.codeEditor.GetTags()
	for _, tagID := range tagIDs {
		snippet.AddTag(tagID)
	}

	err = s.repos.Snippets.Create(snippet)
	if err != nil {
		s.SetError(fmt.Sprintf("Error creating snippet: %v", err))
		return nil
	}

	s.SetSuccess("Snippet created successfully")
	s.mode = snippetViewList
	s.GotoTop()
	s.refreshTable()
	return nil
}

func (s *SnippetsTab) handleEditSnippet(title, language, description, code string) tea.Cmd {
	if title != s.selectedSnippet.Title() {
		err := s.selectedSnippet.SetTitle(title)
		if err != nil {
			s.SetError(fmt.Sprintf("Error: %v", err))
			return nil
		}
	}

	if language != s.selectedSnippet.Language() {
		err := s.selectedSnippet.SetLanguage(language)
		if err != nil {
			s.SetError(fmt.Sprintf("Error: %v", err))
			return nil
		}
	}

	if description != s.selectedSnippet.Description() {
		s.selectedSnippet.SetDescription(description)
	}

	if code != s.selectedSnippet.Code() {
		err := s.selectedSnippet.SetCode(code)
		if err != nil {
			s.SetError(fmt.Sprintf("Error: %v", err))
			return nil
		}
	}

	categoryID := s.codeEditor.GetCategory()
	if categoryID != s.selectedSnippet.CategoryID() {
		s.selectedSnippet.SetCategory(categoryID)
	}

	oldTags := s.selectedSnippet.Tags()
	newTags := s.codeEditor.GetTags()

	// Remove tags that are no longer selected
	for _, oldTag := range oldTags {
		found := false
		for _, newTag := range newTags {
			if oldTag == newTag {
				found = true
				break
			}
		}
		if !found {
			s.selectedSnippet.RemoveTag(oldTag)
		}
	}

	// Add new tags
	for _, newTag := range newTags {
		if !s.selectedSnippet.HasTag(newTag) {
			s.selectedSnippet.AddTag(newTag)
		}
	}

	err := s.repos.Snippets.Update(s.selectedSnippet)
	if err != nil {
		s.SetError(fmt.Sprintf("Error updating snippet: %v", err))
		return nil
	}

	s.SetSuccess("Snippet updated successfully")
	s.mode = snippetViewList
	s.GotoTop()
	s.refreshTable()
	return nil
}
