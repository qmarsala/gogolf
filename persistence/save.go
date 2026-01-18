package persistence

import (
	"encoding/json"
	"fmt"
	"gogolf"
	"gogolf/progression"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const CurrentSaveVersion = 1
const MaxSaveSlots = 5

type SkillData struct {
	Name       string `json:"name"`
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
}

func NewSkillData(skill progression.Skill) SkillData {
	return SkillData{
		Name:       skill.Name,
		Level:      skill.Level,
		Experience: skill.Experience,
	}
}

func (sd SkillData) ToSkill() progression.Skill {
	return progression.Skill{
		Name:       sd.Name,
		Level:      sd.Level,
		Experience: sd.Experience,
	}
}

type AbilityData struct {
	Name       string `json:"name"`
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
}

func NewAbilityData(ability progression.Ability) AbilityData {
	return AbilityData{
		Name:       ability.Name,
		Level:      ability.Level,
		Experience: ability.Experience,
	}
}

func (ad AbilityData) ToAbility() progression.Ability {
	return progression.Ability{
		Name:       ad.Name,
		Level:      ad.Level,
		Experience: ad.Experience,
	}
}

type SaveData struct {
	Version    int                    `json:"version"`
	SavedAt    time.Time              `json:"saved_at"`
	GolferName string                 `json:"golfer_name"`
	Money      int                    `json:"money"`
	Skills     map[string]SkillData   `json:"skills"`
	Abilities  map[string]AbilityData `json:"abilities"`
	Ball       *gogolf.Ball           `json:"ball,omitempty"`
	Glove      *gogolf.Glove          `json:"glove,omitempty"`
	Shoes      *gogolf.Shoes          `json:"shoes,omitempty"`
}

func NewSaveData(golfer gogolf.Golfer) SaveData {
	skills := make(map[string]SkillData)
	for name, skill := range golfer.Skills {
		skills[name] = NewSkillData(skill)
	}

	abilities := make(map[string]AbilityData)
	for name, ability := range golfer.Abilities {
		abilities[name] = NewAbilityData(ability)
	}

	return SaveData{
		Version:    CurrentSaveVersion,
		SavedAt:    time.Now(),
		GolferName: golfer.Name,
		Money:      golfer.Money,
		Skills:     skills,
		Abilities:  abilities,
		Ball:       golfer.Ball,
		Glove:      golfer.Glove,
		Shoes:      golfer.Shoes,
	}
}

func (sd SaveData) ToGolfer() gogolf.Golfer {
	golfer := gogolf.NewGolfer(sd.GolferName)
	golfer.Money = sd.Money

	for name, skillData := range sd.Skills {
		golfer.Skills[name] = skillData.ToSkill()
	}

	for name, abilityData := range sd.Abilities {
		golfer.Abilities[name] = abilityData.ToAbility()
	}

	golfer.Ball = sd.Ball
	golfer.Glove = sd.Glove
	golfer.Shoes = sd.Shoes

	return golfer
}

type SaveSlotInfo struct {
	Slot       int
	GolferName string
	SavedAt    time.Time
}

type SaveManager struct {
	saveDir string
}

func NewSaveManager(saveDir string) *SaveManager {
	return &SaveManager{saveDir: saveDir}
}

func (sm *SaveManager) slotPath(slot int) string {
	return filepath.Join(sm.saveDir, fmt.Sprintf("save_slot_%d.json", slot))
}

func (sm *SaveManager) validateSlot(slot int) error {
	if slot < 1 || slot > MaxSaveSlots {
		return fmt.Errorf("invalid save slot %d: must be between 1 and %d", slot, MaxSaveSlots)
	}
	return nil
}

func (sm *SaveManager) Save(slot int, golfer gogolf.Golfer) error {
	if err := sm.validateSlot(slot); err != nil {
		return err
	}

	if err := os.MkdirAll(sm.saveDir, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	saveData := NewSaveData(golfer)
	jsonBytes, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize save data: %w", err)
	}

	if err := os.WriteFile(sm.slotPath(slot), jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

func (sm *SaveManager) Load(slot int) (gogolf.Golfer, error) {
	if err := sm.validateSlot(slot); err != nil {
		return gogolf.Golfer{}, err
	}

	jsonBytes, err := os.ReadFile(sm.slotPath(slot))
	if err != nil {
		return gogolf.Golfer{}, fmt.Errorf("failed to read save file: %w", err)
	}

	var saveData SaveData
	if err := json.Unmarshal(jsonBytes, &saveData); err != nil {
		return gogolf.Golfer{}, fmt.Errorf("failed to parse save file: %w", err)
	}

	return saveData.ToGolfer(), nil
}

func (sm *SaveManager) Delete(slot int) error {
	if err := sm.validateSlot(slot); err != nil {
		return err
	}

	if err := os.Remove(sm.slotPath(slot)); err != nil {
		return fmt.Errorf("failed to delete save file: %w", err)
	}

	return nil
}

func (sm *SaveManager) SlotExists(slot int) bool {
	if err := sm.validateSlot(slot); err != nil {
		return false
	}

	_, err := os.Stat(sm.slotPath(slot))
	return err == nil
}

func (sm *SaveManager) ListSaveSlots() []SaveSlotInfo {
	var slots []SaveSlotInfo

	for slot := 1; slot <= MaxSaveSlots; slot++ {
		if !sm.SlotExists(slot) {
			continue
		}

		jsonBytes, err := os.ReadFile(sm.slotPath(slot))
		if err != nil {
			continue
		}

		var saveData SaveData
		if err := json.Unmarshal(jsonBytes, &saveData); err != nil {
			continue
		}

		slots = append(slots, SaveSlotInfo{
			Slot:       slot,
			GolferName: saveData.GolferName,
			SavedAt:    saveData.SavedAt,
		})
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].Slot < slots[j].Slot
	})

	return slots
}
