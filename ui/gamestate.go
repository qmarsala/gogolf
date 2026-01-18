package ui

// GameState represents all data needed for rendering the terminal UI
type GameState struct {
	// Player info
	PlayerName     string
	Money          int
	Skills         map[string]SkillDisplay
	Abilities      map[string]AbilityDisplay
	Equipment      EquipmentDisplay

	// Hole info
	HoleNumber int
	TotalHoles int
	Par        int
	HoleDistance float64

	// Ball state
	BallLie           string
	BallLieDifficulty int
	DistanceToHole    float64
	BallLocationX     float64
	BallLocationY     float64
	HoleLocationX     float64
	HoleLocationY     float64

	// Last shot (nil if no shot yet)
	LastShot *ShotDisplay

	// Score
	TotalStrokes     int
	ScoreToPar       int
	StrokesThisHole  int

	// Messages
	StatusMsg string
	ErrorMsg  string
	PromptMsg string
}

// SkillDisplay represents skill information for display
type SkillDisplay struct {
	Name       string
	Level      int
	Value      int
	CurrentXP  int
	XPForNext  int
}

// AbilityDisplay represents ability information for display
type AbilityDisplay struct {
	Name       string
	Level      int
	Value      int
	CurrentXP  int
	XPForNext  int
}

// EquipmentDisplay represents equipped items for display
type EquipmentDisplay struct {
	BallName   string
	BallBonus  string
	GloveName  string
	GloveBonus string
	ShoesName  string
	ShoesBonus string
}

// ShotDisplay represents the result of a shot for display
type ShotDisplay struct {
	ClubName      string
	IntendedShape string
	ActualShape   string
	ShapeSuccess  bool
	Outcome       string // e.g., "Good", "Critical Success"
	Margin        int
	TargetNumber  int
	DiceRolls     []int
	Description   string  // Quality description
	Rotation      float64 // Degrees
	RotationDir   string  // "left" or "right"
	Power         float64 // Percentage 0-1
	Distance      float64 // Yards traveled
	XPEarned      int
	LevelUps      []string // e.g., ["Driver: Level 2!", "Strength: Level 3!"]
}
