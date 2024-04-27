package main

type Skills struct {
	Recovery int
	Driving  int
	Approach int
	Chipping int
	Putting  int
}

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

// func skillCheck(targetNumber int) {
// 	dice := Dice{
// 		Sides: 6,
// 	}
// 	result, rolls := dice.rollN(3)
// 	margin := targetNumber - result
// 	success := margin >= 0
// 	fmt.Println(success, targetNumber, result, margin, rolls)
// }
