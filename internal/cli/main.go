package cli

import (
	"fmt"
	"os"

	"github.com/7-Dany/snip/internal/cli/commands"
	"github.com/7-Dany/snip/internal/cli/components"
	"github.com/7-Dany/snip/internal/cli/config"
	"github.com/7-Dany/snip/internal/cli/tui"
	"github.com/7-Dany/snip/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func CLIMain() {
	config, err := config.LoadConfig()
	if err != nil {
		commands.PrintError("Error loading config!" + err.Error())
		os.Exit(1)
	}

	repos := storage.New(config.StoragePath)
	err = repos.Load()
	if err != nil {
		commands.PrintError("Error loading repos!" + err.Error())
		os.Exit(1)
	}
	defer repos.Save()

	app := commands.NewCLI(repos)

	// If arguments provided, use old CLI
	if len(os.Args) > 1 {
		app.Run(os.Args)
		return
	}

	// Otherwise, launch TUI
	homeTab := tui.NewHomeTab()
	categoriesTab := tui.NewCategoriesTab(repos)

	tabs, err := components.NewTabs(
		[]string{"Home", "Categories", "Tags", "Snippets"},
		[]components.TabModel{
			homeTab,
			categoriesTab,
			homeTab, // Placeholder for Tags
			homeTab, // Placeholder for Snippets
		},
	)
	if err != nil {
		fmt.Println("Error creating tabs:", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(tabs).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
