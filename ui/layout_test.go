package ui

import "testing"

// Test NewLayout creates layout with correct panel positions
func TestNewLayout(t *testing.T) {
	layout := NewLayout(120, 30)

	if layout == nil {
		t.Fatal("NewLayout returned nil")
	}

	if layout.TermWidth != 120 {
		t.Errorf("Layout TermWidth = %d, want 120", layout.TermWidth)
	}

	if layout.TermHeight != 30 {
		t.Errorf("Layout TermHeight = %d, want 30", layout.TermHeight)
	}
}

// Test left panel starts at position 1,1
func TestNewLayout_LeftPanelPosition(t *testing.T) {
	layout := NewLayout(120, 30)

	if layout.LeftPanel.X != 1 {
		t.Errorf("LeftPanel X = %d, want 1", layout.LeftPanel.X)
	}

	if layout.LeftPanel.Y != 1 {
		t.Errorf("LeftPanel Y = %d, want 1", layout.LeftPanel.Y)
	}
}

// Test left panel has expected width
func TestNewLayout_LeftPanelWidth(t *testing.T) {
	layout := NewLayout(120, 30)

	expectedWidth := 60
	if layout.LeftPanel.Width != expectedWidth {
		t.Errorf("LeftPanel Width = %d, want %d", layout.LeftPanel.Width, expectedWidth)
	}
}

// Test right panel starts after left panel plus divider
func TestNewLayout_RightPanelPosition(t *testing.T) {
	layout := NewLayout(120, 30)

	expectedX := layout.LeftPanel.Width + 2 // +1 divider, +1 spacing
	if layout.RightPanel.X != expectedX {
		t.Errorf("RightPanel X = %d, want %d", layout.RightPanel.X, expectedX)
	}

	if layout.RightPanel.Y != 1 {
		t.Errorf("RightPanel Y = %d, want 1", layout.RightPanel.Y)
	}
}

// Test panels have full terminal height
func TestNewLayout_PanelHeight(t *testing.T) {
	layout := NewLayout(120, 30)

	if layout.LeftPanel.Height != 30 {
		t.Errorf("LeftPanel Height = %d, want 30", layout.LeftPanel.Height)
	}

	if layout.RightPanel.Height != 30 {
		t.Errorf("RightPanel Height = %d, want 30", layout.RightPanel.Height)
	}
}

// Test GetDividerColumn returns correct position
func TestLayout_GetDividerColumn(t *testing.T) {
	layout := NewLayout(120, 30)

	dividerCol := layout.GetDividerColumn()
	expectedCol := layout.LeftPanel.Width + 1

	if dividerCol != expectedCol {
		t.Errorf("GetDividerColumn = %d, want %d", dividerCol, expectedCol)
	}
}

// Test SupportsRichUI returns true for large terminals
func TestLayout_SupportsRichUI_Large(t *testing.T) {
	layout := NewLayout(120, 30)

	if !layout.SupportsRichUI() {
		t.Error("SupportsRichUI = false for 120x30 terminal, want true")
	}
}

// Test SupportsRichUI returns false for small terminals
func TestLayout_SupportsRichUI_Small(t *testing.T) {
	layout := NewLayout(70, 20)

	if layout.SupportsRichUI() {
		t.Error("SupportsRichUI = true for 70x20 terminal, want false")
	}
}

// Test layout adjusts for smaller terminals
func TestNewLayout_SmallTerminal(t *testing.T) {
	layout := NewLayout(80, 24)

	// Panels should adjust to smaller size
	if layout.LeftPanel.Width+layout.RightPanel.Width > 80 {
		t.Errorf("Total panel width exceeds terminal width: %d + %d > 80",
			layout.LeftPanel.Width, layout.RightPanel.Width)
	}
}

// Test layout proportions are reasonable
func TestNewLayout_Proportions(t *testing.T) {
	layout := NewLayout(120, 30)

	// Left and right panels should be roughly equal
	ratio := float64(layout.LeftPanel.Width) / float64(layout.RightPanel.Width)
	if ratio < 0.8 || ratio > 1.2 {
		t.Errorf("Panel width ratio = %.2f, expected roughly 1.0", ratio)
	}
}
