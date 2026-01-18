package gogolf

import "math"

type Yard float32
type Inch float32
type Foot float32
type Unit float32

func (y Yard) Inches() Inch {
	return y.Feet().Inches()
}

func (y Yard) Units() Unit {
	return y.Feet().Units()
}

func (y Yard) Feet() Foot {
	return Foot(y * 3)
}

func (f Foot) Inches() Inch {
	return Inch(f * 12)
}

func (f Foot) Units() Unit {
	return f.Inches().Units()
}

func (f Foot) Yards() Yard {
	return Yard(f / 3)
}

func (i Inch) Units() Unit {
	return Unit(i / 5)
}

func (i Inch) Feet() Foot {
	return Foot(i / 12)
}

func (i Inch) Yards() Yard {
	return Yard(i.Feet().Yards())
}

func (u Unit) Feet() Foot {
	return Foot(u.Inches() / 12)
}

func (u Unit) Inches() Inch {
	return Inch(u * 5)
}

func (u Unit) Yards() Yard {
	return Yard(u.Feet() / 3)
}

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

type LieType int

const (
	Tee LieType = iota
	Fairway
	FirstCut
	Rough
	DeepRough
	Bunker
	Green
	PenaltyArea
)

func (l LieType) String() string {
	return [...]string{
		"Tee",
		"Fairway",
		"First Cut",
		"Rough",
		"Deep Rough",
		"Bunker",
		"Green",
		"Penalty Area",
	}[l]
}

func (l LieType) DifficultyModifier() int {
	switch l {
	case Tee:
		return 2
	case Fairway:
		return 0
	case FirstCut:
		return -1
	case Rough:
		return -2
	case DeepRough:
		return -4
	case Bunker:
		return -4
	case Green:
		return 1
	case PenaltyArea:
		return 0
	default:
		return 0
	}
}

type SkillCheckOutcome int

const (
	CriticalFailure SkillCheckOutcome = iota
	Bad
	Poor
	Marginal
	Good
	Excellent
	CriticalSuccess
)

func (o SkillCheckOutcome) String() string {
	return [...]string{
		"Critical Failure",
		"Bad",
		"Poor",
		"Marginal",
		"Good",
		"Excellent",
		"Critical Success",
	}[o]
}

type SkillCheckResult struct {
	Success    bool
	IsCritical bool
	RollTotal  int
	Rolls      []int
	Margin     int
	Outcome    SkillCheckOutcome
}

type ShotShape int

const (
	Straight ShotShape = iota
	Draw
	Fade
	Hook
	Slice
)

func (s ShotShape) String() string {
	return [...]string{
		"Straight",
		"Draw",
		"Fade",
		"Hook",
		"Slice",
	}[s]
}

func (s ShotShape) DifficultyModifier() int {
	switch s {
	case Straight:
		return -2
	case Draw:
		return 1
	case Fade:
		return 1
	case Hook:
		return -1
	case Slice:
		return -1
	default:
		return 0
	}
}
