package main

import (
	"os"
	"time"

	"github.com/7-Dany/snip/internal/cli"
	"github.com/7-Dany/snip/internal/storage"
)

func main() {
	cli.PrintLoading("Loading Config...")
	config, err := cli.LoadConfig()
	if err != nil {
		cli.PrintError("Error loading config!" + err.Error())
		os.Exit(1)
	}
	cli.PrintSuccess("Config loaded")

	cli.PrintLoading("Loading Store...")
	repos := storage.New(config.StoragePath)
	err = repos.Load()
	if err != nil {
		cli.PrintError("Error loading repos!" + err.Error())
		os.Exit(1)
	}
	defer repos.Save()

	cli.PrintSuccess("Store loaded")
	time.Sleep(300 * time.Millisecond)
	cli.ClearScreen()

	cli.PrintLogo()

	cli.Run(os.Args, repos)
}
