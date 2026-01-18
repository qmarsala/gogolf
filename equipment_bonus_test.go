package gogolf_test

import (
	"gogolf"
	"math/rand/v2"
	"testing"
)

func TestGolfer_GetTotalLiePenaltyReduction_NoShoes(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")

	reduction := golfer.GetTotalLiePenaltyReduction()

	if reduction != 0 {
		t.Errorf("GetTotalLiePenaltyReduction with no shoes = %d, want 0", reduction)
	}
}

func TestGolfer_GetTotalLiePenaltyReduction_WithShoes(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	shoes := &gogolf.Shoes{
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

func TestCalculateTargetNumber_WithShoesBonus(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Skills["Driver"] = gogolf.Skill{Name: "Driver", Level: 5, Experience: 0}
	golfer.Abilities["Strength"] = gogolf.Ability{Name: "Strength", Level: 5, Experience: 0}
	club := gogolf.Club{Name: "Driver"}

	baseTarget := golfer.CalculateTargetNumber(club, -2)

	shoes := &gogolf.Shoes{Name: "Test", LiePenaltyReduction: 1, Cost: 30}
	golfer.EquipShoes(shoes)

	newTarget := golfer.CalculateTargetNumber(club, -2)

	expectedIncrease := 1
	if newTarget-baseTarget != expectedIncrease {
		t.Errorf("Target number increase = %d, want %d", newTarget-baseTarget, expectedIncrease)
	}
}

func TestCalculateTargetNumber_PuttingOnGreenWithTourEditionShoes(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	putter := gogolf.Club{Name: "Putter"}

	shoes := &gogolf.Shoes{Name: "Tour Edition", LiePenaltyReduction: 3, Cost: 80}
	golfer.EquipShoes(shoes)

	greenDifficulty := gogolf.Green.DifficultyModifier()

	// In game, putter does not apply shape modifier
	target := golfer.CalculateTargetNumber(putter, greenDifficulty)

	// skill (1) + ability (1) + green (+1) + shoes (+3) = 6
	expectedTarget := 6
	if target != expectedTarget {
		t.Errorf("Putter on green with Tour Edition shoes = %d, want %d (skill=1, ability=1, green=%d, shoes=+3)",
			target, expectedTarget, greenDifficulty)
	}
}

func TestCalculateRotation_WithGloveBonus(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	club := gogolf.Club{Name: "Driver", Accuracy: 0.75, Forgiveness: 0.8}
	random := rand.New(rand.NewPCG(12345, 67890))

	result := gogolf.SkillCheckResult{
		Outcome: gogolf.Good,
		Margin:  3,
	}

	rotationWithout := gogolf.CalculateRotation(club, result, random)

	random = rand.New(rand.NewPCG(12345, 67890))

	glove := &gogolf.Glove{Name: "Test", AccuracyBonus: 0.05, Cost: 30}
	golfer.EquipGlove(glove)

	modifiedClub := golfer.GetModifiedClub(club)
	rotationWith := gogolf.CalculateRotation(modifiedClub, result, random)

	if rotationWith >= rotationWithout {
		t.Errorf("Rotation with glove (%f) should be less than without (%f)", rotationWith, rotationWithout)
	}
}

func TestGolfer_GetModifiedClub(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	baseClub := gogolf.Club{
		Name:        "Driver",
		Distance:    280,
		Accuracy:    0.75,
		Forgiveness: 0.8,
	}

	glove := &gogolf.Glove{Name: "Test", AccuracyBonus: 0.05, Cost: 30}
	golfer.EquipGlove(glove)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	expectedAccuracy := baseClub.Accuracy + 0.05
	if modifiedClub.Accuracy != expectedAccuracy {
		t.Errorf("Modified club accuracy = %f, want %f", modifiedClub.Accuracy, expectedAccuracy)
	}

	if modifiedClub.Distance != baseClub.Distance {
		t.Errorf("Modified club distance = %v, want %v (unchanged)", modifiedClub.Distance, baseClub.Distance)
	}
}

func TestGolfer_GetModifiedClub_WithBall(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	baseClub := gogolf.Club{Name: "Driver", Distance: 280, Accuracy: 0.75, Forgiveness: 0.8}

	ball := &gogolf.Ball{Name: "Test", DistanceBonus: 10, SpinControl: 0.7, Cost: 50}
	golfer.EquipBall(ball)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	expectedDistance := gogolf.Yard(float32(baseClub.Distance) + ball.DistanceBonus)
	if modifiedClub.Distance != expectedDistance {
		t.Errorf("Modified club distance = %v, want %v", modifiedClub.Distance, expectedDistance)
	}
}

func TestGolfer_GetModifiedClub_StacksBonuses(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	baseClub := gogolf.Club{Name: "7 Iron", Distance: 180, Accuracy: 0.9, Forgiveness: 0.8}

	ball := &gogolf.Ball{Name: "Premium", DistanceBonus: 5, SpinControl: 0.8, Cost: 50}
	glove := &gogolf.Glove{Name: "Pro", AccuracyBonus: 0.03, Cost: 40}

	golfer.EquipBall(ball)
	golfer.EquipGlove(glove)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	expectedDistance := gogolf.Yard(float32(baseClub.Distance) + ball.DistanceBonus)
	expectedAccuracy := baseClub.Accuracy + glove.AccuracyBonus

	if modifiedClub.Distance != expectedDistance {
		t.Errorf("Stacked distance = %v, want %v", modifiedClub.Distance, expectedDistance)
	}
	if modifiedClub.Accuracy != expectedAccuracy {
		t.Errorf("Stacked accuracy = %f, want %f", modifiedClub.Accuracy, expectedAccuracy)
	}
}

func TestGolfer_GetModifiedClub_AccuracyCap(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	baseClub := gogolf.Club{Name: "Putter", Distance: 40, Accuracy: 1.0, Forgiveness: 0.95}

	glove := &gogolf.Glove{Name: "Test", AccuracyBonus: 0.05, Cost: 30}
	golfer.EquipGlove(glove)

	modifiedClub := golfer.GetModifiedClub(baseClub)

	if modifiedClub.Accuracy > 1.0 {
		t.Errorf("Modified accuracy = %f, want <= 1.0 (capped)", modifiedClub.Accuracy)
	}
}
