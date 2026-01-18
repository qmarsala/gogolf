package gogolf

import (
	"gogolf/progression"
	"testing"
)

func TestGameLoopIntegration_DynamicTargetNumbers(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	driverClub := Club{Name: "Driver"}
	difficulty := 0

	targetNumber := golfer.CalculateTargetNumber(driverClub, difficulty)
	expectedInitial := 4

	if targetNumber != expectedInitial {
		t.Errorf("Initial target number = %v, want %v", targetNumber, expectedInitial)
	}

	golfer.Skills["Driver"] = progression.Skill{Name: "Driver", Level: 3, Experience: 0}
	golfer.Abilities["Strength"] = progression.Ability{Name: "Strength", Level: 4, Experience: 0}

	targetNumber = golfer.CalculateTargetNumber(driverClub, difficulty)
	expectedLeveled := 14

	if targetNumber != expectedLeveled {
		t.Errorf("After level-up target number = %v, want %v", targetNumber, expectedLeveled)
	}
}

func TestGameLoopIntegration_XPAward(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "Putter"}

	initialSkillXP := golfer.Skills["Putter"].Experience
	initialAbilityXP := golfer.Abilities["Mental"].Experience

	golfer.AwardExperience(club, 10)

	if golfer.Skills["Putter"].Experience != initialSkillXP+10 {
		t.Errorf("Skill XP = %v, want %v", golfer.Skills["Putter"].Experience, initialSkillXP+10)
	}

	if golfer.Abilities["Mental"].Experience != initialAbilityXP+10 {
		t.Errorf("Ability XP = %v, want %v", golfer.Abilities["Mental"].Experience, initialAbilityXP+10)
	}
}

func TestGameLoopIntegration_LevelUpDetection(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	club := Club{Name: "PW"}

	skill := golfer.GetSkillForClub(club)
	ability := golfer.GetAbilityForClub(club)

	prevSkillLevel := golfer.Skills[skill.Name].Level
	prevAbilityLevel := golfer.Abilities[ability.Name].Level

	golfer.AwardExperience(club, 100)

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
