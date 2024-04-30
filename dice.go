package main

import "math/rand/v2"

type Dice struct {
	Sides int
}

func NewD6() Dice {
	return Dice{Sides: 6}
}

func (d Dice) roll() int {
	return rand.IntN(d.Sides) + 1
}

func (d Dice) rollN(dieCount int) (total int, rolls []int) {
	for i := 0; i < dieCount; i++ {
		roll := d.roll()
		rolls = append(rolls, roll)
		total += roll
	}
	return
}

type SkillCheckResult struct {
	Success    bool
	IsCritical bool
	RollTotal  int
	Rolls      []int
	Margin     int
}

func (d Dice) SkillCheck(targetNumber int) SkillCheckResult {
	total, rolls := d.rollN(3)
	margin := targetNumber - total
	return SkillCheckResult{
		Success:    margin >= 0,
		IsCritical: rolls[0] == rolls[1] && rolls[0] == rolls[2],
		RollTotal:  total,
		Rolls:      rolls,
		Margin:     margin,
	}
}
