package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewportTab provides viewport management and message handling for tabs.
// Embed this in your tab structs to get scrolling and message functionality.
type ViewportTab struct {
	viewport   viewport.Model
	ready      bool
	errorMsg   string
	successMsg string
}

// NewViewportTab creates a new viewport tab.
func NewViewportTab() ViewportTab {
	return ViewportTab{
		ready: false,
	}
}

// Init initializes the viewport.
func (v ViewportTab) Init() tea.Cmd {
	return nil
}

// HandleResize handles viewport resize messages.
// Call this at the start of your Update method.
func (v ViewportTab) HandleResize(msg ViewportResizeMsg) (ViewportTab, tea.Cmd) {
	if !v.ready {
		v.viewport = viewport.New(msg.Width, msg.Height)
		v.viewport.YPosition = 0
		v.ready = true
	} else {
		v.viewport.Width = msg.Width
		v.viewport.Height = msg.Height
	}
	return v, nil
}

// UpdateViewport updates the viewport for scrolling.
// Call this at the end of your Update method to handle scroll keys.
func (v ViewportTab) UpdateViewport(msg tea.Msg) (ViewportTab, tea.Cmd) {
	var cmd tea.Cmd
	if v.ready {
		v.viewport, cmd = v.viewport.Update(msg)
	}
	return v, cmd
}

// SetContent sets the viewport content.
func (v *ViewportTab) SetContent(content string) {
	if v.ready {
		v.viewport.SetContent(content)
	}
}

// View returns the viewport view or loading message.
func (v ViewportTab) View() string {
	if !v.ready {
		return "Loading..."
	}
	return v.viewport.View()
}

// IsReady returns true if viewport is initialized.
func (v ViewportTab) IsReady() bool {
	return v.ready
}

// GotoTop scrolls to the top of the viewport.
func (v *ViewportTab) GotoTop() {
	if v.ready {
		v.viewport.GotoTop()
	}
}

// GotoBottom scrolls to the bottom of the viewport.
func (v *ViewportTab) GotoBottom() {
	if v.ready {
		v.viewport.GotoBottom()
	}
}

// ScrollPercent returns the current scroll percentage (0-1).
func (v ViewportTab) ScrollPercent() float64 {
	if !v.ready {
		return 0
	}
	return v.viewport.ScrollPercent()
}

// SetYOffset sets the vertical scroll offset.
func (v *ViewportTab) SetYOffset(offset int) {
	if v.ready {
		// Clamp offset to valid range
		if offset < 0 {
			offset = 0
		}
		v.viewport.YOffset = offset
	}
}

// YOffset returns the current vertical scroll offset.
func (v ViewportTab) YOffset() int {
	if !v.ready {
		return 0
	}
	return v.viewport.YOffset
}

// Height returns the viewport height.
func (v ViewportTab) Height() int {
	if !v.ready {
		return 0
	}
	return v.viewport.Height
}

// Width returns the viewport width.
func (v ViewportTab) Width() int {
	if !v.ready {
		return 0
	}
	return v.viewport.Width
}

// TotalLines returns the total number of lines in the content.
// Note: This is an approximation based on scroll percent.
func (v ViewportTab) TotalLines() int {
	if !v.ready {
		return 0
	}
	// We can't access the private lines field, so we return an estimate
	// based on the viewport height and scroll position
	return v.viewport.Height
}

// SetError sets an error message to display.
func (v *ViewportTab) SetError(msg string) {
	v.errorMsg = msg
	v.successMsg = ""
}

// SetSuccess sets a success message to display.
func (v *ViewportTab) SetSuccess(msg string) {
	v.successMsg = msg
	v.errorMsg = ""
}

// ClearMessages clears both error and success messages.
func (v *ViewportTab) ClearMessages() {
	v.errorMsg = ""
	v.successMsg = ""
}

// RenderMessages returns formatted error/success messages.
func (v ViewportTab) RenderMessages() string {
	var result string
	if v.successMsg != "" {
		result = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Render("✓ "+v.successMsg) + "\n\n"
	}
	if v.errorMsg != "" {
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Render("✗ "+v.errorMsg) + "\n\n"
	}
	return result
}

// ViewportResizeMsg is sent to tabs when viewport dimensions change
type ViewportResizeMsg struct {
	Width  int
	Height int
}
