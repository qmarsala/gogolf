package main

import (
	"gogolf/dice"
	"math"
)

type Club struct {
	Name     string
	Distance Yard
	// how to model this:
	// something like a 'putter' will always go straight on success
	// and deviate a little in a range on failure.  while a 'driver' has room for
	// miss on success, and a greater margin on failure
	Accuracy float32 //? this could represent a percentage. 1 meaning no error, 0 meaning full random
	// 80% accuracy would mean success goes somewhere in a 20 degree, 10 degree angle left or right of target
	// success/failure margin could add to this
	// ie, succeed by 10 would make the accuracy 10 degrees, 5 degree angle left or right of target
	// and failed by 10 would make it 2/3 accuracy - margin = 57 degrees, 28 degrees left or right of target

	// the failure mod of 2/3 could maybe be determined by some 'forgiveness' factor of the club?
	// this represents how much of the accuracy you keep on a failure
	// 1 would mean that a miss hit's accuracy is only affected by the margin
	// 0 would mean that it is basically a random shot
	Forgiveness float32
}

func (c Club) AccuracyDegrees() float32 {
	return ((1 - c.Accuracy) * 100) / 2
}

func DefaultClubs() (clubs []Club) {
	driver := Club{Name: "Driver", Distance: 280, Accuracy: .75, Forgiveness: .8}
	threeWood := Club{Name: "3 Wood", Distance: 250, Accuracy: .8, Forgiveness: .8}
	fiveWood := Club{Name: "5 Wood", Distance: 235, Accuracy: .8, Forgiveness: .8}
	fourIron := Club{Name: "4 Iron", Distance: 215, Accuracy: .85, Forgiveness: .8}
	fiveIron := Club{Name: "5 Iron", Distance: 200, Accuracy: .85, Forgiveness: .8}
	sixIron := Club{Name: "6 Iron", Distance: 190, Accuracy: .85, Forgiveness: .8}
	sevenIron := Club{Name: "7 Iron", Distance: 180, Accuracy: .9, Forgiveness: .8}
	eightIron := Club{Name: "8 Iron", Distance: 170, Accuracy: .9, Forgiveness: .8}
	nineIron := Club{Name: "8 Iron", Distance: 160, Accuracy: .9, Forgiveness: .8}
	pitchingWedge := Club{Name: "PW", Distance: 150, Accuracy: .95, Forgiveness: .8}
	gapWedge := Club{Name: "GW", Distance: 140, Accuracy: .95, Forgiveness: .8}
	sandWedge := Club{Name: "SW", Distance: 125, Accuracy: .95, Forgiveness: .8}
	lobWedge := Club{Name: "LW", Distance: 100, Accuracy: .95, Forgiveness: .8}
	putter := Club{Name: "Putter", Distance: 40, Accuracy: 1, Forgiveness: .95}
	clubs = []Club{driver, threeWood, fiveWood, fourIron, fiveIron, sixIron, sevenIron, eightIron, nineIron, pitchingWedge, gapWedge, sandWedge, lobWedge, putter}
	return
}

type Golfer struct {
	Name      string
	Target    Point
	Clubs     []Club
	Skills    map[string]Skill
	Abilities map[string]Ability
	Money     int // Currency for purchasing equipment
	Ball      *Ball  // Equipped golf ball
	Glove     *Glove // Equipped glove
	Shoes     *Shoes // Equipped shoes
}

// NewGolfer creates a new golfer with default skills and abilities at level 1
func NewGolfer(name string) Golfer {
	return Golfer{
		Name:  name,
		Clubs: DefaultClubs(),
		Skills: map[string]Skill{
			"Driver":       NewSkill("Driver"),
			"Woods":        NewSkill("Woods"),
			"Long Irons":   NewSkill("Long Irons"),
			"Mid Irons":    NewSkill("Mid Irons"),
			"Short Irons":  NewSkill("Short Irons"),
			"Wedges":       NewSkill("Wedges"),
			"Putter":       NewSkill("Putter"),
		},
		Abilities: map[string]Ability{
			"Strength": NewAbility("Strength"),
			"Control":  NewAbility("Control"),
			"Touch":    NewAbility("Touch"),
			"Mental":   NewAbility("Mental"),
		},
		Money: 100, // Starter money for basic equipment
	}
}

//idea: strategy pattern that could be provided by a 'caddie'
func (g Golfer) GetBestClub(distance Yard) Club {
	c := g.Clubs[0]
	closetDiff := float64(1000)
	for _, v := range g.Clubs {
		diff := math.Abs(float64(v.Distance) - float64(distance))
		if diff <= closetDiff && v.Distance >= distance {
			closetDiff = diff
			c = v
		}
	}
	return c
}

func (g Golfer) SkillCheck(d dice.Dice, targetNumber int) dice.SkillCheckResult {
	total, rolls := d.RollN(3)
	margin := targetNumber - total
	isCritical := rolls[0] == rolls[1] && rolls[0] == rolls[2]

	return dice.SkillCheckResult{
		Success:    margin >= 0,
		IsCritical: isCritical,
		RollTotal:  total,
		Rolls:      rolls,
		Margin:     margin,
		Outcome:    dice.DetermineOutcome(margin, isCritical),
	}
}

// GetSkillForClub returns the skill associated with a club
func (g Golfer) GetSkillForClub(club Club) Skill {
	switch club.Name {
	case "Driver":
		return g.Skills["Driver"]
	case "3 Wood", "5 Wood":
		return g.Skills["Woods"]
	case "4 Iron", "5 Iron":
		return g.Skills["Long Irons"]
	case "6 Iron", "7 Iron":
		return g.Skills["Mid Irons"]
	case "8 Iron", "9 Iron":
		return g.Skills["Short Irons"]
	case "PW", "GW", "SW", "LW":
		return g.Skills["Wedges"]
	case "Putter":
		return g.Skills["Putter"]
	default:
		return g.Skills["Driver"] // fallback
	}
}

// GetAbilityForClub returns the ability associated with a club
func (g Golfer) GetAbilityForClub(club Club) Ability {
	switch club.Name {
	case "Driver", "3 Wood", "5 Wood":
		return g.Abilities["Strength"]
	case "4 Iron", "5 Iron", "6 Iron", "7 Iron":
		return g.Abilities["Control"]
	case "8 Iron", "9 Iron", "PW", "GW", "SW", "LW":
		return g.Abilities["Touch"]
	case "Putter":
		return g.Abilities["Mental"]
	default:
		return g.Abilities["Strength"] // fallback
	}
}

// CalculateTargetNumber computes target number for skill check
// Formula: skillValue + abilityValue + difficultyModifier + equipmentBonuses
func (g Golfer) CalculateTargetNumber(club Club, difficulty int) int {
	skill := g.GetSkillForClub(club)
	ability := g.GetAbilityForClub(club)

	// Apply shoes bonus to reduce lie penalties
	adjustedDifficulty := difficulty + g.GetTotalLiePenaltyReduction()

	return skill.Value() + ability.Value() + adjustedDifficulty
}

// AwardExperience adds XP to both the skill and ability for a club
func (g *Golfer) AwardExperience(club Club, xp int) {
	// Get copies of the skill and ability
	skill := g.GetSkillForClub(club)
	ability := g.GetAbilityForClub(club)

	// Call AddExperience on the copies (using pointer receiver)
	// This modifies the copies in place
	(&skill).AddExperience(xp)
	(&ability).AddExperience(xp)

	// Update the maps with the modified copies
	g.Skills[skill.Name] = skill
	g.Abilities[ability.Name] = ability
}

// AddMoney increases the golfer's money by the specified amount
func (g *Golfer) AddMoney(amount int) {
	g.Money += amount
}

// SpendMoney decreases the golfer's money by the specified amount
// Returns true if the transaction succeeded, false if insufficient funds
func (g *Golfer) SpendMoney(amount int) bool {
	if g.Money < amount {
		return false // Insufficient funds
	}
	g.Money -= amount
	return true
}

// EquipBall equips a golf ball
func (g *Golfer) EquipBall(ball *Ball) {
	g.Ball = ball
}

// EquipGlove equips a glove
func (g *Golfer) EquipGlove(glove *Glove) {
	g.Glove = glove
}

// EquipShoes equips shoes
func (g *Golfer) EquipShoes(shoes *Shoes) {
	g.Shoes = shoes
}

// GetEquippedBall returns the currently equipped ball (may be nil)
func (g Golfer) GetEquippedBall() *Ball {
	return g.Ball
}

// GetTotalLiePenaltyReduction returns the total lie penalty reduction from equipment
// Shoes reduce the difficulty of bad lies (e.g., rough, bunkers)
func (g Golfer) GetTotalLiePenaltyReduction() int {
	if g.Shoes == nil {
		return 0
	}
	return g.Shoes.LiePenaltyReduction
}

// GetModifiedClub returns a club with equipment bonuses applied
// Ball adds distance, glove improves accuracy
func (g Golfer) GetModifiedClub(club Club) Club {
	modified := club

	// Apply ball distance bonus
	if g.Ball != nil {
		modified.Distance = Yard(float32(modified.Distance) + g.Ball.DistanceBonus)
	}

	// Apply glove accuracy bonus (capped at 1.0)
	if g.Glove != nil {
		modified.Accuracy += g.Glove.AccuracyBonus
		if modified.Accuracy > 1.0 {
			modified.Accuracy = 1.0
		}
	}

	return modified
}

// how do we want to do this?
// we can have different types of shots
// and use different clubs
// hit it softer or harder
// these types of decisions can affect the target number
// func (g *Golfer) Swing(club Club) {
// 	dice := NewD6()
// 	//todo: Depending on the shot we will use a different skill/attribute
// 	// we will also need to add mods such as lie boost/penalties
// 	// for now, just something
// 	targetNumber := g.Skills[club.DefaultSkill].Value() + g.Abilities[club.DefaultAbility].Value()
// 	result := dice.SkillCheck(targetNumber)
// 	if result.Success {
// 		g.GolfBall.ReceiveHit(Club{}, 1)
// 	} else {
// 		g.GolfBall.ReceiveHit(Club{}, .65)
// 	}
// }

//aiming
// - thinking that by default you will be aiming at the flag
// - to do this, we probably need to calculate the angle
// from the y axis to the hole location and you. so we can do things
// relative to your plane with the hole location.
// - aiming would then rotate this plane?
// - shots then fade/draw/hook/slice etc relative to this plane

// when we hit the ball
// based on shot shape, power, success, etc
// the ball will move to a location
// eventually, the path of the ball needs to be considered
// also need to consider lie (are we in the fairway? a bunker? the rough?)
// skill check success
// where do we end up?
// the failure could be a hook/straight if we were trying to draw
// a slice/straight if we were trying to fade
// the shape/straight ration could be like 80/20 or maybe based on your skill?
// and random if straight, plus maybe some other types of misses...
// golf is about playing your miss
// this would promote picking a shot shape, to know the miss.
