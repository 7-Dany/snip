package commands

import (
	"fmt"

	"github.com/7-Dany/snip/internal/storage"
)

type TagCommand struct {
	repos *storage.Repositories
}

func NewTagCommand(repos *storage.Repositories) *TagCommand {
	return &TagCommand{repos: repos}
}

func (tc *TagCommand) manage(args []string) {
	fmt.Println("Managing Tag Commands")
}
