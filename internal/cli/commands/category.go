// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/7-Dany/snip/internal/domain"
	"github.com/7-Dany/snip/internal/storage"
	"github.com/jedib0t/go-pretty/v6/table"
)

// CategoryCommand handles category-related operations.
type CategoryCommand struct {
	repos *storage.Repositories
}

// NewCategoryCommand creates a new CategoryCommand instance.
func NewCategoryCommand(repos *storage.Repositories) *CategoryCommand {
	return &CategoryCommand{repos: repos}
}

// manage routes category subcommands to the appropriate handler.
func (cc *CategoryCommand) manage(args []string) {
	if len(args) == 0 {
		PrintError("No subcommand provided. Use 'snip help category' for available commands")
		return
	}

	subcommand := strings.ToLower(args[0])
	subcommandArgs := args[1:]

	switch subcommand {
	case "list":
		cc.list()
	case "create":
		cc.create(subcommandArgs)
	case "delete":
		cc.delete(subcommandArgs)
	default:
		PrintError(fmt.Sprintf("Unknown command '%s'. Use 'snip help category' for available commands", args[0]))
	}
}

// list displays all categories in a formatted table.
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
	t.AppendHeader(table.Row{"ID", "Name", "Created", "Updated"})

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

// create creates a new category with the given name or prompts for input.
func (cc *CategoryCommand) create(args []string) {
	var name string

	if len(args) == 0 {
		name = promptForInput(
			"Create a new category",
			"✨",
			"Category name",
			"e.g., algorithms, web-dev, utils...",
			50,
			40,
		)
		if name == "" {
			PrintInfo("Create cancelled")
			return
		}
	} else {
		name = strings.TrimSpace(args[0])
		if name == "" {
			PrintError("category name cannot be empty")
			return
		}
	}

	// Check for duplicates
	existing, err := cc.repos.Categories.FindByName(name)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		PrintError(fmt.Sprintf("failed to check for existing category: %v", err))
		return
	}

	if existing != nil {
		PrintError("category already exists")
		return
	}

	// Create and save the category
	category, err := domain.NewCategory(name)
	if err != nil {
		PrintError(fmt.Sprintf("failed to create category: %v", err))
		return
	}

	if err := cc.repos.Categories.Create(category); err != nil {
		PrintError(fmt.Sprintf("failed to save category: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Created category '%s' (ID: %d)", name, category.ID()))
}

// delete removes a category after user confirmation.
func (cc *CategoryCommand) delete(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'id'. Use 'snip category delete <id>'")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		PrintError(fmt.Sprintf("Invalid ID '%s'. ID must be a number", args[0]))
		return
	}

	// Find the category to confirm deletion
	category, err := cc.repos.Categories.FindByID(id)
	if errors.Is(err, storage.ErrNotFound) {
		PrintError(fmt.Sprintf("Category with ID %d not found", id))
		return
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to find category: %v", err))
		return
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete category '%s'? (y/n): ", category.Name())
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		PrintInfo("Delete cancelled")
		return
	}

	// Delete the category
	if err := cc.repos.Categories.Delete(id); err != nil {
		PrintError(fmt.Sprintf("Failed to delete category: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Deleted category '%s' (ID: %d)", category.Name(), id))
}
