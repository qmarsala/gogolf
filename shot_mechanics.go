package main

import (
	"gogolf/dice"
	"math"
	"math/rand/v2"
)

func CalculateRotation(club Club, result dice.SkillCheckResult, random *rand.Rand) float64 {
	clubAcc := float64(club.AccuracyDegrees())

	switch result.Outcome {
	case dice.CriticalSuccess:
		return random.Float64() * clubAcc * 0.1

	case dice.Excellent:
		possibleRotation := random.Float64() * clubAcc * 0.5
		return math.Max(possibleRotation-float64(result.Margin)*0.5, 0)

	case dice.Good:
		possibleRotation := math.Min(random.Float64()*clubAcc, clubAcc)
		return math.Max(possibleRotation-float64(result.Margin), 0)

	case dice.Marginal:
		return random.Float64() * clubAcc * 0.9

	case dice.Poor:
		baseMisHit := 45 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc, clubAcc)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin)), 1)

	case dice.Bad:
		baseMisHit := 60 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc*1.5, clubAcc*1.5)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin))*1.5, 1)

	case dice.CriticalFailure:
		return 60 + random.Float64()*30
	}

	return clubAcc
}

func CalculatePower(club Club, initialPower float64, result dice.SkillCheckResult) float64 {
	switch result.Outcome {
	case dice.CriticalSuccess:
		return initialPower * 1.05

	case dice.Excellent, dice.Good:
		return initialPower

	case dice.Marginal:
		return initialPower * 0.95

	case dice.Poor:
		return math.Max(initialPower*(float64(club.Forgiveness)-(math.Abs(float64(result.Margin))/100)), 0.1)

	case dice.Bad:
		return math.Max(initialPower*(float64(club.Forgiveness)*0.8-(math.Abs(float64(result.Margin))/100)), 0.1)

	case dice.CriticalFailure:
		return initialPower * 0.2
	}

	return initialPower
}

func GetShotQualityDescription(result dice.SkillCheckResult) string {
	descriptions := map[dice.SkillCheckOutcome]string{
		dice.CriticalSuccess: "PURE STRIKE! Perfect contact.",
		dice.Excellent:       "Great shot! Solid compression.",
		dice.Good:            "Good contact. Ball flights well.",
		dice.Marginal:        "Just caught it. Got away with one.",
		dice.Poor:            "Slight miss. Not quite centered.",
		dice.Bad:             "Poor contact. Significant mishit.",
		dice.CriticalFailure: "DISASTER! Completely topped/chunked it.",
	}
	return descriptions[result.Outcome]
}
