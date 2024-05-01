package main

import (
	"fmt"
)

// holes need a list of good locations to aim at
// sometimes it is best not to aim at the hole.
// and we want to implement 'ob' and 'trees' to force you away form the hole sometimes
// yet inputs like 'aim left 20 degrees' could be clunky as the only option
// having options like: center fairway, center green, hole, etc as quick aim options.
// also, being able to adjust distance aim to gain finer control
// if I am aiming at a point >= club distance, then full power is full power
// but it might be helpful to be able to aim at a point < club distance, and scale power from there.
// ex: aim pw at 100 yrds (making full power 100 instead of 140)

type Hole struct {
	Number       int
	Par          int
	Distance     Yard
	Boundary     Size
	HoleLocation Point
}

func (h Hole) DetectHoleOut(b GolfBall, bPath Vector) bool {
	// eventually when we have 'carry' and 'roll' paths, we will need to make sure
	// it was the roll path, or the carries endpoint that hits the hole
	directHit, distanceFromHole := b.CheckForCollision(bPath, h.HoleLocation)
	hitAndStoppedInHole := directHit && b.Location.Distance(h.HoleLocation) <= Yard(10).Units()
	closeEnough := distanceFromHole <= Unit(2) && b.Location.Distance(h.HoleLocation) <= Yard(1).Units()
	return hitAndStoppedInHole || closeEnough

}

func (h Hole) DetectTapIn(b GolfBall) bool {
	return b.Location.Distance(h.HoleLocation) <= Yard(2).Units()
}

func NewHole(number int, par int, holeLocation Point, boundary Size) *Hole {
	return &Hole{
		Number:       number,
		Par:          par,
		Distance:     holeLocation.Distance(Point{0, 0}).Yards(),
		Boundary:     boundary,
		HoleLocation: holeLocation,
	}
}

func (h Hole) String() string {
	return fmt.Sprintf("Hole: %d Par: %d\nDistance: %f yards",
		h.Number, h.Par, h.Distance)
}

type Course struct {
	Holes []Hole
}

func (c Course) Par() (par int) {
	for _, v := range c.Holes {
		par += v.Par
	}
	return
}

func (c Course) ParUpToHole(holeNumber int) (par int) {
	for _, v := range c.Holes {
		if v.Number <= holeNumber {
			par += v.Par
		}
	}
	return
}
