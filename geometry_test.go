package gogolf

import (
	"math"
	"testing"
)

func TestVector_Subtract(t *testing.T) {
	v1 := Vector{X: 5, Y: 10}
	v2 := Vector{X: 3, Y: 4}

	result := v1.Subtract(v2)

	if result.X != 2 || result.Y != 6 {
		t.Errorf("Vector.Subtract() = %v, want {2, 6}", result)
	}
}

func TestVector_Magnitude(t *testing.T) {
	tests := []struct {
		v        Vector
		expected float64
	}{
		{Vector{X: 3, Y: 4}, 5},
		{Vector{X: 0, Y: 0}, 0},
		{Vector{X: 1, Y: 0}, 1},
		{Vector{X: 0, Y: 1}, 1},
	}

	for _, tt := range tests {
		got := tt.v.Magnitude()
		if got != tt.expected {
			t.Errorf("Vector%v.Magnitude() = %v, want %v", tt.v, got, tt.expected)
		}
	}
}

func TestVector_Normalize(t *testing.T) {
	v := Vector{X: 3, Y: 4}
	result := v.Normalize()

	if result.X != 0.6 || result.Y != 0.8 {
		t.Errorf("Vector.Normalize() = %v, want {0.6, 0.8}", result)
	}

	mag := result.Magnitude()
	if mag != 1 {
		t.Errorf("Normalized vector magnitude = %v, want 1", mag)
	}
}

func TestVector_Dot(t *testing.T) {
	tests := []struct {
		v1, v2   Vector
		expected float64
	}{
		{Vector{X: 1, Y: 0}, Vector{X: 0, Y: 1}, 0},
		{Vector{X: 1, Y: 2}, Vector{X: 3, Y: 4}, 11},
		{Vector{X: 2, Y: 3}, Vector{X: 2, Y: 3}, 13},
	}

	for _, tt := range tests {
		got := tt.v1.Dot(tt.v2)
		if got != tt.expected {
			t.Errorf("%v.Dot(%v) = %v, want %v", tt.v1, tt.v2, got, tt.expected)
		}
	}
}

func TestVector_AngleBetween(t *testing.T) {
	tests := []struct {
		v1, v2   Vector
		expected float64
	}{
		{Vector{X: 1, Y: 0}, Vector{X: 0, Y: 1}, math.Pi / 2},
		{Vector{X: 1, Y: 0}, Vector{X: 1, Y: 0}, 0},
		{Vector{X: 1, Y: 0}, Vector{X: -1, Y: 0}, math.Pi},
		{Vector{X: 0, Y: 0}, Vector{X: 1, Y: 0}, 0},
	}

	for _, tt := range tests {
		got := tt.v1.AngleBetween(tt.v2)
		if math.Abs(got-tt.expected) > 0.0001 {
			t.Errorf("%v.AngleBetween(%v) = %v, want %v", tt.v1, tt.v2, got, tt.expected)
		}
	}
}

func TestVector_Rotate(t *testing.T) {
	v := Vector{X: 1, Y: 0}

	result := v.Rotate(90)
	if math.Abs(result.X) > 0.0001 || math.Abs(result.Y+1) > 0.0001 {
		t.Errorf("Vector.Rotate(90) = %v, want {0, -1}", result)
	}

	result = v.Rotate(180)
	if math.Abs(result.X+1) > 0.0001 || math.Abs(result.Y) > 0.0001 {
		t.Errorf("Vector.Rotate(180) = %v, want {-1, 0}", result)
	}
}

func TestPoint_Move(t *testing.T) {
	p := Point{X: 0, Y: 0}
	direction := Vector{X: 1, Y: 0}

	result := p.Move(direction, 10)

	if result.X != 10 || result.Y != 0 {
		t.Errorf("Point.Move() = %v, want {10, 0}", result)
	}
}

func TestPoint_Vector(t *testing.T) {
	p := Point{X: 5, Y: 10}

	result := p.Vector()

	if result.X != 5 || result.Y != 10 {
		t.Errorf("Point.Vector() = %v, want {5, 10}", result)
	}
}

func TestPoint_Direction(t *testing.T) {
	p1 := Point{X: 0, Y: 0}
	p2 := Point{X: 3, Y: 4}

	result := p1.Direction(p2)

	if result.X != 3 || result.Y != 4 {
		t.Errorf("Point.Direction() = %v, want {3, 4}", result)
	}
}

func TestPoint_Distance(t *testing.T) {
	tests := []struct {
		p1, p2   Point
		expected Unit
	}{
		{Point{X: 0, Y: 0}, Point{X: 3, Y: 4}, 5},
		{Point{X: 0, Y: 0}, Point{X: 0, Y: 0}, 0},
		{Point{X: 1, Y: 1}, Point{X: 4, Y: 5}, 5},
	}

	for _, tt := range tests {
		got := tt.p1.Distance(tt.p2)
		if got != tt.expected {
			t.Errorf("%v.Distance(%v) = %v, want %v", tt.p1, tt.p2, got, tt.expected)
		}
	}
}
