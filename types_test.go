package gogolf

import (
	"testing"
)

func TestInchesToUnits(t *testing.T) {
	inches := float32(5)
	expectedUnits := float32(1)

	actual := Inch(inches).Units()
	if actual != Unit(expectedUnits) {
		t.Error("Expected 5 inches to equal 1 unit, but got", actual)
	}
}

func TestInchesToFeet(t *testing.T) {
	inches := float32(12)
	expectedFeet := float32(1)

	actual := Inch(inches).Feet()
	if actual != Foot(expectedFeet) {
		t.Error("Expected 12 inches to equal 1 feet, but got", actual)
	}
}

func TestInchesToYards(t *testing.T) {
	inches := float32(12)
	expectedYards := float32(0.33333334)

	actual := Inch(inches).Yards()
	if actual != Yard(expectedYards) {
		t.Error("Expected 12 inch to equal 0.33333334 yards, but got", actual)
	}
}

func TestUnitsToInches(t *testing.T) {
	units := float32(1)
	expectedInches := float32(5)

	actual := Unit(units).Inches()
	if actual != Inch(expectedInches) {
		t.Error("Expected 1 unit to equal 5 inches, but got", actual)
	}
}

func TestUnitsToFeet(t *testing.T) {
	units := float32(1)
	expectedFeet := 0.41666666

	actual := Unit(units).Feet()
	if actual != Foot(expectedFeet) {
		t.Error("Expected 1 unit to equal 0.41666666 feet, but got", actual)
	}
}

func TestUnitsToYard(t *testing.T) {
	units := float32(1)
	expectedYards := 0.13888888

	actual := Unit(units).Yards()
	if actual != Yard(expectedYards) {
		t.Error("Expected 1 unit to equal 0.13888888 yards, but got", actual)
	}
}

func TestYardsToUnits(t *testing.T) {
	yards := float32(10)
	expectedUnits := float32(72)

	actual := Yard(yards).Units()
	if actual != Unit(expectedUnits) {
		t.Error("Expected 10 yards to equal 72 units, but got", actual)
	}
}

func TestLieType_String(t *testing.T) {
	tests := []struct {
		lie      LieType
		expected string
	}{
		{Tee, "Tee"},
		{Fairway, "Fairway"},
		{FirstCut, "First Cut"},
		{Rough, "Rough"},
		{DeepRough, "Deep Rough"},
		{Bunker, "Bunker"},
		{Green, "Green"},
		{PenaltyArea, "Penalty Area"},
	}

	for _, tt := range tests {
		got := tt.lie.String()
		if got != tt.expected {
			t.Errorf("LieType.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestLieType_DifficultyModifier(t *testing.T) {
	tests := []struct {
		lie      LieType
		expected int
	}{
		{Tee, 2},
		{Fairway, 0},
		{FirstCut, -1},
		{Rough, -2},
		{DeepRough, -4},
		{Bunker, -4},
		{Green, 1},
		{PenaltyArea, 0},
	}

	for _, tt := range tests {
		got := tt.lie.DifficultyModifier()
		if got != tt.expected {
			t.Errorf("%s.DifficultyModifier() = %v, want %v",
				tt.lie, got, tt.expected)
		}
	}
}

func TestLie_AffectsTargetNumber(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "Driver"}

	teeTarget := golfer.CalculateTargetNumber(club, Tee.DifficultyModifier())
	fairwayTarget := golfer.CalculateTargetNumber(club, Fairway.DifficultyModifier())
	roughTarget := golfer.CalculateTargetNumber(club, Rough.DifficultyModifier())
	bunkerTarget := golfer.CalculateTargetNumber(club, Bunker.DifficultyModifier())

	if teeTarget <= fairwayTarget {
		t.Errorf("Tee target (%v) should be higher than fairway (%v)", teeTarget, fairwayTarget)
	}

	if fairwayTarget <= roughTarget {
		t.Errorf("Fairway target (%v) should be higher than rough (%v)", fairwayTarget, roughTarget)
	}

	if roughTarget <= bunkerTarget {
		t.Errorf("Rough target (%v) should be higher than bunker (%v)", roughTarget, bunkerTarget)
	}

	expectedTee := 6
	expectedFairway := 4
	expectedRough := 2
	expectedBunker := 0

	if teeTarget != expectedTee {
		t.Errorf("Tee target = %v, want %v", teeTarget, expectedTee)
	}
	if fairwayTarget != expectedFairway {
		t.Errorf("Fairway target = %v, want %v", fairwayTarget, expectedFairway)
	}
	if roughTarget != expectedRough {
		t.Errorf("Rough target = %v, want %v", roughTarget, expectedRough)
	}
	if bunkerTarget != expectedBunker {
		t.Errorf("Bunker target = %v, want %v", bunkerTarget, expectedBunker)
	}
}

func TestSkillCheckOutcome_String(t *testing.T) {
	tests := []struct {
		outcome  SkillCheckOutcome
		expected string
	}{
		{CriticalFailure, "Critical Failure"},
		{Bad, "Bad"},
		{Poor, "Poor"},
		{Marginal, "Marginal"},
		{Good, "Good"},
		{Excellent, "Excellent"},
		{CriticalSuccess, "Critical Success"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.outcome.String(); got != tt.expected {
				t.Errorf("SkillCheckOutcome.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
