package main

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

func (pointA Point) Distance(pointB Point) Unit {
	xs := math.Pow(float64(pointB.X-pointA.X), 2)
	ys := math.Pow(float64(pointB.Y-pointA.Y), 2)
	return Unit(math.Sqrt(math.Abs(xs + ys)))
}
