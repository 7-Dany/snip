// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"testing"
)

// TestOutputFunctions tests that output functions don't panic.
// We can't easily test the actual output without capturing stdout,
// but we can ensure they execute without errors.

func TestPrintLogo(t *testing.T) {
	t.Run("prints logo without panic", func(t *testing.T) {
		// Should not panic
		PrintLogo()
	})
}

func TestPrintLoading(t *testing.T) {
	t.Run("prints loading message without panic", func(t *testing.T) {
		// Should not panic
		PrintLoading("Loading...")
	})
}

func TestPrintSuccess(t *testing.T) {
	t.Run("prints success message without panic", func(t *testing.T) {
		// Should not panic
		PrintSuccess("Operation successful")
	})
}

func TestPrintError(t *testing.T) {
	t.Run("prints error message without panic", func(t *testing.T) {
		// Should not panic
		PrintError("Something went wrong")
	})
}

func TestPrintInfo(t *testing.T) {
	t.Run("prints info message without panic", func(t *testing.T) {
		// Should not panic
		PrintInfo("Information message")
	})
}

func TestPrintCommand(t *testing.T) {
	t.Run("prints command without panic", func(t *testing.T) {
		// Should not panic
		PrintCommand("snippet list")
	})
}

func TestClearScreen(t *testing.T) {
	t.Run("clears screen without panic", func(t *testing.T) {
		// Should not panic
		ClearScreen()
	})
}
