# GoGolf Implementation Plan

## Completed Work

Phases 1-4 are complete. See [CORE_MECHANICS.md](CORE_MECHANICS.md) for documentation of implemented systems:
- Skills & Abilities System (PR #2)
- Main Game Loop Integration (PR #3)
- Course Grid & Lie System (PR #4)
- Equipment System & ProShop (PR #5)

**Current test count:** 184 tests passing

---

## Phase 5: Save/Load System
**Priority: HIGH** - Required for RPG progression persistence

### Requirements
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

### Design Considerations
- Save file format should be human-readable (JSON)
- Handle backward compatibility as game evolves
- Graceful error handling for corrupted saves
- Clear save/load confirmation messages

---

## Phase 6: ProShop UI & Browsing
**Priority: MEDIUM** - Enhance equipment shopping experience

### Requirements
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

### Example UI Flow
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

## Sprint Roadmap

### Completed
- Sprint 1: Foundation (Skills System)
- Sprint 2: Game Loop Integration
- Sprint 3: Course Depth (Lie System)
- Sprint 4: Progression (Equipment)

### Next Up
**Sprint 5: Persistence**
- Save/load system
- Auto-save functionality
- Multiple save slots

**Sprint 6: Polish**
- ProShop browsing UI
- Advanced course features
- Additional UI improvements

---

## Key Design Principles

### TDD Requirements (from CLAUDE.md)
- Write tests first (expected input/output)
- Verify test fails
- Commit failing tests
- Implement minimal code to pass
- Commit working code

### RPG Core Loop
```
Action → Skill Check → XP Gain → Level Up → Better Stats → Harder Courses
```

### Golf Realism
- Lie matters (rough is harder than fairway)
- Club selection matters (right tool for the job)
- Course management (risk/reward decisions)
- Equipment helps but skill is primary
