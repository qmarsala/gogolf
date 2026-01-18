package main

import (
	"math/rand/v2"
	"testing"
)

// Test Golfer.GetTotalLiePenaltyReduction returns 0 with no shoes
func TestGolfer_GetTotalLiePenaltyReduction_NoShoes(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	reduction := golfer.GetTotalLiePenaltyReduction()

	if reduction != 0 {
		t.Errorf("GetTotalLiePenaltyReduction with no shoes = %d, want 0", reduction)
	}
}

// Test Golfer.GetTotalLiePenaltyReduction returns shoes bonus
func TestGolfer_GetTotalLiePenaltyReduction_WithShoes(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	shoes := &Shoes{
		Name:                "Test Shoes",
		LiePenaltyReduction: 2,
		Cost:                50,
	}
	golfer.EquipShoes(shoes)

	reduction := golfer.GetTotalLiePenaltyReduction()

	if reduction != 2 {
		t.Errorf("GetTotalLiePenaltyReduction with shoes = %d, want 2", reduction)
	}
}

// Test CalculateTargetNumber includes lie penalty reduction from shoes
func TestCalculateTargetNumber_WithShoesBonus(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "Driver"}

	// Calculate base target number (from rough, difficulty -2)
	baseTarget := golfer.CalculateTargetNumber(club, -2)

	// Equip shoes that reduce penalties by 1
	shoes := &Shoes{Name: "Test", LiePenaltyReduction: 1, Cost: 30}
	golfer.EquipShoes(shoes)

	// Calculate new target number (should be 1 higher)
	newTarget := golfer.CalculateTargetNumber(club, -2)

	expectedIncrease := 1
	if newTarget-baseTarget != expectedIncrease {
		t.Errorf("Target number increase = %d, want %d", newTarget-baseTarget, expectedIncrease)
	}
}

// Test CalculateRotation is affected by glove accuracy bonus
func TestCalculateRotation_WithGloveBonus(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "Driver", Accuracy: 0.75, Forgiveness: 0.8}
	random := rand.New(rand.NewPCG(12345, 67890)) // Fixed seed for reproducibility

	result := SkillCheckResult{
		Outcome: Good,
		Margin:  3,
	}

	// Calculate rotation without glove
	rotationWithout := CalculateRotation(club, result, random)

	// Reset random for fair comparison
	random = rand.New(rand.NewPCG(12345, 67890))

	// Equip glove
	glove := &Glove{Name: "Test", AccuracyBonus: 0.05, Cost: 30}
	golfer.EquipGlove(glove)

	// Calculate rotation with glove (glove improves club accuracy)
	modifiedClub := golfer.GetModifiedClub(club)
	rotationWith := CalculateRotation(modifiedClub, result, random)

	// With better accuracy, rotation should be less
	if rotationWith >= rotationWithout {
		t.Errorf("Rotation with glove (%f) should be less than without (%f)", rotationWith, rotationWithout)
	}
}

// Test Golfer.GetModifiedClub returns modified club with equipment bonuses
func TestGolfer_GetModifiedClub(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	baseClub := Club{
		Name:        "Driver",
		Distance:    280,
		Accuracy:    0.75,
		Forgiveness: 0.8,
	}

	// Equip glove that adds accuracy
	glove := &Glove{Name: "Test", AccuracyBonus: 0.05, Cost: 30}
	golfer.EquipGlove(glove)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	expectedAccuracy := baseClub.Accuracy + 0.05
	if modifiedClub.Accuracy != expectedAccuracy {
		t.Errorf("Modified club accuracy = %f, want %f", modifiedClub.Accuracy, expectedAccuracy)
	}

	// Distance should be unchanged (no ball equipped)
	if modifiedClub.Distance != baseClub.Distance {
		t.Errorf("Modified club distance = %v, want %v (unchanged)", modifiedClub.Distance, baseClub.Distance)
	}
}

// Test Golfer.GetModifiedClub adds ball distance bonus
func TestGolfer_GetModifiedClub_WithBall(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	baseClub := Club{Name: "Driver", Distance: 280, Accuracy: 0.75, Forgiveness: 0.8}

	// Equip ball with distance bonus
	ball := &Ball{Name: "Test", DistanceBonus: 10, SpinControl: 0.7, Cost: 50}
	golfer.EquipBall(ball)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	expectedDistance := Yard(float32(baseClub.Distance) + ball.DistanceBonus)
	if modifiedClub.Distance != expectedDistance {
		t.Errorf("Modified club distance = %v, want %v", modifiedClub.Distance, expectedDistance)
	}
}

// Test Golfer.GetModifiedClub stacks multiple equipment bonuses
func TestGolfer_GetModifiedClub_StacksBonuses(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	baseClub := Club{Name: "7 Iron", Distance: 180, Accuracy: 0.9, Forgiveness: 0.8}

	// Equip both ball and glove
	ball := &Ball{Name: "Premium", DistanceBonus: 5, SpinControl: 0.8, Cost: 50}
	glove := &Glove{Name: "Pro", AccuracyBonus: 0.03, Cost: 40}

	golfer.EquipBall(ball)
	golfer.EquipGlove(glove)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	// Check both bonuses applied
	expectedDistance := Yard(float32(baseClub.Distance) + ball.DistanceBonus)
	expectedAccuracy := baseClub.Accuracy + glove.AccuracyBonus

	if modifiedClub.Distance != expectedDistance {
		t.Errorf("Stacked distance = %v, want %v", modifiedClub.Distance, expectedDistance)
	}
	if modifiedClub.Accuracy != expectedAccuracy {
		t.Errorf("Stacked accuracy = %f, want %f", modifiedClub.Accuracy, expectedAccuracy)
	}
}

// Test equipment bonuses don't break accuracy limits
func TestGolfer_GetModifiedClub_AccuracyCap(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	baseClub := Club{Name: "Putter", Distance: 40, Accuracy: 1.0, Forgiveness: 0.95}

	// Equip glove (shouldn't push accuracy above 1.0)
	glove := &Glove{Name: "Test", AccuracyBonus: 0.05, Cost: 30}
	golfer.EquipGlove(glove)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	// Accuracy should be capped at 1.0
	if modifiedClub.Accuracy > 1.0 {
		t.Errorf("Modified accuracy = %f, want <= 1.0 (capped)", modifiedClub.Accuracy)
	}
}
