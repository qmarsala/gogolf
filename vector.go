package main

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
