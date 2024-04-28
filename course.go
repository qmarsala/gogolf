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
