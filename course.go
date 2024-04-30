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

// for 'hole out' logic, we should scan the path of the ball and the hole location
// for collision. Then if it was not traveling to far past, it could be considered in
// todo: collision detection
// for now the ball's receive hit could return a vector representing the path taken
// but it will eventually need to be an actually line that may curve
func (h Hole) CheckForBall(b GolfBall) bool {
	// this allows a ball to stop short and be counted, we will
	// want to do collision checks on a path to be more accurate
	grace := int(Foot(2).Units())
	holedOut := (b.Location.Y >= (h.HoleLocation.Y-grace) && b.Location.Y <= (h.HoleLocation.Y+grace)) &&
		(b.Location.X >= (h.HoleLocation.X-grace) && b.Location.X <= (h.HoleLocation.X+grace))
	return holedOut
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
