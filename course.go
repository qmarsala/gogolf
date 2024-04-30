package main

import (
	"fmt"
)

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
	holedOut := (b.Location.X == h.HoleLocation.X &&
		b.Location.Y >= (h.HoleLocation.Y-grace) && b.Location.Y <= (h.HoleLocation.Y+grace)) ||
		(b.Location.Y == h.HoleLocation.Y &&
			b.Location.X >= (h.HoleLocation.X-grace) && b.Location.X <= (h.HoleLocation.X+grace))
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
