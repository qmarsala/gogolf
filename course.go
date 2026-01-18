package gogolf

import (
	"fmt"
)

// holes need a list of good locations to aim at
// sometimes it is best not to aim at the hole.
// and we want to implement 'ob' and 'trees' to force you away form the hole sometimes
// yet inputs like 'aim left 20 degrees' could be clunky as the only option
// having options like: center fairway, center green, hole, etc as quick aim options.
// also, being able to adjust distance aim to gain finer control
// if I am aiming at a point >= club distance, then full power is full power
// but it might be helpful to be able to aim at a point < club distance, and scale power from there.
// ex: aim pw at 100 yrds (making full power 100 instead of 140)

type Hole struct {
	Number       int
	Par          int
	Distance     Yard
	Boundary     Size
	HoleLocation Point
	Grid         *CourseGrid // Course grid for lie detection
}

func (h Hole) DetectHoleOut(b GolfBall, bPath Vector) bool {
	// eventually when we have 'carry' and 'roll' paths, we will need to make sure
	// it was the roll path, or the carries endpoint that hits the hole
	directHit, distanceFromHole := b.CheckForCollision(bPath, h.HoleLocation)
	hitAndStoppedInHole := directHit && b.Location.Distance(h.HoleLocation) <= Foot(16).Units()
	closeEnough := distanceFromHole <= Unit(2) && b.Location.Distance(h.HoleLocation) <= Yard(1).Units()
	return hitAndStoppedInHole || closeEnough

}

func (h Hole) DetectTapIn(b GolfBall) bool {
	return b.Location.Distance(h.HoleLocation) <= Foot(4).Units()
}

func NewHole(number int, par int, holeLocation Point, boundary Size) *Hole {
	return &Hole{
		Number:       number,
		Par:          par,
		Distance:     holeLocation.Distance(Point{0, 0}).Yards(),
		Boundary:     boundary,
		HoleLocation: holeLocation,
		Grid:         nil, // No grid by default for backward compatibility
	}
}

// NewHoleWithGrid creates a new hole with an initialized course grid
func NewHoleWithGrid(number int, par int, holeLocation Point, boundary Size, width, length, cellSize Yard) *Hole {
	grid := NewCourseGrid(width, length, cellSize)
	return &Hole{
		Number:       number,
		Par:          par,
		Distance:     holeLocation.Distance(Point{0, 0}).Yards(),
		Boundary:     boundary,
		HoleLocation: holeLocation,
		Grid:         &grid,
	}
}

// GetLieAtPosition returns the lie type at a position on the hole
// Returns Fairway if no grid is configured
func (h Hole) GetLieAtPosition(pos Point) LieType {
	if h.Grid == nil {
		return Fairway // Default to fairway if no grid
	}
	return h.Grid.GetLieAtPosition(pos)
}

func (h Hole) String() string {
	return fmt.Sprintf("Hole: %d Par: %d\nDistance: %f yards",
		h.Number, h.Par, h.Distance)
}

type Course struct {
	Holes []Hole
}

func (c Course) Par() (par int) {
	for _, v := range c.Holes {
		par += v.Par
	}
	return
}

func (c Course) ParUpToHole(holeNumber int) (par int) {
	for _, v := range c.Holes {
		if v.Number <= holeNumber {
			par += v.Par
		}
	}
	return
}

// GenerateSimpleCourse creates a course with grid-based lie detection
// This is a simplified version that sets up basic fairway/rough/green patterns
func GenerateSimpleCourse(holeCount int) (Course, ScoreCard) {
	holes := []Hole{}

	for i := 0; i < holeCount; i++ {
		// Simple hole setup: par 4, ~300 yards
		distance := Yard(300)
		width := Yard(50)
		par := 4

		// Hole location in the green area (last 10 yards)
		// Grid is 300 yards = 2160 units, so use 2110 (about 293 yards)
		holeLocation := Point{X: int(width.Units() / 2), Y: int(distance.Units()) - int(Yard(10).Units())}

		// Create hole with grid (50 yards wide, distance long, 10-yard cells)
		hole := NewHoleWithGrid(i+1, par, holeLocation, Size{}, width, distance, Yard(10))

		// Set up basic lies
		// Tee area (first 10 yards)
		teeAreaSize := Yard(10)
		for y := 0; y < int(teeAreaSize.Units()); y += int(Yard(10).Units()) {
			for x := 0; x < int(width.Units()); x += int(Yard(10).Units()) {
				hole.Grid.SetLieAtPosition(Point{X: x, Y: y}, Tee)
			}
		}

		// Green area (last 20 yards around hole)
		greenStartY := int(distance.Units()) - int(Yard(20).Units())
		for y := greenStartY; y < int(distance.Units()); y += int(Yard(10).Units()) {
			for x := 0; x < int(width.Units()); x += int(Yard(10).Units()) {
				hole.Grid.SetLieAtPosition(Point{X: x, Y: y}, Green)
			}
		}

		// Fairway is the default (already set in NewCourseGrid)
		// Middle area remains as fairway

		holes = append(holes, *hole)
	}

	course := Course{Holes: holes}
	scoreCard := ScoreCard{
		Course: course,
		Scores: map[int]int{},
	}
	return course, scoreCard
}
