package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gogolf"
	"gogolf/game"
	"gogolf/ui"
)

func buildGameState(ctx game.Context, lastShot *ui.ShotDisplay, promptMsg string) ui.GameState {
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

func shotResultToDisplay(result game.ShotResult) *ui.ShotDisplay {
	return &ui.ShotDisplay{
		ClubName:      result.ClubName,
		IntendedShape: result.IntendedShape.String(),
		ActualShape:   result.ActualShape.String(),
		ShapeSuccess:  result.ShapeSuccess,
		Outcome:       result.Outcome.String(),
		Margin:        result.Margin,
		TargetNumber:  result.TargetNumber,
		DiceRolls:     result.DiceRolls,
		Description:   result.Description,
		Rotation:      result.Rotation,
		RotationDir:   result.RotationDir,
		Power:         result.Power,
		Distance:      result.Distance,
		XPEarned:      result.XPEarned,
		LevelUps:      result.LevelUps,
	}
}

func getSaveDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".gogolf_saves"
	}
	return filepath.Join(homeDir, ".gogolf_saves")
}

func showStartupMenu(saveManager *gogolf.SaveManager) *game.Game {
	for {
		options := []ui.MenuOption{
			{Label: "New Game", Value: "new"},
			{Label: "Load Game", Value: "load"},
			{Label: "Quit", Value: "quit"},
		}

		choice := ui.ShowMenu("GoGolf", options)

		switch options[choice].Value {
		case "new":
			name := ui.PromptString("Enter your golfer's name: ")
			if name == "" {
				name = "Player"
			}
			return game.New(name, 3)

		case "load":
			g := showLoadMenu(saveManager)
			if g != nil {
				return g
			}

		case "quit":
			fmt.Println("Goodbye!")
			os.Exit(0)
		}
	}
}

func showLoadMenu(saveManager *gogolf.SaveManager) *game.Game {
	slots := saveManager.ListSaveSlots()

	if len(slots) == 0 {
		fmt.Println("\nNo saved games found.")
		fmt.Println("Press Enter to continue...")
		ui.PromptString("")
		return nil
	}

	options := make([]ui.MenuOption, 0, len(slots)+1)
	for _, slot := range slots {
		label := fmt.Sprintf("Slot %d: %s (saved %s)",
			slot.Slot, slot.GolferName, slot.SavedAt.Format("Jan 2 15:04"))
		options = append(options, ui.MenuOption{Label: label, Value: fmt.Sprintf("%d", slot.Slot)})
	}
	options = append(options, ui.MenuOption{Label: "Back", Value: "back"})

	choice := ui.ShowMenu("Load Game", options)

	if options[choice].Value == "back" {
		return nil
	}

	slot := slots[choice].Slot
	golfer, err := saveManager.Load(slot)
	if err != nil {
		fmt.Printf("\nError loading save: %v\n", err)
		fmt.Println("Press Enter to continue...")
		ui.PromptString("")
		return nil
	}

	fmt.Printf("\nLoaded %s from slot %d\n", golfer.Name, slot)
	return game.NewFromGolfer(golfer, 3)
}

func showSaveMenu(saveManager *gogolf.SaveManager, golfer gogolf.Golfer) {
	options := make([]ui.MenuOption, 0, gogolf.MaxSaveSlots+1)

	for slot := 1; slot <= gogolf.MaxSaveSlots; slot++ {
		label := fmt.Sprintf("Slot %d: ", slot)
		if saveManager.SlotExists(slot) {
			slots := saveManager.ListSaveSlots()
			for _, s := range slots {
				if s.Slot == slot {
					label += fmt.Sprintf("%s (saved %s)", s.GolferName, s.SavedAt.Format("Jan 2 15:04"))
					break
				}
			}
		} else {
			label += "Empty"
		}
		options = append(options, ui.MenuOption{Label: label, Value: fmt.Sprintf("%d", slot)})
	}
	options = append(options, ui.MenuOption{Label: "Back", Value: "back"})

	choice := ui.ShowMenu("Save Game", options)

	if options[choice].Value == "back" {
		return
	}

	slot := choice + 1
	err := saveManager.Save(slot, golfer)
	if err != nil {
		fmt.Printf("\nError saving game: %v\n", err)
	} else {
		fmt.Printf("\nGame saved to slot %d\n", slot)
	}
	fmt.Println("Press Enter to continue...")
	ui.PromptString("")
}

func main() {
	saveManager := gogolf.NewSaveManager(getSaveDir())

	g := showStartupMenu(saveManager)

	renderer := ui.NewRenderer()
	defer renderer.Terminal.ShowCursor()

	if !renderer.Layout.SupportsRichUI() {
		fmt.Println("Terminal too small for rich UI. Minimum 80x24 required.")
		fmt.Println("Falling back to simple mode...")
		return
	}

	renderer.Terminal.HideCursor()

	for !g.IsRoundComplete() {
		g.TeeUp()
		var lastShot *ui.ShotDisplay

		for !g.IsHoleComplete() {
			ctx := g.GetContext()
			state := buildGameState(ctx, lastShot, fmt.Sprintf("Using %s", ctx.CurrentClub.Name))
			renderer.Render(state)

			var shape gogolf.ShotShape
			if ctx.CurrentClub.Name == "Putter" {
				shape = gogolf.Straight
			} else {
				shapeSelector := ui.NewShotShapeSelector(renderer)
				shape = shapeSelector.SelectShotShape()
			}

			modifiedClub := ctx.Golfer.GetModifiedClub(ctx.CurrentClub)
			powerMeter := ui.NewPowerMeter(renderer)
			powerMeter.SetClubDistance(float64(modifiedClub.Distance))
			power := powerMeter.GetPower()

			result := g.TakeShotWithShape(power, shape)

			diceRoller := ui.NewDiceRoller(renderer)
			diceRoller.ShowRoll(result.DiceRolls, result.TargetNumber)

			lastShot = shotResultToDisplay(result)

			if result.TapIn {
				lastShot.Description += " (Tap in)"
			}

			if result.HoledOut {
				break
			}
		}

		reward := g.CompleteHole()
		ctx := g.GetContext()
		statusMsg := fmt.Sprintf("Hole %d Complete! %d strokes (%+d) | +%d money",
			ctx.Hole.Number, g.StrokesThisHole(), ctx.ScoreCard.ScoreThisHole(ctx.Hole), reward)
		state := buildGameState(ctx, lastShot, "Press any key to continue...")
		state.StatusMsg = statusMsg
		renderer.Render(state)

		renderer.Terminal.ShowCursor()
		ui.WaitForAnyKey()
		renderer.Terminal.HideCursor()

		g.NextHole()
	}

	renderer.Terminal.Clear()
	renderer.Terminal.ShowCursor()
	fmt.Println("\n=== Round Complete ===")
	fmt.Printf("Final Score: %d (%+d)\n", g.ScoreCard.TotalStrokes(), g.ScoreCard.Score())
	fmt.Printf("Money: %d\n\n", g.Golfer.Money)
	displayPlayerStats(g.Golfer)

	showPostRoundMenu(saveManager, &g.Golfer)
}

func showPostRoundMenu(saveManager *gogolf.SaveManager, golfer *gogolf.Golfer) {
	proshop := gogolf.NewProShop()

	for {
		options := []ui.MenuOption{
			{Label: "Visit ProShop", Value: "shop"},
			{Label: "Save Game", Value: "save"},
			{Label: "Quit", Value: "quit"},
		}

		choice := ui.ShowMenu("What would you like to do?", options)

		switch options[choice].Value {
		case "shop":
			shopUI := ui.NewShopUI(proshop, os.Stdout, os.Stdin)
			shopUI.Show(golfer)
		case "save":
			showSaveMenu(saveManager, *golfer)
		case "quit":
			fmt.Println("Thanks for playing!")
			return
		}
	}
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
