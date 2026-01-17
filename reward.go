package main

// CalculateHoleReward determines the money reward based on score relative to par
// The reward system encourages good play:
// - Hole-in-one (1 stroke): 100 money
// - Eagle (2 under par): 50 money
// - Birdie (1 under par): 25 money
// - Par: 10 money
// - Bogey (1 over par): 5 money
// - Double bogey or worse: 1 money
func CalculateHoleReward(par, strokes int) int {
	// Special case: hole-in-one is always amazing
	if strokes == 1 {
		return 100
	}

	scoreToPar := strokes - par

	switch {
	case scoreToPar <= -2:
		// Eagle or better
		return 50
	case scoreToPar == -1:
		// Birdie
		return 25
	case scoreToPar == 0:
		// Par
		return 10
	case scoreToPar == 1:
		// Bogey
		return 5
	default:
		// Double bogey or worse
		return 1
	}
}

// AwardHoleReward calculates and awards money based on hole performance
func (g *Golfer) AwardHoleReward(par, strokes int) {
	reward := CalculateHoleReward(par, strokes)
	g.AddMoney(reward)
}
