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
	Unit   string
}

func (s Size) InYards() Size {
	if s.Unit == "yards" {
		return s
	}
	return Size{
		Unit:   "yards",
		Width:  toYards(s.Width),
		Length: toYards(s.Length),
	}
}

func (s Size) InFeet() Size {
	if s.Unit == "feet" {
		return s
	}
	return Size{
		Unit:   "feet",
		Width:  toFeet(s.Width),
		Length: toFeet(s.Length),
	}
}

func toYards(feet float32) float32 {
	return feet / 3
}

func toFeet(yards float32) float32 {
	return yards * 3
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
		Distance: float32(distanceInFeet / 3),
		TeeBox:   teeBox,
		Fairway:  fairway,
		Green:    green,
	}
	return hole
}

func (h Hole) String() string {
	return fmt.Sprintf("Hole: %d Par: %d\nDistance: %f yards",
		h.Number, h.Par, h.Distance)
}

type Course struct {
	Holes []Hole
	Par   int
}
