// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"fmt"
	"strings"

	"github.com/7-Dany/snip/internal/storage"
	"github.com/fatih/color"
)

// HelpCommand handles help-related operations.
type HelpCommand struct {
	repos *storage.Repositories
}

// NewHelpCommand creates a new HelpCommand instance.
func NewHelpCommand(repos *storage.Repositories) *HelpCommand {
	return &HelpCommand{repos: repos}
}

// manage routes help requests to the appropriate handler.
func (hc *HelpCommand) manage(args []string) {
	topic := ""
	if len(args) > 0 {
		topic = strings.ToLower(args[0])
	}
	hc.Print(topic)
}

// Print displays help information for the specified topic.
func (hc *HelpCommand) Print(topic string) {
	cyan := color.New(color.FgCyan, color.Bold)
	white := color.New(color.FgWhite)
	gray := color.New(color.FgHiBlack)

	switch topic {
	case "":
		hc.printGeneralHelp(cyan, white, gray)
	case "snippet":
		hc.printSnippetHelp(cyan, white, gray)
	case "category":
		hc.printCategoryHelp(cyan, white, gray)
	case "tag":
		hc.printTagHelp(cyan, white, gray)
	default:
		PrintError(fmt.Sprintf("Unknown help topic: %s", topic))
		fmt.Println("\nAvailable topics: snippet, category, tag")
	}
}

func (hc *HelpCommand) printGeneralHelp(cyan, white, gray *color.Color) {
	cyan.Println("\nUSAGE")
	white.Println("  snip <command> [arguments]")

	cyan.Println("\nCOMMANDS")
	white.Println("  Snippet Management:")
	fmt.Println("    snippet create                Create a new snippet interactively")
	fmt.Println("    snippet list [--flags]        List all snippets (optional filters)")
	fmt.Println("    snippet show <id>             Display a specific snippet")
	fmt.Println("    snippet update <id>           Update an existing snippet")
	fmt.Println("    snippet delete <id>           Delete a snippet")
	fmt.Println("    snippet search <query>        Search for snippets")

	white.Println("\n  Category Management:")
	fmt.Println("    category create [name]        Create a new category")
	fmt.Println("    category list                 List all categories")
	fmt.Println("    category delete <id>          Delete a category")

	white.Println("\n  Tag Management:")
	fmt.Println("    tag create [name]             Create a new tag")
	fmt.Println("    tag list                      List all tags")
	fmt.Println("    tag delete <id>               Delete a tag")

	white.Println("\n  Other:")
	fmt.Println("    help [topic]                  Show help for a specific topic")

	cyan.Println("\nEXAMPLES")
	fmt.Println("  snip snippet create                       # Interactive snippet creation")
	fmt.Println("  snip snippet list --language go           # List all Go snippets")
	fmt.Println("  snip snippet list --category 1            # List snippets in category 1")
	fmt.Println("  snip snippet search \"quicksort\"           # Search for 'quicksort'")
	fmt.Println("  snip category create algorithms           # Create 'algorithms' category")
	fmt.Println("  snip tag create                           # Interactive tag creation")

	gray.Println("\nFor more information on a specific command, use: snip help <topic>")
	fmt.Println()
}

func (hc *HelpCommand) printSnippetHelp(cyan, white, gray *color.Color) {
	cyan.Println("\nSNIPPET COMMANDS")

	white.Println("\n  snippet create")
	fmt.Println("    Create a new code snippet using an interactive form.")
	gray.Println("    Usage: snip snippet create")

	white.Println("\n  snippet list [--flags]")
	fmt.Println("    List all snippets or filter by category, tag, or language.")
	gray.Println("    Usage: snip snippet list [--category <id>] [--tag <id>] [--language <lang>]")
	gray.Println("    Examples:")
	gray.Println("      snip snippet list")
	gray.Println("      snip snippet list --category 1")
	gray.Println("      snip snippet list --language python")

	white.Println("\n  snippet show <id>")
	fmt.Println("    Display the full details of a specific snippet including code.")
	gray.Println("    Usage: snip snippet show <id>")
	gray.Println("    Example: snip snippet show 5")

	white.Println("\n  snippet update <id>")
	fmt.Println("    Update an existing snippet using an interactive form.")
	gray.Println("    Usage: snip snippet update <id>")
	gray.Println("    Example: snip snippet update 5")

	white.Println("\n  snippet delete <id>")
	fmt.Println("    Delete a snippet after confirmation.")
	gray.Println("    Usage: snip snippet delete <id>")
	gray.Println("    Example: snip snippet delete 5")

	white.Println("\n  snippet search <query>")
	fmt.Println("    Search snippets by title, description, code, or language.")
	gray.Println("    Usage: snip snippet search <query>")
	gray.Println("    Example: snip snippet search \"binary tree\"")

	fmt.Println()
}

func (hc *HelpCommand) printCategoryHelp(cyan, white, gray *color.Color) {
	cyan.Println("\nCATEGORY COMMANDS")

	white.Println("\n  category create [name]")
	fmt.Println("    Create a new category. Provides interactive prompt if name not given.")
	gray.Println("    Usage: snip category create [name]")
	gray.Println("    Examples:")
	gray.Println("      snip category create algorithms")
	gray.Println("      snip category create              # Interactive mode")

	white.Println("\n  category list")
	fmt.Println("    Display all available categories.")
	gray.Println("    Usage: snip category list")

	white.Println("\n  category delete <id>")
	fmt.Println("    Delete a category after confirmation.")
	gray.Println("    Usage: snip category delete <id>")
	gray.Println("    Example: snip category delete 3")

	fmt.Println()
}

func (hc *HelpCommand) printTagHelp(cyan, white, gray *color.Color) {
	cyan.Println("\nTAG COMMANDS")

	white.Println("\n  tag create [name]")
	fmt.Println("    Create a new tag. Provides interactive prompt if name not given.")
	gray.Println("    Usage: snip tag create [name]")
	gray.Println("    Examples:")
	gray.Println("      snip tag create performance")
	gray.Println("      snip tag create                   # Interactive mode")

	white.Println("\n  tag list")
	fmt.Println("    Display all available tags.")
	gray.Println("    Usage: snip tag list")

	white.Println("\n  tag delete <id>")
	fmt.Println("    Delete a tag after confirmation.")
	gray.Println("    Usage: snip tag delete <id>")
	gray.Println("    Example: snip tag delete 7")

	fmt.Println()
}
