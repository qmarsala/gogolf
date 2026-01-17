package main

// Ball represents a golf ball with performance characteristics
type Ball struct {
	Name          string
	DistanceBonus float32 // Extra yards added to shots
	SpinControl   float32 // How well the ball holds spin (0-1, higher is better)
	Cost          int     // Price in shop
}

// Glove represents golf gloves that improve accuracy
type Glove struct {
	Name          string
	AccuracyBonus float32 // Reduces shot dispersion (0-1, higher is better)
	Cost          int     // Price in shop
}

// Shoes represents golf shoes that help with stability on different lies
type Shoes struct {
	Name                string
	LiePenaltyReduction int // Reduces lie difficulty penalties (positive number)
	Cost                int // Price in shop
}
