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

// TagCommand handles tag-related operations.
type TagCommand struct {
	repos *storage.Repositories
}

// NewTagCommand creates a new TagCommand instance.
func NewTagCommand(repos *storage.Repositories) *TagCommand {
	return &TagCommand{repos: repos}
}

// manage routes tag subcommands to the appropriate handler.
func (tc *TagCommand) manage(args []string) {
	if len(args) == 0 {
		PrintError("No subcommand provided. Use 'snip help tag' for available commands")
		return
	}

	subcommand := strings.ToLower(args[0])
	subcommandArgs := args[1:]

	switch subcommand {
	case "list":
		tc.list()
	case "create":
		tc.create(subcommandArgs)
	case "delete":
		tc.delete(subcommandArgs)
	default:
		PrintError(fmt.Sprintf("Unknown command '%s'. Use 'snip help tag' for available commands", args[0]))
	}
}

// list displays all tags in a formatted table.
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
	t.AppendHeader(table.Row{"ID", "Name", "Created", "Updated"})

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

// create creates a new tag with the given name or prompts for input.
func (tc *TagCommand) create(args []string) {
	var name string

	if len(args) == 0 {
		name = promptForInput(
			"Create a new tag",
			"🏷️",
			"Tag name",
			"e.g., performance, security, beginner...",
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
			PrintError("tag name cannot be empty")
			return
		}
	}

	// Check for duplicates
	existing, err := tc.repos.Tags.FindByName(name)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		PrintError(fmt.Sprintf("failed to check for existing tag: %v", err))
		return
	}

	if existing != nil {
		PrintError("tag already exists")
		return
	}

	// Create and save the tag
	tag, err := domain.NewTag(name)
	if err != nil {
		PrintError(fmt.Sprintf("failed to create tag: %v", err))
		return
	}

	if err := tc.repos.Tags.Create(tag); err != nil {
		PrintError(fmt.Sprintf("failed to save tag: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Created tag '%s' (ID: %d)", name, tag.ID()))
}

// delete removes a tag after user confirmation.
func (tc *TagCommand) delete(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'id'. Use 'snip tag delete <id>'")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		PrintError(fmt.Sprintf("Invalid ID '%s'. ID must be a number", args[0]))
		return
	}

	// Find the tag to confirm deletion
	tag, err := tc.repos.Tags.FindByID(id)
	if errors.Is(err, storage.ErrNotFound) {
		PrintError(fmt.Sprintf("Tag with ID %d not found", id))
		return
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to find tag: %v", err))
		return
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete tag '%s'? (y/n): ", tag.Name())
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		PrintInfo("Delete cancelled")
		return
	}

	// Delete the tag
	if err := tc.repos.Tags.Delete(id); err != nil {
		PrintError(fmt.Sprintf("Failed to delete tag: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Deleted tag '%s' (ID: %d)", tag.Name(), id))
}
