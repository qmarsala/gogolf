package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MenuOption struct {
	Label string
	Value string
}

func ShowMenu(title string, options []MenuOption) int {
	fmt.Printf("\n=== %s ===\n\n", title)
	for i, opt := range options {
		fmt.Printf("  %d. %s\n", i+1, opt.Label)
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Printf("Please enter a number between 1 and %d\n", len(options))
			continue
		}

		return choice - 1
	}
}

func PromptString(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func PromptInt(prompt string, min, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		value, err := strconv.Atoi(input)
		if err != nil || value < min || value > max {
			fmt.Printf("Please enter a number between %d and %d\n", min, max)
			continue
		}

		return value
	}
}
