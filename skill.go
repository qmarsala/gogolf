package gogolf

type Skill struct {
	Name       string
	Level      int
	Experience int
}

func NewSkill(name string) Skill {
	return Skill{
		Name:       name,
		Level:      1,
		Experience: 0,
	}
}

func (s Skill) Value() int {
	return s.Level
}

func (s *Skill) AddExperience(xp int) {
	if s.Level >= 9 {
		return
	}

	s.Experience += xp

	for s.CanLevelUp() {
		threshold := s.ExperienceToNextLevel()
		s.Experience -= threshold
		s.Level++

		if s.Level >= 9 {
			s.Level = 9
			s.Experience = 0
			break
		}
	}
}

func (s Skill) ExperienceToNextLevel() int {
	if s.Level >= 9 {
		return 0
	}
	return (s.Level + 1) * 50
}

func (s Skill) CanLevelUp() bool {
	if s.Level >= 9 {
		return false
	}
	return s.Experience >= s.ExperienceToNextLevel()
}
