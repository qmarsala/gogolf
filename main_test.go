package gogolf

import (
	"gogolf/progression"
	"testing"
)

func TestCalculateXP(t *testing.T) {
	tests := []struct {
		outcome  SkillCheckOutcome
		expected int
	}{
		{CriticalSuccess, 15},
		{Excellent, 10},
		{Good, 7},
		{Marginal, 5},
		{Poor, 3},
		{Bad, 2},
		{CriticalFailure, 1},
	}

	for _, tt := range tests {
		got := calculateXP(tt.outcome)

		if got != tt.expected {
			t.Errorf("calculateXP(%v) = %v, want %v",
				tt.outcome, got, tt.expected)
		}
	}
}

// Test that the game loop integration works correctly
// This tests that a golfer with skills/abilities gets different target numbers
func TestGameLoopIntegration_DynamicTargetNumbers(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	// Test with Driver (should use Driver skill + Strength ability)
	driverClub := Club{Name: "Driver"}
	difficulty := 0

	// Initial level 1: skill.Value() = 2, ability.Value() = 2, total = 4
	targetNumber := golfer.CalculateTargetNumber(driverClub, difficulty)
	expectedInitial := 4 // (level 1 * 2) + (level 1 * 2) + 0

	if targetNumber != expectedInitial {
		t.Errorf("Initial target number = %v, want %v", targetNumber, expectedInitial)
	}

	golfer.Skills["Driver"] = progression.Skill{Name: "Driver", Level: 3, Experience: 0}
	golfer.Abilities["Strength"] = progression.Ability{Name: "Strength", Level: 4, Experience: 0}

	// New target: skill.Value() = 6, ability.Value() = 8, total = 14
	targetNumber = golfer.CalculateTargetNumber(driverClub, difficulty)
	expectedLeveled := 14 // (level 3 * 2) + (level 4 * 2) + 0

	if targetNumber != expectedLeveled {
		t.Errorf("After level-up target number = %v, want %v", targetNumber, expectedLeveled)
	}
}

// Test that XP is awarded correctly and triggers level-ups
func TestGameLoopIntegration_XPAward(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "Putter"}

	initialSkillXP := golfer.Skills["Putter"].Experience
	initialAbilityXP := golfer.Abilities["Mental"].Experience

	xp := calculateXP(Excellent)
	golfer.AwardExperience(club, xp)

	if golfer.Skills["Putter"].Experience != initialSkillXP+10 {
		t.Errorf("Skill XP = %v, want %v", golfer.Skills["Putter"].Experience, initialSkillXP+10)
	}

	if golfer.Abilities["Mental"].Experience != initialAbilityXP+10 {
		t.Errorf("Ability XP = %v, want %v", golfer.Abilities["Mental"].Experience, initialAbilityXP+10)
	}
}

// Test level-up detection works correctly
func TestGameLoopIntegration_LevelUpDetection(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "PW"} // Use a specific club name that maps to Wedges

	// Get the skill and ability names for this club
	skill := golfer.GetSkillForClub(club)
	ability := golfer.GetAbilityForClub(club)

	prevSkillLevel := golfer.Skills[skill.Name].Level
	prevAbilityLevel := golfer.Abilities[ability.Name].Level

	t.Logf("Using club %v -> Skill: %v, Ability: %v", club.Name, skill.Name, ability.Name)
	t.Logf("Before: Skill level=%v XP=%v, Ability level=%v XP=%v",
		golfer.Skills[skill.Name].Level, golfer.Skills[skill.Name].Experience,
		golfer.Abilities[ability.Name].Level, golfer.Abilities[ability.Name].Experience)

	// Award enough XP to level up (100 XP for level 1 → 2)
	golfer.AwardExperience(club, 100)

	t.Logf("After: Skill level=%v XP=%v, Ability level=%v XP=%v",
		golfer.Skills[skill.Name].Level, golfer.Skills[skill.Name].Experience,
		golfer.Abilities[ability.Name].Level, golfer.Abilities[ability.Name].Experience)

	// After AwardExperience, the maps should be updated
	newSkillLevel := golfer.Skills[skill.Name].Level
	newAbilityLevel := golfer.Abilities[ability.Name].Level

	if newSkillLevel <= prevSkillLevel {
		t.Errorf("Skill did not level up: prev=%v, new=%v", prevSkillLevel, newSkillLevel)
	}

	if newAbilityLevel <= prevAbilityLevel {
		t.Errorf("Ability did not level up: prev=%v, new=%v", prevAbilityLevel, newAbilityLevel)
	}

	if newSkillLevel != 2 {
		t.Errorf("Skill level = %v, want 2", newSkillLevel)
	}

	if newAbilityLevel != 2 {
		t.Errorf("Ability level = %v, want 2", newAbilityLevel)
	}
}
