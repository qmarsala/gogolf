package main

import "testing"

// Test LieType.String() returns correct names
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

// Test LieType.DifficultyModifier() returns correct values
func TestLieType_DifficultyModifier(t *testing.T) {
	tests := []struct {
		lie      LieType
		expected int
	}{
		{Tee, 2},         // Easiest - teed up perfectly
		{Fairway, 0},     // Normal - good lie
		{FirstCut, -1},   // Slight penalty - ball sitting down a bit
		{Rough, -2},      // Moderate penalty - grass around ball
		{DeepRough, -4},  // Heavy penalty - ball buried
		{Bunker, -4},     // Very hard - sand shot
		{Green, 1},       // Putting is easier than full shots
		{PenaltyArea, 0}, // Same as fairway after drop
	}

	for _, tt := range tests {
		got := tt.lie.DifficultyModifier()
		if got != tt.expected {
			t.Errorf("%s.DifficultyModifier() = %v, want %v",
				tt.lie, got, tt.expected)
		}
	}
}

// Test that different lies produce different target numbers
func TestLie_AffectsTargetNumber(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "Driver"}

	// Test with same golfer/club but different lies
	teeTarget := golfer.CalculateTargetNumber(club, Tee.DifficultyModifier())
	fairwayTarget := golfer.CalculateTargetNumber(club, Fairway.DifficultyModifier())
	roughTarget := golfer.CalculateTargetNumber(club, Rough.DifficultyModifier())
	bunkerTarget := golfer.CalculateTargetNumber(club, Bunker.DifficultyModifier())

	// Tee should be easiest (highest target number)
	if teeTarget <= fairwayTarget {
		t.Errorf("Tee target (%v) should be higher than fairway (%v)", teeTarget, fairwayTarget)
	}

	// Fairway should be easier than rough
	if fairwayTarget <= roughTarget {
		t.Errorf("Fairway target (%v) should be higher than rough (%v)", fairwayTarget, roughTarget)
	}

	// Rough should be easier than bunker
	if roughTarget <= bunkerTarget {
		t.Errorf("Rough target (%v) should be higher than bunker (%v)", roughTarget, bunkerTarget)
	}

	// Verify specific expected values for Level 1 golfer
	// Level 1: Driver skill = 2, Strength ability = 2
	// Base target = 4
	expectedTee := 6      // 4 + 2
	expectedFairway := 4  // 4 + 0
	expectedRough := 2    // 4 + (-2)
	expectedBunker := 0   // 4 + (-4)

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
