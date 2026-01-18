package main

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

func TestDetermineOutcome_MarginBased(t *testing.T) {
	tests := []struct {
		name       string
		margin     int
		isCritical bool
		expected   SkillCheckOutcome
	}{
		{"Margin +7", 7, false, CriticalSuccess},
		{"Margin +8", 8, false, CriticalSuccess},
		{"Margin +10", 10, false, CriticalSuccess},

		{"Margin +6", 6, false, Excellent},
		{"Margin +5", 5, false, Excellent},
		{"Margin +4", 4, false, Excellent},

		{"Margin +3", 3, false, Good},
		{"Margin +2", 2, false, Good},
		{"Margin +1", 1, false, Good},

		{"Margin 0", 0, false, Marginal},

		{"Margin -1", -1, false, Poor},
		{"Margin -2", -2, false, Poor},
		{"Margin -3", -3, false, Poor},

		{"Margin -4", -4, false, Bad},
		{"Margin -5", -5, false, Bad},
		{"Margin -6", -6, false, Bad},

		{"Margin -7", -7, false, CriticalFailure},
		{"Margin -8", -8, false, CriticalFailure},
		{"Margin -10", -10, false, CriticalFailure},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineOutcome(tt.margin, tt.isCritical)
			if got != tt.expected {
				t.Errorf("DetermineOutcome(%d, %v) = %v, want %v",
					tt.margin, tt.isCritical, got, tt.expected)
			}
		})
	}
}

func TestDetermineOutcome_CriticalOverride(t *testing.T) {
	tests := []struct {
		name     string
		margin   int
		expected SkillCheckOutcome
	}{
		{"Critical with margin +7", 7, CriticalSuccess},
		{"Critical with margin +4", 4, CriticalSuccess},
		{"Critical with margin +1", 1, CriticalSuccess},
		{"Critical with margin 0", 0, CriticalSuccess},

		{"Critical with margin -1", -1, CriticalFailure},
		{"Critical with margin -4", -4, CriticalFailure},
		{"Critical with margin -7", -7, CriticalFailure},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineOutcome(tt.margin, true)
			if got != tt.expected {
				t.Errorf("DetermineOutcome(%d, true) = %v, want %v",
					tt.margin, got, tt.expected)
			}
		})
	}
}

func TestDetermineOutcome_TierBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		margin   int
		expected SkillCheckOutcome
	}{
		{"Just below Critical Success", 6, Excellent},
		{"Just at Critical Success", 7, CriticalSuccess},

		{"Just below Excellent", 3, Good},
		{"Just at Excellent", 4, Excellent},

		{"Just below Good", 0, Marginal},
		{"Just at Good", 1, Good},

		{"Just below Marginal", -1, Poor},
		{"At Marginal", 0, Marginal},

		{"Just below Poor", -4, Bad},
		{"Just at Poor", -3, Poor},

		{"Just below Bad", -7, CriticalFailure},
		{"Just at Bad", -6, Bad},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineOutcome(tt.margin, false)
			if got != tt.expected {
				t.Errorf("DetermineOutcome(%d, false) = %v, want %v",
					tt.margin, got, tt.expected)
			}
		})
	}
}
