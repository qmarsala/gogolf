package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("\nWelcome to GoGolf.")
	c := Club{
		Distance:       100,
		DefaultSkill:   SkillApproach,
		DefaultAbility: AbilityControl,
	}
	g := Golfer{
		Name: "",
		Skills: map[string]Stat{
			SkillApproach: {Level: 4},
		},
		Abilities: map[string]Stat{
			AbilityControl: {Level: 4},
		},
		GolfBall:  GolfBall{},
		Clubs:     []Club{c},
		ScoreCard: ScoreCard{},
	}
	g.Swing(g.Clubs[0])
	g.Swing(g.Clubs[0])
	g.Swing(g.Clubs[0])
	fmt.Println(g.GolfBall.Location)

	ballPosition := Point{X: 0, Y: 0} // Example initial position
	direction := Vector{X: 1, Y: 1}   // Example direction vector
	newBallPosition := MovePoint(ballPosition, direction, 200)
	fmt.Printf("New Ball Position: %+v\n", newBallPosition)
}

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
