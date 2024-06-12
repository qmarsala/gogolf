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

const (
	Straight ShotShape = iota
	Fade     ShotShape = iota
	Draw     ShotShape = iota
)

type GolfShot struct {
	Club      Club
	ShotShape ShotShape
	Target    Point
	Power     float32
	Loft      float32
}

type GolfShotResult struct {
}

func (gs *GolfShot) Execute(golfer Golfer) GolfShotResult {
	return GolfShotResult{}
}
