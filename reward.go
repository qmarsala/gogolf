package gogolf

func CalculateHoleReward(par, strokes int) int {
	if strokes == 1 {
		return 100
	}

	scoreToPar := strokes - par

	switch {
	case scoreToPar <= -2:
		return 50
	case scoreToPar == -1:
		return 25
	case scoreToPar == 0:
		return 10
	case scoreToPar == 1:
		return 5
	default:
		return 1
	}
}
