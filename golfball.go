package gogolf

import "fmt"

type GolfBall struct {
	Location     Point
	PrevLocation Point
}

func (ball *GolfBall) TeeUp() {
	ball.Location = Point{0, 0}
	ball.PrevLocation = Point{0, 0}
	fmt.Println("Ball teed up")
}

// GetLie returns the lie type at the ball's current location
func (b GolfBall) GetLie(hole *Hole) LieType {
	return hole.GetLieAtPosition(b.Location)
}

// this is probably the only place we need this logic
// - did the ball hit anything on its way (or the hole)
// - its the only thing 'moving' in the game
// - perhaps when we get to collision detection with hazards
// it will make more sense for this fn to be on the 'boundary' to see if
// anything entered and exited, and event that.
func (b GolfBall) CheckForCollision(bPath Vector, location Point) (bool, Unit) {
	toLocVec := Vector{X: float64(location.X - b.PrevLocation.X), Y: float64(location.Y - b.PrevLocation.Y)}
	dotProduct := toLocVec.Dot(bPath)
	squaredLengthBPath := bPath.Dot(bPath)
	projectionFactor := dotProduct / squaredLengthBPath

	closestPoint := Point{
		X: int(float64(b.PrevLocation.X) + projectionFactor*bPath.X),
		Y: int(float64(b.PrevLocation.Y) + projectionFactor*bPath.Y),
	}
	if projectionFactor < 0 {
		closestPoint = b.PrevLocation
	} else if projectionFactor > 1 {
		closestPoint = b.Location
	}
	distance := closestPoint.Distance(location)
	return distance < 1, distance
}

// todo: eventually we will want to 'shape' the shot
// ie. a fade/draw/slice/hook etc.
// how do we want to message that?
// - I also think it it probably best this func does not actually change the location
// and instead returns a path, so that we can do collision detection on the hole and tall hazards like trees
// the collision detection on the hole is going to be really important.  A simple alternative
// for now, could be to check if it is within a few points of the hole.  Though I think the collision
// will be much better.
func (ball *GolfBall) ReceiveHit(club Club, power float32, direction Vector) (path Vector) {
	yards := Yard(float32(club.Distance) * power)
	ball.PrevLocation = Point{ball.Location.X, ball.Location.Y}
	ball.Location = ball.Location.Move(direction, float64(yards.Units()))
	return Vector{X: float64(ball.Location.X - ball.PrevLocation.X), Y: float64(ball.Location.Y - ball.PrevLocation.Y)}
}
