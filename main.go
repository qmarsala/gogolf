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

func main() {
	fmt.Println("\nWelcome to GoGolf.")
	hole1 := *NewHole(1, 4, Point{X: 20, Y: int(Yard(423).Units())}, Size{})
	hole2 := *NewHole(2, 5, Point{X: -100, Y: int(Yard(523).Units())}, Size{})
	hole3 := *NewHole(3, 3, Point{X: 0, Y: int(Yard(123).Units())}, Size{})
	holes := []Hole{hole1, hole2, hole3}
	ball := GolfBall{Location: Point{X: 0, Y: 0}}
	course := Course{Holes: holes}
	scoreCard := ScoreCard{
		Course: course,
		Scores: map[int]int{},
	}

	driver := Club{Name: "Driver", Distance: 280, Accuracy: .8}
	sevenIron := Club{Name: "7 Iron", Distance: 170, Accuracy: .9}
	pitchingWedge := Club{Name: "PW", Distance: 140, Accuracy: .95}
	lobWedge := Club{Name: "LW", Distance: 100, Accuracy: .95}
	putter := Club{Name: "Putter", Distance: 40, Accuracy: 1}
	clubs := []Club{driver, sevenIron, pitchingWedge, lobWedge, putter}

	for _, h := range course.Holes {
		fmt.Printf("%+v (%+v)\n", scoreCard.TotalStrokes(), scoreCard.ScoreThrough(h.Number-1))
		fmt.Println(h)
		ball.TeeUp()
		for !h.CheckForBall(ball) && scoreCard.TotalStrokesThisHole(h) < 11 {
			fmt.Printf("distance to hole: %f\n", ball.Location.Distance(h.HoleLocation).Yards())
			c := readString("Select a club: ")
			clubChoice, _ := strconv.ParseInt(strings.TrimSpace(c), 10, 8)
			club := clubs[clubChoice]
			p := readString("power: ")
			power, _ := strconv.ParseFloat(strings.TrimSpace(p), 64)
			directionToHole := ball.Location.Direction(h.HoleLocation)

			result := NewD6().SkillCheck(10)
			var rotationDegrees float64 = 0
			rotationDirection := 1
			//how do we want to control the miss direction?
			if (rand.IntN(10)+1)%2 == 0 {
				rotationDirection *= -1
			}
			if result.Success {
				possibleRotation := math.Min(rand.Float64()*100, float64(club.AccuracyDegrees()))
				rotationDegrees = math.Max(possibleRotation-float64(result.Margin), 0)
			} else {
				possibleRotation := math.Min(rand.Float64()*100, float64(club.AccuracyDegrees())*1.3)
				rotationDegrees = math.Max(possibleRotation+float64(result.Margin), 1)
				power = math.Max(power*(.8-(math.Abs(float64(result.Margin))/100)), 0)
			}
			directionToHole.Rotate(float64(rotationDegrees) * float64(rotationDirection))
			ball.ReceiveHit(club, float32(power), directionToHole)
			scoreCard.RecordStroke(h)

			fmt.Println("Success: ", result.Success, " ", result.Margin, " Rotation: ", rotationDegrees, " ", rotationDirection)
			fmt.Printf("ball: %+v | hole: %+v\n", ball.Location, h.HoleLocation)
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
