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

	for ball.Location != hole1.HoleLocation {
		fmt.Printf("distance to hole: %f\n", ball.Location.Distance(hole1.HoleLocation).Yards())
		d := readString("Enter distance: ")
		distance, _ := strconv.ParseFloat(strings.TrimSpace(d), 64)
		directionToHole := ball.Location.Direction(hole1.HoleLocation)
		ball.ReceiveHit(Club{Distance: Yard(distance)}, 1, directionToHole)
		scoreCard.RecordStroke(holes[0])
		fmt.Printf("%+v (%+v)\n", scoreCard.TotalStrokes(), scoreCard.Score())
	}
}

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
