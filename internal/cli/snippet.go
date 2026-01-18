package cli

import (
	"fmt"

	"github.com/7-Dany/snip/internal/storage"
)

type SnippetCommand struct {
	repos *storage.Repositories
}

func NewSnippetCommand(repos *storage.Repositories) *SnippetCommand {
	return &SnippetCommand{repos: repos}
}

func (sc *SnippetCommand) manage(args []string) {
	fmt.Println("Managing Snippets")
}
