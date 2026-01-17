package main

import "testing"

// Test Ability creation with valid initial values
func TestNewAbility(t *testing.T) {
	ability := NewAbility("Strength")

	if ability.Name != "Strength" {
		t.Errorf("NewAbility name = %v, want Strength", ability.Name)
	}

	if ability.Level != 1 {
		t.Errorf("NewAbility level = %v, want 1", ability.Level)
	}

	if ability.Experience != 0 {
		t.Errorf("NewAbility experience = %v, want 0", ability.Experience)
	}
}

// Test Ability.Value() returns correct value based on level
func TestAbility_Value(t *testing.T) {
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
		ability := Ability{Name: "Test", Level: tt.level, Experience: 0}
		got := ability.Value()

		if got != tt.expectedValue {
			t.Errorf("Ability.Value() with level %d = %v, want %v",
				tt.level, got, tt.expectedValue)
		}
	}
}

// Test Ability.AddExperience() increases experience
func TestAbility_AddExperience(t *testing.T) {
	ability := NewAbility("Control")

	ability.AddExperience(10)
	if ability.Experience != 10 {
		t.Errorf("After AddExperience(10), experience = %v, want 10", ability.Experience)
	}

	ability.AddExperience(5)
	if ability.Experience != 15 {
		t.Errorf("After AddExperience(5), experience = %v, want 15", ability.Experience)
	}
}

// Test Ability.AddExperience() triggers level up at threshold
func TestAbility_AddExperience_LevelUp(t *testing.T) {
	ability := NewAbility("Touch")
	// Level 1 → 2 requires 100 XP

	ability.AddExperience(99)
	if ability.Level != 1 {
		t.Errorf("At 99 XP, level = %v, want 1 (no level up yet)", ability.Level)
	}

	ability.AddExperience(1) // Total 100 XP
	if ability.Level != 2 {
		t.Errorf("At 100 XP, level = %v, want 2 (leveled up)", ability.Level)
	}

	if ability.Experience != 0 {
		t.Errorf("After level up, experience = %v, want 0 (reset)", ability.Experience)
	}
}

// Test Ability.AddExperience() respects max level 9
func TestAbility_AddExperience_MaxLevel(t *testing.T) {
	ability := Ability{Name: "Mental", Level: 9, Experience: 0}

	ability.AddExperience(1000)

	if ability.Level != 9 {
		t.Errorf("At max level, level = %v, want 9 (no level up beyond max)", ability.Level)
	}

	// Experience should still accumulate (or cap, depending on design choice)
	// For now, let's assume it caps at 0 when at max level
	if ability.Experience != 0 {
		t.Errorf("At max level, experience = %v, want 0 (no XP accumulation)", ability.Experience)
	}
}

// Test Ability.ExperienceToNextLevel() returns correct threshold
func TestAbility_ExperienceToNextLevel(t *testing.T) {
	tests := []struct {
		level    int
		expected int
	}{
		{1, 100},  // Level 1 → 2 = 100 XP
		{2, 150},  // Level 2 → 3 = 150 XP
		{3, 200},  // Level 3 → 4 = 200 XP
		{8, 450},  // Level 8 → 9 = 450 XP
		{9, 0},    // Level 9 (max) = 0 (no next level)
	}

	for _, tt := range tests {
		ability := Ability{Name: "Test", Level: tt.level, Experience: 0}
		got := ability.ExperienceToNextLevel()

		if got != tt.expected {
			t.Errorf("ExperienceToNextLevel() at level %d = %v, want %v",
				tt.level, got, tt.expected)
		}
	}
}

// Test Ability.CanLevelUp() returns true when enough XP
func TestAbility_CanLevelUp(t *testing.T) {
	ability := NewAbility("Strength")

	if ability.CanLevelUp() {
		t.Errorf("With 0 XP, CanLevelUp() = true, want false")
	}

	ability.Experience = 99
	if ability.CanLevelUp() {
		t.Errorf("With 99 XP, CanLevelUp() = true, want false")
	}

	ability.Experience = 100
	if !ability.CanLevelUp() {
		t.Errorf("With 100 XP, CanLevelUp() = false, want true")
	}

	ability.Level = 9
	if ability.CanLevelUp() {
		t.Errorf("At max level, CanLevelUp() = true, want false")
	}
}

// Test multiple level ups in sequence
func TestAbility_MultipleLevelUps(t *testing.T) {
	ability := NewAbility("Control")

	// Level 1 → 2 (100 XP)
	ability.AddExperience(100)
	if ability.Level != 2 {
		t.Errorf("After 100 XP, level = %v, want 2", ability.Level)
	}

	// Level 2 → 3 (150 XP)
	ability.AddExperience(150)
	if ability.Level != 3 {
		t.Errorf("After 150 more XP, level = %v, want 3", ability.Level)
	}

	// Level 3 → 4 (200 XP)
	ability.AddExperience(200)
	if ability.Level != 4 {
		t.Errorf("After 200 more XP, level = %v, want 4", ability.Level)
	}
}
