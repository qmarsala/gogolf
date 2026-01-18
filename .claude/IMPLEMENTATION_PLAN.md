# GoGolf Implementation Plan

## Current Sprint: Mechanics Alignment

This sprint focuses on aligning the codebase with the updated [CORE_MECHANICS.md](CORE_MECHANICS.md) design document.

---

## Task 1: Skill/Ability Value Calculation

**Problem:** Code uses `level * 2` but design specifies `level`

### Files to Update

| File | Line | Current | Target |
|------|------|---------|--------|
| `skill.go` | 17-19 | `return s.Level * 2` | `return s.Level` |
| `ability.go` | 17-19 | `return a.Level * 2` | `return a.Level` |

### Tests to Update

| File | Lines | Change |
|------|-------|--------|
| `skill_test.go` | 21-42 | Update expected values: Level 1→1, 2→2, 3→3, 5→5, 9→9 |
| `ability_test.go` | 21-42 | Update expected values: Level 1→1, 2→2, 3→3, 5→5, 9→9 |
| `golfer_test.go` | 76-92 | Update CalculateTargetNumber: skill(3)+ability(4)+difficulty(0) = 7 |
| `golfer_test.go` | 273-295 | Update DynamicWithLevels: initial=2, after leveling=7 |
| `golfer_test.go` | 347-373 | Update CalculateTargetNumberWithShape test expectations |

---

## Task 2: Outcome Threshold Alignment

**Problem:** Code thresholds don't match CORE_MECHANICS.md table

### Design Document Says:
| Tier | Condition |
|------|-----------|
| Excellent | Margin +6 or more |
| Good | Margin +3 to +5 |
| Marginal | Margin 0 to +2 |
| Poor | Margin -1 to -3 |
| Bad | Margin -4 to -6 |

### Code Currently Has (golfer.go:122-133):
```go
case margin >= 4:  return Excellent
case margin >= 1:  return Good
case margin == 0:  return Marginal
case margin >= -3: return Poor
default:           return Bad
```

### Decision Required:
- Option A: Update code to match document (stricter thresholds)
- Option B: Update document to match code (current gameplay balance)

---

## Task 3: Verify Roll-Under Logic

**Status:** Already correct ✓

The current implementation in `golfer.go:99-112`:
```go
margin := targetNumber - total
Success: margin >= 0
```

This correctly implements "roll at or below target":
- Roll 8 vs Target 10 → margin = +2 → Success ✓
- Roll 12 vs Target 10 → margin = -2 → Failure ✓

No changes needed for this task.

---

## Implementation Order

Following TDD approach from CLAUDE.md:

### Step 1: Update Tests First
1. Update `skill_test.go` expected values
2. Update `ability_test.go` expected values
3. Update `golfer_test.go` target number expectations
4. Verify tests fail

### Step 2: Commit Failing Tests
```
git commit -m "test: update expectations for level=value formula"
```

### Step 3: Update Implementation
1. Change `skill.go` Value() to return `s.Level`
2. Change `ability.go` Value() to return `a.Level`
3. Run tests to verify pass

### Step 4: Commit Working Code
```
git commit -m "refactor: skill/ability value equals level (not level*2)"
```

### Step 5: Address Outcome Thresholds (if needed)
1. Decide on correct thresholds
2. Update tests or document accordingly
3. Follow same TDD flow

---

## Test Commands

```bash
go test ./...                    # Run all tests
go test -v -run TestSkill_Value  # Run specific test
go test -v -run TestAbility_Value
go test -v -run TestGolfer_CalculateTargetNumber
```

---

## Completion Checklist

- [ ] Skill.Value() returns level (not level*2)
- [ ] Ability.Value() returns level (not level*2)
- [ ] All tests updated and passing
- [ ] Outcome thresholds aligned (code ↔ document)
- [ ] Failing tests committed before implementation
- [ ] Working code committed after implementation
