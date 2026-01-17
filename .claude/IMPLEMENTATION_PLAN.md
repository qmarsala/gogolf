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
- **Equipment System & ProShop** (PR #5 ready)
  - Currency system with tiered hole rewards
  - 3 equipment types (Ball, Glove, Shoes) with stat bonuses
  - ProShop with 10 items across 3 categories
  - Purchase system with auto-equipping
  - Equipment bonuses integrated into shot calculations
  - Money awarded after each hole based on score

### ❌ Missing Core RPG Features
- Save/Load system for persistence
- ProShop UI for in-game browsing/purchasing

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

## ✅ Phase 4: Equipment System & ProShop [COMPLETE]
**Status: READY FOR PR** - RPG progression mechanic fully functional

### Completed Implementation

**Currency System:**
- Added `Money` field to Golfer struct (starts with 100)
- `AddMoney()` method to increase currency
- `SpendMoney()` method with insufficient funds checking
- Currency persisted throughout gameplay

**Hole Reward System:**
- `CalculateHoleReward()` function with tiered rewards:
  - Hole-in-one (1 stroke): 100 money
  - Eagle (2 under par): 50 money
  - Birdie (1 under par): 25 money
  - Par: 10 money
  - Bogey (1 over par): 5 money
  - Double bogey or worse: 1 money
- `AwardHoleReward()` method automatically awards money after each hole
- Rewards displayed to player after hole completion

**Equipment Types:**
- **Ball**: Name, DistanceBonus, SpinControl, Cost
- **Glove**: Name, AccuracyBonus, Cost
- **Shoes**: Name, LiePenaltyReduction, Cost
- Equipment fields added to Golfer (Ball, Glove, Shoes)
- Equip methods: `EquipBall()`, `EquipGlove()`, `EquipShoes()`
- `GetEquippedBall()` method for querying equipped items

**ProShop System:**
- `ProShop` struct with equipment inventory
- `NewProShop()` creates shop with starter inventory:
  - 4 ball options (Budget Ball: 20, Standard: 35, Premium: 50, Pro V1: 75)
  - 3 glove options (Basic: 25, Leather Pro: 45, Precision Grip: 65)
  - 3 shoe options (Casual Spikes: 30, All-Terrain Pro: 55, Tour Edition: 80)
- Purchase methods: `PurchaseBall()`, `PurchaseGlove()`, `PurchaseShoes()`
- Automatic equipment after purchase
- Proper error handling (insufficient funds, item not found)

**Equipment Bonuses:**
- `GetTotalLiePenaltyReduction()` - shoes reduce lie difficulty penalties
- Updated `CalculateTargetNumber()` to apply shoe bonuses to lie difficulty
- `GetModifiedClub()` applies equipment bonuses:
  - Ball adds distance bonus to shots
  - Glove improves club accuracy (capped at 1.0)
  - Multiple bonuses stack correctly
- Equipment bonuses integrated into all shot calculations

**Main Game Loop Integration:**
- Equipment bonuses automatically applied via `GetModifiedClub()`
- Money awarded after each hole completion
- Player sees reward amount and total money
- Full equipment system functional in gameplay

**Test Coverage:**
- Added 44 new tests for Phase 4
- 184 total test cases (all passing)
- Tests cover:
  - Currency system (7 tests)
  - Hole rewards (9 tests)
  - Equipment types (10 tests)
  - ProShop purchasing (10 tests)
  - Equipment bonuses (8 tests)
- 100% test coverage on new code
- Strict TDD methodology followed

**Files Created:**
- currency_test.go
- reward.go, reward_test.go
- equipment.go, equipment_test.go
- proshop.go, proshop_test.go
- equipment_bonus_test.go

**Files Modified:**
- golfer.go (Money, equipment fields, bonus methods)
- main.go (equipment integration, money rewards)

---

## Phase 5: Save/Load System
**Priority: MEDIUM** - Required for RPG progression persistence

### Requirements
Not explicitly in CLAUDE.md but needed for long-term progression

### Design Decisions
- JSON serialization for human-readable save files
- Save golfer state (name, skills, abilities, equipment, currency)
- Auto-save after each round
- Manual save/load options
- Multiple save slots (3-5 slots)

### Implementation Tasks
1. Create `SaveGame` struct with all player data
2. Implement JSON serialization/deserialization
3. Create save file management (write, read, list)
4. Add auto-save after round completion
5. Add manual save/load commands
6. Handle save file versioning for future updates

---

## Phase 6: ProShop UI & Browsing
**Priority: MEDIUM** - Enhance equipment shopping experience

### Requirements
> "Equipment, such as clubs, balls, gloves, can be purchased from the ProShop" (from CLAUDE.md)

### Current State
- ProShop exists with inventory and purchase methods
- Purchasing works programmatically but no in-game UI
- Players need a way to browse and buy equipment during gameplay

### Design Decisions
- Text-based menu system (fits CLI nature)
- Display equipment with stats and prices
- Compare with currently equipped items
- Show affordability based on current money
- Allow browsing without purchasing

### Implementation Tasks
1. Create `DisplayProShop()` function to show inventory
2. Create equipment comparison display (current vs. available)
3. Add interactive menu for browsing categories (Balls/Gloves/Shoes)
4. Add purchase confirmation prompts
5. Integrate ProShop access into main game loop (e.g., between rounds)
6. Display money and current equipment when entering shop
7. Add "Back to Game" option to exit shop

**Example UI Flow:**
```
=== ProShop ===
Money: 150

1. Balls
2. Gloves
3. Shoes
4. Back to Game

> 1

=== Balls ===
Currently equipped: Standard Ball (+3 distance, 0.5 spin)

Available:
1. Budget Ball - 20 money (+0 distance, 0.3 spin)
2. Premium Ball - 50 money (+5 distance, 0.7 spin)
3. Pro V1 - 75 money (+8 distance, 0.9 spin)
4. Back

> 2
Purchase Premium Ball for 50 money? (y/n)
```

---

## Phase 7: Advanced Course Features
**Priority: LOW** - Polish and depth

### Features
- Wind (affects shot difficulty and ball flight)
- Elevation changes (uphill/downhill)
- Greens with break/slope
- Hazards (water, trees, OB)
- Course variety (links, parkland, desert)

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

### Sprint 4: Progression (Equipment) ✅ COMPLETE
13. ✅ Create equipment types (Ball, Glove, Shoes) (PR #5)
14. ✅ Create ProShop with inventory (PR #5)
15. ✅ Add currency and rewards (PR #5)
16. ✅ Equipment bonuses to stats (PR #5)

### Sprint 5: Persistence (NEXT)
17. ⬜ Save/load system
18. ⬜ Auto-save functionality
19. ⬜ Multiple save slots

### Sprint 6: Polish
20. ⬜ ProShop browsing UI
21. ⬜ Advanced course features
22. ⬜ Additional UI improvements

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

**Immediate Priority:** Phase 5 - Save/Load System

With the core RPG mechanics and equipment progression fully functional (Phases 1-4 complete), the next priority is implementing save/load functionality to persist player progress between sessions.

**Why Save/Load is Next:**
- Players can lose hours of progression without saves
- Required before adding more content (players need to keep progress)
- Enables long-term character development
- Natural checkpoint: good time to implement before more features

**Key Features to Implement:**
1. **Save Game Structure** - JSON serialization of all player data
2. **Auto-Save** - Automatic saves after round completion
3. **Manual Save/Load** - Player-controlled save points
4. **Multiple Slots** - Allow different character progressions

**After Save/Load:**
- Phase 6: ProShop UI for in-game equipment browsing
- Phase 7: Advanced course features (wind, elevation, hazards)

**Design Considerations:**
- Save file format should be human-readable (JSON)
- Handle backward compatibility as game evolves
- Graceful error handling for corrupted saves
- Clear save/load confirmation messages
