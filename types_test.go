package gogolf

import (
	"testing"
)

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

	expectedTee := 4
	expectedFairway := 2
	expectedRough := 0
	expectedBunker := -2

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
