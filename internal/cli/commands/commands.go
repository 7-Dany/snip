package commands

import (
	"github.com/7-Dany/snip/internal/storage"
)

type CLI struct {
	snippet  *SnippetCommand
	category *CategoryCommand
	tag      *TagCommand
	help     *HelpCommand
}

func NewCLI(repos *storage.Repositories) *CLI {
	return &CLI{
		snippet:  NewSnippetCommand(repos),
		category: NewCategoryCommand(repos),
		tag:      NewTagCommand(repos),
		help:     NewHelpCommand(repos),
	}
}

func (cli *CLI) Run(args []string) {
	if len(args) < 2 {
		PrintError("no command provided")
		cli.help.Print("")
		return
	}

	switch args[1] {
	case "help":
		cli.help.manage(args[2:])
	case "snippet":
		cli.snippet.manage(args[2:])
	case "category":
		cli.category.manage(args[2:])
	case "tag":
		cli.tag.manage(args[2:])
	default:
		cli.snippet.manage(args[1:])
	}
}
