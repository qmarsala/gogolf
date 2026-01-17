package ui

import (
	"fmt"
	"strings"
)

// colorizeOutcome returns colored text based on shot outcome
func colorizeOutcome(outcome string) string {
	switch outcome {
	case "Critical Success":
		return colorBrightGreen + outcome + colorReset
	case "Excellent", "Good":
		return colorGreen + outcome + colorReset
	case "Marginal":
		return colorYellow + outcome + colorReset
	case "Poor":
		return colorYellow + outcome + colorReset
	case "Bad", "Critical Failure":
		return colorRed + outcome + colorReset
	default:
		return outcome
	}
}

// colorizeMoney returns colored money value
func colorizeMoney(value int) string {
	return colorYellow + fmt.Sprintf("%d", value) + colorReset
}

// colorizeXP returns colored XP value
func colorizeXP(value int) string {
	return colorCyan + fmt.Sprintf("+%d", value) + colorReset
}

// colorizeLevelUp returns colored level up message
func colorizeLevelUp(message string) string {
	return colorBrightGreen + "🎉 " + message + colorReset
}

// Renderer manages the rendering of game state to the terminal
type Renderer struct {
	Terminal *Terminal
	Layout   *Layout
}

// NewRenderer creates a new Renderer
func NewRenderer() *Renderer {
	term := NewTerminal()
	layout := NewLayout(term.Width, term.Height)

	return &Renderer{
		Terminal: term,
		Layout:   layout,
	}
}

// Render renders the complete game state to the terminal
func (r *Renderer) Render(state GameState) {
	if !r.Layout.SupportsRichUI() {
		// Fall back to simple mode (not implemented yet)
		return
	}

	r.Terminal.Clear()
	r.RenderBorder()
	r.RenderLeftPanel(state)
	r.RenderRightPanel(state)

	// Position cursor at prompt location
	promptRow := r.Layout.LeftPanel.Height - 2
	r.Terminal.MoveCursor(promptRow, r.Layout.LeftPanel.X+2)
}

// RenderBorder draws the vertical divider between panels
func (r *Renderer) RenderBorder() {
	dividerCol := r.Layout.GetDividerColumn()

	for row := 1; row <= r.Layout.TermHeight; row++ {
		r.Terminal.MoveCursor(row, dividerCol)
		fmt.Print(PanelDivider)
	}
}

// RenderLeftPanel renders the game area (left panel)
func (r *Renderer) RenderLeftPanel(state GameState) {
	panel := r.Layout.LeftPanel
	row := panel.Y

	// Header
	r.printInPanel(panel, row, fmt.Sprintf("=== HOLE %d - PAR %d ===", state.HoleNumber, state.Par), true)
	row++
	r.printInPanel(panel, row, fmt.Sprintf("Distance: %.0f yards", state.HoleDistance), false)
	row++
	row++ // blank line

	// Current lie
	difficultyStr := ""
	if state.BallLieDifficulty != 0 {
		difficultyStr = fmt.Sprintf(" (difficulty: %+d)", state.BallLieDifficulty)
	}
	r.printInPanel(panel, row, fmt.Sprintf("Current Lie: %s%s", state.BallLie, difficultyStr), false)
	row++
	r.printInPanel(panel, row, fmt.Sprintf("Distance to hole: %.1f yards", state.DistanceToHole), false)
	row++
	row++ // blank line

	// Last shot info
	if state.LastShot != nil {
		shot := state.LastShot
		r.printInPanel(panel, row, "Last Shot:", false)
		row++
		r.printInPanel(panel, row, fmt.Sprintf("├─ Club: %s", shot.ClubName), false)
		row++
		r.printInPanel(panel, row, fmt.Sprintf("├─ Quality: %s (Margin: %+d)", colorizeOutcome(shot.Outcome), shot.Margin), false)
		row++
		r.printInPanel(panel, row, fmt.Sprintf("├─ Description: %s", shot.Description), false)
		row++
		r.printInPanel(panel, row, fmt.Sprintf("├─ Rotation: %.1f° %s", shot.Rotation, shot.RotationDir), false)
		row++
		r.printInPanel(panel, row, fmt.Sprintf("├─ Power: %.0f%%", shot.Power*100), false)
		row++
		r.printInPanel(panel, row, fmt.Sprintf("└─ Distance: %.1f yards", shot.Distance), false)
		row++
		row++ // blank line

		// Ball location
		r.printInPanel(panel, row, fmt.Sprintf("Ball Location: (%.1f, %.1f)", state.BallLocationX, state.BallLocationY), false)
		row++
		r.printInPanel(panel, row, fmt.Sprintf("Hole Location: (%.1f, %.1f)", state.HoleLocationX, state.HoleLocationY), false)
		row++
		row++ // blank line

		// XP and level ups
		if shot.XPEarned > 0 {
			r.printInPanel(panel, row, fmt.Sprintf("XP Earned: %s", colorizeXP(shot.XPEarned)), false)
			row++
			for _, levelUp := range shot.LevelUps {
				r.printInPanel(panel, row, colorizeLevelUp(levelUp), false)
				row++
			}
			row++ // blank line
		}
	}

	// Separator
	r.printInPanel(panel, row, strings.Repeat(HorizontalLine, panel.Width-4), false)
	row++
	row++ // blank line

	// Status/error messages
	if state.ErrorMsg != "" {
		r.printInPanel(panel, row, state.ErrorMsg, false)
		row++
	}
	if state.StatusMsg != "" {
		r.printInPanel(panel, row, state.StatusMsg, false)
		row++
	}

	// Prompt at bottom
	promptRow := panel.Height - 2
	r.printInPanel(panel, promptRow, state.PromptMsg, false)
	r.printInPanel(panel, promptRow+1, "> ", false)
}

// RenderRightPanel renders the player stats (right panel)
func (r *Renderer) RenderRightPanel(state GameState) {
	panel := r.Layout.RightPanel
	row := panel.Y

	// Header
	r.printInPanel(panel, row, "=== PLAYER STATS ===", true)
	row++
	r.printInPanel(panel, row, fmt.Sprintf("Player: %s", state.PlayerName), false)
	row++
	r.printInPanel(panel, row, fmt.Sprintf("Money: %s", colorizeMoney(state.Money)), false)
	row++
	row++ // blank line

	// Skills
	r.printInPanel(panel, row, "--- Skills ---", false)
	row++
	skillOrder := []string{"Driver", "Woods", "Long Irons", "Mid Irons", "Short Irons", "Wedges", "Putter"}
	for _, skillName := range skillOrder {
		if skill, ok := state.Skills[skillName]; ok {
			r.printInPanel(panel, row, fmt.Sprintf("%s: Lvl %d [%d] (%d/%d)",
				skill.Name, skill.Level, skill.Value, skill.CurrentXP, skill.XPForNext), false)
			row++
		}
	}
	row++ // blank line

	// Abilities
	r.printInPanel(panel, row, "--- Abilities ---", false)
	row++
	abilityOrder := []string{"Strength", "Control", "Touch", "Mental"}
	for _, abilityName := range abilityOrder {
		if ability, ok := state.Abilities[abilityName]; ok {
			r.printInPanel(panel, row, fmt.Sprintf("%s: Lvl %d [%d] (%d/%d)",
				ability.Name, ability.Level, ability.Value, ability.CurrentXP, ability.XPForNext), false)
			row++
		}
	}
	row++ // blank line

	// Equipment
	r.printInPanel(panel, row, "--- Equipment ---", false)
	row++
	if state.Equipment.BallName != "" {
		r.printInPanel(panel, row, fmt.Sprintf("Ball: %s (%s)", state.Equipment.BallName, state.Equipment.BallBonus), false)
	} else {
		r.printInPanel(panel, row, "Ball: None", false)
	}
	row++
	if state.Equipment.GloveName != "" {
		r.printInPanel(panel, row, fmt.Sprintf("Glove: %s (%s)", state.Equipment.GloveName, state.Equipment.GloveBonus), false)
	} else {
		r.printInPanel(panel, row, "Glove: None", false)
	}
	row++
	if state.Equipment.ShoesName != "" {
		r.printInPanel(panel, row, fmt.Sprintf("Shoes: %s (%s)", state.Equipment.ShoesName, state.Equipment.ShoesBonus), false)
	} else {
		r.printInPanel(panel, row, "Shoes: None", false)
	}
	row++
	row++ // blank line

	// Score
	r.printInPanel(panel, row, "--- Score ---", false)
	row++
	scoreStr := fmt.Sprintf("%+d", state.ScoreToPar)
	if state.ScoreToPar == 0 {
		scoreStr = "E"
	}
	r.printInPanel(panel, row, fmt.Sprintf("Total: %d (%s)", state.TotalStrokes, scoreStr), false)
	row++
	r.printInPanel(panel, row, fmt.Sprintf("This Hole: %d strokes", state.StrokesThisHole), false)
	row++
	r.printInPanel(panel, row, fmt.Sprintf("Holes: %d/%d", state.HoleNumber, state.TotalHoles), false)
}

// printInPanel prints text at the specified row within a panel, truncating if necessary
func (r *Renderer) printInPanel(panel Panel, row int, text string, centered bool) {
	if row > panel.Height {
		return // Don't print beyond panel height
	}

	maxWidth := panel.Width - 4 // Leave margin on both sides
	if len(text) > maxWidth {
		text = text[:maxWidth-3] + "..."
	}

	if centered {
		padding := (maxWidth - len(text)) / 2
		text = strings.Repeat(" ", padding) + text
	}

	r.Terminal.MoveCursor(row, panel.X+2) // +2 for left margin
	fmt.Print(text)
}
