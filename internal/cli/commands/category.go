package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/7-Dany/snip/internal/domain"
	"github.com/7-Dany/snip/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedib0t/go-pretty/v6/table"
)

type CategoryCommand struct {
	repos *storage.Repositories
}

func NewCategoryCommand(repos *storage.Repositories) *CategoryCommand {
	return &CategoryCommand{repos: repos}
}

func (cc *CategoryCommand) manage(args []string) {
	if len(args) == 0 {
		PrintError("No subcommand provided. Use 'snip help category' for available commands")
		return
	}

	switch strings.ToLower(args[0]) {
	case "list":
		cc.list()
	case "create":
		cc.create(args[1:])
	case "delete":
		cc.delete(args[1:])
	default:
		PrintError(fmt.Sprintf("Unknown command '%s'. Use 'snip help category' for available commands", args[0]))
	}
}

func (cc *CategoryCommand) list() {
	categories, err := cc.repos.Categories.List()
	if err != nil {
		PrintError(fmt.Sprintf("failed to list categories: %v", err))
		return
	}

	if len(categories) == 0 {
		PrintInfo("no categories found, create one with 'snip category create'")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "CreatedAt", "UpdatedAt"})

	for _, category := range categories {
		t.AppendRow(table.Row{
			category.ID(),
			category.Name(),
			category.CreatedAt().Format("2006-01-02 15:04"),
			category.UpdatedAt().Format("2006-01-02 15:04"),
		})
	}

	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

// create supports both argument and interactive modes
func (cc *CategoryCommand) create(args []string) {
	var name string
	if len(args) == 0 {
		name = cc.promptForName()
		if name == "" {
			PrintInfo("Create cancelled")
			return
		}
	} else {
		name = args[0]
	}

	// Check for duplicates
	existed, err := cc.repos.Categories.FindByName(name)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		PrintError(fmt.Sprintf("Internal error, failed to check category by name: %v", err))
		return
	}

	if existed != nil {
		PrintError("category already exists")
		return
	}

	category, err := domain.NewCategory(name)
	if err != nil {
		PrintError(fmt.Sprintf("failed to create category: %v", err))
		return
	}

	err = cc.repos.Categories.Create(category)
	if err != nil {
		PrintError(fmt.Sprintf("failed to register category to storage: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Successfully created category name: '%v'", name))
}

func (cc *CategoryCommand) delete(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'id'. Use 'snip category delete <id>'")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		PrintError(fmt.Sprintf("Validation failed. 'id' must be a number: %v", err))
		return
	}

	category, err := cc.repos.Categories.FindByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		PrintError(fmt.Sprintf("Category with ID %d not found", id))
		return
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to find category: %v", err))
		return
	}

	fmt.Printf("Are you sure you want to delete category '%s'? (y/n): ", category.Name())
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" {
		PrintInfo("Delete cancelled")
		return
	}

	err = cc.repos.Categories.Delete(id)
	if err != nil {
		PrintError(fmt.Sprintf("Delete failed, couldn't delete category: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Deleted category with id: '%v'", id))
}

// Interactive prompt
type categoryInputModel struct {
	textInput textinput.Model
	cancelled bool
}

func newCategoryInputModel() categoryInputModel {
	ti := textinput.New()
	ti.Placeholder = "e.g., algorithms, web-dev, utils..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 40

	return categoryInputModel{
		textInput: ti,
	}
}

func (m categoryInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cc *CategoryCommand) promptForName() string {
	p := tea.NewProgram(newCategoryInputModel())
	finalModel, err := p.Run()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to run input prompt: %v", err))
		return ""
	}

	m := finalModel.(categoryInputModel)
	if m.cancelled {
		return ""
	}

	value := strings.TrimSpace(m.textInput.Value())
	if value == "" {
		PrintError("Category name cannot be empty")
		return ""
	}

	return value
}

func (m categoryInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m categoryInputModel) View() string {
	return fmt.Sprintf(
		"\n✨ Create a new category\n\n"+
			"Category name: %s\n\n"+
			"(press Enter to confirm, Esc to cancel)\n",
		m.textInput.View(),
	)
}
