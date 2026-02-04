// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"
)

func TestNewCLI(t *testing.T) {
	t.Run("creates CLI with all command handlers", func(t *testing.T) {
		repos := setupTestRepos(t)
		cli := NewCLI(repos)

		if cli == nil {
			t.Fatal("NewCLI returned nil")
		}

		if cli.snippet == nil {
			t.Error("snippet command handler is nil")
		}

		if cli.category == nil {
			t.Error("category command handler is nil")
		}

		if cli.tag == nil {
			t.Error("tag command handler is nil")
		}

		if cli.help == nil {
			t.Error("help command handler is nil")
		}
	})
}

func TestCLI_Run(t *testing.T) {
	repos := setupTestRepos(t)
	cli := NewCLI(repos)

	t.Run("shows error when no command provided", func(t *testing.T) {
		// Should show error, not panic
		cli.Run([]string{"snip"})
	})

	t.Run("routes help command", func(t *testing.T) {
		// Should not panic
		cli.Run([]string{"snip", "help"})
	})

	t.Run("routes snippet command", func(t *testing.T) {
		// Should not panic
		cli.Run([]string{"snip", "snippet", "list"})
	})

	t.Run("routes category command", func(t *testing.T) {
		// Should not panic
		cli.Run([]string{"snip", "category", "list"})
	})

	t.Run("routes tag command", func(t *testing.T) {
		// Should not panic
		cli.Run([]string{"snip", "tag", "list"})
	})

	t.Run("routes unknown command to snippet handler", func(t *testing.T) {
		// Should treat as snippet command for backward compatibility
		cli.Run([]string{"snip", "list"})
	})
}
