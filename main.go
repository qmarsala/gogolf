package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("\nWelcome to GoGolf.")
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
}
