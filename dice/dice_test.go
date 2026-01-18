package dice

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

func TestDice_Roll(t *testing.T) {
	d := NewD6()
	for i := 0; i < 100; i++ {
		roll := d.Roll()
		if roll < 1 || roll > 6 {
			t.Errorf("Roll() = %d, want 1-6", roll)
		}
	}
}

func TestDice_RollN(t *testing.T) {
	d := NewD6()
	total, rolls := d.RollN(3)

	if len(rolls) != 3 {
		t.Errorf("RollN(3) returned %d rolls, want 3", len(rolls))
	}

	sum := 0
	for _, r := range rolls {
		if r < 1 || r > 6 {
			t.Errorf("Individual roll = %d, want 1-6", r)
		}
		sum += r
	}

	if total != sum {
		t.Errorf("Total = %d, sum of rolls = %d", total, sum)
	}
}
