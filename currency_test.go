package gogolf

import "testing"

// Test Golfer has a Money field
func TestGolfer_HasMoney(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	if golfer.Money < 0 {
		t.Errorf("Golfer money = %d, want >= 0", golfer.Money)
	}
}

// Test Golfer starts with starter money
func TestGolfer_StarterMoney(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	// Players should start with some money to buy basic equipment
	expectedStarterMoney := 100
	if golfer.Money != expectedStarterMoney {
		t.Errorf("Golfer money = %d, want %d", golfer.Money, expectedStarterMoney)
	}
}

// Test AddMoney increases golfer's money
func TestGolfer_AddMoney(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	initialMoney := golfer.Money

	golfer.AddMoney(50)

	expectedMoney := initialMoney + 50
	if golfer.Money != expectedMoney {
		t.Errorf("After AddMoney(50), money = %d, want %d", golfer.Money, expectedMoney)
	}
}

// Test AddMoney can be called multiple times
func TestGolfer_AddMoney_Multiple(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	initialMoney := golfer.Money

	golfer.AddMoney(25)
	golfer.AddMoney(75)
	golfer.AddMoney(100)

	expectedMoney := initialMoney + 25 + 75 + 100
	if golfer.Money != expectedMoney {
		t.Errorf("After multiple AddMoney calls, money = %d, want %d", golfer.Money, expectedMoney)
	}
}

// Test SpendMoney decreases golfer's money
func TestGolfer_SpendMoney(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	golfer.AddMoney(100) // Ensure enough money
	initialMoney := golfer.Money

	success := golfer.SpendMoney(30)

	if !success {
		t.Error("SpendMoney(30) returned false, want true (enough money)")
	}

	expectedMoney := initialMoney - 30
	if golfer.Money != expectedMoney {
		t.Errorf("After SpendMoney(30), money = %d, want %d", golfer.Money, expectedMoney)
	}
}

// Test SpendMoney fails if not enough money
func TestGolfer_SpendMoney_InsufficientFunds(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	initialMoney := golfer.Money

	// Try to spend more than we have
	success := golfer.SpendMoney(initialMoney + 100)

	if success {
		t.Error("SpendMoney(too much) returned true, want false (insufficient funds)")
	}

	// Money should not change
	if golfer.Money != initialMoney {
		t.Errorf("After failed SpendMoney, money = %d, want %d (unchanged)", golfer.Money, initialMoney)
	}
}

// Test SpendMoney allows spending exact amount
func TestGolfer_SpendMoney_ExactAmount(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	golfer.Money = 100 // Set specific amount

	success := golfer.SpendMoney(100)

	if !success {
		t.Error("SpendMoney(exact amount) returned false, want true")
	}

	if golfer.Money != 0 {
		t.Errorf("After SpendMoney(all money), money = %d, want 0", golfer.Money)
	}
}
