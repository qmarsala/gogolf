package main

// Skill represents club proficiency (Driver, Woods, Irons, Wedges, Putter)
// Skills contribute to the target number for skill checks
type Skill struct {
	Name       string
	Level      int // 1-9 (max level 9)
	Experience int // XP toward next level
}

// NewSkill creates a new Skill at level 1 with 0 experience
func NewSkill(name string) Skill {
	return Skill{
		Name:       name,
		Level:      1,
		Experience: 0,
	}
}

// Value returns the contribution of this skill to target numbers
// Formula: level * 2
func (s Skill) Value() int {
	return s.Level * 2
}

// AddExperience adds XP and handles level-ups
func (s *Skill) AddExperience(xp int) {
	// Don't accumulate XP at max level
	if s.Level >= 9 {
		return
	}

	s.Experience += xp

	// Check for level-ups (may level up multiple times)
	for s.CanLevelUp() {
		threshold := s.ExperienceToNextLevel()
		s.Experience -= threshold
		s.Level++

		// Cap at level 9
		if s.Level >= 9 {
			s.Level = 9
			s.Experience = 0
			break
		}
	}
}

// ExperienceToNextLevel returns XP needed for next level
// Formula: (level + 1) * 50
// Level 1→2: 100, Level 2→3: 150, Level 3→4: 200, ..., Level 8→9: 400
func (s Skill) ExperienceToNextLevel() int {
	if s.Level >= 9 {
		return 0 // Max level, no next level
	}
	return (s.Level + 1) * 50
}

// CanLevelUp returns true if the skill has enough XP to level up
func (s Skill) CanLevelUp() bool {
	if s.Level >= 9 {
		return false
	}
	return s.Experience >= s.ExperienceToNextLevel()
}
