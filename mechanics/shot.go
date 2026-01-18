package mechanics

import (
	"gogolf"
	"math"
)

func CalculateRotation(club gogolf.Club, result gogolf.SkillCheckResult, random gogolf.RandomSource) float64 {
	clubAcc := float64(club.AccuracyDegrees())

	switch result.Outcome {
	case gogolf.CriticalSuccess:
		return random.Float64() * clubAcc * 0.1

	case gogolf.Excellent:
		possibleRotation := random.Float64() * clubAcc * 0.5
		return math.Max(possibleRotation-float64(result.Margin)*0.5, 0)

	case gogolf.Good:
		possibleRotation := math.Min(random.Float64()*clubAcc, clubAcc)
		return math.Max(possibleRotation-float64(result.Margin), 0)

	case gogolf.Marginal:
		return random.Float64() * clubAcc * 0.9

	case gogolf.Poor:
		baseMisHit := 45 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc, clubAcc)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin)), 1)

	case gogolf.Bad:
		baseMisHit := 60 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc*1.5, clubAcc*1.5)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin))*1.5, 1)

	case gogolf.CriticalFailure:
		return 60 + random.Float64()*30
	}

	return clubAcc
}

func CalculatePower(club gogolf.Club, initialPower float64, result gogolf.SkillCheckResult) float64 {
	switch result.Outcome {
	case gogolf.CriticalSuccess:
		return initialPower * 1.05

	case gogolf.Excellent, gogolf.Good:
		return initialPower

	case gogolf.Marginal:
		return initialPower * 0.95

	case gogolf.Poor:
		return math.Max(initialPower*(float64(club.Forgiveness)-(math.Abs(float64(result.Margin))/100)), 0.1)

	case gogolf.Bad:
		return math.Max(initialPower*(float64(club.Forgiveness)*0.8-(math.Abs(float64(result.Margin))/100)), 0.1)

	case gogolf.CriticalFailure:
		return initialPower * 0.2
	}

	return initialPower
}

func GetShotQualityDescription(result gogolf.SkillCheckResult) string {
	descriptions := map[gogolf.SkillCheckOutcome]string{
		gogolf.CriticalSuccess: "PURE STRIKE! Perfect contact.",
		gogolf.Excellent:       "Great shot! Solid compression.",
		gogolf.Good:            "Good contact. Ball flights well.",
		gogolf.Marginal:        "Just caught it. Got away with one.",
		gogolf.Poor:            "Slight miss. Not quite centered.",
		gogolf.Bad:             "Poor contact. Significant mishit.",
		gogolf.CriticalFailure: "DISASTER! Completely topped/chunked it.",
	}
	return descriptions[result.Outcome]
}
