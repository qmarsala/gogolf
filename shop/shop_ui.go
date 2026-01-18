package shop

import (
	"bufio"
	"fmt"
	"gogolf"
	"io"
	"strconv"
	"strings"
)

type ShopUI struct {
	shop   ProShop
	output io.Writer
	reader *bufio.Reader
}

func NewShopUI(shop ProShop, output io.Writer, input io.Reader) *ShopUI {
	return &ShopUI{
		shop:   shop,
		output: output,
		reader: bufio.NewReader(input),
	}
}

func FormatBallDisplay(ball gogolf.Ball) string {
	return fmt.Sprintf("%s - %d money (+%.0f distance, %.1f spin)",
		ball.Name, ball.Cost, ball.DistanceBonus, ball.SpinControl)
}

func FormatGloveDisplay(glove gogolf.Glove) string {
	return fmt.Sprintf("%s - %d money (+%.2f accuracy)",
		glove.Name, glove.Cost, glove.AccuracyBonus)
}

func FormatShoesDisplay(shoes gogolf.Shoes) string {
	return fmt.Sprintf("%s - %d money (-%d lie penalty)",
		shoes.Name, shoes.Cost, shoes.LiePenaltyReduction)
}

func (ui *ShopUI) printf(format string, args ...interface{}) {
	fmt.Fprintf(ui.output, format, args...)
}

func (ui *ShopUI) println(args ...interface{}) {
	fmt.Fprintln(ui.output, args...)
}

func (ui *ShopUI) readLine() string {
	input, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (ui *ShopUI) readInt(min, max int) int {
	for {
		input := ui.readLine()
		value, err := strconv.Atoi(input)
		if err != nil || value < min || value > max {
			ui.printf("Please enter a number between %d and %d\n", min, max)
			continue
		}
		return value
	}
}

func (ui *ShopUI) readYesNo() bool {
	input := strings.ToLower(ui.readLine())
	return input == "y" || input == "yes"
}

func (ui *ShopUI) Show(golfer *gogolf.Golfer) {
	for {
		ui.showMainMenu(golfer)
		choice := ui.readInt(1, 4)

		switch choice {
		case 1:
			ui.showBallsMenu(golfer)
		case 2:
			ui.showGlovesMenu(golfer)
		case 3:
			ui.showShoesMenu(golfer)
		case 4:
			return
		}
	}
}

func (ui *ShopUI) showMainMenu(golfer *gogolf.Golfer) {
	ui.printf("\n=== ProShop ===\n")
	ui.printf("Money: %d\n\n", golfer.Money)
	ui.println("  1. Balls")
	ui.println("  2. Gloves")
	ui.println("  3. Shoes")
	ui.println("  4. Back to Game")
	ui.println()
	ui.printf("> ")
}

func (ui *ShopUI) showBallsMenu(golfer *gogolf.Golfer) {
	for {
		ui.printf("\n=== Balls ===\n")
		ui.printf("Currently equipped: ")
		if golfer.Ball != nil {
			ui.printf("%s (+%.0f distance, %.1f spin)\n", golfer.Ball.Name, golfer.Ball.DistanceBonus, golfer.Ball.SpinControl)
		} else {
			ui.println("None")
		}
		ui.println()
		ui.println("Available:")

		for i, ball := range ui.shop.Balls {
			affordable := golfer.Money >= ball.Cost
			indicator := ""
			if !affordable {
				indicator = " [Cannot afford]"
			}
			ui.printf("  %d. %s%s\n", i+1, FormatBallDisplay(ball), indicator)
		}
		ui.printf("  %d. Back\n", len(ui.shop.Balls)+1)
		ui.println()
		ui.printf("> ")

		choice := ui.readInt(1, len(ui.shop.Balls)+1)

		if choice == len(ui.shop.Balls)+1 {
			return
		}

		selectedBall := ui.shop.Balls[choice-1]
		ui.handleBallPurchase(golfer, selectedBall)
	}
}

func (ui *ShopUI) handleBallPurchase(golfer *gogolf.Golfer, ball gogolf.Ball) {
	if golfer.Money < ball.Cost {
		ui.printf("\nNot enough money! You have %d but need %d.\n", golfer.Money, ball.Cost)
		return
	}

	ui.printf("\nPurchase %s for %d money? (y/n) ", ball.Name, ball.Cost)
	if ui.readYesNo() {
		if ui.shop.PurchaseBall(golfer, ball.Name) {
			ui.printf("Purchased %s!\n", ball.Name)
		}
	}
}

func (ui *ShopUI) showGlovesMenu(golfer *gogolf.Golfer) {
	for {
		ui.printf("\n=== Gloves ===\n")
		ui.printf("Currently equipped: ")
		if golfer.Glove != nil {
			ui.printf("%s (+%.2f accuracy)\n", golfer.Glove.Name, golfer.Glove.AccuracyBonus)
		} else {
			ui.println("None")
		}
		ui.println()
		ui.println("Available:")

		for i, glove := range ui.shop.Gloves {
			affordable := golfer.Money >= glove.Cost
			indicator := ""
			if !affordable {
				indicator = " [Cannot afford]"
			}
			ui.printf("  %d. %s%s\n", i+1, FormatGloveDisplay(glove), indicator)
		}
		ui.printf("  %d. Back\n", len(ui.shop.Gloves)+1)
		ui.println()
		ui.printf("> ")

		choice := ui.readInt(1, len(ui.shop.Gloves)+1)

		if choice == len(ui.shop.Gloves)+1 {
			return
		}

		selectedGlove := ui.shop.Gloves[choice-1]
		ui.handleGlovePurchase(golfer, selectedGlove)
	}
}

func (ui *ShopUI) handleGlovePurchase(golfer *gogolf.Golfer, glove gogolf.Glove) {
	if golfer.Money < glove.Cost {
		ui.printf("\nNot enough money! You have %d but need %d.\n", golfer.Money, glove.Cost)
		return
	}

	ui.printf("\nPurchase %s for %d money? (y/n) ", glove.Name, glove.Cost)
	if ui.readYesNo() {
		if ui.shop.PurchaseGlove(golfer, glove.Name) {
			ui.printf("Purchased %s!\n", glove.Name)
		}
	}
}

func (ui *ShopUI) showShoesMenu(golfer *gogolf.Golfer) {
	for {
		ui.printf("\n=== Shoes ===\n")
		ui.printf("Currently equipped: ")
		if golfer.Shoes != nil {
			ui.printf("%s (-%d lie penalty)\n", golfer.Shoes.Name, golfer.Shoes.LiePenaltyReduction)
		} else {
			ui.println("None")
		}
		ui.println()
		ui.println("Available:")

		for i, shoes := range ui.shop.Shoes {
			affordable := golfer.Money >= shoes.Cost
			indicator := ""
			if !affordable {
				indicator = " [Cannot afford]"
			}
			ui.printf("  %d. %s%s\n", i+1, FormatShoesDisplay(shoes), indicator)
		}
		ui.printf("  %d. Back\n", len(ui.shop.Shoes)+1)
		ui.println()
		ui.printf("> ")

		choice := ui.readInt(1, len(ui.shop.Shoes)+1)

		if choice == len(ui.shop.Shoes)+1 {
			return
		}

		selectedShoes := ui.shop.Shoes[choice-1]
		ui.handleShoesPurchase(golfer, selectedShoes)
	}
}

func (ui *ShopUI) handleShoesPurchase(golfer *gogolf.Golfer, shoes gogolf.Shoes) {
	if golfer.Money < shoes.Cost {
		ui.printf("\nNot enough money! You have %d but need %d.\n", golfer.Money, shoes.Cost)
		return
	}

	ui.printf("\nPurchase %s for %d money? (y/n) ", shoes.Name, shoes.Cost)
	if ui.readYesNo() {
		if ui.shop.PurchaseShoes(golfer, shoes.Name) {
			ui.printf("Purchased %s!\n", shoes.Name)
		}
	}
}
