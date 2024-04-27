package main

// have the following actions available in game
// - inspecting a lie to know its quality (Recovery check)
// - shaping a shot, or adding/removing spin (Any, except putting check)
// - get out of thick rough (Recovery check)
// - Normal shots (appropriate check)
// - Flop shot (chipping check)
// - Pitch (Approach or Chipping)
// - TeeShot (Driving or Approach)
type Skills struct {
	Recovery int
	Driving  int
	Approach int
	Chipping int
	Putting  int
}

// have the following actions available in game
// - inspecting a lie to know its quality (intellect check)
// - shaping a shot, or adding/removing spin (control check)
// - get out of thick rough (strength check)
type Abilities struct {
	Strength  int
	Intellect int
	Control   int
}

type Golfer struct {
	Name      string
	Skills    Skills
	Abilities Abilities
}

type Lie struct {
	Name                 string
	TargetNumberModifier int
	// 'quality' - being based on rsk, do we want
	// 'quality' to actually be more in line
	// with that definition? rather than a multiplier to distance?
	Quality float32
}

type GolfBall struct {
	Lie      Lie
	Location Point
}

// when we hit the ball
// based on shot shape, power, success, etc
// the ball will move to a location
// eventually, the path of the ball needs to be considered
// also need to consider lie (are we in the fairway? a bunker? the rough?)
func (b *GolfBall) Hit() {
	// skill check success
	if true {
		b.Location.Y += int(toFeet(300))
	} else {
		// the failure could be a hook/straight if we were trying to draw
		// a slice/straight if we were trying to fade
		// the shape/straight ration could be like 80/20 or maybe based on your skill?
		// and random if straight, plus maybe some other types of misses...
		// golf is about playing your miss
		// this would promote picking a shot shape, to know the miss.
		b.Location.Y += int(toFeet(250))
		b.Location.X -= int(toFeet(30))
	}
	// where do we end up?
	b.Lie = Lie{
		Name:                 "Perfect",
		TargetNumberModifier: 1,
		Quality:              1,
	}
}

// func skillCheck(targetNumber int) {
// 	dice := Dice{
// 		Sides: 6,
// 	}
// 	result, rolls := dice.rollN(3)
// 	margin := targetNumber - result
// 	success := margin >= 0
// 	fmt.Println(success, targetNumber, result, margin, rolls)
// }
