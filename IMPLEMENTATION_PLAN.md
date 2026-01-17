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
- **Skills & Abilities System** (PR #2 merged)
  - 4 Abilities: Strength, Control, Touch, Mental
  - 7 Skills: Driver, Woods, Long Irons, Mid Irons, Short Irons, Wedges, Putter
  - Experience and leveling (1-9 levels)
  - Club-to-skill/ability mapping
  - Dynamic target number calculation
- **Main Game Loop Integration** (PR #3 merged)
  - Dynamic target numbers based on skill + ability
  - XP award system (1-15 XP per shot based on outcome)
  - Level-up notifications during gameplay
  - Player stats display (skills/abilities with XP progress)
  - Full character progression system active
- **Course Grid & Lie System** (PR #4 merged)
  - 8 LieTypes with difficulty modifiers
  - Spatial grid system for course layout
  - Ball position to lie detection
  - Lie difficulty integrated into target numbers
  - GenerateSimpleCourse with automatic lie setup

### ❌ Missing Core RPG Features

---

## ✅ Phase 1: Skills & Abilities System [COMPLETE]
**Status: MERGED** (PR #2) - Foundation of RPG mechanics established

### Completed Implementation

**Created Types:**
- `Ability` - Character attributes (Strength, Control, Touch, Mental)
- `Skill` - Club proficiencies (Driver, Woods, Long/Mid/Short Irons, Wedges, Putter)

**Key Features:**
- Level system: 1-9 levels with `Value() = level * 2`
- XP progression: `(level + 1) * 50` XP to next level
- Auto-leveling with overflow XP handling
- Map-based skill/ability storage in Golfer

**New Golfer Methods:**
- `NewGolfer(name)` - Constructor with default skills/abilities at level 1
- `GetSkillForClub(club)` - Maps club to relevant skill
- `GetAbilityForClub(club)` - Maps club to relevant ability
- `CalculateTargetNumber(club, difficulty)` - Dynamic target number formula
- `AwardExperience(club, xp)` - Grants XP to both skill and ability

**Test Coverage:**
- 52 tests total (all passing)
- 100% coverage on new code
- TDD methodology followed strictly

---

## ✅ Phase 2: Main Game Loop Integration [COMPLETE]
**Status: MERGED** (PR #3) - RPG mechanics now active in gameplay

### Completed Implementation

**Dynamic Target Numbers:**
- Replaced hardcoded `targetNumber = 10` with `golfer.CalculateTargetNumber(club, difficulty)`
- Target numbers now scale: Level 1 golfer ~4, Level 9 golfer ~36
- Formula: `skill.Value() + ability.Value() + difficulty`

**XP Award System:**
- Created `calculateXP()` helper function with tier-based XP values:
  - Critical Success: 15 XP
  - Excellent: 10 XP
  - Good: 7 XP
  - Marginal: 5 XP
  - Poor: 3 XP
  - Bad: 2 XP
  - Critical Failure: 1 XP
- XP awarded after each shot to both skill and ability

**Level-Up Notifications:**
- Detects level-ups by comparing pre/post levels
- Displays celebration messages: "🎉 Driver leveled up to 2!"
- Shows XP progress after each shot: `XP: +7 (Driver: 45/100, Strength: 32/100)`

**Player Stats Display:**
- Added `displayPlayerStats()` function
- Shows all 7 skills and 4 abilities with:
  - Current level and value contribution
  - XP progress toward next level (e.g., "45/100 XP")
  - MAX indicator for level 9 skills/abilities
- Displayed at start of round and after completion

**Bug Fixes:**
- Fixed `AwardExperience` pointer semantics to properly modify skill/ability copies

**Test Coverage:**
- Added 4 new integration tests
- Total: 68 tests passing
- Tests verify: calculateXP, dynamic target numbers, XP awards, level-ups

---

## ✅ Phase 3: Course Grid & Lie System [COMPLETE]
**Status: MERGED** (PR #4) - Strategic depth added to shot difficulty

### Completed Implementation

**LieType System:**
- Created `LieType` enum with 8 types:
  - **Tee** (+2): Easiest, ball teed up perfectly
  - **Fairway** (0): Normal lie, no penalty
  - **First Cut** (-1): Slight penalty
  - **Rough** (-2): Moderate penalty
  - **Deep Rough** (-4): Heavy penalty, ball buried
  - **Bunker** (-4): Very hard, sand shot
  - **Green** (+1): Putting is easier
  - **Penalty Area** (0): Same as fairway after drop
- Each lie type has `DifficultyModifier()` method

**Course Grid System:**
- Created `GridCell` type with position and lie properties
- Created `CourseGrid` with 2D cell array
- Implemented `NewCourseGrid(width, length, cellSize)` constructor
- Implemented `GetLieAtPosition()` for position-to-lie mapping
- Implemented `SetLieAtPosition()` for grid configuration
- Added bounds checking (returns PenaltyArea for out-of-bounds)
- Configurable cell size (e.g., 10-yard squares)
- Unit conversion handling: 1 yard = 72 units

**Hole Integration:**
- Added `Grid *CourseGrid` field to `Hole` struct
- Created `NewHoleWithGrid()` constructor
- Added `Hole.GetLieAtPosition()` method
- Maintained backward compatibility (nil grid defaults to Fairway)

**Ball Lie Detection:**
- Added `GolfBall.GetLie(hole)` method
- Automatic lie detection based on ball position

**Course Generation:**
- Created `GenerateSimpleCourse()` function
- Automatic lie setup:
  - Tee area: First 10 yards
  - Green area: Last 20 yards around hole
  - Fairway: Middle area (default)

**Main Game Loop Integration:**
- Switched to `GenerateSimpleCourse()`
- Ball lie detected before each shot
- Lie difficulty modifier applied to target number
- Player feedback: `Lie: Tee (difficulty: +2)`

**Gameplay Impact:**
Target numbers now vary by lie (Level 1 golfer with Driver, base ~4):
- From Tee: 6 (easier)
- From Fairway: 4 (normal)
- From Rough: 2 (harder)
- From Bunker: 0 (very hard)

**Test Coverage:**
- Added 13 new tests (140 total test cases)
- All tests passing
- Tests cover:
  - LieType enum and difficulty modifiers (3 tests)
  - CourseGrid position mapping (6 tests)
  - Hole/Ball integration (4 tests)
  - Full integration testing

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

### Sprint 1: Foundation (Skills System) ✅ COMPLETE
1. ✅ Tiered skill checks (PR #1)
2. ✅ Implement Abilities and Skills (PR #2)
3. ✅ Add experience/leveling system (PR #2)
4. ✅ Dynamic target number calculation (PR #2)

### Sprint 2: Game Loop Integration ✅ COMPLETE
5. ✅ Replace hardcoded target numbers with `CalculateTargetNumber()` (PR #3)
6. ✅ Implement XP award system after each shot (PR #3)
7. ✅ Add player stat display (PR #3)
8. ✅ Add level-up notifications (PR #3)

### Sprint 3: Course Depth (Lie System) ✅ COMPLETE
9. ✅ Create course grid system (PR #4)
10. ✅ Implement lie types (PR #4)
11. ✅ Integrate lie difficulty into target numbers (PR #4)
12. ✅ Update ball landing to detect lie (PR #4)

### Sprint 4: Progression (Equipment) (CURRENT)
13. ⬜ Expand equipment system
14. ⬜ Create ProShop
15. ⬜ Add currency and rewards
16. ⬜ Equipment bonuses to stats

### Sprint 5: Polish
17. ⬜ Save/load system
18. ⬜ Advanced course features
19. ⬜ UI improvements

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

**Immediate Priority:** Phase 4 - Equipment System & ProShop

With the core RPG mechanics fully functional (Phases 1-3 complete), the next step is to add equipment progression by implementing a shop system where players can purchase upgraded equipment using currency earned from rounds.

**Key Features to Implement:**
1. **Equipment Types** - Expand beyond clubs to balls, gloves, shoes
2. **ProShop System** - Interface to browse and purchase equipment
3. **Currency System** - Earn money based on round performance
4. **Equipment Bonuses** - Items provide stat bonuses (accuracy, distance, etc.)

**Design Considerations:**
- Equipment should complement skills, not replace them
- Starter equipment vs premium equipment tiers
- Reward good play (par/birdie bonuses) with currency
- Balance equipment costs with earning rate

**Recommended First Task:**
Create the currency system for the Golfer, then add reward logic for completing holes based on score (birdie, par, bogey). Follow TDD approach.
