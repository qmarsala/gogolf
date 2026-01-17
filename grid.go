package main

// GridCell represents a single cell in the course grid
type GridCell struct {
	Position Point   // Center position of the cell
	Lie      LieType // The lie type for this cell
}

// CourseGrid represents a spatial grid of the golf course
// The grid maps positions on the course to lie types
type CourseGrid struct {
	Cells    [][]GridCell // 2D array of grid cells [row][col]
	Width    Yard         // Total width of the grid
	Length   Yard         // Total length of the grid
	CellSize Yard         // Size of each grid cell (square)
}

// NewCourseGrid creates a new course grid with the specified dimensions
// Width and length are in yards, cellSize determines grid resolution
func NewCourseGrid(width, length, cellSize Yard) CourseGrid {
	rows := int(length / cellSize)
	cols := int(width / cellSize)

	// Initialize 2D array
	cells := make([][]GridCell, rows)
	for i := range cells {
		cells[i] = make([]GridCell, cols)
		for j := range cells[i] {
			// Set position to center of each cell
			cells[i][j] = GridCell{
				Position: Point{
					X: int(cellSize.Units()) * j + int(cellSize.Units())/2,
					Y: int(cellSize.Units()) * i + int(cellSize.Units())/2,
				},
				Lie: Fairway, // Default to fairway
			}
		}
	}

	return CourseGrid{
		Cells:    cells,
		Width:    width,
		Length:   length,
		CellSize: cellSize,
	}
}

// GetLieAtPosition returns the lie type at a given position
// Returns PenaltyArea if position is out of bounds
func (g CourseGrid) GetLieAtPosition(pos Point) LieType {
	row, col := g.positionToIndices(pos)

	// Check bounds
	if row < 0 || row >= len(g.Cells) || col < 0 || col >= len(g.Cells[0]) {
		return PenaltyArea // Out of bounds
	}

	return g.Cells[row][col].Lie
}

// SetLieAtPosition updates the lie type at a given position
// Does nothing if position is out of bounds
func (g *CourseGrid) SetLieAtPosition(pos Point, lie LieType) {
	row, col := g.positionToIndices(pos)

	// Check bounds
	if row < 0 || row >= len(g.Cells) || col < 0 || col >= len(g.Cells[0]) {
		return // Out of bounds, do nothing
	}

	g.Cells[row][col].Lie = lie
}

// positionToIndices converts a position to grid cell indices
func (g CourseGrid) positionToIndices(pos Point) (row, col int) {
	// Check for negative coordinates first
	if pos.X < 0 || pos.Y < 0 {
		return -1, -1 // Out of bounds
	}

	cellSizeUnits := int(g.CellSize.Units())
	col = pos.X / cellSizeUnits
	row = pos.Y / cellSizeUnits
	return row, col
}
