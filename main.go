package main

import (
	"os"

	"github.com/7-Dany/snip/internal/cli"
	"github.com/7-Dany/snip/internal/storage"
)

func main() {
	config, err := cli.LoadConfig()
	if err != nil {
		cli.PrintError("Error loading config!" + err.Error())
		os.Exit(1)
	}

	repos := storage.New(config.StoragePath)
	err = repos.Load()
	if err != nil {
		cli.PrintError("Error loading repos!" + err.Error())
		os.Exit(1)
	}
	defer repos.Save()

	cli.PrintLogo()

	app := cli.NewCLI(repos)
	app.Run(os.Args)
}
