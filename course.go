package main

import "fmt"

// we probably want some sort of grid, probably hex based
// to store the holes long term

// for now, for simplicity
// fairway starts at the tee box and goes to the green
// and the only hazard is rough to the sides
type Hole struct {
	Number        int
	Par           int
	Distance      int
	FairwayLength int
	FairwayWidth  int
	GreenLength   int
	GreenWidth    int
	HoleLocation  int
}

func NewHole(number int, par int, distance int) *Hole {
	fairwayLength := .8 * float32(distance)
	greenLength := .2 * float32(distance)
	hole := &Hole{
		Number:        number,
		Par:           par,
		Distance:      distance,
		FairwayLength: int(fairwayLength),
		FairwayWidth:  50,
		GreenLength:   int(greenLength),
		GreenWidth:    50,
		HoleLocation:  1,
	}
	return hole
}

func (h Hole) String() string {
	return fmt.Sprintf("Hole: %d Par: %d\nDistance: %d yards",
		h.Number, h.Par, h.Distance)
}

type Course struct {
	Holes         []Hole
	Par           int
	TotalDistance int
}
