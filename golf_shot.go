package main

//what we need (to start)
// - a way to indicate a 'fade' or 'draw'
// - a way for these shots to limit the failure in some way
//   - instead of it going right or left, if its a draw, it could only go left, by maybe to much left
// - some sort of result, this should probably include
// - carry distance
// - roll distance
// - list of points to make a path that could change direction (fade or draw) to get around something

type ShotShape int

// const (
// 	Straight ShotShape = iota
// 	Fade     ShotShape = iota
// 	Draw     ShotShape = iota
// )

type GolfShot interface {
	Execute(Golfer, GolfBall) GolfShotResult
	// Club      Club
	// ShotShape ShotShape
	// Target    Point
	// Power     float32
}

type Straight struct {
	//what would these things be called?
	Club   Club
	Target Point
	Power  float32
}

type Fade struct {
	Club   Club
	Target Point
	Power  float32
}

type Draw struct {
	Club   Club
	Target Point
	Power  float32
}

type GolfShotResult struct {
	CarryDistance float32
	RollDistance  float32
	Path          []Point
}

func (gs *GolfShotResult) TotalDistance() float32 {
	return gs.RollDistance + gs.CarryDistance
}

func (gs *Straight) Execute(golfer Golfer, ball GolfBall) GolfShotResult {
	//todo:
	// - calculate shot distance
	// - depending on club, determine carry and roll.
	// - plot path, point for start, carry end (and direction change for curved shots), and roll end.
	return GolfShotResult{}
}

func (gs *Fade) Execute(golfer Golfer, ball GolfBall) GolfShotResult {
	//todo:
	// - calculate shot distance
	// - depending on club, determine carry and roll.
	// - plot path, point for start, carry end (and direction change for curved shots), and roll end.
	return GolfShotResult{}
}

func (gs *Draw) Execute(golfer Golfer, ball GolfBall) GolfShotResult {
	//todo:
	// - calculate shot distance
	// - depending on club, determine carry and roll.
	// - plot path, point for start, carry end (and direction change for curved shots), and roll end.
	return GolfShotResult{}
}
