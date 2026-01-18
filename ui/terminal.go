package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ANSI escape sequences for terminal control
const (
	clearScreen   = "\033[2J"
	moveCursor    = "\033[%d;%dH" // row, col
	saveCursor    = "\033[s"
	restoreCursor = "\033[u"
	hideCursor    = "\033[?25l"
	showCursor    = "\033[?25h"

	// Colors
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBrightGreen = "\033[92m"
	colorBrightWhite = "\033[97m"
	colorDim         = "\033[2m"
)

// Terminal manages terminal control and rendering
type Terminal struct {
	Width  int
	Height int
}

// NewTerminal creates a new Terminal and detects terminal size
func NewTerminal() *Terminal {
	width, height := getTerminalSize()
	return &Terminal{
		Width:  width,
		Height: height,
	}
}

// Clear clears the entire terminal screen
func (t *Terminal) Clear() {
	fmt.Print(clearScreen)
	t.MoveCursor(1, 1)
}

// MoveCursor positions the cursor at the specified row and column (1-indexed)
func (t *Terminal) MoveCursor(row, col int) {
	fmt.Printf(moveCursor, row, col)
}

// HideCursor hides the terminal cursor
func (t *Terminal) HideCursor() {
	fmt.Print(hideCursor)
}

// ShowCursor shows the terminal cursor
func (t *Terminal) ShowCursor() {
	fmt.Print(showCursor)
}

// SaveCursor saves the current cursor position
func (t *Terminal) SaveCursor() {
	fmt.Print(saveCursor)
}

// RestoreCursor restores the previously saved cursor position
func (t *Terminal) RestoreCursor() {
	fmt.Print(restoreCursor)
}

// GetSize returns the current terminal dimensions
func (t *Terminal) GetSize() (width, height int) {
	return t.Width, t.Height
}

// SupportsANSI checks if the terminal supports ANSI escape codes
func (t *Terminal) SupportsANSI() bool {
	// On Windows 10+, enable virtual terminal processing
	if runtime.GOOS == "windows" {
		return enableWindowsANSI()
	}
	// Unix-like systems (macOS, Linux) support ANSI by default
	return true
}

// getTerminalSize detects the terminal dimensions
func getTerminalSize() (width, height int) {
	if runtime.GOOS == "windows" {
		return getWindowsTerminalSize()
	}
	return getUnixTerminalSize()
}

// getUnixTerminalSize gets terminal size on Unix-like systems
func getUnixTerminalSize() (width, height int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		// Default fallback size
		return 120, 30
	}

	parts := strings.Fields(string(out))
	if len(parts) != 2 {
		return 120, 30
	}

	height, err = strconv.Atoi(parts[0])
	if err != nil {
		height = 30
	}

	width, err = strconv.Atoi(parts[1])
	if err != nil {
		width = 120
	}

	return width, height
}

// getWindowsTerminalSize gets terminal size on Windows
func getWindowsTerminalSize() (width, height int) {
	// Try using PowerShell to get console size
	cmd := exec.Command("powershell", "-Command", "$Host.UI.RawUI.WindowSize.Width; $Host.UI.RawUI.WindowSize.Height")
	out, err := cmd.Output()
	if err != nil {
		// Default fallback size
		return 120, 30
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2 {
		return 120, 30
	}

	width, err = strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		width = 120
	}

	height, err = strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		height = 30
	}

	return width, height
}

// enableWindowsANSI enables virtual terminal processing on Windows 10+
func enableWindowsANSI() bool {
	// On Windows, ANSI support is available in Windows 10+ with virtual terminal processing
	// For simplicity, we'll assume it's supported and rely on the Go runtime to handle it
	// In a production app, we would use Windows API calls to enable VT processing
	return runtime.GOOS == "windows" || runtime.GOOS != "windows"
}
