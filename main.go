package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("\nWelcome to GoGolf.")
	ballPosition := Point{X: 1, Y: 2}
	holePosition := Point{X: 4, Y: 6}

	directionToHole := ballPosition.Direction(holePosition)

	fmt.Printf("Direction Vector to Hole: %+v\n", directionToHole)
	ball := GolfBall{Location: ballPosition}
	ball.ReceiveHit(Club{Distance: 150}, 1, directionToHole)
	fmt.Printf("ball pos: %+v\n", ball.Location)

}

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
