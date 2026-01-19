package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TabModel represents content that can be displayed in a tab.
type TabModel interface {
	tea.Model
}

// Tabs represents the main tabbed interface component.
type Tabs struct {
	labels     []string
	contents   []TabModel
	activeTab  int
	labelWidth int
	debug      bool
	ready      bool
	width      int
	height     int
	styles     *TabStyles
	focusMode  bool
	xOffset    int
}

// TabStyles contains all styling configuration for the tabs component.
type TabStyles struct {
	InactiveTabBorder lipgloss.Border
	ActiveTabBorder   lipgloss.Border
	Doc               lipgloss.Style
	HighlightColor    lipgloss.AdaptiveColor
	InactiveTab       lipgloss.Style
	ActiveTab         lipgloss.Style
	Window            lipgloss.Style
	ContentPadding    lipgloss.Style
	ModeIndicator     lipgloss.Style
	InteractiveMode   lipgloss.Style
	NavigationMode    lipgloss.Style
}

// DefaultTabStyles creates the default styling configuration.
func DefaultTabStyles() *TabStyles {
	highlightColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")

	inactiveTabStyle := lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)

	activeTabStyle := inactiveTabStyle.
		Border(activeTabBorder, true).
		Foreground(highlightColor).
		Bold(true)

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
		ModeIndicator: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true),
		InteractiveMode: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			Bold(true),
		NavigationMode: lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			Bold(true),
	}
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// NewTabs creates a new tabbed interface.
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
		focusMode:  true,
	}, nil
}

func (t *Tabs) SetLabelWidth(width int) {
	t.labelWidth = width
}

func (t *Tabs) SetStyles(styles *TabStyles) {
	t.styles = styles
}

func (t *Tabs) EnableDebug(enabled bool) {
	t.debug = enabled
}

func (t Tabs) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(t.contents))
	for i, content := range t.contents {
		cmds[i] = content.Init()
	}
	return tea.Batch(cmds...)
}

func (t Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return t.handleWindowResize(msg)
	case tea.KeyMsg:
		// Handle global keys first
		switch msg.String() {
		case "ctrl+c":
			return t, tea.Quit
		case "ctrl+d":
			t.debug = !t.debug
			return t, nil
		case "ctrl+f":
			t.focusMode = !t.focusMode
			return t, nil
		}

		// In focus mode, prioritize tab content for most keys
		if t.focusMode {
			switch msg.String() {
			case "ctrl+left":
				t.xOffset = max(0, t.xOffset-5)
				return t, nil
			case "ctrl+right":
				t.xOffset += 5
				return t, nil
			default:
				// Let tab content handle the key
				updatedContent, cmd := t.contents[t.activeTab].Update(msg)
				t.contents[t.activeTab] = updatedContent
				cmds = append(cmds, cmd)
				return t, tea.Batch(cmds...)
			}
		} else {
			// In navigation mode, handle tab switching and scrolling
			return t.handleKeyPress(msg)
		}
	}

	// Update active tab content for other messages
	updatedContent, cmd := t.contents[t.activeTab].Update(msg)
	t.contents[t.activeTab] = updatedContent
	cmds = append(cmds, cmd)

	return t, tea.Batch(cmds...)
}

func (t Tabs) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	t.width = msg.Width
	t.height = msg.Height

	dimensions := t.calculateDimensions()

	// Forward resize to all tabs with viewport dimensions
	resizeMsg := ViewportResizeMsg{
		Width:  dimensions.viewportWidth,
		Height: dimensions.viewportHeight,
	}

	var cmds []tea.Cmd
	for i := range t.contents {
		updatedContent, cmd := t.contents[i].Update(resizeMsg)
		t.contents[i] = updatedContent
		cmds = append(cmds, cmd)
	}

	// Also send the original WindowSizeMsg to active tab
	updatedContent, cmd := t.contents[t.activeTab].Update(msg)
	t.contents[t.activeTab] = updatedContent
	cmds = append(cmds, cmd)

	t.ready = true
	return t, tea.Batch(cmds...)
}

func (t Tabs) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "right":
		t.activeTab = min(t.activeTab+1, len(t.labels)-1)
		t.xOffset = 0
		return t, nil

	case "left":
		t.activeTab = max(t.activeTab-1, 0)
		t.xOffset = 0
		return t, nil

	case "tab":
		t.activeTab = min(t.activeTab+1, len(t.labels)-1)
		t.xOffset = 0
		return t, nil

	case "shift+tab":
		t.activeTab = max(t.activeTab-1, 0)
		t.xOffset = 0
		return t, nil

	case "ctrl+right":
		t.xOffset += 5
		return t, nil

	case "ctrl+left":
		t.xOffset = max(0, t.xOffset-5)
		return t, nil

	default:
		// Pass keys to active tab
		updatedContent, cmd := t.contents[t.activeTab].Update(msg)
		t.contents[t.activeTab] = updatedContent
		return t, cmd
	}
}

type dimensions struct {
	tabWidth       int
	viewportWidth  int
	viewportHeight int
}

func (t Tabs) calculateDimensions() dimensions {
	// Calculate the actual height used by UI elements
	tabHeaders := t.renderTabHeaders()
	tabHeaderLines := strings.Count(tabHeaders, "\n") + 1

	var footerText string
	if t.debug {
		footerText = "\n\nCtrl+F: toggle mode | Ctrl+D: debug | Ctrl+C: quit"
	} else {
		footerText = "\n\nCtrl+F: toggle mode | Ctrl+C: quit"
	}
	footerLines := strings.Count(footerText, "\n") + 1

	docVPadding := t.styles.Doc.GetVerticalFrameSize()
	windowVPadding := t.styles.Window.GetVerticalFrameSize()
	contentVPadding := t.styles.ContentPadding.GetVerticalFrameSize()

	// Total height - all UI overhead = available viewport height
	usedHeight := tabHeaderLines + footerLines + docVPadding + windowVPadding + contentVPadding
	viewportHeight := t.height - usedHeight

	if viewportHeight < 1 {
		viewportHeight = 1
	}

	// Calculate widths
	tabWidth := lipgloss.Width(tabHeaders)
	windowHPadding := t.styles.Window.GetHorizontalFrameSize()
	contentHPadding := t.styles.ContentPadding.GetHorizontalFrameSize()
	viewportWidth := tabWidth - windowHPadding - contentHPadding

	return dimensions{
		tabWidth:       tabWidth,
		viewportWidth:  viewportWidth,
		viewportHeight: viewportHeight,
	}
}

func (t Tabs) applyHorizontalScroll(content string) string {
	if t.xOffset == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	var scrolledLines []string

	for _, line := range lines {
		runes := []rune(line)
		if t.xOffset >= len(runes) {
			scrolledLines = append(scrolledLines, "")
		} else {
			scrolledLines = append(scrolledLines, string(runes[t.xOffset:]))
		}
	}

	return strings.Join(scrolledLines, "\n")
}

func (t Tabs) buildDebugInfo() string {
	dimensions := t.calculateDimensions()
	focusStatus := "Navigation"
	if t.focusMode {
		focusStatus = "Interactive"
	}

	debugInfo := fmt.Sprintf(
		"Width Debug:\n"+
			"  Terminal height: %d\n"+
			"  Label width setting: %d\n"+
			"  Tab row width: %d\n"+
			"  Viewport width: %d\n"+
			"  Viewport height: %d\n"+
			"  Window frame size: %d\n"+
			"  Content padding size: %d\n"+
			"  Horizontal offset: %d\n"+
			"  Focus mode: %s\n\n",
		t.height,
		t.labelWidth,
		dimensions.tabWidth,
		dimensions.viewportWidth,
		dimensions.viewportHeight,
		t.styles.Window.GetHorizontalFrameSize(),
		t.styles.ContentPadding.GetHorizontalFrameSize(),
		t.xOffset,
		focusStatus,
	)

	return debugInfo + strings.Repeat("─", 40) + "\n\n"
}

func (t Tabs) renderTabHeaders() string {
	var renderedTabs []string

	for i, label := range t.labels {
		renderedTab := t.renderSingleTab(i, label)
		renderedTabs = append(renderedTabs, renderedTab)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (t Tabs) renderSingleTab(index int, label string) string {
	isFirst := index == 0
	isLast := index == len(t.labels)-1
	isActive := index == t.activeTab

	var style lipgloss.Style
	if isActive {
		style = t.styles.ActiveTab
	} else {
		style = t.styles.InactiveTab
	}

	displayLabel := label
	if isActive {
		var modeIcon string
		var iconColor lipgloss.Color
		if t.focusMode {
			modeIcon = "⌨"
			iconColor = lipgloss.Color("10")
		} else {
			modeIcon = "⇄"
			iconColor = lipgloss.Color("12")
		}
		coloredIcon := lipgloss.NewStyle().Foreground(iconColor).Render(modeIcon)
		displayLabel = coloredIcon + " " + label
	}

	border := t.adjustBorder(style, isFirst, isLast, isActive)
	style = style.Border(border).Width(t.labelWidth).Align(lipgloss.Center)

	return style.Render(displayLabel)
}

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

func (t Tabs) View() string {
	if !t.ready {
		return "Initializing..."
	}

	var doc strings.Builder

	tabHeaders := t.renderTabHeaders()
	doc.WriteString(tabHeaders)
	doc.WriteString("\n")

	content := t.renderContent(tabHeaders)
	doc.WriteString(content)

	if t.debug {
		doc.WriteString("\n\nCtrl+F: toggle mode | Ctrl+D: debug | Ctrl+C: quit")
	} else {
		doc.WriteString("\n\nCtrl+F: toggle mode | Ctrl+C: quit")
	}

	return t.styles.Doc.Render(doc.String())
}

func (t Tabs) renderContent(tabHeaders string) string {
	tabWidth := lipgloss.Width(tabHeaders)
	contentWidth := tabWidth - t.styles.Window.GetHorizontalFrameSize()

	content := t.contents[t.activeTab].View()

	if t.debug {
		content = t.buildDebugInfo() + content
	}

	// Apply horizontal scrolling
	content = t.applyHorizontalScroll(content)

	paddedViewport := t.styles.ContentPadding.Render(content)

	return t.styles.Window.Width(contentWidth).Render(paddedViewport)
}
