package main

import (
	"math/rand/v2"
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame("TestPlayer", 3)

	if game.Golfer.Name != "TestPlayer" {
		t.Errorf("expected golfer name 'TestPlayer', got '%s'", game.Golfer.Name)
	}
	if len(game.Course.Holes) != 3 {
		t.Errorf("expected 3 holes, got %d", len(game.Course.Holes))
	}
	if game.CurrentHoleIndex != 0 {
		t.Errorf("expected to start at hole 0, got %d", game.CurrentHoleIndex)
	}
}

func TestGameTeeUp(t *testing.T) {
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	if game.Ball.Location.X != 0 || game.Ball.Location.Y != 0 {
		t.Errorf("expected ball at origin after tee up, got %+v", game.Ball.Location)
	}
}

func TestGetCurrentHole(t *testing.T) {
	game := NewGame("TestPlayer", 3)

	hole := game.GetCurrentHole()
	if hole.Number != 1 {
		t.Errorf("expected hole 1, got %d", hole.Number)
	}
}

func TestGetGameContext(t *testing.T) {
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	ctx := game.GetGameContext()

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
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	initialStrokes := game.ScoreCard.TotalStrokes()
	result := game.TakeShot(0.8)

	if game.ScoreCard.TotalStrokes() != initialStrokes+1 {
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
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	result := game.TakeShot(1.0)

	if result.XPEarned <= 0 {
		t.Error("expected XP to be earned")
	}
}

func TestIsHoleComplete(t *testing.T) {
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	if game.IsHoleComplete() {
		t.Error("hole should not be complete at start")
	}
}

func TestIsRoundComplete(t *testing.T) {
	game := NewGame("TestPlayer", 3)

	if game.IsRoundComplete() {
		t.Error("round should not be complete at start")
	}
}

func TestNextHole(t *testing.T) {
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	game.NextHole()

	if game.CurrentHoleIndex != 1 {
		t.Errorf("expected hole index 1, got %d", game.CurrentHoleIndex)
	}
}

func TestCompleteHoleAwardsMoney(t *testing.T) {
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	initialMoney := game.Golfer.Money
	game.ScoreCard.RecordStroke(game.GetCurrentHole())
	game.CompleteHole()

	if game.Golfer.Money <= initialMoney {
		t.Error("expected money to increase after completing hole")
	}
}


func TestGameWithCustomRandom(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 42))
	game := NewGameWithRandom("TestPlayer", 3, rng)

	if game.random == nil {
		t.Error("expected random to be set")
	}
}

func TestStrokesThisHole(t *testing.T) {
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	if game.StrokesThisHole() != 0 {
		t.Errorf("expected 0 strokes, got %d", game.StrokesThisHole())
	}

	game.TakeShot(0.5)

	if game.StrokesThisHole() != 1 {
		t.Errorf("expected 1 stroke, got %d", game.StrokesThisHole())
	}
}

func TestMaxStrokesPerHole(t *testing.T) {
	game := NewGame("TestPlayer", 3)
	game.TeeUp()

	for i := 0; i < 11; i++ {
		game.TakeShot(0.1)
	}

	if !game.IsHoleComplete() {
		t.Error("hole should be complete after 11 strokes")
	}
}
