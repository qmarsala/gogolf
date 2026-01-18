package game

import (
	"gogolf"
	"math/rand/v2"
	"testing"
)

func TestNew(t *testing.T) {
	g := New("TestPlayer", 3)

	if g.Golfer.Name != "TestPlayer" {
		t.Errorf("expected golfer name 'TestPlayer', got '%s'", g.Golfer.Name)
	}
	if len(g.Course.Holes) != 3 {
		t.Errorf("expected 3 holes, got %d", len(g.Course.Holes))
	}
	if g.CurrentHoleIndex != 0 {
		t.Errorf("expected to start at hole 0, got %d", g.CurrentHoleIndex)
	}
}

func TestGameTeeUp(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	if g.Ball.Location.X != 0 || g.Ball.Location.Y != 0 {
		t.Errorf("expected ball at origin after tee up, got %+v", g.Ball.Location)
	}
}

func TestGetCurrentHole(t *testing.T) {
	g := New("TestPlayer", 3)

	hole := g.GetCurrentHole()
	if hole.Number != 1 {
		t.Errorf("expected hole 1, got %d", hole.Number)
	}
}

func TestGetContext(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	ctx := g.GetContext()

	if ctx.Golfer.Name != "TestPlayer" {
		t.Errorf("expected golfer name 'TestPlayer', got '%s'", ctx.Golfer.Name)
	}
	if ctx.Hole.Number != 1 {
		t.Errorf("expected hole 1, got %d", ctx.Hole.Number)
	}
	if ctx.CurrentClub.Name == "" {
		t.Error("expected a club to be selected")
	}
}

func TestTakeShot(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	initialStrokes := g.ScoreCard.TotalStrokes()
	result := g.TakeShot(0.8)

	if g.ScoreCard.TotalStrokes() != initialStrokes+1 {
		t.Errorf("expected stroke count to increase by 1")
	}
	if result.ClubName == "" {
		t.Error("expected shot result to have club name")
	}
	if result.Power != 0.8 {
		t.Errorf("expected power 0.8, got %f", result.Power)
	}
	if result.Distance <= 0 {
		t.Error("expected positive distance")
	}
}

func TestTakeShotAwardsXP(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	result := g.TakeShot(1.0)

	if result.XPEarned <= 0 {
		t.Error("expected XP to be earned")
	}
}

func TestIsHoleComplete(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	if g.IsHoleComplete() {
		t.Error("hole should not be complete at start")
	}
}

func TestIsRoundComplete(t *testing.T) {
	g := New("TestPlayer", 3)

	if g.IsRoundComplete() {
		t.Error("round should not be complete at start")
	}
}

func TestNextHole(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	g.NextHole()

	if g.CurrentHoleIndex != 1 {
		t.Errorf("expected hole index 1, got %d", g.CurrentHoleIndex)
	}
}

func TestCompleteHoleAwardsMoney(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	initialMoney := g.Golfer.Money
	g.ScoreCard.RecordStroke(g.GetCurrentHole())
	g.CompleteHole()

	if g.Golfer.Money <= initialMoney {
		t.Error("expected money to increase after completing hole")
	}
}

func TestGameWithCustomRandom(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 42))
	g := NewWithRandom("TestPlayer", 3, rng)

	if g.random == nil {
		t.Error("expected random to be set")
	}
}

func TestStrokesThisHole(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	if g.StrokesThisHole() != 0 {
		t.Errorf("expected 0 strokes, got %d", g.StrokesThisHole())
	}

	g.TakeShot(0.5)

	if g.StrokesThisHole() != 1 {
		t.Errorf("expected 1 stroke, got %d", g.StrokesThisHole())
	}
}

func TestMaxStrokesPerHole(t *testing.T) {
	g := New("TestPlayer", 3)
	g.TeeUp()

	for i := 0; i < 11; i++ {
		g.TakeShot(0.1)
	}

	if !g.IsHoleComplete() {
		t.Error("hole should be complete after 11 strokes")
	}
}

func TestCalculateXP(t *testing.T) {
	tests := []struct {
		outcome  gogolf.SkillCheckOutcome
		expected int
	}{
		{gogolf.CriticalSuccess, 15},
		{gogolf.Excellent, 10},
		{gogolf.Good, 7},
		{gogolf.Marginal, 5},
		{gogolf.Poor, 3},
		{gogolf.Bad, 2},
		{gogolf.CriticalFailure, 1},
	}

	for _, tt := range tests {
		got := calculateXP(tt.outcome)

		if got != tt.expected {
			t.Errorf("calculateXP(%v) = %v, want %v",
				tt.outcome, got, tt.expected)
		}
	}
}

func TestNewFromGolfer(t *testing.T) {
	golfer := gogolf.NewGolfer("LoadedPlayer")
	golfer.Money = 500
	golfer.Ball = &gogolf.Ball{Name: "Pro V1", DistanceBonus: 8, SpinControl: 0.9, Cost: 75}

	skill := golfer.Skills["Driver"]
	(&skill).AddExperience(200)
	golfer.Skills["Driver"] = skill

	g := NewFromGolfer(golfer, 3)

	if g.Golfer.Name != "LoadedPlayer" {
		t.Errorf("expected golfer name 'LoadedPlayer', got '%s'", g.Golfer.Name)
	}
	if g.Golfer.Money != 500 {
		t.Errorf("expected money 500, got %d", g.Golfer.Money)
	}
	if g.Golfer.Ball == nil || g.Golfer.Ball.Name != "Pro V1" {
		t.Errorf("expected ball 'Pro V1', got %v", g.Golfer.Ball)
	}
	if g.Golfer.Skills["Driver"].Level != 2 {
		t.Errorf("expected Driver level 2, got %d", g.Golfer.Skills["Driver"].Level)
	}
	if len(g.Course.Holes) != 3 {
		t.Errorf("expected 3 holes, got %d", len(g.Course.Holes))
	}
}

func TestNewFromGolferPreservesAbilities(t *testing.T) {
	golfer := gogolf.NewGolfer("TestPlayer")

	ability := golfer.Abilities["Strength"]
	(&ability).AddExperience(150)
	golfer.Abilities["Strength"] = ability

	g := NewFromGolfer(golfer, 1)

	if g.Golfer.Abilities["Strength"].Level != 2 {
		t.Errorf("expected Strength level 2, got %d", g.Golfer.Abilities["Strength"].Level)
	}
	if g.Golfer.Abilities["Strength"].Experience != 50 {
		t.Errorf("expected Strength XP 50, got %d", g.Golfer.Abilities["Strength"].Experience)
	}
}
