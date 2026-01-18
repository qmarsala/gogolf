package main

import (
	"bufio"
	"fmt"
	"gogolf/ui"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

// things to explore:
// [ ] changing aim target - at least need to aim 'left' or 'right', but aiming shorter would be nice too
// as it would allow flexibility with a full,3/4,1/2,1/4 shot system. (need to remove typing raw power, its cumbersome, annoying, and too easy in a way as you can be very precise)
// [ ] wind - translate the final point like a draw for now, though this is a little different
// [ ] greens adding break

//todo:
// [ ] go through brain storm comments and create some exploration tasks
// [ ] refactor experimental code into longer term solutions

// buildGameState creates a GameState from current game data for UI rendering
func buildGameState(golfer Golfer, hole Hole, ball GolfBall, scoreCard ScoreCard, lastShot *ui.ShotDisplay, promptMsg string) ui.GameState {
	// Build skills map
	skills := make(map[string]ui.SkillDisplay)
	for name, skill := range golfer.Skills {
		skills[name] = ui.SkillDisplay{
			Name:      skill.Name,
			Level:     skill.Level,
			Value:     skill.Value(),
			CurrentXP: skill.Experience,
			XPForNext: skill.ExperienceToNextLevel(),
		}
	}

	// Build abilities map
	abilities := make(map[string]ui.AbilityDisplay)
	for name, ability := range golfer.Abilities {
		abilities[name] = ui.AbilityDisplay{
			Name:      ability.Name,
			Level:     ability.Level,
			Value:     ability.Value(),
			CurrentXP: ability.Experience,
			XPForNext: ability.ExperienceToNextLevel(),
		}
	}

	// Build equipment display
	equipment := ui.EquipmentDisplay{}
	if golfer.Ball != nil {
		equipment.BallName = golfer.Ball.Name
		equipment.BallBonus = fmt.Sprintf("+%.0f dist", golfer.Ball.DistanceBonus)
	}
	if golfer.Glove != nil {
		equipment.GloveName = golfer.Glove.Name
		equipment.GloveBonus = fmt.Sprintf("+%.2f acc", golfer.Glove.AccuracyBonus)
	}
	if golfer.Shoes != nil {
		equipment.ShoesName = golfer.Shoes.Name
		equipment.ShoesBonus = fmt.Sprintf("-%d lie pen", golfer.Shoes.LiePenaltyReduction)
	}

	// Get current lie
	lie := ball.GetLie(&hole)

	return ui.GameState{
		PlayerName:        golfer.Name,
		Money:             golfer.Money,
		Skills:            skills,
		Abilities:         abilities,
		Equipment:         equipment,
		HoleNumber:        hole.Number,
		TotalHoles:        len(scoreCard.Course.Holes),
		Par:               hole.Par,
		HoleDistance:      float64(hole.Distance),
		BallLie:           lie.String(),
		BallLieDifficulty: lie.DifficultyModifier(),
		DistanceToHole:    float64(ball.Location.Distance(hole.HoleLocation).Yards()),
		BallLocationX:     float64(ball.Location.X),
		BallLocationY:     float64(ball.Location.Y),
		HoleLocationX:     float64(hole.HoleLocation.X),
		HoleLocationY:     float64(hole.HoleLocation.Y),
		LastShot:          lastShot,
		TotalStrokes:      scoreCard.TotalStrokes(),
		ScoreToPar:        scoreCard.Score(),
		StrokesThisHole:   scoreCard.TotalStrokesThisHole(hole),
		PromptMsg:         promptMsg,
	}
}

func main() {
	// Initialize UI renderer
	renderer := ui.NewRenderer()
	defer renderer.Terminal.ShowCursor() // Ensure cursor is restored on exit

	// Check if terminal supports rich UI
	if !renderer.Layout.SupportsRichUI() {
		fmt.Println("Terminal too small for rich UI. Minimum 80x24 required.")
		fmt.Println("Falling back to simple mode...")
		// TODO: Implement simple fallback mode
		return
	}

	renderer.Terminal.HideCursor()

	ball := GolfBall{Location: Point{X: 0, Y: 0}}
	course, scoreCard := GenerateSimpleCourse(3) // Use course with lie system
	golfer := NewGolfer("Player")

	random := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	var lastShot *ui.ShotDisplay // Track last shot for display

	for _, h := range course.Holes {
		ball.TeeUp()
		lastShot = nil // Reset shot display for new hole

		for scoreCard.TotalStrokesThisHole(h) < 11 {
			// Render UI with current state
			club := golfer.GetBestClub(ball.Location.Distance(h.HoleLocation).Yards())
			state := buildGameState(golfer, h, ball, scoreCard, lastShot, fmt.Sprintf("Using %s - Enter power (0-1):", club.Name))
			renderer.Render(state)

			// Get power input from user
			// Position cursor after "> " and show it
			promptRow := renderer.Layout.LeftPanel.Height - 1
			renderer.Terminal.MoveCursor(promptRow, renderer.Layout.LeftPanel.X+4)
			renderer.Terminal.ShowCursor()

			reader := bufio.NewReader(os.Stdin)
			powerInput, _ := reader.ReadString('\n')
			renderer.Terminal.HideCursor()
			power, _ := strconv.ParseFloat(strings.TrimSpace(powerInput), 64)

			directionToHole := ball.Location.Direction(h.HoleLocation)

			// Get current lie and calculate difficulty modifier
			lie := ball.GetLie(&h)
			difficulty := lie.DifficultyModifier()

			// Calculate dynamic target number based on club, skill, ability, and lie
			targetNumber := golfer.CalculateTargetNumber(club, difficulty)
			result := golfer.SkillCheck(NewD6(), targetNumber)

			// Determine rotation direction (keep existing logic)
			rotationDirection := float64(1)
			if int(math.Abs(float64(result.Margin)))%2 == 0 {
				rotationDirection *= -1
			}

			// Apply equipment bonuses to club
			modifiedClub := golfer.GetModifiedClub(club)

			// Use tier-based calculations with modified club
			rotationDegrees := CalculateRotation(modifiedClub, result, random)
			power = CalculatePower(modifiedClub, power, result)

			// Award XP and detect level-ups
			skill := golfer.GetSkillForClub(club)
			ability := golfer.GetAbilityForClub(club)
			prevSkillLevel := skill.Level
			prevAbilityLevel := ability.Level

			xpAward := calculateXP(result.Outcome)
			golfer.AwardExperience(club, xpAward)

			// Check for level-ups and collect messages
			newSkill := golfer.GetSkillForClub(club)
			newAbility := golfer.GetAbilityForClub(club)
			var levelUps []string

			if newSkill.Level > prevSkillLevel {
				levelUps = append(levelUps, fmt.Sprintf("%s leveled up to %d!", newSkill.Name, newSkill.Level))
			}
			if newAbility.Level > prevAbilityLevel {
				levelUps = append(levelUps, fmt.Sprintf("%s leveled up to %d!", newAbility.Name, newAbility.Level))
			}

			// Execute shot
			directionToHole.Rotate(rotationDegrees * rotationDirection)
			ballPath := ball.ReceiveHit(modifiedClub, float32(power), directionToHole)
			if club.Name != "Putter" {
				if rand.IntN(100)%2 == 0 {
					experimentWithShotSimpleShapes_Draw(&ball, ballPath, h)
				} else {
					experimentWithShotSimpleShapes_Fade(&ball, ballPath, h)
				}
			}

			scoreCard.RecordStroke(h)

			// Build shot display for next render
			rotationDir := "right"
			if rotationDirection < 0 {
				rotationDir = "left"
			}
			lastShot = &ui.ShotDisplay{
				ClubName:    club.Name,
				Outcome:     result.Outcome.String(),
				Margin:      result.Margin,
				Description: GetShotQualityDescription(result),
				Rotation:    rotationDegrees,
				RotationDir: rotationDir,
				Power:       power,
				Distance:    float64(Unit(ballPath.Magnitude()).Yards()),
				XPEarned:    xpAward,
				LevelUps:    levelUps,
			}
			if h.DetectHoleOut(ball, ballPath) {
				break
			} else if h.DetectTapIn(ball) {
				scoreCard.RecordStroke(h)
				lastShot.Description += " (Tap in)"
				break
			}
		}

		// Award money for hole completion
		golfer.AwardHoleReward(h.Par, scoreCard.TotalStrokesThisHole(h))

		// Show hole completion screen
		reward := CalculateHoleReward(h.Par, scoreCard.TotalStrokesThisHole(h))
		statusMsg := fmt.Sprintf("Hole %d Complete! %d strokes (%+d) | +%d money",
			h.Number, scoreCard.TotalStrokesThisHole(h), scoreCard.ScoreThisHole(h), reward)
		state := buildGameState(golfer, h, ball, scoreCard, lastShot, "Press Enter to continue...")
		state.StatusMsg = statusMsg
		renderer.Render(state)

		// Wait for user to press enter
		// Position cursor after "> " and show it
		promptRow := renderer.Layout.LeftPanel.Height - 1
		renderer.Terminal.MoveCursor(promptRow, renderer.Layout.LeftPanel.X+4)
		renderer.Terminal.ShowCursor()

		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		renderer.Terminal.HideCursor()
	}

	// Show round complete screen
	renderer.Terminal.Clear()
	renderer.Terminal.ShowCursor()
	fmt.Println("\n=== Round Complete ===")
	fmt.Printf("Final Score: %d (%+d)\n", scoreCard.TotalStrokes(), scoreCard.Score())
	fmt.Printf("Money: %d\n\n", golfer.Money)
	displayPlayerStats(golfer)
}

// displayPlayerStats shows the golfer's current skills and abilities
func displayPlayerStats(golfer Golfer) {
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

// calculateXP determines experience points based on shot outcome
func calculateXP(outcome SkillCheckOutcome) int {
	switch outcome {
	case CriticalSuccess:
		return 15
	case Excellent:
		return 10
	case Good:
		return 7
	case Marginal:
		return 5
	case Poor:
		return 3
	case Bad:
		return 2
	case CriticalFailure:
		return 1
	default:
		return 1
	}
}

func experimentWithShotSimpleShapes_Draw(ball *GolfBall, ballPath Vector, h Hole) {
	fmt.Printf("pre draw ball: %+v | hole: %+v\n", ball.Location, h.HoleLocation)
	directionToHole := ball.PrevLocation.Direction(h.HoleLocation)
	drawRotationDegrees := -45
	if directionToHole.Y < 0 {
		drawRotationDegrees = 45
	}
	rotatedPath := ballPath.Rotate(float64(drawRotationDegrees))
	// this should probably be a factor of total distance
	// shorter shots can move as much as longer shots
	translationDistance := Yard(math.Max(rand.Float64()*3, 1)).Units()
	fmt.Println("Draw: ", translationDistance)
	ball.Location = ball.Location.Move(rotatedPath, float64(translationDistance))
}

func experimentWithShotSimpleShapes_Fade(ball *GolfBall, ballPath Vector, h Hole) {
	fmt.Printf("pre fade ball: %+v | hole: %+v\n", ball.Location, h.HoleLocation)
	directionToHole := ball.PrevLocation.Direction(h.HoleLocation)
	drawRotationDegrees := 45
	if directionToHole.Y < 0 {
		drawRotationDegrees = -45
	}
	rotatedPath := ballPath.Rotate(float64(drawRotationDegrees))
	translationDistance := Yard(math.Max(rand.Float64()*3, 1)).Units()
	fmt.Println("Fade: ", translationDistance)
	ball.Location = ball.Location.Move(rotatedPath, float64(translationDistance))
}
