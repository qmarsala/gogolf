package gogolf

type LieType int

const (
	Tee LieType = iota
	Fairway
	FirstCut
	Rough
	DeepRough
	Bunker
	Green
	PenaltyArea
)

func (l LieType) String() string {
	return [...]string{
		"Tee",
		"Fairway",
		"First Cut",
		"Rough",
		"Deep Rough",
		"Bunker",
		"Green",
		"Penalty Area",
	}[l]
}

func (l LieType) DifficultyModifier() int {
	switch l {
	case Tee:
		return 2
	case Fairway:
		return 0
	case FirstCut:
		return -1
	case Rough:
		return -2
	case DeepRough:
		return -4
	case Bunker:
		return -4
	case Green:
		return 1
	case PenaltyArea:
		return 0
	default:
		return 0
	}
}

type SkillCheckOutcome int

const (
	CriticalFailure SkillCheckOutcome = iota
	Bad
	Poor
	Marginal
	Good
	Excellent
	CriticalSuccess
)

func (o SkillCheckOutcome) String() string {
	return [...]string{
		"Critical Failure",
		"Bad",
		"Poor",
		"Marginal",
		"Good",
		"Excellent",
		"Critical Success",
	}[o]
}

type SkillCheckResult struct {
	Success    bool
	IsCritical bool
	RollTotal  int
	Rolls      []int
	Margin     int
	Outcome    SkillCheckOutcome
}

type ShotShape int

const (
	Straight ShotShape = iota
	Draw
	Fade
	Hook
	Slice
)

func (s ShotShape) String() string {
	return [...]string{
		"Straight",
		"Draw",
		"Fade",
		"Hook",
		"Slice",
	}[s]
}

func (s ShotShape) DifficultyModifier() int {
	switch s {
	case Straight:
		return -2
	case Draw:
		return 1
	case Fade:
		return 1
	case Hook:
		return -1
	case Slice:
		return -1
	default:
		return 0
	}
}
