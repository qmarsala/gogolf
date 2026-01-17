package main

import "testing"

// Test GridCell creation and properties
func TestGridCell_Creation(t *testing.T) {
	cell := GridCell{
		Position: Point{X: 100, Y: 200},
		Lie:      Fairway,
	}

	if cell.Position.X != 100 || cell.Position.Y != 200 {
		t.Errorf("GridCell position = %+v, want {100, 200}", cell.Position)
	}

	if cell.Lie != Fairway {
		t.Errorf("GridCell lie = %v, want Fairway", cell.Lie)
	}
}

// Test CourseGrid creation with proper dimensions
func TestCourseGrid_Creation(t *testing.T) {
	// Create a simple 3x3 grid (30x30 yards with 10-yard cells)
	width := Yard(30)
	length := Yard(30)
	cellSize := Yard(10)

	grid := NewCourseGrid(width, length, cellSize)

	if len(grid.Cells) != 3 {
		t.Errorf("Grid rows = %d, want 3", len(grid.Cells))
	}

	if len(grid.Cells[0]) != 3 {
		t.Errorf("Grid columns = %d, want 3", len(grid.Cells[0]))
	}

	if grid.CellSize != cellSize {
		t.Errorf("Grid cell size = %v, want %v", grid.CellSize, cellSize)
	}

	if grid.Width != width {
		t.Errorf("Grid width = %v, want %v", grid.Width, width)
	}

	if grid.Length != length {
		t.Errorf("Grid length = %v, want %v", grid.Length, length)
	}
}

// Test GetLieAtPosition returns correct lie for a position
func TestCourseGrid_GetLieAtPosition(t *testing.T) {
	// Create a simple grid: 30 yards x 30 yards with 10-yard cells = 3x3 grid
	// Each cell is 10 yards = 72 units
	grid := NewCourseGrid(Yard(30), Yard(30), Yard(10))

	// Use SetLieAtPosition to set lies at specific positions (in units)
	// Cell [0][0]: X=0-72, Y=0-72
	// Cell [1][1]: X=72-144, Y=72-144
	// Cell [2][2]: X=144-216, Y=144-216
	grid.SetLieAtPosition(Point{X: 36, Y: 36}, Tee)       // Cell [0][0] - center
	grid.SetLieAtPosition(Point{X: 108, Y: 108}, Fairway) // Cell [1][1] - center
	grid.SetLieAtPosition(Point{X: 180, Y: 180}, Green)   // Cell [2][2] - center
	grid.SetLieAtPosition(Point{X: 180, Y: 36}, Rough)    // Cell [0][2] - center

	tests := []struct {
		name     string
		position Point
		expected LieType
	}{
		{"Tee position", Point{X: 36, Y: 36}, Tee},
		{"Fairway position", Point{X: 108, Y: 108}, Fairway},
		{"Green position", Point{X: 180, Y: 180}, Green},
		{"Rough position", Point{X: 180, Y: 36}, Rough},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lie := grid.GetLieAtPosition(tt.position)
			if lie != tt.expected {
				t.Errorf("GetLieAtPosition(%+v) = %v, want %v",
					tt.position, lie, tt.expected)
			}
		})
	}
}

// Test GetLieAtPosition handles out-of-bounds positions
func TestCourseGrid_GetLieAtPosition_OutOfBounds(t *testing.T) {
	grid := NewCourseGrid(Yard(30), Yard(30), Yard(10))

	tests := []struct {
		name     string
		position Point
	}{
		{"Negative X", Point{X: -10, Y: 15}},
		{"Negative Y", Point{X: 15, Y: -10}},
		{"X too large", Point{X: 1000, Y: 15}},
		{"Y too large", Point{X: 15, Y: 1000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Out of bounds should default to Penalty Area
			lie := grid.GetLieAtPosition(tt.position)
			if lie != PenaltyArea {
				t.Errorf("GetLieAtPosition(%+v) = %v, want PenaltyArea (out of bounds)",
					tt.position, lie)
			}
		})
	}
}

// Test SetLieAtPosition updates grid cells
func TestCourseGrid_SetLieAtPosition(t *testing.T) {
	grid := NewCourseGrid(Yard(30), Yard(30), Yard(10))

	// Initially should be default (Fairway by default in NewCourseGrid)
	initialLie := grid.GetLieAtPosition(Point{X: 15, Y: 15})

	// Set to Rough
	grid.SetLieAtPosition(Point{X: 15, Y: 15}, Rough)
	newLie := grid.GetLieAtPosition(Point{X: 15, Y: 15})

	if newLie == initialLie {
		t.Errorf("SetLieAtPosition did not change lie type")
	}

	if newLie != Rough {
		t.Errorf("After SetLieAtPosition, lie = %v, want Rough", newLie)
	}
}

// Test that grid correctly maps units to cell indices
func TestCourseGrid_CellIndexCalculation(t *testing.T) {
	// 100 yards x 100 yards with 10-yard cells = 10x10 grid
	// Each cell is 10 yards = 72 units
	grid := NewCourseGrid(Yard(100), Yard(100), Yard(10))

	// Set specific cells (in units)
	// Cell [0][0]: X=0-72, Y=0-72
	// Cell [1][1]: X=72-144, Y=72-144
	// Cell [9][9]: X=648-720, Y=648-720
	grid.SetLieAtPosition(Point{X: 0, Y: 0}, Tee)       // Cell [0][0]
	grid.SetLieAtPosition(Point{X: 72, Y: 72}, Rough)   // Cell [1][1]
	grid.SetLieAtPosition(Point{X: 680, Y: 680}, Green) // Cell [9][9]

	// Verify positions within the same cell get the same lie
	tests := []struct {
		name     string
		position Point
		expected LieType
	}{
		{"Tee cell - corner", Point{X: 0, Y: 0}, Tee},
		{"Tee cell - center", Point{X: 36, Y: 36}, Tee},
		{"Tee cell - edge", Point{X: 71, Y: 71}, Tee},
		{"Rough cell - corner", Point{X: 72, Y: 72}, Rough},
		{"Rough cell - center", Point{X: 108, Y: 108}, Rough},
		{"Green cell - anywhere", Point{X: 680, Y: 680}, Green},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lie := grid.GetLieAtPosition(tt.position)
			if lie != tt.expected {
				t.Errorf("GetLieAtPosition(%+v) = %v, want %v",
					tt.position, lie, tt.expected)
			}
		})
	}
}
