package scoring

import (
	"gogolf"
	"testing"
)

func TestRecordScore(t *testing.T) {
	hole := gogolf.Hole{
		Number: 1,
		Par:    3,
	}
	sc := ScoreCard{
		Course: gogolf.Course{Holes: []gogolf.Hole{hole}},
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
	hole := gogolf.Hole{
		Number: 1,
		Par:    3,
	}
	hole2 := gogolf.Hole{
		Number: 2,
		Par:    4,
	}
	sc := ScoreCard{
		Course: gogolf.Course{Holes: []gogolf.Hole{hole, hole2}},
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
	hole := gogolf.Hole{
		Number: 1,
		Par:    3,
	}
	hole2 := gogolf.Hole{
		Number: 2,
		Par:    4,
	}
	sc := ScoreCard{
		Course: gogolf.Course{Holes: []gogolf.Hole{hole}},
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

func TestNewScoreCard(t *testing.T) {
	hole := gogolf.Hole{Number: 1, Par: 4}
	course := gogolf.Course{Holes: []gogolf.Hole{hole}}

	sc := NewScoreCard(course)

	if sc.Course.Holes[0].Number != 1 {
		t.Error("Expected course hole 1, but got", sc.Course.Holes[0].Number)
	}
	if sc.Scores == nil {
		t.Error("Expected Scores map to be initialized")
	}
}
