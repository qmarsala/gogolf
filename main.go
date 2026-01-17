package main

import (
	"bufio"
	"fmt"
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

func main() {
	fmt.Println("\nWelcome to GoGolf.")
	ball := GolfBall{Location: Point{X: 0, Y: 0}}
	course, scoreCard := GenerateCourse(3)
	golfer := NewGolfer("Player")

	// Display initial player stats
	displayPlayerStats(golfer)

	random := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	for _, h := range course.Holes {
		fmt.Printf("%+v (%+v)\n", scoreCard.TotalStrokes(), scoreCard.ScoreThrough(h.Number-1))
		fmt.Println(h)
		ball.TeeUp()
		for scoreCard.TotalStrokesThisHole(h) < 11 {
			distance := ball.Location.Distance(h.HoleLocation).Yards()
			fmt.Printf("distance to hole: %f\n", distance)
			//this can be a good default, but would still want a way to change it in the game
			club := golfer.GetBestClub(distance)
			fmt.Println("Using ", club.Name)
			//still need a better way to do this, something like a 'select a shot' and have a few options
			// like full, 3/4, 1/2, 1/4. as well as things like 'draw' and 'fade' or 'straight'
			// though, when putting, we may need an option like 'tap in' that just adds a stroke and finishes the hole
			// when the ball is within a certain range.  perhaps this could be part of hole out logic. and it auto taps in if the ball is close.
			p := readString("power: ")
			power, _ := strconv.ParseFloat(strings.TrimSpace(p), 64)
			directionToHole := ball.Location.Direction(h.HoleLocation)

			// Calculate dynamic target number based on club, skill, and ability
			difficulty := 0 // TODO: Will come from lie system in Phase 3
			targetNumber := golfer.CalculateTargetNumber(club, difficulty)
			result := golfer.SkillCheck(NewD6(), targetNumber)

			// Determine rotation direction (keep existing logic)
			rotationDirection := float64(1)
			if int(math.Abs(float64(result.Margin)))%2 == 0 {
				rotationDirection *= -1
			}

			// Use tier-based calculations
			rotationDegrees := CalculateRotation(club, result, random)
			power = CalculatePower(club, power, result)

			// Enhanced feedback
			fmt.Printf("\nShot Quality: %s (Margin: %+d)\n", result.Outcome, result.Margin)
			fmt.Printf("├─ %s\n", GetShotQualityDescription(result))
			fmt.Printf("├─ Rotation: %.1f° %s\n", rotationDegrees, map[bool]string{true: "left", false: "right"}[rotationDirection < 0])
			fmt.Printf("└─ Power: %.0f%%\n\n", power*100)

			// Award XP and detect level-ups
			skill := golfer.GetSkillForClub(club)
			ability := golfer.GetAbilityForClub(club)
			prevSkillLevel := skill.Level
			prevAbilityLevel := ability.Level

			xpAward := calculateXP(result.Outcome)
			golfer.AwardExperience(club, xpAward)

			// Check for level-ups
			newSkill := golfer.GetSkillForClub(club)
			newAbility := golfer.GetAbilityForClub(club)

			if newSkill.Level > prevSkillLevel {
				fmt.Printf("🎉 %s leveled up to %d!\n", newSkill.Name, newSkill.Level)
			}
			if newAbility.Level > prevAbilityLevel {
				fmt.Printf("🎉 %s leveled up to %d!\n", newAbility.Name, newAbility.Level)
			}

			fmt.Printf("XP: +%d (%s: %d/%d, %s: %d/%d)\n\n",
				xpAward,
				newSkill.Name, newSkill.Experience, newSkill.ExperienceToNextLevel(),
				newAbility.Name, newAbility.Experience, newAbility.ExperienceToNextLevel())

			directionToHole.Rotate(rotationDegrees * rotationDirection)
			ballPath := ball.ReceiveHit(club, float32(power), directionToHole)
			if club.Name != "Putter" {
				if rand.IntN(100)%2 == 0 {
					experimentWithShotSimpleShapes_Draw(&ball, ballPath, h)
				} else {
					experimentWithShotSimpleShapes_Fade(&ball, ballPath, h)
				}
			}
			fmt.Printf("Ball traveled %.1f yards\n", Unit(ballPath.Magnitude()).Yards())
			scoreCard.RecordStroke(h)
			fmt.Printf("Result: %s | Rotation: %.1f° %s\n",
				result.Outcome,
				rotationDegrees,
				map[bool]string{true: "left", false: "right"}[rotationDirection < 0])
			fmt.Printf("ball: %+v | hole: %+v\n", ball.Location, h.HoleLocation)
			if h.DetectHoleOut(ball, ballPath) {
				break
			} else if h.DetectTapIn(ball) {
				scoreCard.RecordStroke(h)
				fmt.Println("tap in")
				break
			}
		}
		fmt.Println("Hole Completed: ", scoreCard.TotalStrokesThisHole(h), " (", scoreCard.ScoreThisHole(h), ")")
	}
	fmt.Println("Score: ", scoreCard.TotalStrokes(), "(", scoreCard.Score(), ")")

	// Display final player stats
	fmt.Println("\n=== Round Complete ===")
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

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
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
