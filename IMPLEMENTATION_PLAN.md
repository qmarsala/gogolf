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

### Remaining Work for Future Phases
- [ ] Integrate lie difficulty into target number calculation (Phase 3)

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

### Sprint 3: Course Depth (Lie System) (CURRENT)
9. ⬜ Create course grid system
10. ⬜ Implement lie types
11. ⬜ Integrate lie difficulty into target numbers
12. ⬜ Update ball landing to detect lie

### Sprint 4: Progression (Equipment)
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

**Immediate Priority:** Phase 3 - Course Grid & Lie System

With the RPG mechanics fully functional (Phases 1 & 2 complete), the next step is to add strategic depth by implementing a course grid system where ball position affects difficulty.

**Key Features to Implement:**
1. **Lie Types** - Tee, Fairway, Rough, Bunker, Green, etc.
2. **Grid System** - Spatial grid that maps ball position to lie type
3. **Difficulty Modifiers** - Each lie type affects target number differently:
   - Tee: +2 (easier)
   - Fairway: 0 (normal)
   - Rough: -2 (harder)
   - Bunker: -4 (very hard)
4. **Ball Position Detection** - Determine lie after each shot lands

**Recommended First Task:**
Create the `LieType` enum and difficulty modifier system. Follow TDD approach by writing tests that verify different lies produce different target numbers.
