package main

import (
	"slices"
)

type ScoreCard struct {
	Holes  []Hole
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
		if holeIndex := slices.IndexFunc(sc.Holes, func(h Hole) bool {
			return h.Number == k
		}); holeIndex > -1 {
			score += v - sc.Holes[holeIndex].Par
		}
	}
	return
}
