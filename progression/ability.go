package progression

type Ability struct {
	Name       string
	Level      int
	Experience int
}

func NewAbility(name string) Ability {
	return Ability{
		Name:       name,
		Level:      1,
		Experience: 0,
	}
}

func (a Ability) Value() int {
	return a.Level * 2
}

func (a *Ability) AddExperience(xp int) {
	if a.Level >= 9 {
		return
	}

	a.Experience += xp

	for a.CanLevelUp() {
		threshold := a.ExperienceToNextLevel()
		a.Experience -= threshold
		a.Level++

		if a.Level >= 9 {
			a.Level = 9
			a.Experience = 0
			break
		}
	}
}

func (a Ability) ExperienceToNextLevel() int {
	if a.Level >= 9 {
		return 0
	}
	return (a.Level + 1) * 50
}

func (a Ability) CanLevelUp() bool {
	if a.Level >= 9 {
		return false
	}
	return a.Experience >= a.ExperienceToNextLevel()
}
