package cli

import (
	"fmt"
	"strings"

	"github.com/7-Dany/snip/internal/storage"
	"github.com/fatih/color"
)

type HelpCommand struct {
	repos *storage.Repositories
}

func NewHelpCommand(repos *storage.Repositories) *HelpCommand {
	return &HelpCommand{repos: repos}
}

func (hc *HelpCommand) manage(args []string) {
	if len(args) == 0 {
		hc.Print("")
		return
	}

	hc.Print(strings.ToLower(args[0]))
}

// PrintHelp displays general help or specific command help
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
	fmt.Println("    create                    Create a new snippet interactively")
	fmt.Println("    list                      Display all snippets in a table")
	fmt.Println("    search <query>            Search for snippets")
	fmt.Println("    show <id>                 Display a specific snippet")
	fmt.Println("    update <id>               Update an existing snippet")
	fmt.Println("    delete <id>               Delete a snippet")

	white.Println("\n  Category Management:")
	fmt.Println("    category create           Create a new category")
	fmt.Println("    category list             List all categories")
	fmt.Println("    category delete <id>      Delete a category")

	white.Println("\n  Tag Management:")
	fmt.Println("    tag create                Create a new tag")
	fmt.Println("    tag list                  List all tags")
	fmt.Println("    tag delete <id>           Delete a tag")

	white.Println("\n  Other:")
	fmt.Println("    help [topic]              Show help for a specific topic")

	cyan.Println("\nEXAMPLES")
	fmt.Println("  snip create                        # Start interactive snippet creation")
	fmt.Println("  snip search \"quicksort\"            # Search for snippets containing 'quicksort'")
	fmt.Println("  snip show 5                        # Display snippet with ID 5")
	fmt.Println("  snip help snippet                  # Show detailed help for snippet commands")

	gray.Println("\nFor more information on a specific command, use: snip help <topic>")
	fmt.Println()
}

func (hc *HelpCommand) printSnippetHelp(cyan, white, gray *color.Color) {
	cyan.Println("\nSNIPPET COMMANDS")

	white.Println("\n  create")
	fmt.Println("    Create a new code snippet interactively using a form interface.")
	gray.Println("    Usage: snip create")

	white.Println("\n  list")
	fmt.Println("    Display all snippets in a formatted table.")
	gray.Println("    Usage: snip list")

	white.Println("\n  search <query>")
	fmt.Println("    Search snippets by title, description, code, or language.")
	gray.Println("    Usage: snip search <query>")
	gray.Println("    Example: snip search \"binary tree\"")

	white.Println("\n  show <id>")
	fmt.Println("    Display the full details of a specific snippet.")
	gray.Println("    Usage: snip show <id>")
	gray.Println("    Example: snip show 5")

	white.Println("\n  update <id>")
	fmt.Println("    Update an existing snippet using an interactive form.")
	gray.Println("    Usage: snip update <id>")
	gray.Println("    Example: snip update 5")

	white.Println("\n  delete <id>")
	fmt.Println("    Delete a snippet with confirmation prompt.")
	gray.Println("    Usage: snip delete <id>")
	gray.Println("    Example: snip delete 5")

	fmt.Println()
}

func (hc *HelpCommand) printCategoryHelp(cyan, white, gray *color.Color) {
	cyan.Println("\nCATEGORY COMMANDS")

	white.Println("\n  category create")
	fmt.Println("    Create a new category for organizing snippets.")
	gray.Println("    Usage: snip category create <name>")
	gray.Println("    Example: snip category create C++")

	white.Println("\n  category list")
	fmt.Println("    Display all available categories.")
	gray.Println("    Usage: snip category list")

	white.Println("\n  category delete <id>")
	fmt.Println("    Delete a category. Snippets in this category will have no category.")
	gray.Println("    Usage: snip category delete <id>")
	gray.Println("    Example: snip category delete 3")

	fmt.Println()
}

func (hc *HelpCommand) printTagHelp(cyan, white, gray *color.Color) {
	cyan.Println("\nTAG COMMANDS")

	white.Println("\n  tag create")
	fmt.Println("    Create a new tag for labeling snippets.")
	gray.Println("    Usage: snip tag create")

	white.Println("\n  tag list")
	fmt.Println("    Display all available tags.")
	gray.Println("    Usage: snip tag list")

	white.Println("\n  tag delete <id>")
	fmt.Println("    Delete a tag. This will remove the tag from all snippets.")
	gray.Println("    Usage: snip tag delete <id>")
	gray.Println("    Example: snip tag delete 7")

	fmt.Println()
}
