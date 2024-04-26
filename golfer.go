package main

type Skills struct {
	Woods   int
	Irons   int
	Wedges  int
	Putting int
}

type Abilities struct {
	Power    int
	Accuracy int
	Control  int
}

type Golfer struct {
	Name      string
	Skills    Skills
	Abilities Abilities
}
