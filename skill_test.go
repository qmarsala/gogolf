package main

import "testing"

// Test Skill creation with valid initial values
func TestNewSkill(t *testing.T) {
	skill := NewSkill("Driver")

	if skill.Name != "Driver" {
		t.Errorf("NewSkill name = %v, want Driver", skill.Name)
	}

	if skill.Level != 1 {
		t.Errorf("NewSkill level = %v, want 1", skill.Level)
	}

	if skill.Experience != 0 {
		t.Errorf("NewSkill experience = %v, want 0", skill.Experience)
	}
}

// Test Skill.Value() returns correct value based on level
func TestSkill_Value(t *testing.T) {
	tests := []struct {
		level         int
		expectedValue int
	}{
		{1, 2},  // Level 1 = 2 value
		{2, 4},  // Level 2 = 4 value
		{3, 6},  // Level 3 = 6 value
		{5, 10}, // Level 5 = 10 value
		{9, 18}, // Level 9 = 18 value (max level)
	}

	for _, tt := range tests {
		skill := Skill{Name: "Test", Level: tt.level, Experience: 0}
		got := skill.Value()

		if got != tt.expectedValue {
			t.Errorf("Skill.Value() with level %d = %v, want %v",
				tt.level, got, tt.expectedValue)
		}
	}
}

// Test Skill.AddExperience() increases experience
func TestSkill_AddExperience(t *testing.T) {
	skill := NewSkill("Putter")

	skill.AddExperience(10)
	if skill.Experience != 10 {
		t.Errorf("After AddExperience(10), experience = %v, want 10", skill.Experience)
	}

	skill.AddExperience(5)
	if skill.Experience != 15 {
		t.Errorf("After AddExperience(5), experience = %v, want 15", skill.Experience)
	}
}

// Test Skill.AddExperience() triggers level up at threshold
func TestSkill_AddExperience_LevelUp(t *testing.T) {
	skill := NewSkill("Woods")
	// Level 1 → 2 requires 100 XP

	skill.AddExperience(99)
	if skill.Level != 1 {
		t.Errorf("At 99 XP, level = %v, want 1 (no level up yet)", skill.Level)
	}

	skill.AddExperience(1) // Total 100 XP
	if skill.Level != 2 {
		t.Errorf("At 100 XP, level = %v, want 2 (leveled up)", skill.Level)
	}

	if skill.Experience != 0 {
		t.Errorf("After level up, experience = %v, want 0 (reset)", skill.Experience)
	}
}

// Test Skill.AddExperience() respects max level 9
func TestSkill_AddExperience_MaxLevel(t *testing.T) {
	skill := Skill{Name: "Wedges", Level: 9, Experience: 0}

	skill.AddExperience(1000)

	if skill.Level != 9 {
		t.Errorf("At max level, level = %v, want 9 (no level up beyond max)", skill.Level)
	}

	// Experience should not accumulate at max level
	if skill.Experience != 0 {
		t.Errorf("At max level, experience = %v, want 0 (no XP accumulation)", skill.Experience)
	}
}

// Test Skill.ExperienceToNextLevel() returns correct threshold
func TestSkill_ExperienceToNextLevel(t *testing.T) {
	tests := []struct {
		level    int
		expected int
	}{
		{1, 100}, // Level 1 → 2 = 100 XP
		{2, 150}, // Level 2 → 3 = 150 XP
		{3, 200}, // Level 3 → 4 = 200 XP
		{8, 400}, // Level 8 → 9 = 400 XP
		{9, 0},   // Level 9 (max) = 0 (no next level)
	}

	for _, tt := range tests {
		skill := Skill{Name: "Test", Level: tt.level, Experience: 0}
		got := skill.ExperienceToNextLevel()

		if got != tt.expected {
			t.Errorf("ExperienceToNextLevel() at level %d = %v, want %v",
				tt.level, got, tt.expected)
		}
	}
}

// Test Skill.CanLevelUp() returns true when enough XP
func TestSkill_CanLevelUp(t *testing.T) {
	skill := NewSkill("Long Irons")

	if skill.CanLevelUp() {
		t.Errorf("With 0 XP, CanLevelUp() = true, want false")
	}

	skill.Experience = 99
	if skill.CanLevelUp() {
		t.Errorf("With 99 XP, CanLevelUp() = true, want false")
	}

	skill.Experience = 100
	if !skill.CanLevelUp() {
		t.Errorf("With 100 XP, CanLevelUp() = false, want true")
	}

	skill.Level = 9
	if skill.CanLevelUp() {
		t.Errorf("At max level, CanLevelUp() = true, want false")
	}
}

// Test skill type names match expected club categories
func TestSkillTypes(t *testing.T) {
	expectedSkills := []string{
		"Driver",
		"Woods",
		"Long Irons",
		"Mid Irons",
		"Short Irons",
		"Wedges",
		"Putter",
	}

	for _, skillName := range expectedSkills {
		skill := NewSkill(skillName)
		if skill.Name != skillName {
			t.Errorf("NewSkill(%v) name = %v, want %v", skillName, skill.Name, skillName)
		}
	}
}

// Test multiple level ups in sequence
func TestSkill_MultipleLevelUps(t *testing.T) {
	skill := NewSkill("Mid Irons")

	// Level 1 → 2 (100 XP)
	skill.AddExperience(100)
	if skill.Level != 2 {
		t.Errorf("After 100 XP, level = %v, want 2", skill.Level)
	}

	// Level 2 → 3 (150 XP)
	skill.AddExperience(150)
	if skill.Level != 3 {
		t.Errorf("After 150 more XP, level = %v, want 3", skill.Level)
	}

	// Level 3 → 4 (200 XP)
	skill.AddExperience(200)
	if skill.Level != 4 {
		t.Errorf("After 200 more XP, level = %v, want 4", skill.Level)
	}
}
