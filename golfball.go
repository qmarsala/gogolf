package main

import "fmt"

type GolfBall struct {
	Location Point
}

func (ball *GolfBall) TeeUp() {
	ball.Location = Point{0, 0}
	fmt.Println("Ball teed up")
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
	startingLocation := Point{ball.Location.X, ball.Location.Y}
	ball.Location = ball.Location.Move(direction, float64(yards.Units()))
	fmt.Printf("Ball traveled %f\n", startingLocation.Distance(ball.Location).Yards())
}
