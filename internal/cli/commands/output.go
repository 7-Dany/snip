// Package commands provides the CLI command handlers for SNIP.
package commands

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// PrintLogo displays the application logo with branding.
func PrintLogo() {
	magenta := color.New(color.FgMagenta, color.Bold)
	yellow := color.New(color.FgYellow)

	ghost := []string{
		"      ████████      ",
		"    ██        ██    ",
		"  ██   ██  ██   ██  ",
		"  ██   ██  ██   ██  ",
		"  ██            ██  ",
		"  ██  ████████  ██  ",
		"  ████        ████  ",
	}

	snip := []string{
		"   _____ _   _ _____ _____  ",
		"  / ____| \\ | |_   _|  __ \\ ",
		" | (___ |  \\| | | | | |__) |",
		"  \\___ \\| . ` | | | |  ___/ ",
		"  ____) | |\\  |_| |_| |     ",
		" |_____/|_| \\_|_____|_|     ",
		"  Code Snippet Manager v1.0 ",
	}

	fmt.Println()
	for i := range snip {
		if i < len(ghost) {
			magenta.Print(ghost[i])
		} else {
			magenta.Print("                    ")
		}
		magenta.Print("    ")
		magenta.Println(snip[i])
	}

	yellow.Println("\n    ✨ Manage your code snippets with style! ✨")
	fmt.Println()
}

// PrintLoading displays a loading message with animation.
func PrintLoading(msg string) {
	yellow := color.New(color.FgYellow)
	yellow.Printf("⚡ %s", msg)
	time.Sleep(150 * time.Millisecond)
	fmt.Println()
}

// PrintSuccess displays a success message in green.
func PrintSuccess(msg string) {
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✓ %s\n", msg)
}

// PrintError displays an error message in red.
func PrintError(msg string) {
	red := color.New(color.FgRed, color.Bold)
	red.Printf("✗ %s\n", msg)
}

// PrintInfo displays an informational message in cyan.
func PrintInfo(msg string) {
	cyan := color.New(color.FgCyan)
	cyan.Println(msg)
}

// PrintCommand displays the command being executed.
func PrintCommand(command string) {
	gray := color.New(color.FgHiWhite)
	gray.Printf("\n→ Executing command: %s\n\n", command)
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
