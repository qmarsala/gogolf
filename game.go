package main

import (
	"gogolf/dice"
	"gogolf/progression"
	"math"
	"math/rand/v2"
)

type Game struct {
	Golfer           Golfer
	Course           Course
	Ball             GolfBall
	ScoreCard        ScoreCard
	CurrentHoleIndex int
	random           *rand.Rand
	lastShotResult   *ShotResult
}

type GameContext struct {
	Golfer      Golfer
	Hole        Hole
	Ball        GolfBall
	ScoreCard   ScoreCard
	CurrentClub Club
	Lie         LieType
}

type ShotResult struct {
	ClubName    string
	Outcome     dice.SkillCheckOutcome
	Margin      int
	Description string
	Rotation    float64
	RotationDir string
	Power       float64
	Distance    float64
	XPEarned    int
	LevelUps    []string
	HoledOut    bool
	TapIn       bool
}

func NewGame(playerName string, holeCount int) *Game {
	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	return NewGameWithRandom(playerName, holeCount, rng)
}

func NewGameWithRandom(playerName string, holeCount int, rng *rand.Rand) *Game {
	course, scoreCard := GenerateSimpleCourse(holeCount)
	return &Game{
		Golfer:           NewGolfer(playerName),
		Course:           course,
		Ball:             GolfBall{Location: Point{X: 0, Y: 0}},
		ScoreCard:        scoreCard,
		CurrentHoleIndex: 0,
		random:           rng,
	}
}

func (g *Game) TeeUp() {
	g.Ball.TeeUp()
	g.lastShotResult = nil
}

func (g *Game) GetCurrentHole() Hole {
	return g.Course.Holes[g.CurrentHoleIndex]
}

func (g *Game) GetGameContext() GameContext {
	hole := g.GetCurrentHole()
	club := g.Golfer.GetBestClub(g.Ball.Location.Distance(hole.HoleLocation).Yards())
	lie := g.Ball.GetLie(&hole)

	return GameContext{
		Golfer:      g.Golfer,
		Hole:        hole,
		Ball:        g.Ball,
		ScoreCard:   g.ScoreCard,
		CurrentClub: club,
		Lie:         lie,
	}
}

func (g *Game) TakeShot(power float64) ShotResult {
	hole := g.GetCurrentHole()
	club := g.Golfer.GetBestClub(g.Ball.Location.Distance(hole.HoleLocation).Yards())
	directionToHole := g.Ball.Location.Direction(hole.HoleLocation)

	lie := g.Ball.GetLie(&hole)
	difficulty := lie.DifficultyModifier()

	targetNumber := g.Golfer.CalculateTargetNumber(club, difficulty)
	result := g.Golfer.SkillCheck(dice.NewD6(), targetNumber)

	rotationDirection := float64(1)
	if int(math.Abs(float64(result.Margin)))%2 == 0 {
		rotationDirection *= -1
	}

	modifiedClub := g.Golfer.GetModifiedClub(club)

	rotationDegrees := CalculateRotation(modifiedClub, result, g.random)
	adjustedPower := CalculatePower(modifiedClub, power, result)

	skill := g.Golfer.GetSkillForClub(club)
	ability := g.Golfer.GetAbilityForClub(club)
	prevSkillLevel := skill.Level
	prevAbilityLevel := ability.Level

	xpAward := calculateXP(result.Outcome)
	g.Golfer.AwardExperience(club, xpAward)

	newSkill := g.Golfer.GetSkillForClub(club)
	newAbility := g.Golfer.GetAbilityForClub(club)
	var levelUps []string

	if newSkill.Level > prevSkillLevel {
		levelUps = append(levelUps, newSkill.Name+" leveled up!")
	}
	if newAbility.Level > prevAbilityLevel {
		levelUps = append(levelUps, newAbility.Name+" leveled up!")
	}

	directionToHole.Rotate(rotationDegrees * rotationDirection)
	ballPath := g.Ball.ReceiveHit(modifiedClub, float32(adjustedPower), directionToHole)

	if club.Name != "Putter" {
		if g.random.IntN(100)%2 == 0 {
			g.applyDraw(ballPath, hole)
		} else {
			g.applyFade(ballPath, hole)
		}
	}

	g.ScoreCard.RecordStroke(hole)

	rotationDir := "right"
	if rotationDirection < 0 {
		rotationDir = "left"
	}

	holedOut := hole.DetectHoleOut(g.Ball, ballPath)
	tapIn := !holedOut && hole.DetectTapIn(g.Ball)

	if tapIn {
		holedOut = true
		g.ScoreCard.RecordStroke(hole)
	}

	shotResult := ShotResult{
		ClubName:    club.Name,
		Outcome:     result.Outcome,
		Margin:      result.Margin,
		Description: GetShotQualityDescription(result),
		Rotation:    rotationDegrees,
		RotationDir: rotationDir,
		Power:       power,
		Distance:    float64(Unit(ballPath.Magnitude()).Yards()),
		XPEarned:    xpAward,
		LevelUps:    levelUps,
		HoledOut:    holedOut,
		TapIn:       tapIn,
	}

	g.lastShotResult = &shotResult
	return shotResult
}

func (g *Game) applyDraw(ballPath Vector, h Hole) {
	directionToHole := g.Ball.PrevLocation.Direction(h.HoleLocation)
	drawRotationDegrees := -45
	if directionToHole.Y < 0 {
		drawRotationDegrees = 45
	}
	rotatedPath := ballPath.Rotate(float64(drawRotationDegrees))
	translationDistance := Yard(math.Max(g.random.Float64()*3, 1)).Units()
	g.Ball.Location = g.Ball.Location.Move(rotatedPath, float64(translationDistance))
}

func (g *Game) applyFade(ballPath Vector, h Hole) {
	directionToHole := g.Ball.PrevLocation.Direction(h.HoleLocation)
	drawRotationDegrees := 45
	if directionToHole.Y < 0 {
		drawRotationDegrees = -45
	}
	rotatedPath := ballPath.Rotate(float64(drawRotationDegrees))
	translationDistance := Yard(math.Max(g.random.Float64()*3, 1)).Units()
	g.Ball.Location = g.Ball.Location.Move(rotatedPath, float64(translationDistance))
}

func (g *Game) IsHoleComplete() bool {
	if g.lastShotResult != nil && g.lastShotResult.HoledOut {
		return true
	}
	return g.StrokesThisHole() >= 11
}

func (g *Game) IsRoundComplete() bool {
	return g.CurrentHoleIndex >= len(g.Course.Holes)
}

func (g *Game) StrokesThisHole() int {
	return g.ScoreCard.TotalStrokesThisHole(g.GetCurrentHole())
}

func (g *Game) NextHole() {
	g.CurrentHoleIndex++
	if !g.IsRoundComplete() {
		g.TeeUp()
	}
}

func (g *Game) CompleteHole() int {
	hole := g.GetCurrentHole()
	strokes := g.StrokesThisHole()
	g.Golfer.AwardHoleReward(hole.Par, strokes)
	return progression.CalculateHoleReward(hole.Par, strokes)
}

func (g *Game) GetLastShotResult() *ShotResult {
	return g.lastShotResult
}

func calculateXP(outcome dice.SkillCheckOutcome) int {
	switch outcome {
	case dice.CriticalSuccess:
		return 15
	case dice.Excellent:
		return 10
	case dice.Good:
		return 7
	case dice.Marginal:
		return 5
	case dice.Poor:
		return 3
	case dice.Bad:
		return 2
	case dice.CriticalFailure:
		return 1
	default:
		return 1
	}
}
