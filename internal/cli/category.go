// Package cli provides command-line interface functionality for the snippet manager.
// This file implements category management commands with interactive TUI support.
package cli

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

// CategoryCommand handles all category-related CLI operations.
// This type wraps the storage repositories and provides command routing
// for category management.
//
// Design: Command pattern with dependency injection
// - Wraps *storage.Repositories for data access
// - Routes subcommands to handler methods
// - Supports both interactive (TUI) and argument-based modes
//
// Supported operations:
//   - list: Display all categories in formatted table
//   - create: Add new category (interactive or via argument)
//   - delete: Remove category with confirmation prompt
type CategoryCommand struct {
	repos *storage.Repositories
}

// NewCategoryCommand creates a new CategoryCommand with the given repositories.
// This is the public constructor used by the CLI coordinator.
//
// Parameters:
//
//	repos - Storage repositories for data access
func NewCategoryCommand(repos *storage.Repositories) *CategoryCommand {
	return &CategoryCommand{repos: repos}
}

// manage routes category subcommands to their respective handlers.
// This is the main entry point called by the CLI coordinator.
//
// Performance: O(1) for routing, actual performance depends on handler
//
// Parameters:
//
//	args - Command arguments (first element is subcommand)
//
// Supported subcommands:
//   - "list": List all categories
//   - "create": Create new category
//   - "delete": Delete category by ID
//
// Error handling:
//   - Shows error if no subcommand provided
//   - Shows error for unknown subcommands
//   - Case-insensitive subcommand matching
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

// list retrieves and displays all categories in a formatted table.
// Uses go-pretty/table with StyleColoredBright for visual output.
//
// Performance: O(n) where n = total categories
//
// Output format:
//
//	Colored table with columns: ID, Name, CreatedAt, UpdatedAt
//	Shows helpful message if no categories exist
//
// Error handling:
//   - Prints error if repository access fails
//   - Shows info message for empty result set
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

// create adds a new category to the repository.
// Supports both interactive (TUI) and argument-based modes.
//
// Performance: O(n) for duplicate check, O(1) for creation
//
// Parameters:
//
//	args - Optional category name as first element
//
// Behavior:
//   - If args empty: Launch interactive Bubbletea prompt
//   - If args[0] provided: Use as category name
//
// Validation:
//   - Checks for duplicate names (case-insensitive)
//   - Validates domain constraints via domain.NewCategory()
//
// Error handling:
//   - Shows error for duplicate names
//   - Shows error if creation fails
//   - Shows success message with category name
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

// delete removes a category by ID after user confirmation.
// Always prompts for confirmation before deletion.
//
// Performance: O(1) for lookup and deletion
//
// Parameters:
//
//	args - Category ID as first element (must be numeric string)
//
// Validation:
//   - Requires ID argument
//   - Validates ID is numeric
//   - Checks category exists before prompting
//
// User interaction:
//   - Displays category name in confirmation prompt
//   - Accepts 'y' (case-insensitive) to confirm
//   - Any other input cancels deletion
//
// Error handling:
//   - Shows error for missing/invalid ID
//   - Shows error if category not found
//   - Shows info message if cancelled
//   - Shows success message if deleted
//
// Note: Does not cascade delete snippets in this category.
// Snippets will have category_id but category won't exist (orphaned reference).
// CLI layer should handle this when displaying snippets.
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

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
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

// categoryInputModel is the Bubbletea model for interactive category name input.
// Implements tea.Model interface (Init, Update, View).
//
// Design: Single-field form using bubbles/textinput
// - Focused by default for immediate typing
// - Shows placeholder text for guidance
// - Character limit enforced (50 chars)
type categoryInputModel struct {
	textInput textinput.Model
	cancelled bool
}

// newCategoryInputModel initializes a new category input form with default settings.
//
// Configuration:
//   - Placeholder: "e.g., algorithms, web-dev, utils..."
//   - Auto-focused for immediate input
//   - Character limit: 50
//   - Display width: 40 characters
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

// Init initializes the Bubbletea program and returns the blink command.
// Required by tea.Model interface.
//
// Returns:
//
//	textinput.Blink - Makes cursor blink for better UX
func (m categoryInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// promptForName launches an interactive Bubbletea TUI for category name input.
// This provides a user-friendly alternative to command-line arguments.
//
// Performance: Blocks until user input or cancellation
//
// Returns:
//
//	Trimmed category name, or empty string if cancelled/empty
//
// User interaction:
//   - Enter: Confirm input
//   - Esc/Ctrl+C: Cancel (returns empty string)
//   - Character limit: 50 chars
//
// Note: Returns empty string for both cancellation and empty input.
// Caller should handle empty string appropriately.
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

// Update handles keyboard events and updates the model state.
// Required by tea.Model interface.
//
// Keyboard handling:
//   - Enter: Accept input and quit
//   - Esc/Ctrl+C: Cancel and quit
//   - Other keys: Handled by textinput component
//
// Returns:
//
//	Updated model and command (tea.Quit for termination)
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

// View renders the current state of the input form as a string.
// Required by tea.Model interface.
//
// Output format:
//
//	Title with emoji
//	Input field with current value
//	Help text for keyboard shortcuts
//
// Returns:
//
//	Formatted string for terminal display
func (m categoryInputModel) View() string {
	return fmt.Sprintf(
		"\n✨ Create a new category\n\n"+
			"Category name: %s\n\n"+
			"(press Enter to confirm, Esc to cancel)\n",
		m.textInput.View(),
	)
}
