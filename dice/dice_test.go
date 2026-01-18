package dice

import "testing"

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
