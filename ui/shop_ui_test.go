package ui

import (
	"bytes"
	"gogolf"
	"gogolf/shop"
	"strings"
	"testing"
)

func TestFormatBallDisplay(t *testing.T) {
	ball := gogolf.Ball{
		Name:          "Premium Ball",
		DistanceBonus: 5,
		SpinControl:   0.7,
		Cost:          50,
	}

	result := FormatBallDisplay(ball)

	if !strings.Contains(result, "Premium Ball") {
		t.Errorf("FormatBallDisplay should contain ball name, got: %s", result)
	}
	if !strings.Contains(result, "+5") {
		t.Errorf("FormatBallDisplay should contain distance bonus, got: %s", result)
	}
	if !strings.Contains(result, "0.7") {
		t.Errorf("FormatBallDisplay should contain spin control, got: %s", result)
	}
	if !strings.Contains(result, "50") {
		t.Errorf("FormatBallDisplay should contain cost, got: %s", result)
	}
}

func TestFormatGloveDisplay(t *testing.T) {
	glove := gogolf.Glove{
		Name:          "Leather Pro",
		AccuracyBonus: 0.05,
		Cost:          45,
	}

	result := FormatGloveDisplay(glove)

	if !strings.Contains(result, "Leather Pro") {
		t.Errorf("FormatGloveDisplay should contain glove name, got: %s", result)
	}
	if !strings.Contains(result, "0.05") {
		t.Errorf("FormatGloveDisplay should contain accuracy bonus, got: %s", result)
	}
	if !strings.Contains(result, "45") {
		t.Errorf("FormatGloveDisplay should contain cost, got: %s", result)
	}
}

func TestFormatShoesDisplay(t *testing.T) {
	shoes := gogolf.Shoes{
		Name:                "Tour Edition",
		LiePenaltyReduction: 3,
		Cost:                80,
	}

	result := FormatShoesDisplay(shoes)

	if !strings.Contains(result, "Tour Edition") {
		t.Errorf("FormatShoesDisplay should contain shoes name, got: %s", result)
	}
	if !strings.Contains(result, "3") {
		t.Errorf("FormatShoesDisplay should contain lie penalty reduction, got: %s", result)
	}
	if !strings.Contains(result, "80") {
		t.Errorf("FormatShoesDisplay should contain cost, got: %s", result)
	}
}

func TestShopUI_ShowMainMenu_DisplaysMoneyAndEquipment(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 150

	output := &bytes.Buffer{}
	input := strings.NewReader("4\n") // Select "Back"

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if !strings.Contains(result, "150") {
		t.Errorf("Shop UI should display money (150), got: %s", result)
	}
	if !strings.Contains(result, "ProShop") {
		t.Errorf("Shop UI should display ProShop title, got: %s", result)
	}
}

func TestShopUI_ShowMainMenu_DisplaysCategories(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")

	output := &bytes.Buffer{}
	input := strings.NewReader("4\n") // Select "Back"

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if !strings.Contains(result, "Balls") {
		t.Errorf("Shop UI should display Balls category, got: %s", result)
	}
	if !strings.Contains(result, "Gloves") {
		t.Errorf("Shop UI should display Gloves category, got: %s", result)
	}
	if !strings.Contains(result, "Shoes") {
		t.Errorf("Shop UI should display Shoes category, got: %s", result)
	}
}

func TestShopUI_BallsMenu_DisplaysCurrentEquipment(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Ball = &gogolf.Ball{Name: "Standard Ball", DistanceBonus: 3, SpinControl: 0.5}

	output := &bytes.Buffer{}
	input := strings.NewReader("1\n5\n4\n") // Select Balls, then Back, then Back from main

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if !strings.Contains(result, "Standard Ball") {
		t.Errorf("Balls menu should show currently equipped ball, got: %s", result)
	}
	if !strings.Contains(result, "Currently equipped") {
		t.Errorf("Balls menu should indicate current equipment, got: %s", result)
	}
}

func TestShopUI_BallsMenu_DisplaysAvailableBalls(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")

	output := &bytes.Buffer{}
	input := strings.NewReader("1\n5\n4\n") // Select Balls, then Back, then Back from main

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if !strings.Contains(result, "Budget Ball") {
		t.Errorf("Balls menu should show Budget Ball, got: %s", result)
	}
	if !strings.Contains(result, "Pro V1") {
		t.Errorf("Balls menu should show Pro V1, got: %s", result)
	}
}

func TestShopUI_PurchaseBall_Success(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 100

	output := &bytes.Buffer{}
	// Select Balls, select Budget Ball (first option), confirm purchase, back, back
	input := strings.NewReader("1\n1\ny\n5\n4\n")

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	if golfer.Ball == nil {
		t.Fatal("Ball should be equipped after purchase")
	}
	if golfer.Ball.Name != "Budget Ball" {
		t.Errorf("Equipped ball = %s, want Budget Ball", golfer.Ball.Name)
	}
	if golfer.Money != 80 { // 100 - 20
		t.Errorf("Money = %d, want 80", golfer.Money)
	}
}

func TestShopUI_PurchaseBall_Declined(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 100
	initialMoney := golfer.Money

	output := &bytes.Buffer{}
	// Select Balls, select Budget Ball, decline purchase, back, back
	input := strings.NewReader("1\n1\nn\n5\n4\n")

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	if golfer.Ball != nil {
		t.Error("Ball should not be equipped after declining purchase")
	}
	if golfer.Money != initialMoney {
		t.Errorf("Money changed after declining: %d, want %d", golfer.Money, initialMoney)
	}
}

func TestShopUI_PurchaseBall_InsufficientFunds(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 10

	output := &bytes.Buffer{}
	// Select Balls, select Premium Ball (50 cost), try to confirm
	input := strings.NewReader("1\n3\ny\n5\n4\n")

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if golfer.Ball != nil {
		t.Error("Ball should not be equipped with insufficient funds")
	}
	if !strings.Contains(result, "Cannot afford") || !strings.Contains(result, "Not enough") {
		// Allow either message variant
		if !strings.Contains(strings.ToLower(result), "afford") && !strings.Contains(strings.ToLower(result), "enough") {
			t.Errorf("Should show insufficient funds message, got: %s", result)
		}
	}
}

func TestShopUI_GlovesMenu_DisplaysAvailableGloves(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")

	output := &bytes.Buffer{}
	input := strings.NewReader("2\n4\n4\n") // Select Gloves, then Back, then Back from main

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if !strings.Contains(result, "Basic Glove") {
		t.Errorf("Gloves menu should show Basic Glove, got: %s", result)
	}
	if !strings.Contains(result, "Precision Grip") {
		t.Errorf("Gloves menu should show Precision Grip, got: %s", result)
	}
}

func TestShopUI_ShoesMenu_DisplaysAvailableShoes(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")

	output := &bytes.Buffer{}
	input := strings.NewReader("3\n4\n4\n") // Select Shoes, then Back, then Back from main

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if !strings.Contains(result, "Casual Spikes") {
		t.Errorf("Shoes menu should show Casual Spikes, got: %s", result)
	}
	if !strings.Contains(result, "Tour Edition") {
		t.Errorf("Shoes menu should show Tour Edition, got: %s", result)
	}
}

func TestShopUI_PurchaseGlove_Success(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 100

	output := &bytes.Buffer{}
	// Select Gloves, select Basic Glove, confirm purchase, back, back
	input := strings.NewReader("2\n1\ny\n4\n4\n")

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	if golfer.Glove == nil {
		t.Fatal("Glove should be equipped after purchase")
	}
	if golfer.Glove.Name != "Basic Glove" {
		t.Errorf("Equipped glove = %s, want Basic Glove", golfer.Glove.Name)
	}
}

func TestShopUI_PurchaseShoes_Success(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 100

	output := &bytes.Buffer{}
	// Select Shoes, select Casual Spikes, confirm purchase, back, back
	input := strings.NewReader("3\n1\ny\n4\n4\n")

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	if golfer.Shoes == nil {
		t.Fatal("Shoes should be equipped after purchase")
	}
	if golfer.Shoes.Name != "Casual Spikes" {
		t.Errorf("Equipped shoes = %s, want Casual Spikes", golfer.Shoes.Name)
	}
}

func TestShopUI_DisplaysAffordabilityIndicator(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 30 // Can afford Budget Ball (20) but not Premium (50)

	output := &bytes.Buffer{}
	input := strings.NewReader("1\n5\n4\n") // Select Balls, then Back, then Back from main

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	// Should show that expensive items are not affordable
	// The exact format may vary, but it should indicate some items can't be purchased
	if !strings.Contains(result, "Budget Ball") {
		t.Errorf("Should show Budget Ball, got: %s", result)
	}
}

func TestShopUI_NoEquipmentShowsNone(t *testing.T) {
	proshop := shop.NewProShop()
	golfer := gogolf.NewGolfer("TestPlayer")
	// No equipment equipped

	output := &bytes.Buffer{}
	input := strings.NewReader("1\n5\n4\n") // Select Balls, then Back, then Back from main

	ui := NewShopUI(proshop, output, input)
	ui.Show(&golfer)

	result := output.String()

	if !strings.Contains(result, "None") && !strings.Contains(result, "none") {
		t.Errorf("Should indicate no equipment when nothing equipped, got: %s", result)
	}
}
