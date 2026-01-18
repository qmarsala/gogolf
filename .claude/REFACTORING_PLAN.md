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

### Phase 3: Move Shot Mechanics to `mechanics/` Subpackage
- Created `mechanics/` package with `CalculateRotation`, `CalculatePower`, `GetShotQualityDescription`
- Updated `game/game.go` to use mechanics package
- Updated `equipment_bonus_test.go` to use external test package

### Phase 4: Move ProShop to `shop/` Subpackage
- Created `shop/` package with `ProShop` and purchase methods
- Updated types to reference `gogolf.Ball`, `gogolf.Glove`, `gogolf.Shoes`, `gogolf.Golfer`

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
├── scorecard.go         # ScoreCard entity (domain struct)
├── grid.go              # CourseGrid for lie detection
│
├── game/                # Game orchestration
│   └── game.go
├── mechanics/           # Shot physics calculations
│   └── shot.go
├── shop/                # Equipment purchasing
│   └── proshop.go
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
- Root package retains domain entities (Golfer, Club, Course, ScoreCard, etc.)
- Subpackages implement behavior and orchestration
