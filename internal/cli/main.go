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

func Run() {
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
	tagsTab := tui.NewTagsTab(repos)         // You'll need to implement this
	snippetsTab := tui.NewSnippetsTab(repos) // You'll need to implement this

	tabs, err := components.NewTabs(
		[]string{"Home", "Categories", "Tags", "Snippets"},
		[]components.TabModel{
			homeTab,
			categoriesTab,
			tagsTab,
			snippetsTab,
		},
	)
	if err != nil {
		fmt.Println("Error creating tabs:", err)
		os.Exit(1)
	}

	// Configure tabs
	tabs.SetLabelWidth(20)

	// Run with alt screen for clean display
	p := tea.NewProgram(tabs, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
