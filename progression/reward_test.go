package progression

import "testing"

func TestCalculateHoleReward_Birdie(t *testing.T) {
	par := 4
	strokes := 3

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 25
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (birdie)",
			par, strokes, reward, expectedReward)
	}
}

func TestCalculateHoleReward_Par(t *testing.T) {
	par := 4
	strokes := 4

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 10
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (par)",
			par, strokes, reward, expectedReward)
	}
}

func TestCalculateHoleReward_Bogey(t *testing.T) {
	par := 4
	strokes := 5

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 5
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (bogey)",
			par, strokes, reward, expectedReward)
	}
}

func TestCalculateHoleReward_Eagle(t *testing.T) {
	par := 5
	strokes := 3

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 50
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (eagle)",
			par, strokes, reward, expectedReward)
	}
}

func TestCalculateHoleReward_DoubleBogey(t *testing.T) {
	par := 4
	strokes := 6

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 1
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (double bogey)",
			par, strokes, reward, expectedReward)
	}
}

func TestCalculateHoleReward_TripleBogey(t *testing.T) {
	par := 4
	strokes := 7

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 1
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (triple bogey)",
			par, strokes, reward, expectedReward)
	}
}

func TestCalculateHoleReward_HoleInOne(t *testing.T) {
	par := 3
	strokes := 1

	reward := CalculateHoleReward(par, strokes)

	expectedReward := 100
	if reward != expectedReward {
		t.Errorf("CalculateHoleReward(par=%d, strokes=%d) = %d, want %d (hole in one)",
			par, strokes, reward, expectedReward)
	}
}
