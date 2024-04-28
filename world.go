package main

import (
	"math"
)

type Point struct {
	X int
	Y int
}

type Size struct {
	Width  Unit
	Length Unit
}

type Vector struct {
	X, Y float64
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector) Normalize() Vector {
	mag := v.Magnitude()
	return Vector{X: v.X / mag, Y: v.Y / mag}
}

func MovePoint(p Point, direction Vector, distance float64) Point {
	unit := direction.Normalize()
	return Point{
		X: int(float64(p.X) + unit.X*distance),
		Y: int(float64(p.Y) + unit.Y*distance),
	}
}
