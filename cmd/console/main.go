package main

import (
	"fmt"

	"gogolf"
	"gogolf/ui"
)

func buildGameState(ctx gogolf.GameContext, lastShot *ui.ShotDisplay, promptMsg string) ui.GameState {
	skills := make(map[string]ui.SkillDisplay)
	for name, skill := range ctx.Golfer.Skills {
		skills[name] = ui.SkillDisplay{
			Name:      skill.Name,
			Level:     skill.Level,
			Value:     skill.Value(),
			CurrentXP: skill.Experience,
			XPForNext: skill.ExperienceToNextLevel(),
		}
	}

	abilities := make(map[string]ui.AbilityDisplay)
	for name, ability := range ctx.Golfer.Abilities {
		abilities[name] = ui.AbilityDisplay{
			Name:      ability.Name,
			Level:     ability.Level,
			Value:     ability.Value(),
			CurrentXP: ability.Experience,
			XPForNext: ability.ExperienceToNextLevel(),
		}
	}

	equipment := ui.EquipmentDisplay{}
	if ctx.Golfer.Ball != nil {
		equipment.BallName = ctx.Golfer.Ball.Name
		equipment.BallBonus = fmt.Sprintf("+%.0f dist", ctx.Golfer.Ball.DistanceBonus)
	}
	if ctx.Golfer.Glove != nil {
		equipment.GloveName = ctx.Golfer.Glove.Name
		equipment.GloveBonus = fmt.Sprintf("+%.2f acc", ctx.Golfer.Glove.AccuracyBonus)
	}
	if ctx.Golfer.Shoes != nil {
		equipment.ShoesName = ctx.Golfer.Shoes.Name
		equipment.ShoesBonus = fmt.Sprintf("-%d lie pen", ctx.Golfer.Shoes.LiePenaltyReduction)
	}

	return ui.GameState{
		PlayerName:        ctx.Golfer.Name,
		Money:             ctx.Golfer.Money,
		Skills:            skills,
		Abilities:         abilities,
		Equipment:         equipment,
		HoleNumber:        ctx.Hole.Number,
		TotalHoles:        len(ctx.ScoreCard.Course.Holes),
		Par:               ctx.Hole.Par,
		HoleDistance:      float64(ctx.Hole.Distance),
		BallLie:           ctx.Lie.String(),
		BallLieDifficulty: ctx.Lie.DifficultyModifier(),
		DistanceToHole:    float64(ctx.Ball.Location.Distance(ctx.Hole.HoleLocation).Yards()),
		BallLocationX:     float64(ctx.Ball.Location.X),
		BallLocationY:     float64(ctx.Ball.Location.Y),
		HoleLocationX:     float64(ctx.Hole.HoleLocation.X),
		HoleLocationY:     float64(ctx.Hole.HoleLocation.Y),
		LastShot:          lastShot,
		TotalStrokes:      ctx.ScoreCard.TotalStrokes(),
		ScoreToPar:        ctx.ScoreCard.Score(),
		StrokesThisHole:   ctx.ScoreCard.TotalStrokesThisHole(ctx.Hole),
		PromptMsg:         promptMsg,
	}
}

func shotResultToDisplay(result gogolf.ShotResult) *ui.ShotDisplay {
	return &ui.ShotDisplay{
		ClubName:    result.ClubName,
		Outcome:     result.Outcome.String(),
		Margin:      result.Margin,
		Description: result.Description,
		Rotation:    result.Rotation,
		RotationDir: result.RotationDir,
		Power:       result.Power,
		Distance:    result.Distance,
		XPEarned:    result.XPEarned,
		LevelUps:    result.LevelUps,
	}
}

func main() {
	renderer := ui.NewRenderer()
	defer renderer.Terminal.ShowCursor()

	if !renderer.Layout.SupportsRichUI() {
		fmt.Println("Terminal too small for rich UI. Minimum 80x24 required.")
		fmt.Println("Falling back to simple mode...")
		return
	}

	renderer.Terminal.HideCursor()

	game := gogolf.NewGame("Player", 3)

	for !game.IsRoundComplete() {
		game.TeeUp()
		var lastShot *ui.ShotDisplay

		for !game.IsHoleComplete() {
			ctx := game.GetGameContext()
			state := buildGameState(ctx, lastShot, fmt.Sprintf("Using %s", ctx.CurrentClub.Name))
			renderer.Render(state)

			powerMeter := ui.NewPowerMeter(renderer)
			power := powerMeter.GetPower()

			result := game.TakeShot(power)
			lastShot = shotResultToDisplay(result)

			if result.TapIn {
				lastShot.Description += " (Tap in)"
			}

			if result.HoledOut {
				break
			}
		}

		reward := game.CompleteHole()
		ctx := game.GetGameContext()
		statusMsg := fmt.Sprintf("Hole %d Complete! %d strokes (%+d) | +%d money",
			ctx.Hole.Number, game.StrokesThisHole(), ctx.ScoreCard.ScoreThisHole(ctx.Hole), reward)
		state := buildGameState(ctx, lastShot, "Press any key to continue...")
		state.StatusMsg = statusMsg
		renderer.Render(state)

		renderer.Terminal.ShowCursor()
		ui.WaitForAnyKey()
		renderer.Terminal.HideCursor()

		game.NextHole()
	}

	renderer.Terminal.Clear()
	renderer.Terminal.ShowCursor()
	fmt.Println("\n=== Round Complete ===")
	fmt.Printf("Final Score: %d (%+d)\n", game.ScoreCard.TotalStrokes(), game.ScoreCard.Score())
	fmt.Printf("Money: %d\n\n", game.Golfer.Money)
	displayPlayerStats(game.Golfer)
}

func displayPlayerStats(golfer gogolf.Golfer) {
	fmt.Println("\n=== Player Stats ===")
	fmt.Println("Skills:")
	for _, skillName := range []string{"Driver", "Woods", "Long Irons", "Mid Irons", "Short Irons", "Wedges", "Putter"} {
		skill := golfer.Skills[skillName]
		xpToNext := skill.ExperienceToNextLevel()
		if xpToNext == 0 {
			fmt.Printf("  %s: Level %d (MAX) [Value: %d]\n", skill.Name, skill.Level, skill.Value())
		} else {
			fmt.Printf("  %s: Level %d [Value: %d] (%d/%d XP)\n", skill.Name, skill.Level, skill.Value(), skill.Experience, xpToNext)
		}
	}
	fmt.Println("\nAbilities:")
	for _, abilityName := range []string{"Strength", "Control", "Touch", "Mental"} {
		ability := golfer.Abilities[abilityName]
		xpToNext := ability.ExperienceToNextLevel()
		if xpToNext == 0 {
			fmt.Printf("  %s: Level %d (MAX) [Value: %d]\n", ability.Name, ability.Level, ability.Value())
		} else {
			fmt.Printf("  %s: Level %d [Value: %d] (%d/%d XP)\n", ability.Name, ability.Level, ability.Value(), ability.Experience, xpToNext)
		}
	}
	fmt.Println("===================")
}
