package ui

import (
	"testing"
)

// Test NewTerminal creates a terminal with detected size
func TestNewTerminal(t *testing.T) {
	term := NewTerminal()

	if term == nil {
		t.Fatal("NewTerminal returned nil")
	}

	// Terminal should have some reasonable dimensions
	if term.Width < 10 || term.Width > 1000 {
		t.Errorf("Terminal width = %d, expected reasonable value (10-1000)", term.Width)
	}

	if term.Height < 10 || term.Height > 1000 {
		t.Errorf("Terminal height = %d, expected reasonable value (10-1000)", term.Height)
	}
}

// Test GetSize returns terminal dimensions
func TestTerminal_GetSize(t *testing.T) {
	term := NewTerminal()

	width, height := term.GetSize()

	if width != term.Width {
		t.Errorf("GetSize width = %d, want %d", width, term.Width)
	}

	if height != term.Height {
		t.Errorf("GetSize height = %d, want %d", height, term.Height)
	}
}

// Test SupportsANSI returns true on modern systems
func TestTerminal_SupportsANSI(t *testing.T) {
	term := NewTerminal()

	supports := term.SupportsANSI()

	// All modern systems should support ANSI
	if !supports {
		t.Error("SupportsANSI = false, expected true on modern systems")
	}
}

// Test terminal control methods don't panic
func TestTerminal_ControlMethods(t *testing.T) {
	term := NewTerminal()

	// These should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Terminal control method panicked: %v", r)
		}
	}()

	term.HideCursor()
	term.ShowCursor()
	term.SaveCursor()
	term.RestoreCursor()
	term.MoveCursor(1, 1)
	term.Clear()
}

// Test getTerminalSize returns reasonable values
func TestGetTerminalSize(t *testing.T) {
	width, height := getTerminalSize()

	if width < 10 || width > 1000 {
		t.Errorf("getTerminalSize width = %d, expected reasonable value (10-1000)", width)
	}

	if height < 10 || height > 1000 {
		t.Errorf("getTerminalSize height = %d, expected reasonable value (10-1000)", height)
	}
}
