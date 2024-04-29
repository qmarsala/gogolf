package main

type Point struct {
	X int
	Y int
}

type Size struct {
	Width  Unit
	Length Unit
}

func (p Point) Move(direction Vector, distance float64) Point {
	normalizedDirection := direction.Normalize()
	return Point{
		X: int(float64(p.X) + normalizedDirection.X*distance),
		Y: int(float64(p.Y) + normalizedDirection.Y*distance),
	}
}

func (p Point) Vector() Vector {
	return Vector{float64(p.X), float64(p.Y)}
}

func (p1 Point) Direction(p2 Point) Vector {
	return p2.Vector().Subtract(p1.Vector())
}
