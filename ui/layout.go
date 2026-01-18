package ui

// Layout constants for terminal UI organization
const (
	MinTerminalWidth  = 80
	MinTerminalHeight = 24

	// Panel dimensions (assuming 120 column terminal)
	LeftPanelWidth  = 60
	RightPanelWidth = 60

	// Box-drawing characters
	PanelDivider   = "│"
	HorizontalLine = "─"
	TopLeftCorner  = "┌"
	TopRightCorner = "┐"
	BotLeftCorner  = "└"
	BotRightCorner = "┘"
	LeftTee        = "├"
	RightTee       = "┤"
	TopTee         = "┬"
	BotTee         = "┴"
)

// Layout defines the screen layout with distinct areas
type Layout struct {
	LeftPanel  Panel
	RightPanel Panel
	TermWidth  int
	TermHeight int
}

// Panel represents a rectangular area of the terminal
type Panel struct {
	X      int // Column position (1-indexed)
	Y      int // Row position (1-indexed)
	Width  int // Width in columns
	Height int // Height in rows
}

// NewLayout creates a layout based on terminal dimensions
func NewLayout(termWidth, termHeight int) *Layout {
	// Calculate panel dimensions based on terminal size
	// For now, use fixed 60/60 split if terminal is wide enough
	leftWidth := LeftPanelWidth
	rightWidth := RightPanelWidth

	// If terminal is smaller, adjust proportionally
	if termWidth < MinTerminalWidth*1.5 {
		leftWidth = termWidth / 2
		rightWidth = termWidth - leftWidth - 1 // -1 for divider
	}

	return &Layout{
		LeftPanel: Panel{
			X:      1,
			Y:      1,
			Width:  leftWidth,
			Height: termHeight,
		},
		RightPanel: Panel{
			X:      leftWidth + 2, // +1 for divider, +1 for spacing
			Y:      1,
			Width:  rightWidth,
			Height: termHeight,
		},
		TermWidth:  termWidth,
		TermHeight: termHeight,
	}
}

// GetDividerColumn returns the column where the vertical divider should be drawn
func (l *Layout) GetDividerColumn() int {
	return l.LeftPanel.Width + 1
}

// SupportsRichUI returns true if the terminal is large enough for the rich UI
func (l *Layout) SupportsRichUI() bool {
	return l.TermWidth >= MinTerminalWidth && l.TermHeight >= MinTerminalHeight
}
