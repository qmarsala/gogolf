# GoGolf Implementation Plan
## Addressing CLAUDE.md Shortcomings

Based on the project vision in CLAUDE.md, this document outlines the gaps between the current implementation and the target RPG golf experience.

---

## Current State Analysis

### ✅ Implemented
- 3D6 dice roll mechanics with 7-tier outcomes (Critical Success → Critical Failure)
- Basic club system (14 clubs with Accuracy/Forgiveness)
- Course and hole generation
- Ball physics with rotation and power
- Score tracking system
- Basic game loop

### ❌ Missing Core RPG Features

---

## Phase 1: Skills & Abilities System
**Priority: HIGH** - This is the foundation of the RPG mechanics

### Requirements (from CLAUDE.md)
> "A character has skills and abilities, such as strength, accuracy, green reading, irons, drivers, etc..."
> "Skills/Abilities gain experience when used for a skill check"
> "Skills/Abilities can be leveled to a maximum level of 9"

### Design Decisions

**Abilities (Attributes):**
- Strength (power/distance)
- Control (accuracy/consistency)
- Touch (finesse shots, putting)
- Mental (course management, green reading)

**Skills (Club Proficiencies):**
- Driver
- Woods
- Long Irons (3-5)
- Mid Irons (6-7)
- Short Irons (8-9)
- Wedges
- Putter

**Leveling System:**
- Levels: 1-9
- Experience gained per skill check
- Level determines skill value contribution to target number

### Implementation Tasks
1. Create `Ability` type with name, level (1-9), experience
2. Create `Skill` type with name, level (1-9), experience
3. Add `Skills` and `Abilities` maps to `Golfer` struct
4. Create level → value conversion formula (e.g., level * 2 = bonus to target number)
5. Create experience gain system (award XP on skill check)
6. Create level-up logic (threshold-based)
7. Update SkillCheck to accept skill/ability and calculate dynamic target number

**Formula for Target Number:**
```
targetNumber = baseSkillValue + relevantAbilityValue + difficultyModifier
```

---

## Phase 2: Dynamic Target Numbers
**Priority: HIGH** - Directly tied to Phase 1

### Requirements (from CLAUDE.md)
> "target number is determined by the skills involved in the check as well as the difficulty of the action"

### Current Issue
Target number is hardcoded to 10 in `main.go:48`

### Implementation Tasks
1. Remove hardcoded target number
2. Create method: `Golfer.CalculateTargetNumber(skill Skill, ability Ability, difficulty int) int`
3. Integrate lie difficulty into target number
4. Update all SkillCheck calls to use calculated target numbers

**Example:**
```go
// Driver shot from fairway
skill := golfer.Skills["Driver"]
ability := golfer.Abilities["Strength"]
lieDifficulty := currentLie.Difficulty // 0 = fairway, -2 = rough, -4 = bunker
targetNumber := skill.Value() + ability.Value() + lieDifficulty
```

---

## Phase 3: Course Grid & Lie System
**Priority: MEDIUM** - Adds strategic depth

### Requirements (from CLAUDE.md)
> "The course is a grid of locations the ball can be within"
> "Each location will have properties that determine the lie of the ball"
> "ex: a location may be considered to be 'in a bunker', 'in the rough', 'or in the fairway'"

### Current Issue
- Course exists but not as a spatial grid
- No lie detection
- Ball location is just X/Y coordinates without context

### Design Decisions

**Lie Types:**
- Tee (easiest)
- Fairway (normal)
- First Cut (slight penalty)
- Rough (moderate penalty)
- Deep Rough (heavy penalty)
- Bunker (very heavy penalty, special mechanics)
- Green (putting only)
- Penalty Area (stroke penalty)

**Grid System:**
- Divide hole into grid cells (e.g., 10-yard squares)
- Each cell has a `LieType`
- Ball location maps to grid cell
- Lie affects target number difficulty

### Implementation Tasks
1. Create `LieType` enum
2. Create `GridCell` type with LieType, position, properties
3. Create `CourseGrid` structure (2D array of GridCells)
4. Update `Hole` to include `CourseGrid`
5. Update ball landing logic to determine lie
6. Create lie → difficulty modifier mapping
7. Integrate lie difficulty into target number calculation

**Example Grid:**
```
[Tee][Fairway][Fairway][Rough][Fairway][Green]
```

---

## Phase 4: Equipment System & ProShop
**Priority: MEDIUM** - RPG progression mechanic

### Requirements (from CLAUDE.md)
> "Equipment, such as clubs, balls, gloves, can be purchased from the ProShop"
> "Both equipment and skills/abilities are used to determine the outcome of actions"

### Current Issue
- Only clubs exist
- No equipment variation
- No purchasing system
- No currency/economy

### Design Decisions

**Equipment Types:**
- Clubs (already exist, need shop integration)
- Balls (distance, spin, control modifiers)
- Gloves (accuracy bonus)
- Shoes (lie penalty reduction)

**Currency:**
- Earn money through rounds (par/birdie bonuses)
- Starter equipment vs premium equipment

**Equipment Stats:**
```go
type Ball struct {
    Name string
    DistanceBonus float32
    SpinControl float32
    Cost int
}

type Glove struct {
    Name string
    AccuracyBonus float32
    Cost int
}
```

### Implementation Tasks
1. Create equipment types (Ball, Glove, Shoes)
2. Create `ProShop` with inventory
3. Add currency to `Golfer`
4. Create purchase interface
5. Update shot calculations to include equipment bonuses
6. Create reward system (earn money per round)

---

## Phase 5: Save/Load System
**Priority: LOW** - Nice to have for persistence

### Requirements
Not explicitly in CLAUDE.md but needed for RPG progression

### Implementation Tasks
1. Create `SaveGame` serialization (JSON)
2. Save golfer state (skills, equipment, currency)
3. Load game on startup
4. Multiple save slots

---

## Phase 6: Advanced Course Features
**Priority: LOW** - Polish and depth

### Features
- Wind (affects shot difficulty and ball flight)
- Elevation changes (uphill/downhill)
- Greens with break/slope
- Hazards (water, trees, OB)

---

## Recommended Implementation Order

### Sprint 1: Foundation (Skills System)
1. ✅ Tiered skill checks (DONE)
2. ⬜ Implement Abilities and Skills
3. ⬜ Add experience/leveling system
4. ⬜ Dynamic target numbers

### Sprint 2: Course Depth (Lie System)
5. ⬜ Create course grid system
6. ⬜ Implement lie types
7. ⬜ Integrate lie difficulty into target numbers

### Sprint 3: Progression (Equipment)
8. ⬜ Expand equipment system
9. ⬜ Create ProShop
10. ⬜ Add currency and rewards

### Sprint 4: Polish
11. ⬜ Save/load system
12. ⬜ Advanced course features
13. ⬜ UI improvements

---

## Key Design Principles

### From CLAUDE.md TDD Requirements:
- Write tests first (expected input/output)
- Verify test fails
- Commit failing tests
- Implement minimal code to pass
- Commit working code

### RPG Core Loop:
```
Action → Skill Check → XP Gain → Level Up → Better Stats → Harder Courses
```

### Golf Realism:
- Lie matters (rough is harder than fairway)
- Club selection matters (right tool for the job)
- Course management (risk/reward decisions)
- Equipment helps but skill is primary

---

## Next Steps

**Immediate Priority:** Phase 1 - Skills & Abilities System

This is the foundation that enables dynamic target numbers, experience progression, and the core RPG loop.

**Recommended First Task:**
Create the Skills/Abilities data structures and leveling system, starting with TDD approach.
