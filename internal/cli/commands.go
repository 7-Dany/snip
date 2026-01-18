package cli

import (
	"github.com/7-Dany/snip/internal/storage"
)

func Run(args []string, repos *storage.Repositories) {
	if len(args) < 2 {
		PrintError("no command provided")
		PrintHelp("")
		return
	}

	PrintInfo(args[1])
}
