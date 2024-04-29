package main

type Club struct {
	Name     string
	Distance Yard
}

type Golfer struct {
	Name   string
	Target Point
}

// how do we want to do this?
// we can have different types of shots
// and use different clubs
// hit it softer or harder
// these types of decisions can affect the target number
// func (g *Golfer) Swing(club Club) {
// 	dice := NewD6()
// 	//todo: Depending on the shot we will use a different skill/attribute
// 	// we will also need to add mods such as lie boost/penalties
// 	// for now, just something
// 	targetNumber := g.Skills[club.DefaultSkill].Value() + g.Abilities[club.DefaultAbility].Value()
// 	result := dice.SkillCheck(targetNumber)
// 	if result.Success {
// 		g.GolfBall.ReceiveHit(Club{}, 1)
// 	} else {
// 		g.GolfBall.ReceiveHit(Club{}, .65)
// 	}
// }

//aiming
// - thinking that by default you will be aiming at the flag
// - to do this, we probably need to calculate the angle
// from the y axis to the hole location and you. so we can do things
// relative to your plane with the hole location.
// - aiming would then rotate this plane?
// - shots then fade/draw/hook/slice etc relative to this plane

// when we hit the ball
// based on shot shape, power, success, etc
// the ball will move to a location
// eventually, the path of the ball needs to be considered
// also need to consider lie (are we in the fairway? a bunker? the rough?)
// skill check success
// where do we end up?
// the failure could be a hook/straight if we were trying to draw
// a slice/straight if we were trying to fade
// the shape/straight ration could be like 80/20 or maybe based on your skill?
// and random if straight, plus maybe some other types of misses...
// golf is about playing your miss
// this would promote picking a shot shape, to know the miss.
