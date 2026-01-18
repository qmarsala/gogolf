package gogolf

import "math"

type Vector struct {
	X, Y float64
}

func (v1 Vector) Subtract(v2 Vector) Vector {
	return Vector{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector) Normalize() Vector {
	mag := v.Magnitude()
	return Vector{X: v.X / mag, Y: v.Y / mag}
}

func (v1 Vector) Dot(v2 Vector) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v1 Vector) AngleBetween(v2 Vector) float64 {
	dot := v1.Dot(v2)
	magV1 := v1.Magnitude()
	magV2 := v2.Magnitude()
	if magV1 == 0 || magV2 == 0 {
		return 0
	}
	cosine := dot / (magV1 * magV2)
	if cosine < -1 {
		cosine = -1
	} else if cosine > 1 {
		cosine = 1
	}
	return math.Acos(cosine)
}

func (v Vector) Rotate(angleDegrees float64) Vector {
	angleRadians := angleDegrees * math.Pi / 180
	cosTheta := math.Cos(angleRadians)
	sinTheta := math.Sin(angleRadians)
	return Vector{
		X: v.X*cosTheta + v.Y*sinTheta,
		Y: -v.X*sinTheta + v.Y*cosTheta,
	}
}

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

func (pointA Point) Distance(pointB Point) Unit {
	xs := math.Pow(float64(pointB.X-pointA.X), 2)
	ys := math.Pow(float64(pointB.Y-pointA.Y), 2)
	return Unit(math.Sqrt(math.Abs(xs + ys)))
}
