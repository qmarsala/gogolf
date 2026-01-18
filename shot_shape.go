package gogolf

type ShapeResult struct {
	Intended ShotShape
	Actual   ShotShape
	Success  bool
}

func DetermineActualShape(intended ShotShape, result SkillCheckResult, random RandomSource) ShapeResult {
	if result.Outcome >= Marginal {
		return ShapeResult{Intended: intended, Actual: intended, Success: true}
	}

	actual := determineFailedShape(intended, random)
	return ShapeResult{Intended: intended, Actual: actual, Success: false}
}

func determineFailedShape(intended ShotShape, random RandomSource) ShotShape {
	switch intended {
	case Straight:
		if random.IntN(2) == 0 {
			return Draw
		}
		return Fade

	case Draw:
		if random.IntN(10) < 8 {
			return Hook
		}
		return Straight

	case Fade:
		if random.IntN(10) < 8 {
			return Slice
		}
		return Straight

	case Hook:
		if random.IntN(10) < 7 {
			return Hook
		}
		return Draw

	case Slice:
		if random.IntN(10) < 7 {
			return Slice
		}
		return Fade
	}

	return intended
}
