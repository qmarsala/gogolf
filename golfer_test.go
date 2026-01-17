package main

import "testing"

// Test NewGolfer creates a golfer with default skills and abilities
func TestNewGolfer(t *testing.T) {
	golfer := NewGolfer("Tiger")

	if golfer.Name != "Tiger" {
		t.Errorf("NewGolfer name = %v, want Tiger", golfer.Name)
	}

	// Should have 4 abilities
	expectedAbilities := []string{"Strength", "Control", "Touch", "Mental"}
	for _, abilityName := range expectedAbilities {
		ability, exists := golfer.Abilities[abilityName]
		if !exists {
			t.Errorf("NewGolfer missing ability: %v", abilityName)
		}
		if ability.Level != 1 {
			t.Errorf("Ability %v level = %v, want 1", abilityName, ability.Level)
		}
	}

	// Should have 7 skills
	expectedSkills := []string{"Driver", "Woods", "Long Irons", "Mid Irons", "Short Irons", "Wedges", "Putter"}
	for _, skillName := range expectedSkills {
		skill, exists := golfer.Skills[skillName]
		if !exists {
			t.Errorf("NewGolfer missing skill: %v", skillName)
		}
		if skill.Level != 1 {
			t.Errorf("Skill %v level = %v, want 1", skillName, skill.Level)
		}
	}
}

// Test Golfer.GetSkillForClub maps clubs to correct skills
func TestGolfer_GetSkillForClub(t *testing.T) {
	golfer := NewGolfer("Test")

	tests := []struct {
		clubName     string
		expectedSkill string
	}{
		{"Driver", "Driver"},
		{"3 Wood", "Woods"},
		{"5 Wood", "Woods"},
		{"4 Iron", "Long Irons"},
		{"5 Iron", "Long Irons"},
		{"6 Iron", "Mid Irons"},
		{"7 Iron", "Mid Irons"},
		{"8 Iron", "Short Irons"},
		{"9 Iron", "Short Irons"},
		{"PW", "Wedges"},
		{"GW", "Wedges"},
		{"SW", "Wedges"},
		{"LW", "Wedges"},
		{"Putter", "Putter"},
	}

	for _, tt := range tests {
		club := Club{Name: tt.clubName}
		skill := golfer.GetSkillForClub(club)

		if skill.Name != tt.expectedSkill {
			t.Errorf("GetSkillForClub(%v) = %v, want %v",
				tt.clubName, skill.Name, tt.expectedSkill)
		}
	}
}

// Test Golfer.CalculateTargetNumber with skill, ability, and difficulty
func TestGolfer_CalculateTargetNumber(t *testing.T) {
	golfer := NewGolfer("Test")

	// Set specific levels for testing
	golfer.Skills["Driver"] = Skill{Name: "Driver", Level: 3, Experience: 0}      // Value: 6
	golfer.Abilities["Strength"] = Ability{Name: "Strength", Level: 4, Experience: 0} // Value: 8

	club := Club{Name: "Driver"}
	difficulty := 0 // Fairway lie

	targetNumber := golfer.CalculateTargetNumber(club, difficulty)

	// Expected: skill(6) + ability(8) + difficulty(0) = 14
	expected := 14
	if targetNumber != expected {
		t.Errorf("CalculateTargetNumber() = %v, want %v", targetNumber, expected)
	}
}

// Test Golfer.CalculateTargetNumber with difficulty modifiers
func TestGolfer_CalculateTargetNumber_WithDifficulty(t *testing.T) {
	golfer := NewGolfer("Test")
	golfer.Skills["Driver"] = Skill{Name: "Driver", Level: 3, Experience: 0}      // Value: 6
	golfer.Abilities["Strength"] = Ability{Name: "Strength", Level: 4, Experience: 0} // Value: 8

	club := Club{Name: "Driver"}

	tests := []struct {
		difficulty int
		expected   int
	}{
		{0, 14},   // Fairway: 6 + 8 + 0 = 14
		{-2, 12},  // Rough: 6 + 8 + (-2) = 12
		{-4, 10},  // Bunker: 6 + 8 + (-4) = 10
		{2, 16},   // Tee: 6 + 8 + 2 = 16
	}

	for _, tt := range tests {
		targetNumber := golfer.CalculateTargetNumber(club, tt.difficulty)

		if targetNumber != tt.expected {
			t.Errorf("CalculateTargetNumber(difficulty=%v) = %v, want %v",
				tt.difficulty, targetNumber, tt.expected)
		}
	}
}

// Test Golfer.GetAbilityForClub maps clubs to correct abilities
func TestGolfer_GetAbilityForClub(t *testing.T) {
	golfer := NewGolfer("Test")

	tests := []struct {
		clubName        string
		expectedAbility string
	}{
		{"Driver", "Strength"},      // Woods use Strength
		{"3 Wood", "Strength"},
		{"4 Iron", "Control"},       // Long irons use Control
		{"5 Iron", "Control"},
		{"6 Iron", "Control"},       // Mid irons use Control
		{"7 Iron", "Control"},
		{"8 Iron", "Touch"},         // Short irons use Touch
		{"9 Iron", "Touch"},
		{"PW", "Touch"},             // Wedges use Touch
		{"GW", "Touch"},
		{"SW", "Touch"},
		{"LW", "Touch"},
		{"Putter", "Mental"},        // Putter uses Mental
	}

	for _, tt := range tests {
		club := Club{Name: tt.clubName}
		ability := golfer.GetAbilityForClub(club)

		if ability.Name != tt.expectedAbility {
			t.Errorf("GetAbilityForClub(%v) = %v, want %v",
				tt.clubName, ability.Name, tt.expectedAbility)
		}
	}
}

// Test Golfer.AwardExperience increases skill and ability XP
func TestGolfer_AwardExperience(t *testing.T) {
	golfer := NewGolfer("Test")
	club := Club{Name: "Driver"}
	xp := 10

	initialSkillXP := golfer.Skills["Driver"].Experience
	initialAbilityXP := golfer.Abilities["Strength"].Experience

	golfer.AwardExperience(club, xp)

	finalSkillXP := golfer.Skills["Driver"].Experience
	finalAbilityXP := golfer.Abilities["Strength"].Experience

	if finalSkillXP != initialSkillXP+xp {
		t.Errorf("Skill XP after award = %v, want %v", finalSkillXP, initialSkillXP+xp)
	}

	if finalAbilityXP != initialAbilityXP+xp {
		t.Errorf("Ability XP after award = %v, want %v", finalAbilityXP, initialAbilityXP+xp)
	}
}

// Test Golfer.AwardExperience triggers level ups
func TestGolfer_AwardExperience_LevelUp(t *testing.T) {
	golfer := NewGolfer("Test")
	club := Club{Name: "Putter"}

	// Award enough XP to level up (100 XP for level 1 → 2)
	golfer.AwardExperience(club, 100)

	if golfer.Skills["Putter"].Level != 2 {
		t.Errorf("Putter skill level = %v, want 2 (leveled up)", golfer.Skills["Putter"].Level)
	}

	if golfer.Abilities["Mental"].Level != 2 {
		t.Errorf("Mental ability level = %v, want 2 (leveled up)", golfer.Abilities["Mental"].Level)
	}
}
