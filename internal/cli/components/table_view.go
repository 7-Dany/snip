package components

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableView represents a reusable table component with actions
type TableView struct {
	table      table.Model
	title      string
	emptyMsg   string
	actionHint string
	height     int
}

// NewTableView creates a new table view
func NewTableView(columns []table.Column, title, emptyMsg, actionHint string, height int) TableView {
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

	return TableView{
		table:      t,
		title:      title,
		emptyMsg:   emptyMsg,
		actionHint: actionHint,
		height:     height,
	}
}

func (tv *TableView) SetRows(rows []table.Row) {
	tv.table.SetRows(rows)
}

func (tv *TableView) SelectedRow() table.Row {
	return tv.table.SelectedRow()
}

func (tv *TableView) Update(msg tea.Msg) (TableView, tea.Cmd) {
	var cmd tea.Cmd
	tv.table, cmd = tv.table.Update(msg)
	return *tv, cmd
}

func (tv TableView) View() string {
	if len(tv.table.Rows()) == 0 {
		return tv.emptyMsg + "\n\n" +
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render(tv.actionHint)
	}

	var b string
	b += tv.table.View()
	b += "\n\n"
	b += lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(tv.actionHint)

	return b
}
