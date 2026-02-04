// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"
)

func TestHelpCommand_manage(t *testing.T) {
	repos := setupTestRepos(t)
	hc := NewHelpCommand(repos)

	t.Run("shows general help with no arguments", func(t *testing.T) {
		// Should not panic
		hc.manage([]string{})
	})

	t.Run("shows snippet help", func(t *testing.T) {
		// Should not panic
		hc.manage([]string{"snippet"})
	})

	t.Run("shows category help", func(t *testing.T) {
		// Should not panic
		hc.manage([]string{"category"})
	})

	t.Run("shows tag help", func(t *testing.T) {
		// Should not panic
		hc.manage([]string{"tag"})
	})

	t.Run("handles case insensitive topics", func(t *testing.T) {
		// Should work with different cases
		hc.manage([]string{"SNIPPET"})
		hc.manage([]string{"Category"})
		hc.manage([]string{"tAg"})
	})
}

func TestHelpCommand_Print(t *testing.T) {
	repos := setupTestRepos(t)
	hc := NewHelpCommand(repos)

	t.Run("prints general help", func(t *testing.T) {
		// Should not panic
		hc.Print("")
	})

	t.Run("prints snippet help", func(t *testing.T) {
		// Should not panic
		hc.Print("snippet")
	})

	t.Run("prints category help", func(t *testing.T) {
		// Should not panic
		hc.Print("category")
	})

	t.Run("prints tag help", func(t *testing.T) {
		// Should not panic
		hc.Print("tag")
	})

	t.Run("handles unknown topic", func(t *testing.T) {
		// Should show error, not panic
		hc.Print("unknown")
	})
}

func TestNewHelpCommand(t *testing.T) {
	t.Run("creates help command with repos", func(t *testing.T) {
		repos := setupTestRepos(t)
		hc := NewHelpCommand(repos)

		if hc == nil {
			t.Fatal("NewHelpCommand returned nil")
		}

		if hc.repos == nil {
			t.Error("repos is nil")
		}
	})
}
