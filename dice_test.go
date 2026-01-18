package gogolf

import (
	"testing"
)

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

func TestSkillCheck_OutcomePopulated(t *testing.T) {
	golfer := Golfer{Clubs: DefaultClubs()}
	d := NewD6()

	for i := 0; i < 10; i++ {
		result := golfer.SkillCheck(d, 10)

		if result.Outcome < CriticalFailure || result.Outcome > CriticalSuccess {
			t.Errorf("SkillCheck returned invalid outcome: %v", result.Outcome)
		}

		if result.Success && result.Outcome < Marginal {
			t.Errorf("Success=true but outcome=%v is a failure tier", result.Outcome)
		}
		if !result.Success && result.Outcome >= Marginal {
			t.Errorf("Success=false but outcome=%v is a success tier", result.Outcome)
		}
	}
}

func TestSkillCheck_CriticalDetection(t *testing.T) {
	golfer := Golfer{Clubs: DefaultClubs()}
	d := NewD6()

	for i := 0; i < 100; i++ {
		result := golfer.SkillCheck(d, 10)

		if result.IsCritical {
			if !(result.Rolls[0] == result.Rolls[1] && result.Rolls[0] == result.Rolls[2]) {
				t.Errorf("IsCritical=true but rolls not equal: %v", result.Rolls)
			}

			if result.Outcome != CriticalSuccess && result.Outcome != CriticalFailure {
				t.Errorf("IsCritical=true but outcome=%v is not critical", result.Outcome)
			}

			if result.Margin >= 0 && result.Outcome != CriticalSuccess {
				t.Errorf("Critical with margin=%d should be CriticalSuccess, got %v",
					result.Margin, result.Outcome)
			}

			if result.Margin < 0 && result.Outcome != CriticalFailure {
				t.Errorf("Critical with margin=%d should be CriticalFailure, got %v",
					result.Margin, result.Outcome)
			}
		}
	}
}

func TestSkillCheck_MarginCalculation(t *testing.T) {
	golfer := Golfer{Clubs: DefaultClubs()}
	d := NewD6()
	targetNumber := 10

	for i := 0; i < 50; i++ {
		result := golfer.SkillCheck(d, targetNumber)

		expectedMargin := targetNumber - result.RollTotal
		if result.Margin != expectedMargin {
			t.Errorf("Margin calculation error: got %d, expected %d (target=%d, roll=%d)",
				result.Margin, expectedMargin, targetNumber, result.RollTotal)
		}

		if result.Success != (result.Margin >= 0) {
			t.Errorf("Success/Margin mismatch: Success=%v, Margin=%d",
				result.Success, result.Margin)
		}
	}
}
