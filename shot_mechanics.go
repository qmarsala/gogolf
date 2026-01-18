package main

import (
	"math"
	"math/rand/v2"
)

func CalculateRotation(club Club, result SkillCheckResult, random *rand.Rand) float64 {
	clubAcc := float64(club.AccuracyDegrees())

	switch result.Outcome {
	case CriticalSuccess:
		return random.Float64() * clubAcc * 0.1

	case Excellent:
		possibleRotation := random.Float64() * clubAcc * 0.5
		return math.Max(possibleRotation-float64(result.Margin)*0.5, 0)

	case Good:
		possibleRotation := math.Min(random.Float64()*clubAcc, clubAcc)
		return math.Max(possibleRotation-float64(result.Margin), 0)

	case Marginal:
		return random.Float64() * clubAcc * 0.9

	case Poor:
		baseMisHit := 45 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc, clubAcc)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin)), 1)

	case Bad:
		baseMisHit := 60 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc*1.5, clubAcc*1.5)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin))*1.5, 1)

	case CriticalFailure:
		return 60 + random.Float64()*30
	}

	return clubAcc
}

func CalculatePower(club Club, initialPower float64, result SkillCheckResult) float64 {
	switch result.Outcome {
	case CriticalSuccess:
		return initialPower * 1.05

	case Excellent, Good:
		return initialPower

	case Marginal:
		return initialPower * 0.95

	case Poor:
		return math.Max(initialPower*(float64(club.Forgiveness)-(math.Abs(float64(result.Margin))/100)), 0.1)

	case Bad:
		return math.Max(initialPower*(float64(club.Forgiveness)*0.8-(math.Abs(float64(result.Margin))/100)), 0.1)

	case CriticalFailure:
		return initialPower * 0.2
	}

	return initialPower
}

func GetShotQualityDescription(result SkillCheckResult) string {
	descriptions := map[SkillCheckOutcome]string{
		CriticalSuccess: "PURE STRIKE! Perfect contact.",
		Excellent:       "Great shot! Solid compression.",
		Good:            "Good contact. Ball flights well.",
		Marginal:        "Just caught it. Got away with one.",
		Poor:            "Slight miss. Not quite centered.",
		Bad:             "Poor contact. Significant mishit.",
		CriticalFailure: "DISASTER! Completely topped/chunked it.",
	}
	return descriptions[result.Outcome]
}
