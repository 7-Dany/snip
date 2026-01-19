package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SearchableTableView is a table component with integrated search functionality.
// It provides real-time filtering as the user types.
type SearchableTableView struct {
	table        table.Model
	searchInput  textinput.Model
	allRows      []table.Row // Store all rows for filtering
	filteredRows []table.Row // Currently displayed rows
	title        string
	emptyMsg     string
	actionHint   string
	searchActive bool
	height       int
}

// NewSearchableTableView creates a new searchable table view.
func NewSearchableTableView(columns []table.Column, title, emptyMsg, actionHint string, height int) SearchableTableView {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Create search input
	si := textinput.New()
	si.Placeholder = "Search..."
	si.CharLimit = 50
	si.Width = 40

	return SearchableTableView{
		table:        t,
		searchInput:  si,
		allRows:      []table.Row{},
		filteredRows: []table.Row{},
		title:        title,
		emptyMsg:     emptyMsg,
		actionHint:   actionHint,
		searchActive: false,
		height:       height,
	}
}

// SetRows sets all rows and resets search filter.
func (stv *SearchableTableView) SetRows(rows []table.Row) {
	stv.allRows = rows
	stv.filteredRows = rows
	stv.table.SetRows(rows)
}

// SelectedRow returns the currently selected row.
func (stv SearchableTableView) SelectedRow() table.Row {
	return stv.table.SelectedRow()
}

// ToggleSearch enables/disables search mode.
func (stv *SearchableTableView) ToggleSearch() {
	stv.searchActive = !stv.searchActive
	if stv.searchActive {
		stv.searchInput.Focus()
		stv.table.Blur()
	} else {
		stv.searchInput.Blur()
		stv.searchInput.SetValue("")
		stv.table.Focus()
		stv.filterRows() // Reset to show all rows
	}
}

// IsSearchActive returns whether search mode is active.
func (stv SearchableTableView) IsSearchActive() bool {
	return stv.searchActive
}

// GetSearchQuery returns the current search query.
func (stv SearchableTableView) GetSearchQuery() string {
	return stv.searchInput.Value()
}

// filterRows filters the rows based on search query.
func (stv *SearchableTableView) filterRows() {
	query := strings.ToLower(strings.TrimSpace(stv.searchInput.Value()))

	if query == "" {
		stv.filteredRows = stv.allRows
		stv.table.SetRows(stv.allRows)
		return
	}

	var filtered []table.Row
	for _, row := range stv.allRows {
		// Search across all columns
		match := false
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), query) {
				match = true
				break
			}
		}
		if match {
			filtered = append(filtered, row)
		}
	}

	stv.filteredRows = filtered
	stv.table.SetRows(filtered)
}

// Update handles messages and updates state.
func (stv *SearchableTableView) Update(msg tea.Msg) (SearchableTableView, tea.Cmd) {
	var cmd tea.Cmd

	if stv.searchActive {
		// Handle search input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				// Exit search mode
				stv.ToggleSearch()
				return *stv, nil
			case "enter":
				// Exit search mode but keep filter
				stv.searchActive = false
				stv.searchInput.Blur()
				stv.table.Focus()
				return *stv, nil
			}
		}

		// Update search input
		stv.searchInput, cmd = stv.searchInput.Update(msg)
		stv.filterRows() // Real-time filtering
		return *stv, cmd
	}

	// Handle table navigation
	stv.table, cmd = stv.table.Update(msg)
	return *stv, cmd
}

// View renders the searchable table.
func (stv SearchableTableView) View() string {
	var b strings.Builder

	// Show search input if active
	if stv.searchActive {
		searchStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Bold(true)

		b.WriteString(searchStyle.Render("🔍 Search: "))
		b.WriteString(stv.searchInput.View())
		b.WriteString("\n\n")

		// Show results count
		countStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

		countText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Bold(true).
			Render(formatInt(len(stv.filteredRows)))

		b.WriteString(countStyle.Render("Found " + countText + " of " + formatInt(len(stv.allRows)) + " items"))
		b.WriteString("\n\n")
	}

	// Show table or empty message
	if len(stv.filteredRows) == 0 {
		emptyMsg := stv.emptyMsg
		if stv.searchActive && stv.searchInput.Value() != "" {
			emptyMsg = "No results found for \"" + stv.searchInput.Value() + "\""
		}
		b.WriteString(emptyMsg)
		b.WriteString("\n\n")
	} else {
		b.WriteString(stv.table.View())
		b.WriteString("\n\n")
	}

	// Action hints
	actionHint := stv.actionHint
	if stv.searchActive {
		actionHint = "Type to search | Enter: apply filter | Esc: cancel search"
	} else if stv.searchInput.Value() != "" {
		// Show that filter is active
		actionHint = "Filtered: \"" + stv.searchInput.Value() + "\" | " + stv.actionHint
	}

	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(actionHint))

	return b.String()
}

func formatInt(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	// For larger numbers, use standard formatting
	result := ""
	for n > 0 {
		result = string(rune('0'+(n%10))) + result
		n /= 10
	}
	if result == "" {
		return "0"
	}
	return result
}
