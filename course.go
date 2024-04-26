package main

// we probably want some sort of grid, probably hex based
// to store the holes long term

// for now, for simplicity
// fairway starts at the tee box and goes to the green
// and the only hazard is rough to the sides
type Hole struct {
	Number        int
	Par           int
	Distance      int
	FairwayLength int
	FairwayWidth  int
	GreenLength   int
	GreenWidth    int
	HoleLocation  int
}

type Course struct {
	Holes         []Hole
	Par           int
	TotalDistance int
}

type Score struct {
	Hole    Hole
	Strokes int
}

type ScoreCard struct {
	Scores map[int]Score
}
