# GoGolf UX Improvement Plan
## Terminal UI Redesign

### Current Problems

**Issue 1: Scrolling Text Chaos**
- Terminal constantly scrolls with new output
- Previous information disappears off screen
- Hard to track game state
- No visual continuity between shots
- Player must scroll back to see stats/equipment

**Issue 2: Information Overload**
- All information mixed together chronologically
- No clear separation between persistent state and temporary actions
- Important info (stats, equipment, money) buried in output

**Issue 3: Poor Spatial Organization**
- Everything in a single column
- No dedicated areas for different types of information
- Cognitive load to parse text wall

---

## Solution: Split-Panel Terminal UI

### Design Concept

Use terminal control sequences to create a **persistent, non-scrolling interface** with distinct areas:

```
┌─────────────────────────────────┬─────────────────────────────────┐
│         GAME AREA (LEFT)        │    PLAYER INFO (RIGHT)          │
│                                 │                                 │
│  Hole 1 - Par 4                 │  === PLAYER STATS ===           │
│  Distance: 300 yards            │  Player: Jordan                 │
│  Ball: Tee                      │  Money: 150                     │
│                                 │                                 │
│  Your shot:                     │  --- Skills ---                 │
│  ├─ Quality: Good               │  Driver: Lvl 2 [4]              │
│  ├─ Rotation: 3.2° right        │  Woods: Lvl 1 [2]               │
│  └─ Power: 95%                  │  ...                            │
│                                 │                                 │
│  Ball traveled 285 yards        │  --- Abilities ---              │
│  New lie: Fairway               │  Strength: Lvl 2 [4]            │
│  Distance to hole: 15 yards     │  Control: Lvl 1 [2]             │
│                                 │  ...                            │
│  [Enter club selection]         │                                 │
│  >                              │  --- Equipment ---              │
│                                 │  Ball: Standard (+3 dist)       │
│                                 │  Glove: None                    │
│                                 │  Shoes: None                    │
│                                 │                                 │
│                                 │  --- Score ---                  │
│                                 │  Total: 5 (+2)                  │
│                                 │  This Hole: 2 strokes           │
└─────────────────────────────────┴─────────────────────────────────┘
```

### Key Principles

1. **Static Layout**: Terminal screen is cleared and redrawn, not scrolled
2. **Persistent Right Panel**: Player stats always visible
3. **Dynamic Left Panel**: Game actions and current state
4. **Clean Separation**: Clear visual boundaries between areas
5. **Minimal Re-rendering**: Only update changed sections when possible

---

## Technical Implementation

### Option 1: ANSI Escape Sequences (Recommended)
**Pros:**
- Built into all modern terminals
- Cross-platform (Windows, macOS, Linux)
- No external dependencies
- Full control over cursor positioning and colors

**Cons:**
- Need to handle terminal size detection
- Manual layout calculations
- More complex than simple `fmt.Println()`

**Implementation:**
```go
// Terminal control codes
const (
    clearScreen   = "\033[2J"
    moveCursor    = "\033[%d;%dH"  // row, col
    saveCursor    = "\033[s"
    restoreCursor = "\033[u"
    hideCursor    = "\033[?25l"
    showCursor    = "\033[?25h"

    // Colors
    colorReset    = "\033[0m"
    colorGreen    = "\033[32m"
    colorYellow   = "\033[33m"
    colorRed      = "\033[31m"
    colorCyan     = "\033[36m"
)

// Get terminal size
func getTerminalSize() (width, height int) {
    // Unix: use syscall
    // Windows: use Windows API
}

// Position cursor at row, col
func moveCursorTo(row, col int) {
    fmt.Printf("\033[%d;%dH", row, col)
}
```

### Option 2: Third-party Library (tcell, termbox-go, bubbletea)
**Pros:**
- Higher-level abstractions
- Built-in layout helpers
- Event handling included

**Cons:**
- External dependency
- Learning curve
- May be overkill for our needs

**Recommendation**: Start with ANSI escape sequences (Option 1) for maximum control and minimal dependencies.

---

## Layout Specification

### Terminal Size Assumptions
- **Minimum**: 80 columns × 24 rows
- **Recommended**: 120 columns × 30 rows
- **Graceful degradation** if terminal is smaller

### Panel Division
```
┌─────────────────────────┬────────────────────────┐
│   LEFT: 60 cols         │   RIGHT: 60 cols      │
│   (Game Area)           │   (Player Info)        │
│   Lines 1-30            │   Lines 1-30           │
└─────────────────────────┴────────────────────────┘
```

### Right Panel Layout (Static)
```
Row  1: === PLAYER STATS ===
Row  2: Player: [Name]
Row  3: Money: [Amount]
Row  4:
Row  5: --- Skills ---
Row  6: Driver: Lvl X [Value] (XP/NextXP)
Row  7: Woods: ...
Row  8: Long Irons: ...
Row  9: Mid Irons: ...
Row 10: Short Irons: ...
Row 11: Wedges: ...
Row 12: Putter: ...
Row 13:
Row 14: --- Abilities ---
Row 15: Strength: Lvl X [Value] (XP/NextXP)
Row 16: Control: ...
Row 17: Touch: ...
Row 18: Mental: ...
Row 19:
Row 20: --- Equipment ---
Row 21: Ball: [Name] ([Bonus])
Row 22: Glove: [Name] ([Bonus])
Row 23: Shoes: [Name] ([Bonus])
Row 24:
Row 25: --- Score ---
Row 26: Total: [Strokes] ([Score])
Row 27: This Hole: [Strokes]
Row 28: Holes: [Current]/[Total]
Row 29:
Row 30: [Status messages]
```

### Left Panel Layout (Dynamic)
```
Row  1: === HOLE [N] - PAR [X] ===
Row  2: Distance: [Y] yards
Row  3:
Row  4: Current Lie: [LieType] (difficulty: +/-N)
Row  5: Distance to hole: [Z] yards
Row  6:
Row  7: Last Shot:
Row  8: ├─ Club: [ClubName]
Row  9: ├─ Quality: [Outcome] (Margin: +/-N)
Row 10: ├─ Description: [Quality text]
Row 11: ├─ Rotation: [Degrees]° [Direction]
Row 12: ├─ Power: [Percent]%
Row 13: └─ Distance: [Yards] yards
Row 14:
Row 15: Ball Location: (X, Y)
Row 16: Hole Location: (X, Y)
Row 17:
Row 18: XP Earned: +[N]
Row 19: [Level up messages]
Row 20:
Row 21: Money Earned: +[N] ([Reason])
Row 22:
Row 23: ────────────────────────────
Row 24:
Row 25: SELECT ACTION:
Row 26: > [Input prompt]
Row 27:
Row 28: [Error/Info messages]
Row 29:
Row 30: [Help: Press ? for commands]
```

---

## Rendering Strategy

### Approach 1: Full Re-render (Simpler)
**How it works:**
1. Clear screen
2. Redraw entire UI
3. Position cursor at input prompt

**Pros:**
- Simple to implement
- No state tracking needed
- Always consistent

**Cons:**
- Potential flicker
- More expensive (doesn't matter for turn-based game)

**Recommendation**: Use this approach initially.

### Approach 2: Partial Re-render (Optimized)
**How it works:**
1. Track dirty regions
2. Only update changed sections
3. More complex but smoother

**When to use**: If full re-render causes noticeable flicker.

---

## File Structure

### New Files to Create

**`ui/terminal.go`**
```go
package ui

// Terminal control and rendering
type Terminal struct {
    Width  int
    Height int
}

func NewTerminal() *Terminal
func (t *Terminal) Clear()
func (t *Terminal) MoveCursor(row, col int)
func (t *Terminal) GetSize() (width, height int)
func (t *Terminal) HideCursor()
func (t *Terminal) ShowCursor()
```

**`ui/layout.go`**
```go
package ui

// Layout constants and calculations
const (
    MinTerminalWidth  = 80
    MinTerminalHeight = 24

    LeftPanelWidth  = 60
    RightPanelWidth = 60

    PanelDivider = "│"
)

type Layout struct {
    LeftPanel  Panel
    RightPanel Panel
}

type Panel struct {
    X      int
    Y      int
    Width  int
    Height int
}

func NewLayout(termWidth, termHeight int) *Layout
```

**`ui/renderer.go`**
```go
package ui

// Main rendering logic
type Renderer struct {
    Terminal *Terminal
    Layout   *Layout
}

func NewRenderer() *Renderer
func (r *Renderer) Render(state GameState)
func (r *Renderer) RenderLeftPanel(state GameState)
func (r *Renderer) RenderRightPanel(state GameState)
func (r *Renderer) RenderBorder()
```

**`ui/gamestate.go`**
```go
package ui

// Data structure for rendering
type GameState struct {
    // Player info
    Player      Golfer

    // Hole info
    CurrentHole Hole
    HoleNumber  int
    TotalHoles  int

    // Shot info
    LastShot    *ShotResult

    // Score info
    ScoreCard   ScoreCard

    // Messages
    StatusMsg   string
    ErrorMsg    string
}

type ShotResult struct {
    Club        Club
    Outcome     SkillCheckOutcome
    Margin      int
    Rotation    float64
    Power       float64
    Distance    float64
    XPEarned    int
    LevelUps    []string
}
```

### Modified Files

**`main.go`**
- Replace all `fmt.Println()` with UI calls
- Create GameState after each action
- Call `renderer.Render(gameState)`
- Handle input with cursor positioning

---

## Color Scheme

Use colors sparingly for important information:

**Game Outcomes:**
- Critical Success: Bright Green
- Excellent/Good: Green
- Marginal: Yellow
- Poor: Yellow
- Bad/Critical Failure: Red

**Stats:**
- Skill/Ability names: Cyan
- Level up notifications: Bright Green
- Money earned: Yellow
- Current values: White (default)

**UI Elements:**
- Panel borders: Dim White
- Headers: Bright White
- Error messages: Red
- Success messages: Green

---

## Implementation Phases

### Phase 1: Core Terminal Control
**Goal**: Basic ANSI escape sequence support
1. Create `ui/terminal.go` with basic functions
2. Test clear screen, cursor movement
3. Detect terminal size
4. Handle Windows vs Unix differences

### Phase 2: Layout System
**Goal**: Define and render static layout
1. Create `ui/layout.go` with panel definitions
2. Implement border rendering
3. Test two-column layout
4. Ensure proper alignment

### Phase 3: Game State Rendering
**Goal**: Render actual game data
1. Create `ui/gamestate.go` with data structures
2. Implement `RenderRightPanel()` for player stats
3. Implement `RenderLeftPanel()` for game area
4. Test with real game data

### Phase 4: Main Loop Integration
**Goal**: Replace text output with UI
1. Modify `main.go` to use renderer
2. Replace all fmt.Println() calls
3. Update input handling for positioned cursor
4. Test full game flow

### Phase 5: Polish
**Goal**: Improve visual experience
1. Add color coding
2. Add box-drawing characters
3. Handle terminal resize gracefully
4. Add help text and keyboard shortcuts

---

## Testing Strategy

### Manual Testing
1. **Different terminal sizes**: 80×24, 120×30, 160×40
2. **Different terminals**: Windows Terminal, iTerm2, GNOME Terminal
3. **Different OS**: Windows, macOS, Linux
4. **Edge cases**: Very long names, max level stats

### Automated Testing
Focus on:
- Layout calculations (panel positions)
- State formatting (text alignment, truncation)
- Color code generation

Not needed:
- Actual terminal rendering (hard to test)

---

## Fallback Mode

If terminal doesn't support ANSI codes or is too small:
```go
func (r *Renderer) SupportsRichUI() bool {
    width, height := r.Terminal.GetSize()
    return width >= MinTerminalWidth &&
           height >= MinTerminalHeight &&
           r.Terminal.SupportsANSI()
}

// If false, fall back to simple text output
if !renderer.SupportsRichUI() {
    // Use old fmt.Println() style
}
```

---

## Example Usage

```go
// In main.go
func main() {
    renderer := ui.NewRenderer()
    defer renderer.Terminal.ShowCursor() // Ensure cursor restored

    renderer.Terminal.HideCursor()
    renderer.Terminal.Clear()

    // Game loop
    for _, hole := range course.Holes {
        for !holeComplete {
            // Update game state
            state := ui.GameState{
                Player:      golfer,
                CurrentHole: hole,
                HoleNumber:  holeNum,
                TotalHoles:  len(course.Holes),
                LastShot:    &lastShotResult,
                ScoreCard:   scoreCard,
            }

            // Render
            renderer.Render(state)

            // Get input (cursor already positioned)
            input := readInput()

            // Process action...
        }
    }

    renderer.Terminal.Clear()
    renderer.Terminal.ShowCursor()
}
```

---

## Success Criteria

The UX improvement is successful when:

1. ✅ **No scrolling during gameplay** - Screen updates in place
2. ✅ **Stats always visible** - Right panel never changes position
3. ✅ **Clear information hierarchy** - Game area vs player info separated
4. ✅ **Better readability** - Aligned columns, clear sections
5. ✅ **Smooth experience** - No flicker, responsive updates
6. ✅ **Cross-platform** - Works on Windows, macOS, Linux
7. ✅ **Graceful degradation** - Falls back to simple mode if needed

---

## Risk Mitigation

**Risk 1: Terminal Compatibility Issues**
- **Mitigation**: Detect and fall back to simple mode
- **Test on**: Windows Terminal, iTerm2, GNOME Terminal, macOS Terminal

**Risk 2: Windows ANSI Support**
- **Mitigation**: Enable virtual terminal processing on Windows 10+
- **Fallback**: Use simple mode on older Windows versions

**Risk 3: Terminal Size Variations**
- **Mitigation**: Detect size, adjust layout dynamically
- **Minimum size**: 80×24 enforced, warn if smaller

**Risk 4: Increased Complexity**
- **Mitigation**: Keep simple text mode as backup
- **Refactor**: UI code isolated in `ui/` package

---

## Future Enhancements

After basic split-panel works:

1. **Animations**: Smooth transitions for ball movement
2. **Course Visualization**: ASCII art representation of hole layout
3. **Shot Tracer**: Show ball path with ASCII graphics
4. **Color Themes**: Different color schemes (light/dark mode)
5. **Mouse Support**: Click on options instead of keyboard only
6. **ProShop UI**: Dedicated screen for equipment browsing
7. **Minimap**: Bird's eye view of hole in corner

---

## Next Steps

1. Create new branch: `feature/terminal-ui`
2. Implement Phase 1 (Core Terminal Control)
3. Test on multiple platforms
4. Implement Phase 2 (Layout System)
5. Continue through phases with TDD where applicable
6. Create PR when basic functionality works

**Estimated Effort**: 4-6 hours for basic implementation
**Priority**: HIGH - Significantly improves playability
