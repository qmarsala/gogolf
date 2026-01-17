package main

import "math/rand/v2"

type Dice struct {
	Sides int
}

func NewD6() Dice {
	return Dice{Sides: 6}
}

func (d Dice) Roll() int {
	return rand.IntN(d.Sides) + 1
}

func (d Dice) RollN(dieCount int) (total int, rolls []int) {
	for i := 0; i < dieCount; i++ {
		roll := d.Roll()
		rolls = append(rolls, roll)
		total += roll
	}
	return
}

type SkillCheckOutcome int

const (
	CriticalFailure SkillCheckOutcome = iota
	Bad
	Poor
	Marginal
	Good
	Excellent
	CriticalSuccess
)

func (o SkillCheckOutcome) String() string {
	return [...]string{
		"Critical Failure",
		"Bad",
		"Poor",
		"Marginal",
		"Good",
		"Excellent",
		"Critical Success",
	}[o]
}

type SkillCheckResult struct {
	Success    bool
	IsCritical bool
	RollTotal  int
	Rolls      []int
	Margin     int
	Outcome    SkillCheckOutcome
}

func determineOutcome(margin int, isCritical bool) SkillCheckOutcome {
	// Critical rolls override margin-based tiers
	if isCritical {
		if margin >= 0 {
			return CriticalSuccess
		}
		return CriticalFailure
	}

	// Margin-based tier determination
	switch {
	case margin >= 7:
		return CriticalSuccess
	case margin >= 4:
		return Excellent
	case margin >= 1:
		return Good
	case margin == 0:
		return Marginal
	case margin >= -3:
		return Poor
	case margin >= -6:
		return Bad
	default:
		return CriticalFailure
	}
}
