package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectorView provides a searchable selection interface for categories or tags.
// Can be used for single selection (category) or multi-selection (tags).
type SelectorView struct {
	title       string
	subtitle    string
	table       table.Model
	multiSelect bool
	selectedIDs map[int]bool
	width       int
	height      int
}

// NewSelectorView creates a new selector view.
// multiSelect: true for tags (can select multiple), false for category (single selection)
func NewSelectorView(title, subtitle string, multiSelect bool, width, height int) SelectorView {
	columns := []table.Column{
		{Title: "ID", Width: 8},
		{Title: "Name", Width: 40},
		{Title: "Count", Width: 10},
	}

	if multiSelect {
		columns = []table.Column{
			{Title: "✓", Width: 3},
			{Title: "ID", Width: 8},
			{Title: "Name", Width: 35},
			{Title: "Count", Width: 10},
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(height-10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("13")).
		Bold(false)
	t.SetStyles(s)

	return SelectorView{
		title:       title,
		subtitle:    subtitle,
		table:       t,
		multiSelect: multiSelect,
		selectedIDs: make(map[int]bool),
		width:       width,
		height:      height,
	}
}

// SetRows sets the data rows for the selector.
func (sv *SelectorView) SetRows(rows []table.Row) {
	sv.table.SetRows(rows)
}

// SetSelected pre-selects items (for editing existing snippet)
func (sv *SelectorView) SetSelected(ids []int) {
	sv.selectedIDs = make(map[int]bool)
	for _, id := range ids {
		sv.selectedIDs[id] = true
	}
}

// ToggleSelection toggles the selection state of the current row (multi-select only)
func (sv *SelectorView) ToggleSelection(id int) {
	if !sv.multiSelect {
		return
	}
	if sv.selectedIDs[id] {
		delete(sv.selectedIDs, id)
	} else {
		sv.selectedIDs[id] = true
	}
}

// GetSelectedID returns the currently focused ID (single-select only)
func (sv SelectorView) GetSelectedID() int {
	if sv.multiSelect || len(sv.table.Rows()) == 0 {
		return 0
	}

	selected := sv.table.SelectedRow()
	if len(selected) < 1 {
		return 0
	}

	var id int
	fmt.Sscanf(selected[0], "%d", &id)
	return id
}

// GetSelectedIDs returns all selected IDs (multi-select only)
func (sv SelectorView) GetSelectedIDs() []int {
	if !sv.multiSelect {
		return nil
	}

	ids := []int{}
	for id, selected := range sv.selectedIDs {
		if selected {
			ids = append(ids, id)
		}
	}

	sort.Ints(ids)
	return ids
}

// IsSelected checks if an ID is selected
func (sv SelectorView) IsSelected(id int) bool {
	return sv.selectedIDs[id]
}

// Update handles input for the selector.
func (sv *SelectorView) Update(msg tea.Msg) (SelectorView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ": // Use " " instead of "space" to catch the space key
			if sv.multiSelect {
				selected := sv.table.SelectedRow()
				if len(selected) >= 2 { // Changed from >= 3 to >= 2
					var id int
					n, err := fmt.Sscanf(selected[1], "%d", &id)

					if err == nil && n == 1 && id > 0 {
						sv.ToggleSelection(id)
						sv.UpdateCheckboxDisplay()
					}
				}
				// Return early to prevent the table from processing the space key
				return *sv, nil
			}
		}
	}

	// Only update the table if we didn't handle the key above
	var cmd tea.Cmd
	sv.table, cmd = sv.table.Update(msg)
	return *sv, cmd
}

// UpdateCheckboxDisplay updates the checkbox column for all rows
func (sv *SelectorView) UpdateCheckboxDisplay() {
	if !sv.multiSelect {
		return
	}

	rows := sv.table.Rows()
	newRows := make([]table.Row, len(rows))

	for i, row := range rows {
		if len(row) < 3 {
			newRows[i] = row
			continue
		}

		var id int
		fmt.Sscanf(row[1], "%d", &id)

		newRow := make(table.Row, len(row))
		copy(newRow, row)

		if sv.IsSelected(id) {
			newRow[0] = "✓"
		} else {
			newRow[0] = " "
		}

		newRows[i] = newRow
	}

	sv.table.SetRows(newRows)
}

// View renders the selector.
func (sv SelectorView) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13")).
		MarginBottom(1)

	b.WriteString(titleStyle.Render(sv.title))
	b.WriteString("\n")

	if sv.subtitle != "" {
		subtitleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			MarginBottom(1)
		b.WriteString(subtitleStyle.Render(sv.subtitle))
		b.WriteString("\n")
	}

	if sv.multiSelect {
		selectedIDs := []int{}
		for id, selected := range sv.selectedIDs {
			if selected {
				selectedIDs = append(selectedIDs, id)
			}
		}
		sort.Ints(selectedIDs)

		countStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)
		b.WriteString(countStyle.Render(fmt.Sprintf("Selected: %d", len(selectedIDs))))
		b.WriteString("\n\n")
	}

	b.WriteString(sv.table.View())
	b.WriteString("\n\n")

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	if sv.multiSelect {
		b.WriteString(footerStyle.Render(
			"↑↓: navigate | Space: toggle selection | Enter: confirm | Esc: cancel | n: clear all",
		))
	} else {
		b.WriteString(footerStyle.Render(
			"↑↓: navigate | Enter: select | Esc: cancel | n: clear selection",
		))
	}

	return b.String()
}
