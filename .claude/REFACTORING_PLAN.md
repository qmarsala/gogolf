# Codebase Refactoring Plan

## Goal
Define clear types, structs, and interfaces at the root level (the WHAT) and move implementations to subpackages (the HOW).

## Completed

### Phase 1: Define Root Interfaces
- Added `interfaces.go` with `DiceRoller`, `RandomSource`, `ClubSelector`
- Updated `Golfer.SkillCheck` to accept `DiceRoller` interface
- Updated `Game` and `CalculateRotation` to use `RandomSource` interface

### Phase 2: Move Game Orchestration
- Created `game/` subpackage with `Game`, `Context`, `ShotResult`
- Moved game loop orchestration, shot execution, XP calculation
- Updated `cmd/console/main.go` to use `game.New()`

---

## Remaining Steps

### Phase 3: Move Shot Mechanics to `mechanics/` Subpackage

**Files to move:**
- `shot_mechanics.go` → `mechanics/shot.go`

**Functions:**
- `CalculateRotation(club, result, random)` - rotation based on skill check
- `CalculatePower(club, power, result)` - power adjustment based on outcome
- `GetShotQualityDescription(result)` - description text for outcomes

**Steps:**
1. Create branch `refactor/mechanics-package`
2. Create `mechanics/` directory
3. Move functions to `mechanics/shot.go`
4. Update imports in `game/game.go`
5. Move tests to `mechanics/shot_test.go`
6. Verify all tests pass
7. Create PR

---

### Phase 4: Move ProShop to `shop/` Subpackage

**Files to move:**
- `proshop.go` → `shop/proshop.go`

**Types and Functions:**
- `ProShop` struct with inventory
- `NewProShop()` factory
- `PurchaseBall()`, `PurchaseGlove()`, `PurchaseShoes()`

**Steps:**
1. Create branch `refactor/shop-package`
2. Create `shop/` directory
3. Move `ProShop` to `shop/proshop.go`
4. Update any imports (currently no consumers in codebase)
5. Move tests to `shop/proshop_test.go`
6. Verify all tests pass
7. Create PR

---

### Phase 5: Move ScoreCard to `scoring/` Subpackage

**Files to move:**
- `scorecard.go` → `scoring/scorecard.go`

**Types and Functions:**
- `ScoreCard` struct
- `RecordStroke()`, `TotalStrokes()`, `Score()`, etc.

**Steps:**
1. Create branch `refactor/scoring-package`
2. Create `scoring/` directory
3. Move `ScoreCard` to `scoring/scorecard.go`
4. Update imports in `game/game.go`
5. Move tests to `scoring/scorecard_test.go`
6. Verify all tests pass
7. Create PR

---

## Final Structure

```
gogolf/
├── interfaces.go        # DiceRoller, RandomSource, ClubSelector
├── types.go             # Yard, Point, Vector, LieType, SkillCheckOutcome
├── golfer.go            # Golfer entity
├── course.go            # Course, Hole entities
├── golfball.go          # GolfBall entity
├── equipment.go         # Ball, Glove, Shoes structs
├── grid.go              # CourseGrid for lie detection
│
├── game/                # Game orchestration
│   └── game.go
├── mechanics/           # Shot physics calculations
│   └── shot.go
├── shop/                # Equipment purchasing
│   └── proshop.go
├── scoring/             # Score tracking
│   └── scorecard.go
├── dice/                # Dice rolling implementation
│   └── dice.go
├── progression/         # Skills & abilities with XP
│   ├── skill.go
│   ├── ability.go
│   └── reward.go
├── ui/                  # Terminal rendering
│   └── ...
└── cmd/console/         # Entry point
    └── main.go
```

## Notes
- Each phase should follow TDD practices per CLAUDE.md
- Create feature branches and PRs for each phase
- Root package retains domain entities (Golfer, Club, Course, etc.)
- Subpackages implement behavior and orchestration
