package main

import "fmt"

type Skills struct {
	Woods   int
	Irons   int
	Wedges  int
	Putting int
}

type Abilities struct {
	Power    int
	Accuracy int
	Control  int
}

type ShotShape int

func (ss ShotShape) String() string {
	switch ss {
	case -1:
		return "fade"
	case 1:
		return "draw"
	default:
		return "straight"
	}
}

type Shot struct {
	Club  Club
	Shape ShotShape
	Power int
}

func (s Shot) String() string {
	return fmt.Sprintf("club: %s\nshot shape: %s\npower: %d percent",
		s.Club, s.Shape, s.Power)
}

type Club struct {
	Name          string
	Id            int
	StockDistance int
	Skill         string
}

func (c Club) String() string {
	return c.Name
}

type GolfBag struct {
	Clubs []Club
}

type Golfer struct {
	Name      string
	Skills    Skills
	Abilities Abilities
}

func (g *Golfer) PlayShot(shot Shot) {
	fmt.Println(shot.String())
	targetNumber := g.Skills.Wedges + g.Abilities.Accuracy
	dice := Dice{
		Sides: 6,
	}
	result, rolls := dice.rollN(3)
	margin := targetNumber - result
	success := margin >= 0

	fmt.Println(success, targetNumber, result, margin, rolls)
	if success {
		// state changes somewhere
		// we are now shot.distance closer to the hole
	} else {
		// how bad did we fail?
		// critical fail? big margin fail? normal fail?
		// how did we fail? alignment? chunky? thin?
		// what is the consequence for the result?
	}
}

type Score struct {
	Hole    Hole
	Strokes int
}

type ScoreCard struct {
	Scores map[int]Score
}

type GolfBall struct {
	Lie            string
	DistanceToHole int
}
