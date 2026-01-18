# GoGolf Implementation Plan

## Completed: Power Meter & Putting Improvements

### Summary

Implemented putting-specific mechanic (Option C):
- **Green-only putting**: Putter only available when on Green lie
- **Feet-based display**: Distance shown in feet when on green
- **Auto-scaled power meter**: Meter scales to putt distance + 50% buffer (min 10 feet)
- Other clubs unchanged (yards, full power range)

### Changes Made

1. **GetBestClubForLie()** in [golfer.go](../golfer.go) - restricts putter to green only
2. **SetPuttingMode()** in [ui/input.go](../ui/input.go) - auto-scales meter for putting
3. **IsOnGreen** in [ui/gamestate.go](../ui/gamestate.go) - tracks when on green
4. **Renderer** in [ui/renderer.go](../ui/renderer.go) - shows feet when on green
5. **Game loop** in [cmd/console/main.go](../cmd/console/main.go) - uses putting mode when putter selected

---

## Future Consideration: Grip-Down for Other Clubs

The short-shot problem (needing very fast spacebar timing) still affects wedges for short chips. If this becomes an issue, consider adding a grip-down/power-scaling option for all clubs (Option A or D from original plan).