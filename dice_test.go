package main

import "testing"

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
		// Critical Success tier (margin >= 7)
		{"Margin +7", 7, false, CriticalSuccess},
		{"Margin +8", 8, false, CriticalSuccess},
		{"Margin +10", 10, false, CriticalSuccess},

		// Excellent tier (margin 4-6)
		{"Margin +6", 6, false, Excellent},
		{"Margin +5", 5, false, Excellent},
		{"Margin +4", 4, false, Excellent},

		// Good tier (margin 1-3)
		{"Margin +3", 3, false, Good},
		{"Margin +2", 2, false, Good},
		{"Margin +1", 1, false, Good},

		// Marginal tier (margin 0)
		{"Margin 0", 0, false, Marginal},

		// Poor tier (margin -1 to -3)
		{"Margin -1", -1, false, Poor},
		{"Margin -2", -2, false, Poor},
		{"Margin -3", -3, false, Poor},

		// Bad tier (margin -4 to -6)
		{"Margin -4", -4, false, Bad},
		{"Margin -5", -5, false, Bad},
		{"Margin -6", -6, false, Bad},

		// Critical Failure tier (margin <= -7)
		{"Margin -7", -7, false, CriticalFailure},
		{"Margin -8", -8, false, CriticalFailure},
		{"Margin -10", -10, false, CriticalFailure},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineOutcome(tt.margin, tt.isCritical)
			if got != tt.expected {
				t.Errorf("determineOutcome(%d, %v) = %v, want %v",
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
		// Critical rolls with positive margin should be Critical Success
		{"Critical with margin +7", 7, CriticalSuccess},
		{"Critical with margin +4", 4, CriticalSuccess},
		{"Critical with margin +1", 1, CriticalSuccess},
		{"Critical with margin 0", 0, CriticalSuccess},

		// Critical rolls with negative margin should be Critical Failure
		{"Critical with margin -1", -1, CriticalFailure},
		{"Critical with margin -4", -4, CriticalFailure},
		{"Critical with margin -7", -7, CriticalFailure},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineOutcome(tt.margin, true)
			if got != tt.expected {
				t.Errorf("determineOutcome(%d, true) = %v, want %v",
					tt.margin, got, tt.expected)
			}
		})
	}
}

func TestDetermineOutcome_TierBoundaries(t *testing.T) {
	// Test boundary transitions between tiers
	tests := []struct {
		name     string
		margin   int
		expected SkillCheckOutcome
	}{
		// Excellent/Critical Success boundary
		{"Just below Critical Success", 6, Excellent},
		{"Just at Critical Success", 7, CriticalSuccess},

		// Good/Excellent boundary
		{"Just below Excellent", 3, Good},
		{"Just at Excellent", 4, Excellent},

		// Marginal/Good boundary
		{"Just below Good", 0, Marginal},
		{"Just at Good", 1, Good},

		// Poor/Marginal boundary
		{"Just below Marginal", -1, Poor},
		{"At Marginal", 0, Marginal},

		// Bad/Poor boundary
		{"Just below Poor", -4, Bad},
		{"Just at Poor", -3, Poor},

		// Critical Failure/Bad boundary
		{"Just below Bad", -7, CriticalFailure},
		{"Just at Bad", -6, Bad},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineOutcome(tt.margin, false)
			if got != tt.expected {
				t.Errorf("determineOutcome(%d, false) = %v, want %v",
					tt.margin, got, tt.expected)
			}
		})
	}
}

func TestSkillCheck_OutcomePopulated(t *testing.T) {
	golfer := Golfer{Clubs: DefaultClubs()}
	dice := NewD6()

	// Run multiple skill checks to verify Outcome is always populated
	for i := 0; i < 10; i++ {
		result := golfer.SkillCheck(dice, 10)

		// Verify outcome is set
		if result.Outcome < CriticalFailure || result.Outcome > CriticalSuccess {
			t.Errorf("SkillCheck returned invalid outcome: %v", result.Outcome)
		}

		// Verify outcome matches success/failure
		if result.Success && result.Outcome < Marginal {
			t.Errorf("Success=true but outcome=%v is a failure tier", result.Outcome)
		}
		if !result.Success && result.Outcome >= Marginal {
			t.Errorf("Success=false but outcome=%v is a success tier", result.Outcome)
		}
	}
}

func TestSkillCheck_CriticalDetection(t *testing.T) {
	// This test verifies that critical detection works correctly
	// We can't easily force specific rolls, but we can verify the logic is sound
	golfer := Golfer{Clubs: DefaultClubs()}
	dice := NewD6()

	for i := 0; i < 100; i++ {
		result := golfer.SkillCheck(dice, 10)

		if result.IsCritical {
			// All three rolls should be equal
			if !(result.Rolls[0] == result.Rolls[1] && result.Rolls[0] == result.Rolls[2]) {
				t.Errorf("IsCritical=true but rolls not equal: %v", result.Rolls)
			}

			// Outcome should be critical
			if result.Outcome != CriticalSuccess && result.Outcome != CriticalFailure {
				t.Errorf("IsCritical=true but outcome=%v is not critical", result.Outcome)
			}

			// If margin >= 0, should be Critical Success
			if result.Margin >= 0 && result.Outcome != CriticalSuccess {
				t.Errorf("Critical with margin=%d should be CriticalSuccess, got %v",
					result.Margin, result.Outcome)
			}

			// If margin < 0, should be Critical Failure
			if result.Margin < 0 && result.Outcome != CriticalFailure {
				t.Errorf("Critical with margin=%d should be CriticalFailure, got %v",
					result.Margin, result.Outcome)
			}
		}
	}
}

func TestSkillCheck_MarginCalculation(t *testing.T) {
	golfer := Golfer{Clubs: DefaultClubs()}
	dice := NewD6()
	targetNumber := 10

	for i := 0; i < 50; i++ {
		result := golfer.SkillCheck(dice, targetNumber)

		// Verify margin calculation
		expectedMargin := targetNumber - result.RollTotal
		if result.Margin != expectedMargin {
			t.Errorf("Margin calculation error: got %d, expected %d (target=%d, roll=%d)",
				result.Margin, expectedMargin, targetNumber, result.RollTotal)
		}

		// Verify success/failure matches margin
		if result.Success != (result.Margin >= 0) {
			t.Errorf("Success/Margin mismatch: Success=%v, Margin=%d",
				result.Success, result.Margin)
		}
	}
}
