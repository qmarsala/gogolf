package main

// Ability represents a character attribute (Strength, Control, Touch, Mental)
// Abilities contribute to the target number for skill checks
type Ability struct {
	Name       string
	Level      int // 1-9 (max level 9)
	Experience int // XP toward next level
}

// NewAbility creates a new Ability at level 1 with 0 experience
func NewAbility(name string) Ability {
	return Ability{
		Name:       name,
		Level:      1,
		Experience: 0,
	}
}

// Value returns the contribution of this ability to target numbers
// Formula: level * 2
func (a Ability) Value() int {
	return a.Level * 2
}

// AddExperience adds XP and handles level-ups
func (a *Ability) AddExperience(xp int) {
	// Don't accumulate XP at max level
	if a.Level >= 9 {
		return
	}

	a.Experience += xp

	// Check for level-ups (may level up multiple times)
	for a.CanLevelUp() {
		threshold := a.ExperienceToNextLevel()
		a.Experience -= threshold
		a.Level++

		// Cap at level 9
		if a.Level >= 9 {
			a.Level = 9
			a.Experience = 0
			break
		}
	}
}

// ExperienceToNextLevel returns XP needed for next level
// Formula: (level + 1) * 50
// Level 1→2: 100, Level 2→3: 150, Level 3→4: 200, ..., Level 8→9: 400
func (a Ability) ExperienceToNextLevel() int {
	if a.Level >= 9 {
		return 0 // Max level, no next level
	}
	return (a.Level + 1) * 50
}

// CanLevelUp returns true if the ability has enough XP to level up
func (a Ability) CanLevelUp() bool {
	if a.Level >= 9 {
		return false
	}
	return a.Experience >= a.ExperienceToNextLevel()
}
