# GoGolf Layered Architecture Refactoring Plan

## Current Structure

All game logic lives in the root package (`main`) with UI in `ui/`:

```
gogolf/
├── main.go              # Entry point + orchestration
├── game.go              # Game struct, TakeShot, state management
├── golfer.go            # Golfer, Club structs
├── skill.go             # Skill struct + leveling
├── ability.go           # Ability struct + leveling
├── equipment.go         # Ball, Glove, Shoes structs
├── proshop.go           # ProShop struct + inventory
├── course.go            # Course, Hole structs
├── grid.go              # CourseGrid, GridCell structs
├── golfball.go          # GolfBall struct
├── lie.go               # LieType enum
├── dice.go              # Dice, SkillCheckResult structs
├── shot_mechanics.go    # Shot calculation functions
├── scorecard.go         # ScoreCard struct
├── reward.go            # Reward calculation functions
├── vector.go            # Vector struct (2D math)
├── world.go             # Point, Size, Unit types
├── measurements.go      # Yard, Foot, Inch conversions
├── *_test.go            # Test files
└── ui/
    ├── gamestate.go     # GameState display struct
    ├── renderer.go      # Renderer struct
    ├── terminal.go      # Terminal control
    ├── layout.go        # Panel layout
    └── input.go         # Power meter input
```

**Problems:**
1. Everything in root makes it hard to navigate
2. No clear separation between domain types and implementation
3. Adding Phase 5 (save/load) will add more files to root
4. Difficult to swap implementations (e.g., file-based vs cloud saves)

---

## Proposed Layered Structure

```
gogolf/
├── main.go                    # Entry point, dependency injection, orchestration
│
├── game.go                    # Game interface/struct (core game loop)
├── golfer.go                  # Golfer interface + struct
├── course.go                  # Course, Hole interfaces/structs
├── scorecard.go               # ScoreCard interface/struct
├── equipment.go               # Ball, Glove, Shoes interfaces/structs
├── types.go                   # Shared domain types: Point, Vector, Unit, LieType
│
├── dice/                      # Dice rolling implementation
│   ├── dice.go                # Dice struct, Roll(), RollN()
│   ├── skillcheck.go          # SkillCheckResult, determineOutcome()
│   └── dice_test.go
│
├── shot/                      # Shot mechanics implementation
│   ├── mechanics.go           # CalculateRotation(), CalculatePower()
│   ├── ball.go                # GolfBall struct, ReceiveHit(), collision
│   └── *_test.go
│
├── course/                    # Course implementation details
│   ├── grid.go                # CourseGrid, GridCell implementation
│   ├── generator.go           # GenerateSimpleCourse() and future generators
│   └── *_test.go
│
├── progression/               # RPG progression implementation
│   ├── skill.go               # Skill leveling, XP calculation
│   ├── ability.go             # Ability leveling, XP calculation
│   ├── reward.go              # Hole reward calculation
│   └── *_test.go
│
├── shop/                      # Pro shop implementation
│   ├── proshop.go             # ProShop inventory, purchasing
│   └── proshop_test.go
│
├── persistence/               # Phase 5: Save/Load (future)
│   ├── save_data.go           # SaveData struct, serialization
│   └── persistence_test.go
│
├── filedb/                    # Phase 5: File-based persistence (future)
│   ├── repository.go          # File read/write implementation
│   └── repository_test.go
│
└── ui/                        # UI implementation (already separated)
    ├── gamestate.go
    ├── renderer.go
    ├── terminal.go
    ├── layout.go
    └── input.go
```

---

## Root Level Files (Interfaces & Core Types)

### types.go (NEW)
Consolidate shared domain types:
- `Point` (from world.go)
- `Size` (from world.go)
- `Vector` (from vector.go)
- `Unit`, `Yard`, `Foot`, `Inch` + conversions (from measurements.go)
- `LieType` enum + `DifficultyModifier()` (from lie.go)

### game.go (MODIFY)
Keep the `Game` struct and core game loop logic:
- `Game` struct
- `GameContext` struct
- `ShotResult` struct
- `NewGame()`, `TakeShot()`, `IsHoleComplete()`, `IsRoundComplete()`

### golfer.go (MODIFY)
Keep the `Golfer` and `Club` structs:
- `Golfer` struct (references Skill, Ability, Equipment)
- `Club` struct
- `GetBestClub()`, `SkillCheck()`, `AwardExperience()`
- Equipment binding methods

### course.go (MODIFY)
Keep the `Course` and `Hole` structs:
- `Course` struct
- `Hole` struct
- `DetectHoleOut()`, `DetectTapIn()`, `GetLieAtPosition()`

### scorecard.go (KEEP)
Already well-scoped, keep as-is.

### equipment.go (KEEP)
Already well-scoped with `Ball`, `Glove`, `Shoes` structs.

---

## Subpackages

### dice/
Move dice rolling and skill check logic:
- `Dice` struct
- `SkillCheckResult`, `SkillCheckOutcome`
- `Roll()`, `RollN()`, `determineOutcome()`

### shot/
Move shot mechanics:
- `GolfBall` struct (from golfball.go)
- `CalculateRotation()`, `CalculatePower()` (from shot_mechanics.go)
- `GetShotQualityDescription()`

### course/
Move course implementation details:
- `CourseGrid`, `GridCell` (from grid.go)
- `GenerateSimpleCourse()` (from course.go)

### progression/
Move RPG leveling system:
- `Skill` struct (from skill.go)
- `Ability` struct (from ability.go)
- `CalculateHoleReward()`, `AwardHoleReward()` (from reward.go)
- Leveling formulas, XP calculations

### shop/
Move pro shop:
- `ProShop` struct (from proshop.go)
- Inventory, purchasing logic

### persistence/ (Phase 5)
Define interfaces for save/load:
- `SaveData` struct
- `Repository` interface with `Save()`, `Load()`, `List()`, `Delete()`

### filedb/ (Phase 5)
File-based implementation:
- `FileRepository` implementing `Repository` interface
- JSON serialization to local files

---

## Migration Steps

### Step 1: Create types.go
1. Create `types.go` with consolidated types
2. Update imports across codebase
3. Delete `world.go`, `vector.go`, `measurements.go`, `lie.go`
4. Run tests

### Step 2: Create dice/ package
1. Create `dice/dice.go` and `dice/skillcheck.go`
2. Move tests to `dice/dice_test.go`
3. Update imports in `golfer.go`, `game.go`
4. Delete root `dice.go`
5. Run tests

### Step 3: Create shot/ package
1. Create `shot/mechanics.go` and `shot/ball.go`
2. Move tests
3. Update imports in `game.go`
4. Delete root `shot_mechanics.go`, `golfball.go`
5. Run tests

### Step 4: Create course/ package
1. Create `course/grid.go` and `course/generator.go`
2. Move grid tests
3. Update imports in `course.go`
4. Delete root `grid.go`
5. Run tests

### Step 5: Create progression/ package
1. Create `progression/skill.go`, `progression/ability.go`, `progression/reward.go`
2. Move tests
3. Update imports in `golfer.go`
4. Delete root `skill.go`, `ability.go`, `reward.go`
5. Run tests

### Step 6: Create shop/ package
1. Create `shop/proshop.go`
2. Move tests
3. Update imports
4. Delete root `proshop.go`
5. Run tests

### Step 7: Clean up test files
1. Consolidate test files into appropriate packages
2. Remove orphaned test files from root
3. Ensure all 184 tests pass

---

## Package Dependencies (After Refactor)

```
main.go
  ├── game.go (uses: dice/, shot/, course/, progression/)
  ├── golfer.go (uses: dice/, progression/, equipment.go)
  ├── course.go (uses: course/, types.go)
  ├── scorecard.go
  ├── equipment.go
  └── ui/

game.go
  ├── dice/
  ├── shot/
  ├── course/
  └── progression/

golfer.go
  ├── dice/
  ├── progression/
  └── equipment.go

course.go
  └── course/

shop/
  └── equipment.go
```

---

## Benefits

1. **Clear separation**: Domain types at root, implementations in packages
2. **Extensibility**: Easy to add `clouddb/` alongside `filedb/` for Phase 5
3. **Testability**: Each package has focused tests
4. **Navigation**: Logical grouping makes codebase easier to understand
5. **Dependency injection**: Interfaces at root allow swapping implementations

---

## Phase 5 Integration Example

With this structure, Phase 5 (Save/Load) would look like:

**Root: game_data.go**
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

func (r *FileRepository) Save(slot int, data SaveData) error {
    // JSON serialize to basePath/slot_N.json
}

func (r *FileRepository) Load(slot int) (SaveData, error) {
    // JSON deserialize from basePath/slot_N.json
}
```

**main.go**
```go
func main() {
    repo := filedb.NewRepository("./saves")
    // Inject repo into game setup
}
```

---

## Notes

- Each step should be a separate commit following TDD principles
- Run `go test ./...` after each step to ensure no regressions
- Consider adding a `go.mod` package name update if needed
- The `ui/` package is already well-structured and can remain as-is
