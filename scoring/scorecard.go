package scoring

import "gogolf"

type ScoreCard struct {
	Course gogolf.Course
	Scores map[int]int
}

func NewScoreCard(course gogolf.Course) ScoreCard {
	return ScoreCard{
		Course: course,
		Scores: map[int]int{},
	}
}

func (sc *ScoreCard) RecordStroke(h gogolf.Hole) {
	sc.Scores[h.Number]++
}

func (sc ScoreCard) TotalStrokesThrough(holeNumber int) (score int) {
	for k, v := range sc.Scores {
		if k <= holeNumber {
			score += v
		}
	}
	return
}

func (sc ScoreCard) TotalStrokes() (score int) {
	for _, v := range sc.Scores {
		score += v
	}
	return
}

func (sc ScoreCard) TotalStrokesThisHole(h gogolf.Hole) (score int) {
	return sc.Scores[h.Number]
}

func (sc ScoreCard) ScoreThisHole(h gogolf.Hole) (score int) {
	return sc.TotalStrokesThisHole(h) - h.Par
}

func (sc ScoreCard) Score() (score int) {
	return sc.TotalStrokes() - sc.Course.Par()
}

func (sc ScoreCard) ScoreThrough(holeNumber int) (score int) {
	return sc.TotalStrokesThrough(holeNumber) - sc.Course.ParUpToHole(holeNumber)
}
