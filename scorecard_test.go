package main

import (
	"testing"
)

func TestRecordScore(t *testing.T) {
	hole := Hole{
		Number: 1,
		Par:    3,
	}
	sc := ScoreCard{
		Course: Course{[]Hole{hole}},
		Scores: map[int]int{},
	}

	sc.RecordStroke(hole)
	sc.RecordStroke(hole)

	strokes := sc.Scores[hole.Number]
	if strokes != 2 {
		t.Error("Expected strokes for hole 1 to be 2, but got", strokes)
	}
}

func TestScore(t *testing.T) {
	hole := Hole{
		Number: 1,
		Par:    3,
	}
	hole2 := Hole{
		Number: 2,
		Par:    4,
	}
	sc := ScoreCard{
		Course: Course{[]Hole{hole}},
		Scores: map[int]int{},
	}

	sc.RecordStroke(hole)
	sc.RecordStroke(hole)

	sc.RecordStroke(hole2)
	sc.RecordStroke(hole2)
	sc.RecordStroke(hole2)
	sc.RecordStroke(hole2)

	score := sc.Score()
	if score != -1 {
		t.Error("Expected score to be -1, but got", score)
	}
}

func TestTotalStrokes(t *testing.T) {
	hole := Hole{
		Number: 1,
		Par:    3,
	}
	hole2 := Hole{
		Number: 2,
		Par:    4,
	}
	sc := ScoreCard{
		Course: Course{[]Hole{hole}},
		Scores: map[int]int{},
	}

	sc.RecordStroke(hole)
	sc.RecordStroke(hole)

	sc.RecordStroke(hole2)
	sc.RecordStroke(hole2)
	sc.RecordStroke(hole2)
	sc.RecordStroke(hole2)

	strokes := sc.TotalStrokes()
	if strokes != 6 {
		t.Error("Expected strokes to be 6, but got", strokes)
	}
}
