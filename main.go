package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("\nWelcome to GoGolf.")
	fmt.Println("Please select a skill package:")
	skill := readString(`
	1) W 5 - I 3 - W 2 - P 2
	2) W 3 - I 3 - W 3 - P 3
	3) W 2 - I 5 - W 2 - P 3
	4) W 2 - I 2 - W 4 - P 4 
	`)

	fmt.Println("\nPlease select an ability package:")
	ability := readString(`
	1) D 3 - A 2 - C 1
	2) D 2 - A 2 - C 2
	3) D 1 - A 3 - C 2
	4) D 1 - A 1 - C 4 
	`)

	println(skill)
	println(ability)

	golfer := Golfer{}
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
	golfer.PlayShot(
		Shot{
			Club: Club{
				Name:          "Pitching Wedge",
				Id:            10,
				StockDistance: 100,
			},
			Shape: 0,
			Power: 50,
		})
}
