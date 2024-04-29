package main

type GolfBall struct {
	Location Point
}

func (ball *GolfBall) ReceiveHit(club Club, power float32, direction Vector) {
	yards := Yard(float32(club.Distance) * power)
	ball.Location = ball.Location.Move(direction, float64(yards.Units()))
}
