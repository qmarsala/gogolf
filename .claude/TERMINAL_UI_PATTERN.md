# Terminal UI Pattern: Two-Column Layout with In-Place Animation

A reusable pattern for building terminal UIs that display content in columns, update in place without scrolling, and support smooth animations.

## Overview

This pattern enables:
- **Fixed-position rendering**: Content stays in place, no scrolling
- **Multi-column layouts**: Split terminal into independent panels
- **Smooth animations**: Dice rolls, progress bars, timers that update in place
- **Cross-platform support**: Works on Windows, macOS, and Linux

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Terminal                                   │
│  ┌─────────────────────────┐ │ ┌─────────────────────────────────┐  │
│  │      Left Panel         │ │ │         Right Panel             │  │
│  │                         │ │ │                                 │  │
│  │  - Dynamic content      │ │ │  - Static info (stats)          │  │
│  │  - Animations           │ │ │  - Updates on state change      │  │
│  │  - User prompts         │ │ │                                 │  │
│  │                         │ │ │                                 │  │
│  └─────────────────────────┘ │ └─────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

## Components

### 1. Terminal Controller

Handles low-level terminal operations using ANSI escape sequences.

```go
package ui

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strconv"
    "strings"
)

// ANSI escape sequences
const (
    clearScreen   = "\033[2J"
    moveCursor    = "\033[%d;%dH" // row, col (1-indexed)
    hideCursor    = "\033[?25l"
    showCursor    = "\033[?25h"

    // Colors
    colorReset       = "\033[0m"
    colorGreen       = "\033[32m"
    colorYellow      = "\033[33m"
    colorRed         = "\033[31m"
    colorCyan        = "\033[36m"
    colorBrightWhite = "\033[97m"
    colorDim         = "\033[2m"
)

type Terminal struct {
    Width  int
    Height int
}

func NewTerminal() *Terminal {
    width, height := getTerminalSize()
    return &Terminal{Width: width, Height: height}
}

func (t *Terminal) Clear() {
    fmt.Print(clearScreen)
    t.MoveCursor(1, 1)
}

func (t *Terminal) MoveCursor(row, col int) {
    fmt.Printf(moveCursor, row, col)
}

func (t *Terminal) HideCursor() {
    fmt.Print(hideCursor)
}

func (t *Terminal) ShowCursor() {
    fmt.Print(showCursor)
}
```

### 2. Terminal Size Detection

Cross-platform terminal dimension detection.

```go
func getTerminalSize() (width, height int) {
    if runtime.GOOS == "windows" {
        return getWindowsTerminalSize()
    }
    return getUnixTerminalSize()
}

func getUnixTerminalSize() (width, height int) {
    cmd := exec.Command("stty", "size")
    cmd.Stdin = os.Stdin
    out, err := cmd.Output()
    if err != nil {
        return 120, 30 // fallback
    }

    parts := strings.Fields(string(out))
    if len(parts) != 2 {
        return 120, 30
    }

    height, _ = strconv.Atoi(parts[0])
    width, _ = strconv.Atoi(parts[1])
    return width, height
}

func getWindowsTerminalSize() (width, height int) {
    cmd := exec.Command("powershell", "-Command",
        "$Host.UI.RawUI.WindowSize.Width; $Host.UI.RawUI.WindowSize.Height")
    out, err := cmd.Output()
    if err != nil {
        return 120, 30
    }

    lines := strings.Split(strings.TrimSpace(string(out)), "\n")
    if len(lines) != 2 {
        return 120, 30
    }

    width, _ = strconv.Atoi(strings.TrimSpace(lines[0]))
    height, _ = strconv.Atoi(strings.TrimSpace(lines[1]))
    return width, height
}
```

### 3. Layout System

Defines panels as rectangular regions of the terminal.

```go
const (
    MinTerminalWidth  = 80
    MinTerminalHeight = 24
    LeftPanelWidth    = 60
    RightPanelWidth   = 60
    PanelDivider      = "│"
)

type Panel struct {
    X      int // Column position (1-indexed)
    Y      int // Row position (1-indexed)
    Width  int
    Height int
}

type Layout struct {
    LeftPanel  Panel
    RightPanel Panel
    TermWidth  int
    TermHeight int
}

func NewLayout(termWidth, termHeight int) *Layout {
    leftWidth := LeftPanelWidth
    rightWidth := RightPanelWidth

    // Adjust for smaller terminals
    if termWidth < MinTerminalWidth*1.5 {
        leftWidth = termWidth / 2
        rightWidth = termWidth - leftWidth - 1
    }

    return &Layout{
        LeftPanel: Panel{
            X:      1,
            Y:      1,
            Width:  leftWidth,
            Height: termHeight,
        },
        RightPanel: Panel{
            X:      leftWidth + 2, // +1 divider, +1 spacing
            Y:      1,
            Width:  rightWidth,
            Height: termHeight,
        },
        TermWidth:  termWidth,
        TermHeight: termHeight,
    }
}

func (l *Layout) GetDividerColumn() int {
    return l.LeftPanel.Width + 1
}
```

### 4. Renderer

Manages rendering content to panels.

```go
type Renderer struct {
    Terminal *Terminal
    Layout   *Layout
}

func NewRenderer() *Renderer {
    term := NewTerminal()
    layout := NewLayout(term.Width, term.Height)
    return &Renderer{Terminal: term, Layout: layout}
}

// RenderBorder draws the vertical divider between panels
func (r *Renderer) RenderBorder() {
    dividerCol := r.Layout.GetDividerColumn()
    for row := 1; row <= r.Layout.TermHeight; row++ {
        r.Terminal.MoveCursor(row, dividerCol)
        fmt.Print(PanelDivider)
    }
}

// PrintInPanel writes text at a specific row within a panel
func (r *Renderer) PrintInPanel(panel Panel, row int, text string) {
    if row > panel.Height {
        return
    }

    maxWidth := panel.Width - 4 // margins
    if len(text) > maxWidth {
        text = text[:maxWidth-3] + "..."
    }

    r.Terminal.MoveCursor(row, panel.X+2)
    fmt.Print(text)
}

// ClearLine clears a line within a panel (for animations)
func (r *Renderer) ClearLine(panel Panel, row int) {
    r.Terminal.MoveCursor(row, panel.X+2)
    fmt.Print(strings.Repeat(" ", panel.Width-4))
}
```

### 5. Animation: Dice Roller

Animates dice that stop one at a time, updating in place.

```go
import (
    "math/rand"
    "time"
)

type DiceRoller struct {
    renderer *Renderer
}

func NewDiceRoller(renderer *Renderer) *DiceRoller {
    return &DiceRoller{renderer: renderer}
}

func (dr *DiceRoller) ShowRoll(finalRolls []int, targetNumber int) {
    panel := dr.renderer.Layout.LeftPanel
    row := panel.Height - 6

    // Show target
    dr.renderer.Terminal.MoveCursor(row-1, panel.X+2)
    fmt.Printf("Target: %d", targetNumber)

    stopped := [3]bool{false, false, false}
    displayed := [3]int{1, 1, 1}

    animationDuration := 1500 * time.Millisecond
    stopInterval := animationDuration / 4
    frameInterval := 50 * time.Millisecond

    startTime := time.Now()
    ticker := time.NewTicker(frameInterval)
    defer ticker.Stop()

    for {
        elapsed := time.Since(startTime)

        // Stop dice sequentially
        if elapsed > stopInterval && !stopped[0] {
            stopped[0] = true
            displayed[0] = finalRolls[0]
        }
        if elapsed > stopInterval*2 && !stopped[1] {
            stopped[1] = true
            displayed[1] = finalRolls[1]
        }
        if elapsed > stopInterval*3 && !stopped[2] {
            stopped[2] = true
            displayed[2] = finalRolls[2]
        }

        // Randomize non-stopped dice
        for i := 0; i < 3; i++ {
            if !stopped[i] {
                displayed[i] = rand.Intn(6) + 1
            }
        }

        // Render dice at fixed position
        dr.renderer.Terminal.MoveCursor(row, panel.X+2)
        dr.drawDice(displayed, stopped)

        if stopped[0] && stopped[1] && stopped[2] {
            break
        }

        <-ticker.C
    }

    // Show total
    total := finalRolls[0] + finalRolls[1] + finalRolls[2]
    dr.renderer.Terminal.MoveCursor(row+1, panel.X+2)
    fmt.Printf("Total: %d", total)

    time.Sleep(500 * time.Millisecond)
}

func (dr *DiceRoller) drawDice(values [3]int, stopped [3]bool) {
    fmt.Print("Dice: ")
    for i, val := range values {
        if stopped[i] {
            fmt.Printf("%s[%d]%s ", colorBrightWhite, val, colorReset)
        } else {
            fmt.Printf("%s[%d]%s ", colorDim, val, colorReset)
        }
    }
}
```

### 6. Animation: Power Meter

Real-time progress bar with goroutines for concurrent input handling.

```go
type PowerMeter struct {
    renderer       *Renderer
    maxTime        time.Duration
    sweetSpotStart float64
    sweetSpotEnd   float64
}

func NewPowerMeter(renderer *Renderer) *PowerMeter {
    return &PowerMeter{
        renderer:       renderer,
        maxTime:        2 * time.Second,
        sweetSpotStart: 0.75,
        sweetSpotEnd:   0.85,
    }
}

func (pm *PowerMeter) GetPower() float64 {
    panel := pm.renderer.Layout.LeftPanel
    meterRow := panel.Height - 3

    // Initial prompt
    pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
    fmt.Print("Press SPACE to start power meter...")

    pm.waitForSpacebar()

    startTime := time.Now()
    ticker := time.NewTicker(50 * time.Millisecond)
    defer ticker.Stop()

    stopChan := make(chan bool, 1)
    powerChan := make(chan float64, 1)

    // Wait for second spacebar in goroutine
    go func() {
        pm.waitForSpacebar()
        stopChan <- true
    }()

    // Update display in goroutine
    var finalPower float64
    go func() {
        for {
            select {
            case <-ticker.C:
                elapsed := time.Since(startTime)
                power := pm.calculatePower(elapsed)

                pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
                fmt.Printf("Power: %s %.0f%%   ", pm.drawMeterBar(elapsed), power*100)

                if elapsed >= pm.maxTime {
                    powerChan <- power
                    return
                }
            case <-stopChan:
                elapsed := time.Since(startTime)
                powerChan <- pm.calculatePower(elapsed)
                return
            }
        }
    }()

    finalPower = <-powerChan
    time.Sleep(500 * time.Millisecond)

    return finalPower
}

func (pm *PowerMeter) calculatePower(elapsed time.Duration) float64 {
    ratio := float64(elapsed) / float64(pm.maxTime)
    if ratio > 1.0 {
        ratio = 1.0
    }

    if ratio < pm.sweetSpotStart {
        return (ratio / pm.sweetSpotStart) * 0.95
    }
    if ratio <= pm.sweetSpotEnd {
        return 1.0
    }

    // Decay past sweet spot
    overshoot := (ratio - pm.sweetSpotEnd) / (1.0 - pm.sweetSpotEnd)
    return 0.95 - (overshoot * 0.45)
}

func (pm *PowerMeter) drawMeterBar(elapsed time.Duration) string {
    barWidth := 20
    ratio := float64(elapsed) / float64(pm.maxTime)
    if ratio > 1.0 {
        ratio = 1.0
    }

    sweetStart := int(pm.sweetSpotStart * float64(barWidth))
    sweetEnd := int(pm.sweetSpotEnd * float64(barWidth))
    currentPos := int(ratio * float64(barWidth))

    bar := "["
    for i := 0; i < barWidth; i++ {
        if i == sweetStart {
            bar += "("
        } else if i == sweetEnd {
            bar += ")"
        } else if i < currentPos {
            bar += "="
        } else {
            bar += " "
        }
    }
    bar += "]"

    return bar
}
```

### 7. Platform-Specific Input

#### Windows (input_windows.go)

```go
// +build windows

package ui

import (
    "syscall"
    "unsafe"
)

var (
    kernel32         = syscall.NewLazyDLL("kernel32.dll")
    procReadConsole  = kernel32.NewProc("ReadConsoleInputW")
    procGetStdHandle = kernel32.NewProc("GetStdHandle")
)

const (
    stdInputHandle = uint32(4294967286) // -10
    keyEvent       = 0x0001
    keyDown        = 0x0001
)

type inputRecord struct {
    EventType uint16
    _         uint16
    KeyEvent  keyEventRecord
}

type keyEventRecord struct {
    KeyDown         uint32
    RepeatCount     uint16
    VirtualKeyCode  uint16
    VirtualScanCode uint16
    UnicodeChar     uint16
    ControlKeyState uint32
}

func (pm *PowerMeter) waitForSpacebar() {
    handle, _, _ := procGetStdHandle.Call(uintptr(stdInputHandle))

    for {
        var record inputRecord
        var numRead uint32

        procReadConsole.Call(
            handle,
            uintptr(unsafe.Pointer(&record)),
            1,
            uintptr(unsafe.Pointer(&numRead)),
            0,
        )

        if record.EventType == keyEvent &&
            record.KeyEvent.KeyDown == keyDown &&
            record.KeyEvent.VirtualKeyCode == 0x20 { // VK_SPACE
            return
        }
    }
}

func WaitForAnyKey() {
    handle, _, _ := procGetStdHandle.Call(uintptr(stdInputHandle))

    for {
        var record inputRecord
        var numRead uint32

        procReadConsole.Call(
            handle,
            uintptr(unsafe.Pointer(&record)),
            1,
            uintptr(unsafe.Pointer(&numRead)),
            0,
        )

        if record.EventType == keyEvent && record.KeyEvent.KeyDown == keyDown {
            return
        }
    }
}

func readSingleKey() byte {
    handle, _, _ := procGetStdHandle.Call(uintptr(stdInputHandle))

    for {
        var record inputRecord
        var numRead uint32

        procReadConsole.Call(
            handle,
            uintptr(unsafe.Pointer(&record)),
            1,
            uintptr(unsafe.Pointer(&numRead)),
            0,
        )

        if record.EventType == keyEvent && record.KeyEvent.KeyDown == keyDown {
            return byte(record.KeyEvent.UnicodeChar)
        }
    }
}
```

#### Unix (input_unix.go)

```go
// +build !windows

package ui

import (
    "os"
    "syscall"
    "unsafe"
)

const (
    tcgets = 0x5401
    tcsets = 0x5402
)

type termios struct {
    Iflag  uint32
    Oflag  uint32
    Cflag  uint32
    Lflag  uint32
    Line   uint8
    Cc     [32]uint8
    Ispeed uint32
    Ospeed uint32
}

func (pm *PowerMeter) waitForSpacebar() {
    fd := int(os.Stdin.Fd())
    var oldState termios
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcgets, uintptr(unsafe.Pointer(&oldState)))

    // Raw mode: no echo, no line buffering
    newState := oldState
    newState.Lflag &^= syscall.ECHO | syscall.ICANON
    newState.Cc[syscall.VMIN] = 1
    newState.Cc[syscall.VTIME] = 0
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&newState)))

    defer func() {
        syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&oldState)))
    }()

    buf := make([]byte, 1)
    for {
        os.Stdin.Read(buf)
        if buf[0] == ' ' {
            return
        }
    }
}

func WaitForAnyKey() {
    fd := int(os.Stdin.Fd())
    var oldState termios
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcgets, uintptr(unsafe.Pointer(&oldState)))

    newState := oldState
    newState.Lflag &^= syscall.ECHO | syscall.ICANON
    newState.Cc[syscall.VMIN] = 1
    newState.Cc[syscall.VTIME] = 0
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&newState)))

    defer func() {
        syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&oldState)))
    }()

    buf := make([]byte, 1)
    os.Stdin.Read(buf)
}

func readSingleKey() byte {
    fd := int(os.Stdin.Fd())
    var oldState termios
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcgets, uintptr(unsafe.Pointer(&oldState)))

    newState := oldState
    newState.Lflag &^= syscall.ECHO | syscall.ICANON
    newState.Cc[syscall.VMIN] = 1
    newState.Cc[syscall.VTIME] = 0
    syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&newState)))

    defer func() {
        syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&oldState)))
    }()

    buf := make([]byte, 1)
    os.Stdin.Read(buf)
    return buf[0]
}
```

## Key Techniques

### In-Place Updates

The core technique is using `MoveCursor()` to position the cursor at a fixed location before printing:

```go
// Always print at the same row - overwrites previous content
renderer.Terminal.MoveCursor(row, col)
fmt.Printf("Value: %d   ", value) // trailing spaces clear old content
```

### Ticker-Based Animation

Use `time.Ticker` for consistent frame rates:

```go
ticker := time.NewTicker(50 * time.Millisecond) // 20 FPS
defer ticker.Stop()

for {
    // Update display
    <-ticker.C // Wait for next frame
}
```

### Goroutine Coordination

Use channels to coordinate between input handling and display updates:

```go
stopChan := make(chan bool, 1)
resultChan := make(chan float64, 1)

go func() {
    waitForInput()
    stopChan <- true
}()

go func() {
    for {
        select {
        case <-ticker.C:
            updateDisplay()
        case <-stopChan:
            resultChan <- calculateResult()
            return
        }
    }
}()

result := <-resultChan
```

### Clearing Lines

When updating in place, use trailing spaces or explicit clearing:

```go
// Option 1: Trailing spaces
fmt.Printf("Value: %d      ", value)

// Option 2: Clear entire line within panel
func (r *Renderer) ClearLine(panel Panel, row int) {
    r.Terminal.MoveCursor(row, panel.X+2)
    fmt.Print(strings.Repeat(" ", panel.Width-4))
}
```

## Usage Example

```go
func main() {
    renderer := NewRenderer()
    defer renderer.Terminal.ShowCursor()

    renderer.Terminal.HideCursor()
    renderer.Terminal.Clear()
    renderer.RenderBorder()

    // Render static content
    renderer.PrintInPanel(renderer.Layout.LeftPanel, 1, "=== GAME ===")
    renderer.PrintInPanel(renderer.Layout.RightPanel, 1, "=== STATS ===")

    // Run animation
    diceRoller := NewDiceRoller(renderer)
    diceRoller.ShowRoll([]int{4, 5, 6}, 12)

    // Get timed input
    powerMeter := NewPowerMeter(renderer)
    power := powerMeter.GetPower()

    renderer.Terminal.ShowCursor()
}
```

## Design Principles

1. **Separation of concerns**: Terminal control, layout, rendering, and input are separate components
2. **Panel isolation**: Each panel has its own coordinate space, simplifying layout logic
3. **State-based rendering**: Pass complete state to renderer, which handles all display logic
4. **Platform abstraction**: Platform-specific code isolated to separate files with build tags
5. **Graceful degradation**: Check terminal size and fall back to simpler UI if needed

## File Structure

```
ui/
├── terminal.go      # ANSI control, cursor, colors
├── layout.go        # Panel definitions, layout calculation
├── renderer.go      # Render state to panels
├── input.go         # Animation components (DiceRoller, PowerMeter)
├── input_windows.go # Windows keyboard input
├── input_unix.go    # Unix keyboard input
└── gamestate.go     # State struct for rendering (domain-specific)
```
