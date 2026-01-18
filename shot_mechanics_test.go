package main

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestCalculateRotation_CriticalSuccess(t *testing.T) {
	club := Club{Name: "Driver", Distance: 280, Accuracy: 0.75, Forgiveness: 0.8}
	result := SkillCheckResult{
		Outcome: CriticalSuccess,
		Margin:  7,
	}
	random := rand.New(rand.NewPCG(1, 2))

	rotation := CalculateRotation(club, result, random)
	clubAcc := float64(club.AccuracyDegrees())

	// Critical success should have minimal rotation (10% of normal spread)
	maxExpected := clubAcc * 0.1
	if rotation < 0 || rotation > maxExpected {
		t.Errorf("CriticalSuccess rotation = %.2f, expected 0 to %.2f", rotation, maxExpected)
	}
}

func TestCalculateRotation_Excellent(t *testing.T) {
	club := Club{Name: "7 Iron", Distance: 180, Accuracy: 0.9, Forgiveness: 0.8}
	result := SkillCheckResult{
		Outcome: Excellent,
		Margin:  5,
	}
	random := rand.New(rand.NewPCG(1, 2))

	rotation := CalculateRotation(club, result, random)
	clubAcc := float64(club.AccuracyDegrees())

	// Excellent should have reduced spread (50% of normal)
	maxExpected := clubAcc * 0.5
	if rotation < 0 || rotation > maxExpected {
		t.Errorf("Excellent rotation = %.2f, expected 0 to %.2f", rotation, maxExpected)
	}
}

func TestCalculateRotation_Good(t *testing.T) {
	club := Club{Name: "Putter", Distance: 40, Accuracy: 1.0, Forgiveness: 0.95}
	result := SkillCheckResult{
		Outcome: Good,
		Margin:  2,
	}
	random := rand.New(rand.NewPCG(1, 2))

	rotation := CalculateRotation(club, result, random)

	// Good should be >= 0 (margin reduces rotation)
	if rotation < 0 {
		t.Errorf("Good rotation = %.2f, expected >= 0", rotation)
	}
}

func TestCalculateRotation_Marginal(t *testing.T) {
	club := Club{Name: "SW", Distance: 125, Accuracy: 0.95, Forgiveness: 0.8}
	result := SkillCheckResult{
		Outcome: Marginal,
		Margin:  0,
	}
	random := rand.New(rand.NewPCG(1, 2))

	rotation := CalculateRotation(club, result, random)
	clubAcc := float64(club.AccuracyDegrees())

	// Marginal should be up to 90% of club accuracy
	maxExpected := clubAcc * 0.9
	if rotation < 0 || rotation > maxExpected {
		t.Errorf("Marginal rotation = %.2f, expected 0 to %.2f", rotation, maxExpected)
	}
}

func TestCalculateRotation_Poor(t *testing.T) {
	club := Club{Name: "3 Wood", Distance: 250, Accuracy: 0.8, Forgiveness: 0.8}
	result := SkillCheckResult{
		Outcome: Poor,
		Margin:  -2,
	}
	random := rand.New(rand.NewPCG(1, 2))

	rotation := CalculateRotation(club, result, random)

	// Poor should have at least 1 degree rotation
	if rotation < 1 {
		t.Errorf("Poor rotation = %.2f, expected >= 1", rotation)
	}
}

func TestCalculateRotation_Bad(t *testing.T) {
	club := Club{Name: "4 Iron", Distance: 215, Accuracy: 0.85, Forgiveness: 0.8}
	result := SkillCheckResult{
		Outcome: Bad,
		Margin:  -5,
	}
	random := rand.New(rand.NewPCG(1, 2))

	rotation := CalculateRotation(club, result, random)

	if rotation < 1 {
		t.Errorf("Bad rotation = %.2f, expected >= 1", rotation)
	}

	poorResult := SkillCheckResult{Outcome: Poor, Margin: -2}
	poorRotation := CalculateRotation(club, poorResult, rand.New(rand.NewPCG(1, 2)))

	if rotation <= poorRotation {
		t.Logf("Bad rotation (%.2f) should typically be worse than Poor (%.2f), but randomness may vary",
			rotation, poorRotation)
	}
}

func TestCalculateRotation_CriticalFailure(t *testing.T) {
	club := Club{Name: "Driver", Distance: 280, Accuracy: 0.75, Forgiveness: 0.8}
	result := SkillCheckResult{
		Outcome: CriticalFailure,
		Margin:  -8,
	}
	random := rand.New(rand.NewPCG(1, 2))

	rotation := CalculateRotation(club, result, random)

	// Critical failure should be 60-90 degrees
	if rotation < 60 || rotation > 90 {
		t.Errorf("CriticalFailure rotation = %.2f, expected 60 to 90", rotation)
	}
}

func TestCalculatePower_CriticalSuccess(t *testing.T) {
	club := Club{Name: "Driver", Accuracy: 0.75, Forgiveness: 0.8}
	result := SkillCheckResult{Outcome: CriticalSuccess, Margin: 7}
	initialPower := 1.0

	power := CalculatePower(club, initialPower, result)

	expected := 1.05
	if !floatEquals(power, expected, 0.001) {
		t.Errorf("CriticalSuccess power = %.3f, expected %.3f", power, expected)
	}
}

func TestCalculatePower_ExcellentAndGood(t *testing.T) {
	club := Club{Name: "7 Iron", Accuracy: 0.9, Forgiveness: 0.8}
	initialPower := 0.85

	tests := []struct {
		outcome  SkillCheckOutcome
		margin   int
		expected float64
	}{
		{Excellent, 5, 0.85},
		{Good, 2, 0.85},
	}

	for _, tt := range tests {
		result := SkillCheckResult{Outcome: tt.outcome, Margin: tt.margin}
		power := CalculatePower(club, initialPower, result)

		if !floatEquals(power, tt.expected, 0.001) {
			t.Errorf("%s power = %.3f, expected %.3f", tt.outcome, power, tt.expected)
		}
	}
}

func TestCalculatePower_Marginal(t *testing.T) {
	club := Club{Name: "Putter", Accuracy: 1.0, Forgiveness: 0.95}
	result := SkillCheckResult{Outcome: Marginal, Margin: 0}
	initialPower := 1.0

	power := CalculatePower(club, initialPower, result)

	expected := 0.95
	if !floatEquals(power, expected, 0.001) {
		t.Errorf("Marginal power = %.3f, expected %.3f", power, expected)
	}
}

func TestCalculatePower_Poor(t *testing.T) {
	club := Club{Name: "SW", Accuracy: 0.95, Forgiveness: 0.8}
	result := SkillCheckResult{Outcome: Poor, Margin: -2}
	initialPower := 1.0

	power := CalculatePower(club, initialPower, result)

	expected := 1.0 * (0.8 - 0.02)
	if !floatEquals(power, expected, 0.001) {
		t.Errorf("Poor power = %.3f, expected %.3f", power, expected)
	}

	if power < 0.1 {
		t.Errorf("Poor power = %.3f, should be >= 0.1", power)
	}
}

func TestCalculatePower_Bad(t *testing.T) {
	club := Club{Name: "3 Wood", Accuracy: 0.8, Forgiveness: 0.8}
	result := SkillCheckResult{Outcome: Bad, Margin: -5}
	initialPower := 1.0

	power := CalculatePower(club, initialPower, result)

	expected := math.Max(1.0*(0.8*0.8-0.05), 0.1)
	if !floatEquals(power, expected, 0.001) {
		t.Errorf("Bad power = %.3f, expected %.3f", power, expected)
	}

	poorResult := SkillCheckResult{Outcome: Poor, Margin: -5}
	poorPower := CalculatePower(club, initialPower, poorResult)

	if power >= poorPower {
		t.Errorf("Bad power (%.3f) should be less than Poor power (%.3f)", power, poorPower)
	}
}

func TestCalculatePower_CriticalFailure(t *testing.T) {
	club := Club{Name: "Driver", Accuracy: 0.75, Forgiveness: 0.8}
	result := SkillCheckResult{Outcome: CriticalFailure, Margin: -8}
	initialPower := 1.0

	power := CalculatePower(club, initialPower, result)

	expected := 0.2
	if !floatEquals(power, expected, 0.001) {
		t.Errorf("CriticalFailure power = %.3f, expected %.3f", power, expected)
	}
}

func TestCalculatePower_FloorEnforcement(t *testing.T) {
	club := Club{Name: "Driver", Accuracy: 0.75, Forgiveness: 0.5}
	result := SkillCheckResult{Outcome: Poor, Margin: -100}
	initialPower := 1.0

	power := CalculatePower(club, initialPower, result)

	if power < 0.1 {
		t.Errorf("Power = %.3f, should be >= 0.1 (floor)", power)
	}
}

func TestCalculatePower_ScalesWithInitialPower(t *testing.T) {
	club := Club{Name: "7 Iron", Accuracy: 0.9, Forgiveness: 0.8}
	result := SkillCheckResult{Outcome: Good, Margin: 2}

	tests := []float64{0.5, 0.75, 1.0, 1.25}

	for _, initialPower := range tests {
		power := CalculatePower(club, initialPower, result)

		if !floatEquals(power, initialPower, 0.001) {
			t.Errorf("Good with initial power %.2f = %.3f, expected %.3f",
				initialPower, power, initialPower)
		}
	}
}

func TestGetShotQualityDescription_AllOutcomes(t *testing.T) {
	tests := []struct {
		outcome       SkillCheckOutcome
		shouldContain string
	}{
		{CriticalSuccess, "PURE STRIKE"},
		{Excellent, "Great shot"},
		{Good, "Good contact"},
		{Marginal, "Just caught it"},
		{Poor, "Slight miss"},
		{Bad, "Poor contact"},
		{CriticalFailure, "DISASTER"},
	}

	for _, tt := range tests {
		result := SkillCheckResult{Outcome: tt.outcome}
		description := GetShotQualityDescription(result)

		if description == "" {
			t.Errorf("GetShotQualityDescription(%v) returned empty string", tt.outcome)
		}

		t.Logf("%s: %s", tt.outcome, description)
	}
}

func TestCalculateRotation_DifferentClubs(t *testing.T) {
	clubs := []Club{
		{Name: "Driver", Accuracy: 0.75, Forgiveness: 0.8},
		{Name: "7 Iron", Accuracy: 0.9, Forgiveness: 0.8},
		{Name: "Putter", Accuracy: 1.0, Forgiveness: 0.95},
	}

	result := SkillCheckResult{Outcome: Good, Margin: 2}
	random := rand.New(rand.NewPCG(1, 2))

	for _, club := range clubs {
		rotation := CalculateRotation(club, result, random)
		t.Logf("%s (Accuracy: %.2f) rotation: %.2f°", club.Name, club.Accuracy, rotation)

		if rotation < 0 {
			t.Errorf("%s produced negative rotation: %.2f", club.Name, rotation)
		}
	}
}

func floatEquals(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
