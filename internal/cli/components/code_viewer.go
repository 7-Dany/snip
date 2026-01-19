package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CodeViewer displays a snippet's code with syntax highlighting and line numbers.
type CodeViewer struct {
	title       string
	language    string
	description string
	code        string
	lineCount   int
	width       int
}

// NewCodeViewer creates a new code viewer.
func NewCodeViewer(title, language, description, code string, width int) CodeViewer {
	lines := strings.Split(code, "\n")
	return CodeViewer{
		title:       title,
		language:    language,
		description: description,
		code:        code,
		lineCount:   len(lines),
		width:       width,
	}
}

// View renders the code viewer.
func (cv CodeViewer) View() string {
	var b strings.Builder

	// Header styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("13")).
		Bold(true).
		Underline(true)

	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	langBadgeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("13")).
		Padding(0, 1).
		Bold(true)

	// Title
	b.WriteString(titleStyle.Render(cv.title))
	b.WriteString("  ")
	b.WriteString(langBadgeStyle.Render(cv.language))
	b.WriteString("\n\n")

	// Description
	if cv.description != "" {
		b.WriteString(metaStyle.Render(cv.description))
		b.WriteString("\n\n")
	}

	// Code box with full width
	codeWidth := cv.width - 8 // Leave margin for borders and padding
	if codeWidth < 40 {
		codeWidth = 40
	}

	codeStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(codeWidth)

	b.WriteString(codeStyle.Render(cv.renderCodeWithLineNumbers()))

	// Footer
	b.WriteString("\n\n")
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	b.WriteString(footerStyle.Render(
		"Lines: " +
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).Render(fmt.Sprintf("%d", cv.lineCount)) +
			" | Characters: " +
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).Render(fmt.Sprintf("%d", len(cv.code)))))

	return b.String()
}

// renderCodeWithLineNumbers adds line numbers to the code.
func (cv CodeViewer) renderCodeWithLineNumbers() string {
	lines := strings.Split(cv.code, "\n")
	var b strings.Builder

	lineNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(4).
		Align(lipgloss.Right)

	codeLineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	for i, line := range lines {
		lineNum := lineNumStyle.Render(fmt.Sprintf("%d", i+1))
		codeLine := codeLineStyle.Render(line)
		b.WriteString(lineNum + " │ " + codeLine)
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
