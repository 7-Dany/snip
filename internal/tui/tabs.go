package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TabModel is the interface that tab content must implement.
// Any Bubble Tea model can be used as tab content.
type TabModel interface {
	tea.Model
}

// Tabs manages a tabbed interface with scrollable content.
//
// Design: Fixed headers with scrollable viewport
// - Tab headers remain fixed at the top
// - Only content area scrolls using bubbles/viewport
// - Width automatically matches tab headers
// - Supports customizable styling via TabStyles
//
// Usage:
//
//	tabs, _ := NewTabs(labels, contents)
//	tabs.SetLabelWidth(20)  // Optional: adjust tab width
//	tabs.EnableDebug(true)  // Optional: show debug info
type Tabs struct {
	labels     []string
	contents   []TabModel
	activeTab  int
	labelWidth int
	debug      bool
	viewport   viewport.Model
	ready      bool
	width      int
	height     int
	styles     *TabStyles
}

// TabStyles holds all the styling configuration for the tabs component.
//
// Design: Separation of concerns
// - All visual styling is centralized here
// - Can be customized via SetStyles()
// - Default styles provided by DefaultTabStyles()
type TabStyles struct {
	InactiveTabBorder lipgloss.Border
	ActiveTabBorder   lipgloss.Border
	Doc               lipgloss.Style
	HighlightColor    lipgloss.AdaptiveColor
	InactiveTab       lipgloss.Style
	ActiveTab         lipgloss.Style
	Window            lipgloss.Style
	ContentPadding    lipgloss.Style
}

// DefaultTabStyles returns the default styling configuration.
//
// Design: Purple/violet theme with rounded borders
// - Highlight color: #874BFD (light) / #7D56F4 (dark)
// - Active tabs use special bottom border to connect with content
// - Inactive tabs use standard bottom border
// - Content padding: 2 spaces horizontal, 0 vertical
func DefaultTabStyles() *TabStyles {
	highlightColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")

	inactiveTabStyle := lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)

	activeTabStyle := inactiveTabStyle.
		Border(activeTabBorder, true)

	return &TabStyles{
		InactiveTabBorder: inactiveTabBorder,
		ActiveTabBorder:   activeTabBorder,
		Doc:               lipgloss.NewStyle().Padding(1, 2, 1, 2),
		HighlightColor:    highlightColor,
		InactiveTab:       inactiveTabStyle,
		ActiveTab:         activeTabStyle,
		Window: lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(0).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop(),
		ContentPadding: lipgloss.NewStyle().Padding(0, 2),
	}
}

// tabBorderWithBottom creates a custom border for tabs.
// This is used to create the seamless connection between active tabs and content.
//
// Design: Modified rounded border
// - Uses lipgloss.RoundedBorder() as base
// - Customizes bottom border characters for tab connection
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// NewTabs creates a new Tabs component with default settings.
//
// Args:
//
//	labels: Tab header labels (must match length of contents)
//	contents: Tab content models (must implement TabModel interface)
//
// Returns error if labels and contents lengths don't match.
//
// Default settings:
//   - Label width: 14 characters
//   - Debug mode: disabled
//   - Styles: DefaultTabStyles()
func NewTabs(labels []string, contents []TabModel) (Tabs, error) {
	if len(labels) != len(contents) {
		return Tabs{}, fmt.Errorf("labels length must be same as contents length")
	}
	return Tabs{
		labels:     labels,
		contents:   contents,
		activeTab:  0,
		labelWidth: 14,
		debug:      false,
		ready:      false,
		styles:     DefaultTabStyles(),
	}, nil
}

// SetLabelWidth sets the width for each tab label
func (t *Tabs) SetLabelWidth(width int) {
	t.labelWidth = width
}

// SetStyles allows customizing the tab styles
func (t *Tabs) SetStyles(styles *TabStyles) {
	t.styles = styles
}

// EnableDebug turns on debug mode to show width calculations
func (t *Tabs) EnableDebug(enabled bool) {
	t.debug = enabled
}

// Init initializes the tabs component
func (t Tabs) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(t.contents))
	for i, content := range t.contents {
		cmds[i] = content.Init()
	}
	return tea.Batch(cmds...)
}

// Update handles messages and updates the component state
func (t Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return t.handleWindowResize(msg)
	case tea.KeyMsg:
		return t.handleKeyPress(msg)
	}

	// Update active tab content
	updatedContent, cmd := t.contents[t.activeTab].Update(msg)
	t.contents[t.activeTab] = updatedContent

	// Update viewport content
	t.updateViewportContent()

	return t, cmd
}

// handleWindowResize handles window resize events and updates viewport dimensions.
//
// Design: Responsive layout
// - Calculates viewport size based on available terminal space
// - Reserves space for headers (5 lines) and footer (2 lines)
// - Viewport width = tab width - borders (2) - padding (4)
// - Initializes viewport on first resize event
func (t Tabs) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	t.width = msg.Width
	t.height = msg.Height

	dimensions := t.calculateDimensions()

	if !t.ready {
		t.viewport = viewport.New(dimensions.viewportWidth, dimensions.viewportHeight)
		t.viewport.YPosition = 0
		t.ready = true
		t.updateViewportContent()
	} else {
		t.viewport.Width = dimensions.viewportWidth
		t.viewport.Height = dimensions.viewportHeight
	}

	return t, nil
}

// handleKeyPress handles keyboard input for tab navigation and scrolling.
//
// Key bindings:
//   - Tab, Right: Next tab
//   - Shift+Tab, Left: Previous tab
//   - d: Toggle debug mode
//   - q, Ctrl+C: Quit
//   - Up/Down, PgUp/PgDown: Scroll content (passed to viewport)
//
// Design: Resets scroll position to top when switching tabs
func (t Tabs) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return t, tea.Quit

	case "right", "tab":
		t.activeTab = min(t.activeTab+1, len(t.labels)-1)
		t.viewport.GotoTop()
		t.updateViewportContent()
		return t, nil

	case "left", "shift+tab":
		t.activeTab = max(t.activeTab-1, 0)
		t.viewport.GotoTop()
		t.updateViewportContent()
		return t, nil

	case "d":
		t.debug = !t.debug
		t.updateViewportContent()
		return t, nil

	default:
		// Pass keys to viewport for scrolling
		if t.ready {
			t.viewport, cmd = t.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return t, tea.Batch(cmds...)
}

// dimensions holds calculated layout dimensions for the tabs component.
//
// Design: Centralized dimension calculations
// - Calculated once per render to avoid repeated computation
// - All dimensions account for borders, padding, and spacing
type dimensions struct {
	tabWidth       int // Total width of all tab headers combined
	viewportWidth  int // Width available for viewport content
	viewportHeight int // Height available for viewport content
}

// calculateDimensions calculates the layout dimensions based on current state.
//
// Design: Layout calculation logic
// - Header: 5 lines (tabs + border)
// - Footer: 2 lines (padding)
// - Viewport width = tabWidth - borders (2) - content padding (4)
// - Viewport height = total height - header - footer
func (t Tabs) calculateDimensions() dimensions {
	const (
		headerHeight = 5
		footerHeight = 2
		borderSize   = 2
		paddingSize  = 4
	)

	tabHeaders := t.renderTabHeaders()
	tabWidth := lipgloss.Width(tabHeaders)
	viewportWidth := tabWidth - borderSize - paddingSize
	viewportHeight := t.height - headerHeight - footerHeight

	return dimensions{
		tabWidth:       tabWidth,
		viewportWidth:  viewportWidth,
		viewportHeight: viewportHeight,
	}
}

// updateViewportContent updates the viewport with the current tab's content.
//
// Design: Content management
// - Fetches content from active tab's Model
// - Prepends debug info if debug mode is enabled
// - Updates viewport to trigger re-render
func (t *Tabs) updateViewportContent() {
	if !t.ready {
		return
	}

	content := t.contents[t.activeTab].View()

	if t.debug {
		content = t.buildDebugInfo() + content
	}

	t.viewport.SetContent(content)
}

// buildDebugInfo creates debug information string showing layout calculations.
//
// Output includes:
//   - Label width setting
//   - Calculated tab row width
//   - Viewport dimensions
//   - Current scroll position percentage
func (t Tabs) buildDebugInfo() string {
	dimensions := t.calculateDimensions()

	debugInfo := fmt.Sprintf(
		"Width Debug:\n"+
			"  Label width setting: %d\n"+
			"  Tab row width: %d\n"+
			"  Viewport width: %d\n"+
			"  Viewport height: %d\n"+
			"  Scroll position: %.0f%%\n\n",
		t.labelWidth,
		dimensions.tabWidth,
		t.viewport.Width,
		t.viewport.Height,
		t.viewport.ScrollPercent()*100,
	)

	return debugInfo + strings.Repeat("─", 40) + "\n\n"
}

// renderTabHeaders renders all tab headers as a horizontal row.
//
// Design: Horizontal tab bar
// - Each tab is rendered individually via renderSingleTab
// - Tabs are joined horizontally at the top alignment
// - Returns complete header row as single string
func (t Tabs) renderTabHeaders() string {
	var renderedTabs []string

	for i, label := range t.labels {
		renderedTab := t.renderSingleTab(i, label)
		renderedTabs = append(renderedTabs, renderedTab)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// renderSingleTab renders a single tab with appropriate styling.
//
// Design: Dynamic border adjustment
// - Active tab: Uses special bottom border to connect with content
// - Inactive tab: Uses standard bottom border with separator
// - First/last tabs: Adjusted borders for seamless edge connection
// - All tabs: Centered text with fixed width
//
// Note: lipgloss methods automatically return new styles (no .Copy() needed)
func (t Tabs) renderSingleTab(index int, label string) string {
	isFirst := index == 0
	isLast := index == len(t.labels)-1
	isActive := index == t.activeTab

	// Choose base style
	var style lipgloss.Style
	if isActive {
		style = t.styles.ActiveTab
	} else {
		style = t.styles.InactiveTab
	}

	// Adjust borders for first/last tabs
	border := t.adjustBorder(style, isFirst, isLast, isActive)
	style = style.Border(border).Width(t.labelWidth).Align(lipgloss.Center)

	return style.Render(label)
}

// adjustBorder adjusts the border characters for edge tabs.
//
// Design: Seamless border connections
// - First active tab: Uses │ for bottom-left (connects to content)
// - First inactive tab: Uses ├ for bottom-left (connects to window border)
// - Last active tab: Uses │ for bottom-right (connects to content)
// - Last inactive tab: Uses ┤ for bottom-right (connects to window border)
// - Middle tabs: Uses default border characters
func (t Tabs) adjustBorder(style lipgloss.Style, isFirst, isLast, isActive bool) lipgloss.Border {
	border, _, _, _, _ := style.GetBorder()

	if isFirst && isActive {
		border.BottomLeft = "│"
	} else if isFirst && !isActive {
		border.BottomLeft = "├"
	} else if isLast && isActive {
		border.BottomRight = "│"
	} else if isLast && !isActive {
		border.BottomRight = "┤"
	}

	return border
}

// View renders the complete tabs UI (headers + scrollable content).
//
// Design: Two-part layout
// 1. Fixed tab headers at top
// 2. Scrollable content area below
//
// The content width automatically matches the tab headers width.
// Debug footer is shown when debug mode is enabled.
func (t Tabs) View() string {
	if !t.ready {
		return "Initializing..."
	}

	var doc strings.Builder

	// Render fixed tab headers
	tabHeaders := t.renderTabHeaders()
	doc.WriteString(tabHeaders)
	doc.WriteString("\n")

	// Render scrollable content area
	content := t.renderContent(tabHeaders)
	doc.WriteString(content)

	// Add debug footer if enabled
	if t.debug {
		doc.WriteString("\n\nPress 'd' to toggle debug | ↑↓ to scroll | Tab/Shift+Tab to switch tabs")
	}

	return t.styles.Doc.Render(doc.String())
}

// renderContent renders the content area with viewport.
//
// Design: Padded viewport in bordered window
// - Window width matches tab headers exactly
// - Content has horizontal padding (2 spaces) inside viewport
// - Viewport handles scrolling for long content
//
// Note: lipgloss methods automatically return new styles (no .Copy() needed)
func (t Tabs) renderContent(tabHeaders string) string {
	tabWidth := lipgloss.Width(tabHeaders)
	contentWidth := tabWidth - t.styles.Window.GetHorizontalFrameSize()

	// Add padding to viewport content
	paddedViewport := t.styles.ContentPadding.Render(t.viewport.View())

	return t.styles.Window.Width(contentWidth).Render(paddedViewport)
}
