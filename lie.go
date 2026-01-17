package main

// LieType represents the condition of the ball's position on the course
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

// String returns the display name for the lie type
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

// DifficultyModifier returns the target number adjustment for this lie
// Positive values make shots easier (higher target number)
// Negative values make shots harder (lower target number)
func (l LieType) DifficultyModifier() int {
	switch l {
	case Tee:
		return 2 // Easiest - ball teed up perfectly
	case Fairway:
		return 0 // Normal - good lie, no penalty
	case FirstCut:
		return -1 // Slight penalty - ball sitting down a bit
	case Rough:
		return -2 // Moderate penalty - grass around ball
	case DeepRough:
		return -4 // Heavy penalty - ball buried in thick grass
	case Bunker:
		return -4 // Very hard - sand shot requires special technique
	case Green:
		return 1 // Putting is easier than full shots
	case PenaltyArea:
		return 0 // Same as fairway after drop
	default:
		return 0
	}
}
