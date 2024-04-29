package main

type GolfBall struct {
	Location Point
}

// todo: eventually we will want to 'shape' the shot
// ie. a fade/draw/slice/hook etc.
// how do we want to message that?
// - I also think it it probably best this func does not actually change the location
// and instead returns a path, so that we can do collision detection on the hole and tall hazards like trees
// the collision detection on the hole is going to be really important.  A simple alternative
// for now, could be to check if it is within a few points of the hole.  Though I think the collision
// will be much better.
func (ball *GolfBall) ReceiveHit(club Club, power float32, direction Vector) {
	yards := Yard(float32(club.Distance) * power)
	ball.Location = ball.Location.Move(direction, float64(yards.Units()))
}
