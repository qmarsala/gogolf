package persistence

import (
	"encoding/json"
	"gogolf"
	"gogolf/progression"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveDataContainsGolferState(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 250
	golfer.Ball = &gogolf.Ball{Name: "Pro V1", DistanceBonus: 8, SpinControl: 0.9, Cost: 75}
	golfer.Glove = &gogolf.Glove{Name: "Leather Pro", AccuracyBonus: 0.05, Cost: 45}
	golfer.Shoes = &gogolf.Shoes{Name: "Tour Edition", LiePenaltyReduction: 3, Cost: 80}

	saveData := NewSaveData(golfer)

	if saveData.GolferName != "TestPlayer" {
		t.Errorf("expected golfer name 'TestPlayer', got '%s'", saveData.GolferName)
	}
	if saveData.Money != 250 {
		t.Errorf("expected money 250, got %d", saveData.Money)
	}
	if saveData.Ball == nil || saveData.Ball.Name != "Pro V1" {
		t.Errorf("expected ball 'Pro V1', got %v", saveData.Ball)
	}
	if saveData.Glove == nil || saveData.Glove.Name != "Leather Pro" {
		t.Errorf("expected glove 'Leather Pro', got %v", saveData.Glove)
	}
	if saveData.Shoes == nil || saveData.Shoes.Name != "Tour Edition" {
		t.Errorf("expected shoes 'Tour Edition', got %v", saveData.Shoes)
	}
}

func TestSaveDataContainsSkillsAndAbilities(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	skill := golfer.Skills["Driver"]
	(&skill).AddExperience(150)
	golfer.Skills["Driver"] = skill

	ability := golfer.Abilities["Strength"]
	(&ability).AddExperience(200)
	golfer.Abilities["Strength"] = ability

	saveData := NewSaveData(golfer)

	if len(saveData.Skills) != 7 {
		t.Errorf("expected 7 skills, got %d", len(saveData.Skills))
	}
	driverSkill := saveData.Skills["Driver"]
	if driverSkill.Level != 2 || driverSkill.Experience != 50 {
		t.Errorf("expected Driver skill level 2 with 50 XP, got level %d with %d XP",
			driverSkill.Level, driverSkill.Experience)
	}

	if len(saveData.Abilities) != 4 {
		t.Errorf("expected 4 abilities, got %d", len(saveData.Abilities))
	}
	strengthAbility := saveData.Abilities["Strength"]
	if strengthAbility.Level != 3 || strengthAbility.Experience != 0 {
		t.Errorf("expected Strength ability level 3 with 0 XP, got level %d with %d XP",
			strengthAbility.Level, strengthAbility.Experience)
	}
}

func TestSaveDataHasVersionAndTimestamp(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	before := time.Now()
	saveData := NewSaveData(golfer)
	after := time.Now()

	if saveData.Version != CurrentSaveVersion {
		t.Errorf("expected version %d, got %d", CurrentSaveVersion, saveData.Version)
	}
	if saveData.SavedAt.Before(before) || saveData.SavedAt.After(after) {
		t.Errorf("expected timestamp between %v and %v, got %v", before, after, saveData.SavedAt)
	}
}

func TestSaveDataSerializesToJSON(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")
	golfer.Money = 150
	saveData := NewSaveData(golfer)

	jsonBytes, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		t.Fatalf("failed to serialize to JSON: %v", err)
	}

	jsonStr := string(jsonBytes)
	if len(jsonStr) == 0 {
		t.Error("expected non-empty JSON output")
	}

	var loaded SaveData
	err = json.Unmarshal(jsonBytes, &loaded)
	if err != nil {
		t.Fatalf("failed to deserialize from JSON: %v", err)
	}

	if loaded.GolferName != "TestPlayer" {
		t.Errorf("expected golfer name 'TestPlayer', got '%s'", loaded.GolferName)
	}
	if loaded.Money != 150 {
		t.Errorf("expected money 150, got %d", loaded.Money)
	}
}

func TestRestoreGolferFromSaveData(t *testing.T) {
	original := gogolf.NewGolfer("TestPlayer")
	original.Money = 300
	original.Ball = &gogolf.Ball{Name: "Premium Ball", DistanceBonus: 5, SpinControl: 0.7, Cost: 50}

	skill := original.Skills["Putter"]
	(&skill).AddExperience(100)
	original.Skills["Putter"] = skill

	saveData := NewSaveData(original)

	restored := saveData.ToGolfer()

	if restored.Name != original.Name {
		t.Errorf("expected name '%s', got '%s'", original.Name, restored.Name)
	}
	if restored.Money != original.Money {
		t.Errorf("expected money %d, got %d", original.Money, restored.Money)
	}
	if restored.Ball == nil || restored.Ball.Name != "Premium Ball" {
		t.Errorf("expected ball 'Premium Ball', got %v", restored.Ball)
	}
	if restored.Skills["Putter"].Level != 2 {
		t.Errorf("expected Putter skill level 2, got %d", restored.Skills["Putter"].Level)
	}
	if len(restored.Clubs) != 14 {
		t.Errorf("expected 14 clubs, got %d", len(restored.Clubs))
	}
}

func TestSaveAndLoadFile(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSaveManager(tempDir)

	golfer := gogolf.NewGolfer("FileTestPlayer")
	golfer.Money = 500

	err := manager.Save(1, golfer)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	expectedPath := filepath.Join(tempDir, "save_slot_1.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("save file was not created")
	}

	loaded, err := manager.Load(1)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.Name != "FileTestPlayer" {
		t.Errorf("expected name 'FileTestPlayer', got '%s'", loaded.Name)
	}
	if loaded.Money != 500 {
		t.Errorf("expected money 500, got %d", loaded.Money)
	}
}

func TestListSaveSlots(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSaveManager(tempDir)

	golfer1 := gogolf.NewGolfer("Player1")
	golfer2 := gogolf.NewGolfer("Player2")

	manager.Save(1, golfer1)
	manager.Save(3, golfer2)

	slots := manager.ListSaveSlots()

	if len(slots) != 2 {
		t.Fatalf("expected 2 save slots, got %d", len(slots))
	}

	foundSlot1 := false
	foundSlot3 := false
	for _, slot := range slots {
		if slot.Slot == 1 && slot.GolferName == "Player1" {
			foundSlot1 = true
		}
		if slot.Slot == 3 && slot.GolferName == "Player2" {
			foundSlot3 = true
		}
	}

	if !foundSlot1 {
		t.Error("expected to find slot 1 with Player1")
	}
	if !foundSlot3 {
		t.Error("expected to find slot 3 with Player2")
	}
}

func TestDeleteSaveSlot(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSaveManager(tempDir)

	golfer := gogolf.NewGolfer("ToDelete")
	manager.Save(2, golfer)

	err := manager.Delete(2)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	slots := manager.ListSaveSlots()
	if len(slots) != 0 {
		t.Errorf("expected 0 save slots after delete, got %d", len(slots))
	}
}

func TestLoadNonExistentSlotReturnsError(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSaveManager(tempDir)

	_, err := manager.Load(99)
	if err == nil {
		t.Error("expected error when loading non-existent slot")
	}
}

func TestSaveSlotValidation(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSaveManager(tempDir)
	golfer := gogolf.NewGolfer("Test")

	err := manager.Save(0, golfer)
	if err == nil {
		t.Error("expected error for slot 0")
	}

	err = manager.Save(6, golfer)
	if err == nil {
		t.Error("expected error for slot 6 (max is 5)")
	}
}

func TestSaveDataWithNilEquipment(t *testing.T) {
	golfer := gogolf.NewGolfer("NoEquipment")

	saveData := NewSaveData(golfer)

	if saveData.Ball != nil {
		t.Error("expected nil ball")
	}
	if saveData.Glove != nil {
		t.Error("expected nil glove")
	}
	if saveData.Shoes != nil {
		t.Error("expected nil shoes")
	}

	jsonBytes, err := json.Marshal(saveData)
	if err != nil {
		t.Fatalf("failed to serialize: %v", err)
	}

	var loaded SaveData
	json.Unmarshal(jsonBytes, &loaded)

	restored := loaded.ToGolfer()
	if restored.Ball != nil {
		t.Error("expected nil ball after restore")
	}
}

func TestSavePreservesAllSkillProgress(t *testing.T) {
	golfer := gogolf.NewGolfer("SkillTest")

	skillNames := []string{"Driver", "Woods", "Long Irons", "Mid Irons", "Short Irons", "Wedges", "Putter"}
	for i, name := range skillNames {
		skill := golfer.Skills[name]
		(&skill).AddExperience((i + 1) * 50)
		golfer.Skills[name] = skill
	}

	saveData := NewSaveData(golfer)
	restored := saveData.ToGolfer()

	for i, name := range skillNames {
		originalSkill := golfer.Skills[name]
		restoredSkill := restored.Skills[name]

		if originalSkill.Level != restoredSkill.Level {
			t.Errorf("skill %s: expected level %d, got %d", name, originalSkill.Level, restoredSkill.Level)
		}
		if originalSkill.Experience != restoredSkill.Experience {
			t.Errorf("skill %s: expected XP %d, got %d", name, originalSkill.Experience, restoredSkill.Experience)
		}
	}
}

func TestSkillData_Serialization(t *testing.T) {
	skillData := SkillData{
		Name:       "Driver",
		Level:      5,
		Experience: 125,
	}

	jsonBytes, _ := json.Marshal(skillData)
	var loaded SkillData
	json.Unmarshal(jsonBytes, &loaded)

	if loaded.Name != "Driver" || loaded.Level != 5 || loaded.Experience != 125 {
		t.Errorf("skill data not properly serialized: %+v", loaded)
	}
}

func TestSkillDataToSkill(t *testing.T) {
	skillData := SkillData{
		Name:       "Putter",
		Level:      3,
		Experience: 75,
	}

	skill := skillData.ToSkill()

	if skill.Name != "Putter" {
		t.Errorf("expected name 'Putter', got '%s'", skill.Name)
	}
	if skill.Level != 3 {
		t.Errorf("expected level 3, got %d", skill.Level)
	}
	if skill.Experience != 75 {
		t.Errorf("expected experience 75, got %d", skill.Experience)
	}
}

func TestAbilityDataToAbility(t *testing.T) {
	abilityData := AbilityData{
		Name:       "Strength",
		Level:      4,
		Experience: 100,
	}

	ability := abilityData.ToAbility()

	if ability.Name != "Strength" {
		t.Errorf("expected name 'Strength', got '%s'", ability.Name)
	}
	if ability.Level != 4 {
		t.Errorf("expected level 4, got %d", ability.Level)
	}
	if ability.Experience != 100 {
		t.Errorf("expected experience 100, got %d", ability.Experience)
	}
}

func TestNewSkillDataFromSkill(t *testing.T) {
	skill := progression.Skill{
		Name:       "Woods",
		Level:      2,
		Experience: 50,
	}

	skillData := NewSkillData(skill)

	if skillData.Name != "Woods" {
		t.Errorf("expected name 'Woods', got '%s'", skillData.Name)
	}
	if skillData.Level != 2 {
		t.Errorf("expected level 2, got %d", skillData.Level)
	}
	if skillData.Experience != 50 {
		t.Errorf("expected experience 50, got %d", skillData.Experience)
	}
}

func TestNewAbilityDataFromAbility(t *testing.T) {
	ability := progression.Ability{
		Name:       "Touch",
		Level:      5,
		Experience: 0,
	}

	abilityData := NewAbilityData(ability)

	if abilityData.Name != "Touch" {
		t.Errorf("expected name 'Touch', got '%s'", abilityData.Name)
	}
	if abilityData.Level != 5 {
		t.Errorf("expected level 5, got %d", abilityData.Level)
	}
	if abilityData.Experience != 0 {
		t.Errorf("expected experience 0, got %d", abilityData.Experience)
	}
}

func TestSlotExistsCheck(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSaveManager(tempDir)

	if manager.SlotExists(1) {
		t.Error("slot 1 should not exist initially")
	}

	golfer := gogolf.NewGolfer("Test")
	manager.Save(1, golfer)

	if !manager.SlotExists(1) {
		t.Error("slot 1 should exist after save")
	}
}
