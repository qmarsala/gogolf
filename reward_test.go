package main

import "testing"

// Test CalculateHoleReward returns correct money for birdie
func TestCalculateHoleReward_Birdie(t *testing.T) {
	par := 4
	strokes := 3 // Birdie = 1 under par

	reward := CalculateHoleReward(par, strokes)

	// Birdie should give a good reward (25 money)
	expectedReward := 25
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (birdie)",
			par, strokes, reward, expectedReward)
	}
}

// Test CalculateHoleReward returns correct money for par
func TestCalculateHoleReward_Par(t *testing.T) {
	par := 4
	strokes := 4 // Par

	reward := CalculateHoleReward(par, strokes)

	// Par should give decent reward (10 money)
	expectedReward := 10
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (par)",
			par, strokes, reward, expectedReward)
	}
}

// Test CalculateHoleReward returns correct money for bogey
func TestCalculateHoleReward_Bogey(t *testing.T) {
	par := 4
	strokes := 5 // Bogey = 1 over par

	reward := CalculateHoleReward(par, strokes)

	// Bogey should give small reward (5 money)
	expectedReward := 5
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (bogey)",
			par, strokes, reward, expectedReward)
	}
}

// Test CalculateHoleReward returns correct money for eagle
func TestCalculateHoleReward_Eagle(t *testing.T) {
	par := 5
	strokes := 3 // Eagle = 2 under par

	reward := CalculateHoleReward(par, strokes)

	// Eagle should give great reward (50 money)
	expectedReward := 50
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (eagle)",
			par, strokes, reward, expectedReward)
	}
}

// Test CalculateHoleReward returns minimal money for double bogey or worse
func TestCalculateHoleReward_DoubleBogey(t *testing.T) {
	par := 4
	strokes := 6 // Double bogey = 2 over par

	reward := CalculateHoleReward(par, strokes)

	// Double bogey or worse gives minimal reward (1 money)
	expectedReward := 1
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (double bogey)",
			par, strokes, reward, expectedReward)
	}
}

// Test CalculateHoleReward for triple bogey also gives minimal reward
func TestCalculateHoleReward_TripleBogey(t *testing.T) {
	par := 4
	strokes := 7 // Triple bogey

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 1
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (triple bogey)",
			par, strokes, reward, expectedReward)
	}
}

// Test CalculateHoleReward for hole-in-one (albatross)
func TestCalculateHoleReward_HoleInOne(t *testing.T) {
	par := 3
	strokes := 1 // Hole in one on par 3

	reward := CalculateHoleReward(par, strokes)

	// Hole in one should give massive reward (100 money)
	expectedReward := 100
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (hole in one)",
			par, strokes, reward, expectedReward)
	}
}

// Test AwardHoleReward increases golfer's money
func TestGolfer_AwardHoleReward(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	initialMoney := golfer.Money

	par := 4
	strokes := 4 // Par

	golfer.AwardHoleReward(par, strokes)

	expectedMoney := initialMoney + 10 // Par gives 10 money
	if golfer.Money != expectedMoney {
		t.Errorf("After AwardHoleReward(par), money = %d, want %d", golfer.Money, expectedMoney)
	}
}

// Test AwardHoleReward works for different scores
func TestGolfer_AwardHoleReward_MultipleHoles(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	initialMoney := golfer.Money

	// Hole 1: Birdie (25 money)
	golfer.AwardHoleReward(4, 3)

	// Hole 2: Par (10 money)
	golfer.AwardHoleReward(3, 3)

	// Hole 3: Bogey (5 money)
	golfer.AwardHoleReward(4, 5)

	expectedMoney := initialMoney + 25 + 10 + 5
	if golfer.Money != expectedMoney {
		t.Errorf("After 3 holes, money = %d, want %d", golfer.Money, expectedMoney)
	}
}
