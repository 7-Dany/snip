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

type TagCommand struct {
	repos *storage.Repositories
}

func NewTagCommand(repos *storage.Repositories) *TagCommand {
	return &TagCommand{repos: repos}
}

func (tc *TagCommand) manage(args []string) {
	if len(args) == 0 {
		PrintError("No subcommand provided. Use 'snip help tag' for available commands")
		return
	}

	switch strings.ToLower(args[0]) {
	case "list":
		tc.list()
	case "create":
		tc.create(args[1:])
	case "delete":
		tc.delete(args[1:])
	default:
		PrintError(fmt.Sprintf("Unknown command '%s'. Use 'snip help tag' for available commands", args[0]))
	}
}

func (tc *TagCommand) list() {
	tags, err := tc.repos.Tags.List()
	if err != nil {
		PrintError(fmt.Sprintf("failed to list tags: %v", err))
		return
	}

	if len(tags) == 0 {
		PrintInfo("no tags found, create one with 'snip tag create'")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "CreatedAt", "UpdatedAt"})

	for _, tag := range tags {
		t.AppendRow(table.Row{
			tag.ID(),
			tag.Name(),
			tag.CreatedAt().Format("2006-01-02 15:04"),
			tag.UpdatedAt().Format("2006-01-02 15:04"),
		})
	}

	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

// create supports both argument and interactive modes
func (tc *TagCommand) create(args []string) {
	var name string
	if len(args) == 0 {
		name = tc.promptForName()
		if name == "" {
			PrintInfo("Create cancelled")
			return
		}
	} else {
		name = args[0]
	}

	// Check for duplicates
	existed, err := tc.repos.Tags.FindByName(name)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		PrintError(fmt.Sprintf("Internal error, failed to check tag by name: %v", err))
		return
	}

	if existed != nil {
		PrintError("tag already exists")
		return
	}

	tag, err := domain.NewTag(name)
	if err != nil {
		PrintError(fmt.Sprintf("failed to create tag: %v", err))
		return
	}

	err = tc.repos.Tags.Create(tag)
	if err != nil {
		PrintError(fmt.Sprintf("failed to register tag to storage: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Successfully created tag name: '%v'", name))
}

func (tc *TagCommand) delete(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'id'. Use 'snip tag delete <id>'")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		PrintError(fmt.Sprintf("Validation failed. 'id' must be a number: %v", err))
		return
	}

	tag, err := tc.repos.Tags.FindByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		PrintError(fmt.Sprintf("Tag with ID %d not found", id))
		return
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to find tag: %v", err))
		return
	}

	fmt.Printf("Are you sure you want to delete tag '%s'? (y/n): ", tag.Name())
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" {
		PrintInfo("Delete cancelled")
		return
	}

	err = tc.repos.Tags.Delete(id)
	if err != nil {
		PrintError(fmt.Sprintf("Delete failed, couldn't delete tag: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Deleted tag with id: '%v'", id))
}

// Interactive prompt
type tagInputModel struct {
	textInput textinput.Model
	cancelled bool
}

func newTagInputModel() tagInputModel {
	ti := textinput.New()
	ti.Placeholder = "e.g., performance, security, beginner..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 40

	return tagInputModel{
		textInput: ti,
	}
}

func (m tagInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (tc *TagCommand) promptForName() string {
	p := tea.NewProgram(newTagInputModel())
	finalModel, err := p.Run()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to run input prompt: %v", err))
		return ""
	}

	m := finalModel.(tagInputModel)
	if m.cancelled {
		return ""
	}

	value := strings.TrimSpace(m.textInput.Value())
	if value == "" {
		PrintError("Tag name cannot be empty")
		return ""
	}

	return value
}

func (m tagInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m tagInputModel) View() string {
	return fmt.Sprintf(
		"\n🏷️  Create a new tag\n\n"+
			"Tag name: %s\n\n"+
			"(press Enter to confirm, Esc to cancel)\n",
		m.textInput.View(),
	)
}
