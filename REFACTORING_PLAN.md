# GoGolf Layered Architecture Refactoring

## Completed Structure

```
gogolf/
├── main.go                    # Entry point, dependency injection, orchestration
│
├── types.go                   # Shared domain types: Point, Vector, Unit, LieType
├── game.go                    # Game struct, TakeShot, state management
├── golfer.go                  # Golfer, Club structs
├── course.go                  # Course, Hole structs
├── scorecard.go               # ScoreCard struct
├── equipment.go               # Ball, Glove, Shoes structs
├── proshop.go                 # ProShop inventory and purchasing
├── golfball.go                # GolfBall movement and collision
├── grid.go                    # CourseGrid for lie detection
├── shot_mechanics.go          # Shot calculation functions
│
├── dice/                      # Dice rolling and skill checks
│   ├── dice.go                # Dice, SkillCheckResult, SkillCheckOutcome
│   └── dice_test.go
│
├── progression/               # RPG leveling system
│   ├── skill.go               # Skill struct and leveling
│   ├── ability.go             # Ability struct and leveling
│   ├── reward.go              # CalculateHoleReward
│   └── *_test.go
│
└── ui/                        # Terminal UI (unchanged)
    ├── gamestate.go
    ├── renderer.go
    ├── terminal.go
    ├── layout.go
    └── input.go
```

## Changes Made

### Step 1: Consolidated types.go
Merged `world.go`, `vector.go`, `measurements.go`, and `lie.go` into a single `types.go`:
- `Point`, `Size` structs
- `Vector` struct with math operations
- `Unit`, `Yard`, `Foot`, `Inch` types with conversions
- `LieType` enum with `DifficultyModifier()`

### Step 2: Created dice/ package
Extracted dice rolling and skill check logic:
- `Dice` struct with `Roll()`, `RollN()`
- `SkillCheckOutcome` enum (CriticalFailure → CriticalSuccess)
- `SkillCheckResult` struct
- `DetermineOutcome()` function

### Step 3: Created progression/ package
Extracted RPG progression system:
- `Skill` struct with leveling (1-9)
- `Ability` struct with leveling (1-9)
- `CalculateHoleReward()` function
- XP formulas: `(level + 1) * 50` to next level
- Value formula: `level * 2`

## Package Dependencies

```
main
├── dice/           # Skill check mechanics
├── progression/    # Leveling and rewards
└── ui/            # Terminal rendering

game.go, golfer.go, shot_mechanics.go
├── dice/
└── progression/
```

## Not Extracted (Due to Tight Coupling)

The following were considered but not extracted due to circular dependencies:

- **shot/** - `GolfBall.GetLie()` depends on `Hole`, shot mechanics use `Club`
- **course/** - `CourseGrid` tightly coupled with `Hole` and `LieType`
- **shop/** - `ProShop.Purchase*()` methods call `Golfer` methods directly

These could be extracted in the future using interfaces, but the current structure is simpler.

## Phase 5 Integration (Future)

For the Save/Load system, define interfaces at root:

**game_data.go**
```go
type SaveData struct {
    GolferName string
    Skills     map[string]SkillData
    Abilities  map[string]AbilityData
    Equipment  EquipmentData
    Money      int
    Version    int
}

type Repository interface {
    Save(slot int, data SaveData) error
    Load(slot int) (SaveData, error)
    List() ([]SaveSlot, error)
    Delete(slot int) error
}
```

**filedb/repository.go**
```go
type FileRepository struct {
    basePath string
}

func (r *FileRepository) Save(slot int, data SaveData) error { ... }
func (r *FileRepository) Load(slot int) (SaveData, error) { ... }
```

## Benefits Achieved

1. **Clear type organization** - Domain types consolidated in `types.go`
2. **Isolated dice mechanics** - Easy to test and modify independently
3. **Separated progression system** - RPG formulas in one place
4. **Ready for Phase 5** - Clean separation for adding persistence
5. **Reduced root clutter** - From 20+ files to logical groupings
