package main

import (
	"fmt"
	"math"
)

// we probably want some sort of grid, probably hex based
// to store the holes long term

type Point struct {
	X int
	Y int
}

type Size struct {
	Width  float32
	Length float32
}

type Yards float32
type Feet float32

func toYards(feet float32) Yards {
	return Yards(feet / 3)
}

func toFeet(yards float32) Feet {
	return Feet(yards * 3)
}

type Green struct {
	Size         Size
	Location     Point
	HoleLocation Point
}

type Fairway struct {
	Size     Size
	Location Point
}

type TeeBox struct {
	Size     Size
	Location Point
}

type Hole struct {
	Number   int
	Par      int
	Distance float32
	TeeBox   TeeBox
	Fairway  Fairway
	Green    Green
}

func NewHole(number int, par int, teeBox TeeBox, fairway Fairway, green Green) *Hole {
	distanceInFeet := math.Sqrt(math.Pow(float64(green.HoleLocation.X), 2) + math.Pow(float64(green.HoleLocation.Y), 2))
	hole := &Hole{
		Number:   number,
		Par:      par,
		Distance: float32(distanceInFeet),
		TeeBox:   teeBox,
		Fairway:  fairway,
		Green:    green,
	}
	return hole
}

func (h Hole) DistanceToHole(b GolfBall) Yards {
	xs := math.Pow(float64(b.Location.X-h.Green.HoleLocation.X), 2)
	ys := math.Pow(float64(b.Location.Y-h.Green.HoleLocation.Y), 2)
	distanceInFeet := math.Sqrt(math.Abs(xs + ys))
	return toYards(float32(distanceInFeet))
}

func (h Hole) String() string {
	return fmt.Sprintf("Hole: %d Par: %d\nDistance: %f yards",
		h.Number, h.Par, toYards(h.Distance))
}

type Course struct {
	Holes []Hole
	Par   int
}
