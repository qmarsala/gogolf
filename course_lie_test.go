package gogolf

import "testing"

// Test Hole includes CourseGrid
func TestHole_WithGrid(t *testing.T) {
	hole := NewHoleWithGrid(1, 4, Point{X: 360, Y: 720}, Size{}, Yard(100), Yard(100), Yard(10))

	if hole.Grid == nil {
		t.Error("Hole grid is nil, expected initialized grid")
	}

	if hole.Number != 1 {
		t.Errorf("Hole number = %d, want 1", hole.Number)
	}

	if hole.Par != 4 {
		t.Errorf("Hole par = %d, want 4", hole.Par)
	}
}

// Test Hole.GetLieAtPosition delegates to grid
func TestHole_GetLieAtPosition(t *testing.T) {
	// Grid is 100 yards = 720 units, so valid positions are 0-719
	holeLocation := Point{X: 360, Y: 680}
	hole := NewHoleWithGrid(1, 4, holeLocation, Size{}, Yard(100), Yard(100), Yard(10))

	// Set up different lies on the course
	hole.Grid.SetLieAtPosition(Point{X: 0, Y: 0}, Tee)
	hole.Grid.SetLieAtPosition(Point{X: 180, Y: 360}, Fairway)
	hole.Grid.SetLieAtPosition(Point{X: 360, Y: 680}, Green) // At hole location

	tests := []struct {
		name     string
		position Point
		expected LieType
	}{
		{"Tee box", Point{X: 0, Y: 0}, Tee},
		{"Fairway", Point{X: 180, Y: 360}, Fairway},
		{"Green", Point{X: 360, Y: 680}, Green},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lie := hole.GetLieAtPosition(tt.position)
			if lie != tt.expected {
				t.Errorf("GetLieAtPosition(%+v) = %v, want %v",
					tt.position, lie, tt.expected)
			}
		})
	}
}

// Test GolfBall can detect its lie
func TestGolfBall_GetLie(t *testing.T) {
	holeLocation := Point{X: 360, Y: 680}
	hole := NewHoleWithGrid(1, 4, holeLocation, Size{}, Yard(100), Yard(100), Yard(10))

	// Set up lies
	hole.Grid.SetLieAtPosition(Point{X: 0, Y: 0}, Tee)
	hole.Grid.SetLieAtPosition(Point{X: 180, Y: 360}, Rough)
	hole.Grid.SetLieAtPosition(Point{X: 360, Y: 680}, Green)

	tests := []struct {
		name     string
		ballPos  Point
		expected LieType
	}{
		{"Ball on tee", Point{X: 0, Y: 0}, Tee},
		{"Ball in rough", Point{X: 180, Y: 360}, Rough},
		{"Ball on green", Point{X: 360, Y: 680}, Green},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ball := GolfBall{Location: tt.ballPos}
			lie := ball.GetLie(hole)

			if lie != tt.expected {
				t.Errorf("Ball at %+v got lie %v, want %v",
					tt.ballPos, lie, tt.expected)
			}
		})
	}
}

func TestGenerateSimpleCourse(t *testing.T) {
	course := GenerateSimpleCourse(3)

	if len(course.Holes) != 3 {
		t.Errorf("Course has %d holes, want 3", len(course.Holes))
	}

	for i, hole := range course.Holes {
		if hole.Grid == nil {
			t.Errorf("Hole %d has nil grid", i+1)
		}

		// Verify tee and green are set
		teeLie := hole.GetLieAtPosition(Point{X: 0, Y: 0})
		if teeLie != Tee {
			t.Errorf("Hole %d tee area lie = %v, want Tee", i+1, teeLie)
		}

		greenLie := hole.GetLieAtPosition(hole.HoleLocation)
		if greenLie != Green {
			t.Errorf("Hole %d green lie = %v, want Green", i+1, greenLie)
		}
	}
}
