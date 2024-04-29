package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("\nWelcome to GoGolf.")
	hole1 := *NewHole(1, 4, Point{X: 20, Y: int(Yard(423).Units())}, Size{})
	holes := []Hole{hole1}
	ball := GolfBall{Location: Point{X: 0, Y: 0}}
	scoreCard := ScoreCard{
		Holes:  holes,
		Scores: map[int]int{},
	}

	driver := Club{Name: "Driver", Distance: 280}
	sevenIron := Club{Name: "7 Iron", Distance: 170}
	pitchingWedge := Club{Name: "PW", Distance: 140}
	lobWedge := Club{Name: "LW", Distance: 100}
	putter := Club{Name: "LW", Distance: 40}
	clubs := []Club{driver, sevenIron, pitchingWedge, lobWedge, putter}

	//for 'hole out' logic, we should scan the path of the ball and the hole location
	// for collision. Then if it was not traveling to far past, it could be considered in
	//todo: collision detection
	// for now the ball's receive hit could return a vector representing the path taken
	// but it will eventually need to be an actually line that may curve
	for !hole1.CheckForBall(ball) {
		fmt.Printf("distance to hole: %f\n", ball.Location.Distance(hole1.HoleLocation).Yards())
		c := readString("Select a club: ")
		clubChoice, _ := strconv.ParseInt(strings.TrimSpace(c), 10, 8)
		club := clubs[clubChoice]
		p := readString("power: ")
		power, _ := strconv.ParseFloat(strings.TrimSpace(p), 64)
		directionToHole := ball.Location.Direction(hole1.HoleLocation)
		ball.ReceiveHit(club, float32(power), directionToHole)
		scoreCard.RecordStroke(holes[0])
		fmt.Printf("%+v (%+v)\n", scoreCard.TotalStrokes(), scoreCard.Score())
		fmt.Printf("ball: %+v | hole: %+v\n", ball.Location, hole1.HoleLocation)
	}
}

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
