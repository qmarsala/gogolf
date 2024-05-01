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
	golfer := Golfer{Clubs: DefaultClubs()}

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

			result := NewD6().SkillCheck(10)
			//how do we want to control the miss direction?
			rotationDirection := float64(1)
			if int(math.Abs(float64(result.Margin)))%2 == 0 {
				rotationDirection *= -1
			}
			rotationDegrees := float64(0)
			clubAcc := float64(club.AccuracyDegrees())
			if result.Success {
				possibleRotation := math.Min(random.Float64()*clubAcc, clubAcc)
				fmt.Println("Possible Rotation: ", possibleRotation)
				rotationDegrees = math.Max(possibleRotation-float64(result.Margin), 0)
			} else {
				baseMisHit := 45 * (1 - float64(club.Forgiveness))
				minimumMisHitRotation := random.Float64() * baseMisHit
				possibleRotation := math.Max(minimumMisHitRotation+clubAcc, clubAcc)
				fmt.Println("Possible Rotation: ", possibleRotation)
				rotationDegrees = math.Max(possibleRotation+math.Abs(float64(result.Margin)), 1)
				power = math.Max(power*(float64(club.Forgiveness)-(math.Abs(float64(result.Margin))/100)), 0.1)
			}

			directionToHole.Rotate(rotationDegrees * rotationDirection)
			ballPath := ball.ReceiveHit(club, float32(power), directionToHole)
			if club.Name != "Putter" {
				if rand.IntN(100)%2 == 0 {
					experimentWithShotSimpleShapes_Draw(&ball, ballPath, h)
				} else {
					experimentWithShotSimpleShapes_Fade(&ball, ballPath, h)
				}
			}
			fmt.Printf("Ball traveled %f\n", Unit(ballPath.Magnitude()).Yards())
			scoreCard.RecordStroke(h)
			fmt.Println("Success: ", result.Success, " ", result.Margin, " Rotation: ", rotationDegrees, " ", rotationDirection)
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
