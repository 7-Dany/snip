// Package commands provides the CLI command handlers for SNIP.
// It implements the command-line interface for managing snippets, categories, and tags.
package commands

import (
	"github.com/7-Dany/snip/internal/storage"
)

// CLI coordinates all command handlers and provides the main entry point
// for command execution.
type CLI struct {
	snippet  *SnippetCommand
	category *CategoryCommand
	tag      *TagCommand
	help     *HelpCommand
}

// NewCLI creates a new CLI instance with all command handlers initialized.
func NewCLI(repos *storage.Repositories) *CLI {
	return &CLI{
		snippet:  NewSnippetCommand(repos),
		category: NewCategoryCommand(repos),
		tag:      NewTagCommand(repos),
		help:     NewHelpCommand(repos),
	}
}

// Run executes the appropriate command based on the provided arguments.
// args[0] is expected to be the program name, args[1] is the command.
func (cli *CLI) Run(args []string) {
	if len(args) < 2 {
		PrintError("no command provided")
		cli.help.Print("")
		return
	}

	command := args[1]
	commandArgs := args[2:]

	switch command {
	case "help":
		cli.help.manage(commandArgs)
	case "snippet":
		cli.snippet.manage(commandArgs)
	case "category":
		cli.category.manage(commandArgs)
	case "tag":
		cli.tag.manage(commandArgs)
	default:
		// Assume it's a snippet command for backward compatibility
		cli.snippet.manage(args[1:])
	}
}
