package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("\nWelcome to GoGolf.")

	golfer := Golfer{
		Name: "Joe",
		Skills: Skills{
			Recovery: 4,
			Driving:  4,
			Approach: 4,
			Chipping: 4,
			Putting:  4,
		},
		Abilities: Abilities{
			Strength:  4,
			Intellect: 4,
			Control:   4,
		},
	}

	playGolf(golfer)
}

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func playGolf(golfer Golfer) {
	//display current hole info
	//display current shot info
	// get input on shot
	// execute shot
	//when hole is complete, move to next hole
	//when course is complete, end game
	hole := NewHole(1, 4, TeeBox{
		Size:     Size{Length: 3, Width: 3},
		Location: Point{X: 0, Y: 0},
	}, Fairway{
		Size:     Size{Length: 300, Width: 30},
		Location: Point{X: 0, Y: 10},
	}, Green{
		Size:         Size{Length: 30, Width: 30},
		Location:     Point{X: 0, Y: 910},
		HoleLocation: Point{X: 0, Y: 916},
	})

	fmt.Println(hole)
	fmt.Println(hole.Green.HoleLocation)
	fmt.Println(hole.DistanceToHole(golfer.GolfBall))

	for golfer.GolfBall.Location.Y <= hole.Green.HoleLocation.Y {
		golfer.Swing()
		fmt.Println(golfer.GolfBall.Lie, golfer.GolfBall.Location)
		fmt.Println(hole.DistanceToHole(golfer.GolfBall))
	}
}
