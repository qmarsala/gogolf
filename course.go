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

func (h Hole) CheckForBall(b GolfBall) bool {
	// this allows a ball to stop short and be counted, we will
	// want to do collision checks on a path to be more accurate
	grace := int(Foot(2).Units())
	holedOut := (b.Location.X == h.HoleLocation.X &&
		b.Location.Y >= (h.HoleLocation.Y-grace) && b.Location.Y <= (h.HoleLocation.Y+grace)) ||
		(b.Location.Y == h.HoleLocation.Y &&
			b.Location.X >= (h.HoleLocation.X-grace) && b.Location.X <= (h.HoleLocation.X+grace))
	fmt.Printf("holed out: %+v\n", holedOut)
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
	Par   int
}
