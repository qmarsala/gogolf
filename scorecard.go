package main

import (
	"slices"
)

type ScoreCard struct {
	Course Course
	Scores map[int]int
}

func (sc *ScoreCard) RecordStroke(h Hole) {
	sc.Scores[h.Number]++
}

func (sc ScoreCard) TotalStrokes() (score int) {
	for _, v := range sc.Scores {
		score += v
	}
	return
}

func (sc ScoreCard) Score() (score int) {
	for k, v := range sc.Scores {
		if holeIndex := slices.IndexFunc(sc.Course.Holes, func(h Hole) bool {
			return h.Number == k
		}); holeIndex > -1 {
			score += v - sc.Course.Holes[holeIndex].Par
		}
	}
	return
}
