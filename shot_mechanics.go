package main

import (
	"math"
	"math/rand/v2"
)

// CalculateRotation determines shot rotation based on skill check outcome tier
func CalculateRotation(club Club, result SkillCheckResult, random *rand.Rand) float64 {
	clubAcc := float64(club.AccuracyDegrees())

	switch result.Outcome {
	case CriticalSuccess:
		// Near-perfect: 10% of normal spread
		return random.Float64() * clubAcc * 0.1

	case Excellent:
		// Great contact: reduced spread
		possibleRotation := random.Float64() * clubAcc * 0.5
		return math.Max(possibleRotation-float64(result.Margin)*0.5, 0)

	case Good:
		// Solid shot: current success formula
		possibleRotation := math.Min(random.Float64()*clubAcc, clubAcc)
		return math.Max(possibleRotation-float64(result.Margin), 0)

	case Marginal:
		// Just barely made it: minimal benefit
		return random.Float64() * clubAcc * 0.9

	case Poor:
		// Minor miss: current failure formula
		baseMisHit := 45 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc, clubAcc)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin)), 1)

	case Bad:
		// Bad contact: increased spray (60° base, 1.5x multipliers)
		baseMisHit := 60 * (1 - float64(club.Forgiveness))
		minimumMisHitRotation := random.Float64() * baseMisHit
		possibleRotation := math.Max(minimumMisHitRotation+clubAcc*1.5, clubAcc*1.5)
		return math.Max(possibleRotation+math.Abs(float64(result.Margin))*1.5, 1)

	case CriticalFailure:
		// Catastrophic: 60-90° off target
		return 60 + random.Float64()*30
	}

	return clubAcc // fallback
}

// CalculatePower determines shot power based on skill check outcome tier
func CalculatePower(club Club, initialPower float64, result SkillCheckResult) float64 {
	switch result.Outcome {
	case CriticalSuccess:
		// Perfect compression: +5% bonus
		return initialPower * 1.05

	case Excellent, Good:
		// Great/solid contact: full power
		return initialPower

	case Marginal:
		// Barely caught it: -5% loss
		return initialPower * 0.95

	case Poor:
		// Current formula
		return math.Max(initialPower*(float64(club.Forgiveness)-(math.Abs(float64(result.Margin))/100)), 0.1)

	case Bad:
		// Worse contact: enhanced power loss
		return math.Max(initialPower*(float64(club.Forgiveness)*0.8-(math.Abs(float64(result.Margin))/100)), 0.1)

	case CriticalFailure:
		// Topped/chunked: 20% distance
		return initialPower * 0.2
	}

	return initialPower
}

// GetShotQualityDescription provides flavor text for each outcome tier
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
