# GoGolf Core Mechanics

This document captures the implemented game mechanics for reference when adding new features.

---

## Dice Roll System

**3D6 Dice Mechanics:**
- Roll 3 six-sided dice and sum the result (range: 3-18)
- Compare roll + modifiers against target number
- Critical when all 3 dice match

**7-Tier Outcome System:**
| Tier | Condition | XP Award |
|------|-----------|----------|
| Critical Success | All dice match | 15 |
| Excellent | Margin +6 or more | 10 |
| Good | Margin +3 to +5 | 7 |
| Marginal | Margin 0 to +2 | 5 |
| Poor | Margin -1 to -3 | 3 |
| Bad | Margin -4 to -6 | 2 |
| Critical Failure | All dice match | 1 |

---

## Character System

### Abilities (4 total)
Core attributes that influence shot outcomes:
- **Strength** - Power for Driver, Woods
- **Control** - Precision for Long/Mid Irons
- **Touch** - Finesse for Short Irons, Wedges
- **Mental** - Focus for Putter

### Skills (7 total)
Club proficiencies:
- **Driver** - Tee shots
- **Woods** - Fairway woods
- **Long Irons** - 3-4 irons
- **Mid Irons** - 5-7 irons
- **Short Irons** - 8-9 irons
- **Wedges** - PW, SW, LW
- **Putter** - Putting

### Leveling System
- Levels: 1-9 (max)
- Value contribution: `level`
- XP to next level: `(level + 1) * 50`
- XP awarded after each shot to relevant skill and ability
- Overflow XP carries to next level

### Club Mappings
| Club Type | Skill | Ability |
|-----------|-------|---------|
| Driver | Driver | Strength |
| 3-Wood, 5-Wood | Woods | Strength |
| 3-Iron, 4-Iron | Long Irons | Control |
| 5-Iron, 6-Iron, 7-Iron | Mid Irons | Control |
| 8-Iron, 9-Iron | Short Irons | Touch |
| PW, SW, LW | Wedges | Touch |
| Putter | Putter | Mental |

---

## Target Number Calculation

```
targetNumber = skill.Value() + ability.Value() + lieDifficulty - shoePenaltyReduction
```

**Example (Level 1 golfer, Driver from Tee):**
- Skill value: 1
- Ability value: 1
- Tee lie bonus: +2
- Target: 4

Successful rolls are when you roll at or below your target number

---

## Course & Lie System

### Lie Types (8 total)
| Lie Type | Difficulty Modifier |
|----------|-------------------|
| Tee | +2 (easiest) |
| Green | +1 |
| Fairway | 0 (baseline) |
| Penalty Area | 0 |
| First Cut | -1 |
| Rough | -2 |
| Deep Rough | -4 |
| Bunker | -4 (hardest) |

### Course Grid
- 2D grid of cells covering the course
- Each cell has a position and lie type
- Configurable cell size (default: 10 yards)
- Unit conversion: 1 yard = 72 units
- Out-of-bounds returns Penalty Area lie

### Simple Course Generation
- Tee area: First 10 yards
- Green area: Last 20 yards around hole
- Fairway: Middle area (default)

---

## Equipment System

### Equipment Types
**Ball:**
- DistanceBonus: Added to shot distance
- SpinControl: Affects ball behavior
- Cost: Purchase price

**Glove:**
- AccuracyBonus: Added to club accuracy (capped at 1.0)
- Cost: Purchase price

**Shoes:**
- LiePenaltyReduction: Reduces negative lie modifiers
- Cost: Purchase price

### ProShop Inventory
**Balls:**
| Name | Distance | Spin | Cost |
|------|----------|------|------|
| Budget Ball | +0 | 0.3 | 20 |
| Standard Ball | +3 | 0.5 | 35 |
| Premium Ball | +5 | 0.7 | 50 |
| Pro V1 | +8 | 0.9 | 75 |

**Gloves:**
| Name | Accuracy | Cost |
|------|----------|------|
| Basic Glove | +0.02 | 25 |
| Leather Pro | +0.05 | 45 |
| Precision Grip | +0.08 | 65 |

**Shoes:**
| Name | Lie Reduction | Cost |
|------|---------------|------|
| Casual Spikes | 1 | 30 |
| All-Terrain Pro | 2 | 55 |
| Tour Edition | 3 | 80 |

---

## Currency & Rewards

### Starting Money
- New golfer starts with 100 money

### Hole Rewards
| Score vs Par | Reward |
|--------------|--------|
| Hole-in-one (1 stroke) | 100 |
| Eagle (-2) | 50 |
| Birdie (-1) | 25 |
| Par (0) | 10 |
| Bogey (+1) | 5 |
| Double bogey+ | 1 |

---

## Club System

### 14-Club Bag
Each club has:
- **Distance**: Maximum yards
- **Accuracy**: Probability modifier (0.0-1.0)
- **Forgiveness**: Error reduction

Standard club distances and usage are mapped to skills/abilities as shown above.

---

## Key Formulas Reference

**XP to next level:**
```
xpRequired = (currentLevel + 1) * 50
```

**Skill/Ability value:**
```
value = level
```

**Target number:**
```
target = skillValue + abilityValue + lieDifficulty - shoeReduction
```

**Modified club accuracy:**
```
accuracy = min(baseAccuracy + gloveBonus, 1.0)
```
