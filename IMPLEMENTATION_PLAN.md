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

### Remaining Integration Work
While the system is built, it's not yet integrated into the main game loop:
- [ ] Update main.go to use `CalculateTargetNumber()` instead of hardcoded `10`
- [ ] Award XP after each shot based on outcome
- [ ] Display skill/ability progression to player
- [ ] Integrate lie difficulty into target number calculation

---

## Phase 2: Main Game Loop Integration
**Priority: HIGH** - Connect Phase 1 to gameplay

### Requirements
Make the Skills & Abilities system affect actual gameplay and award experience for player progression.

### Current Issue
[main.go:48](main.go#L48) still uses hardcoded `targetNumber = 10`

### Implementation Tasks

**1. Dynamic Target Numbers**
```go
// Replace this (main.go line 48):
result := golfer.SkillCheck(NewD6(), 10)

// With this:
club := golfer.GetBestClub(targetDistance)
difficulty := 0 // Will come from lie system in Phase 3
targetNumber := golfer.CalculateTargetNumber(club, difficulty)
result := golfer.SkillCheck(NewD6(), targetNumber)
```

**2. XP Award System**
```go
// After each shot, award XP based on outcome
xpAward := calculateXP(result.Outcome)
golfer.AwardExperience(club, xpAward)

// Suggested XP values:
// CriticalSuccess: 15 XP
// Excellent: 10 XP
// Good: 7 XP
// Marginal: 5 XP
// Poor: 3 XP
// Bad: 2 XP
// CriticalFailure: 1 XP
```

**3. Progression Feedback**
```go
// Notify player of level-ups
if leveledUp {
    fmt.Printf("🎉 %s leveled up to %d!\n", skillName, newLevel)
}
```

**4. Player Stats Display**
- Show current skills/abilities in game UI
- Display XP progress toward next level
- Show value contributions to target numbers

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

### Sprint 2: Game Loop Integration (CURRENT)
5. ⬜ Replace hardcoded target numbers with `CalculateTargetNumber()`
6. ⬜ Implement XP award system after each shot
7. ⬜ Add player stat display
8. ⬜ Add level-up notifications

### Sprint 3: Course Depth (Lie System)
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

**Immediate Priority:** Phase 2 - Main Game Loop Integration

With the Skills & Abilities foundation complete (Phase 1), we need to integrate it into actual gameplay so players can use and improve their skills.

**Recommended First Task:**
Update [main.go](main.go) to use dynamic target numbers based on golfer skills/abilities instead of the hardcoded value of 10. Follow TDD approach by writing tests that verify target numbers change based on skill levels.
