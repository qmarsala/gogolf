package gogolf

type DiceRoller interface {
	RollN(count int) (total int, rolls []int)
}

type RandomSource interface {
	Float64() float64
	IntN(n int) int
}

type ClubSelector interface {
	SelectClub(clubs []Club, distance Yard) Club
}
