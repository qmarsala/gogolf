package gogolf

import (
	"testing"
)

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
		clubName      string
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

	golfer.Skills["Driver"] = Skill{Name: "Driver", Level: 3, Experience: 0}
	golfer.Abilities["Strength"] = Ability{Name: "Strength", Level: 4, Experience: 0}

	club := Club{Name: "Driver"}
	difficulty := 0 // Fairway lie

	targetNumber := golfer.CalculateTargetNumber(club, difficulty)

	// Expected: skill(3) + ability(4) + difficulty(0) = 7
	expected := 7
	if targetNumber != expected {
		t.Errorf("CalculateTargetNumber() = %v, want %v", targetNumber, expected)
	}
}

// Test Golfer.CalculateTargetNumber with difficulty modifiers
func TestGolfer_CalculateTargetNumber_WithDifficulty(t *testing.T) {
	golfer := NewGolfer("Test")
	golfer.Skills["Driver"] = Skill{Name: "Driver", Level: 3, Experience: 0}
	golfer.Abilities["Strength"] = Ability{Name: "Strength", Level: 4, Experience: 0}

	club := Club{Name: "Driver"}

	tests := []struct {
		difficulty int
		expected   int
	}{
		{0, 7},  // Fairway: 3 + 4 + 0 = 7
		{-2, 5}, // Rough: 3 + 4 + (-2) = 5
		{-4, 3}, // Bunker: 3 + 4 + (-4) = 3
		{2, 9},  // Tee: 3 + 4 + 2 = 9
	}

	for _, tt := range tests {
		targetNumber := golfer.CalculateTargetNumber(club, tt.difficulty)

		if targetNumber != tt.expected {
			t.Errorf("CalculateTargetNumber(difficulty=%v) = %v, want %v",
				tt.difficulty, targetNumber, tt.expected)
		}
	}
}

func TestGolfer_CalculateTargetNumber_MinimumOfThree(t *testing.T) {
	golfer := NewGolfer("Test")
	golfer.Skills["Driver"] = Skill{Name: "Driver", Level: 1, Experience: 0}
	golfer.Abilities["Strength"] = Ability{Name: "Strength", Level: 1, Experience: 0}

	club := Club{Name: "Driver"}

	tests := []struct {
		difficulty int
		expected   int
		desc       string
	}{
		{-4, 3, "Deep Rough: 1 + 1 + (-4) = -2, clamped to 3"},
		{-6, 3, "Extreme penalty: 1 + 1 + (-6) = -4, clamped to 3"},
		{-1, 3, "First Cut: 1 + 1 + (-1) = 1, clamped to 3"},
		{0, 3, "Fairway: 1 + 1 + 0 = 2, clamped to 3"},
		{2, 4, "Tee: 1 + 1 + 2 = 4, no clamping needed"},
	}

	for _, tt := range tests {
		targetNumber := golfer.CalculateTargetNumber(club, tt.difficulty)

		if targetNumber != tt.expected {
			t.Errorf("%s: CalculateTargetNumber(difficulty=%v) = %v, want %v",
				tt.desc, tt.difficulty, targetNumber, tt.expected)
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
		{"Driver", "Strength"}, // Woods use Strength
		{"3 Wood", "Strength"},
		{"4 Iron", "Control"}, // Long irons use Control
		{"5 Iron", "Control"},
		{"6 Iron", "Control"}, // Mid irons use Control
		{"7 Iron", "Control"},
		{"8 Iron", "Touch"}, // Short irons use Touch
		{"9 Iron", "Touch"},
		{"PW", "Touch"}, // Wedges use Touch
		{"GW", "Touch"},
		{"SW", "Touch"},
		{"LW", "Touch"},
		{"Putter", "Mental"}, // Putter uses Mental
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

func TestDetermineOutcome_MarginBased(t *testing.T) {
	tests := []struct {
		name       string
		margin     int
		isCritical bool
		expected   SkillCheckOutcome
	}{
		{"Margin +10", 10, false, Excellent},
		{"Margin +8", 8, false, Excellent},
		{"Margin +7", 7, false, Excellent},
		{"Margin +6", 6, false, Excellent},

		{"Margin +5", 5, false, Good},
		{"Margin +4", 4, false, Good},
		{"Margin +3", 3, false, Good},

		{"Margin +2", 2, false, Marginal},
		{"Margin +1", 1, false, Marginal},
		{"Margin 0", 0, false, Marginal},

		{"Margin -1", -1, false, Poor},
		{"Margin -2", -2, false, Poor},
		{"Margin -3", -3, false, Poor},

		{"Margin -4", -4, false, Bad},
		{"Margin -5", -5, false, Bad},
		{"Margin -6", -6, false, Bad},
		{"Margin -7", -7, false, Bad},
		{"Margin -8", -8, false, Bad},
		{"Margin -10", -10, false, Bad},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineOutcome(tt.margin, tt.isCritical)
			if got != tt.expected {
				t.Errorf("determineOutcome(%d, %v) = %v, want %v",
					tt.margin, tt.isCritical, got, tt.expected)
			}
		})
	}
}

func TestDetermineOutcome_TierBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		margin   int
		expected SkillCheckOutcome
	}{
		{"Just below Excellent", 5, Good},
		{"Just at Excellent", 6, Excellent},
		{"High margin still Excellent", 10, Excellent},

		{"Just below Good", 2, Marginal},
		{"Just at Good", 3, Good},
		{"Top of Good", 5, Good},

		{"Just below Marginal", -1, Poor},
		{"At Marginal", 0, Marginal},
		{"Top of Marginal", 2, Marginal},

		{"Just below Poor", -4, Bad},
		{"Just at Poor", -3, Poor},
		{"Top of Poor", -1, Poor},

		{"Just below Bad", -7, Bad},
		{"Just at Bad", -6, Bad},
		{"Top of Bad", -4, Bad},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineOutcome(tt.margin, false)
			if got != tt.expected {
				t.Errorf("determineOutcome(%d, false) = %v, want %v",
					tt.margin, got, tt.expected)
			}
		})
	}
}

func TestGolfer_CalculateTargetNumber_DynamicWithLevels(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	driverClub := Club{Name: "Driver"}
	difficulty := 0

	targetNumber := golfer.CalculateTargetNumber(driverClub, difficulty)
	expectedInitial := 2

	if targetNumber != expectedInitial {
		t.Errorf("Initial target number = %v, want %v", targetNumber, expectedInitial)
	}

	golfer.Skills["Driver"] = Skill{Name: "Driver", Level: 3, Experience: 0}
	golfer.Abilities["Strength"] = Ability{Name: "Strength", Level: 4, Experience: 0}

	targetNumber = golfer.CalculateTargetNumber(driverClub, difficulty)
	expectedLeveled := 7

	if targetNumber != expectedLeveled {
		t.Errorf("After level-up target number = %v, want %v", targetNumber, expectedLeveled)
	}
}

func TestGolfer_AwardExperience_XPDistribution(t *testing.T) {
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

func TestGolfer_AwardExperience_LevelUpDetection(t *testing.T) {
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

func TestGolfer_CalculateTargetNumberWithShape(t *testing.T) {
	golfer := NewGolfer("Test")
	golfer.Skills["Driver"] = Skill{Name: "Driver", Level: 3, Experience: 0}
	golfer.Abilities["Strength"] = Ability{Name: "Strength", Level: 4, Experience: 0}

	club := Club{Name: "Driver"}
	lieDifficulty := 0

	tests := []struct {
		shape    ShotShape
		expected int
	}{
		{Straight, 5},
		{Draw, 8},
		{Fade, 8},
		{Hook, 6},
		{Slice, 6},
	}

	for _, tt := range tests {
		targetNumber := golfer.CalculateTargetNumberWithShape(club, lieDifficulty, tt.shape)
		if targetNumber != tt.expected {
			t.Errorf("CalculateTargetNumberWithShape(%v) = %v, want %v",
				tt.shape, targetNumber, tt.expected)
		}
	}
}

func TestGolfer_DrawFadeEasierThanStraight(t *testing.T) {
	golfer := NewGolfer("Test")
	club := Club{Name: "7-Iron"}

	straightTarget := golfer.CalculateTargetNumberWithShape(club, 0, Straight)
	drawTarget := golfer.CalculateTargetNumberWithShape(club, 0, Draw)
	fadeTarget := golfer.CalculateTargetNumberWithShape(club, 0, Fade)

	if drawTarget <= straightTarget {
		t.Errorf("Draw target (%d) should be higher than Straight (%d)", drawTarget, straightTarget)
	}
	if fadeTarget <= straightTarget {
		t.Errorf("Fade target (%d) should be higher than Straight (%d)", fadeTarget, straightTarget)
	}
}
