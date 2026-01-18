package gogolf

import (
	"math/rand/v2"
	"testing"
)

func TestShotShape_String(t *testing.T) {
	tests := []struct {
		shape    ShotShape
		expected string
	}{
		{Straight, "Straight"},
		{Draw, "Draw"},
		{Fade, "Fade"},
		{Hook, "Hook"},
		{Slice, "Slice"},
	}

	for _, tt := range tests {
		if got := tt.shape.String(); got != tt.expected {
			t.Errorf("%v.String() = %v, want %v", tt.shape, got, tt.expected)
		}
	}
}

func TestShotShape_DifficultyModifier(t *testing.T) {
	tests := []struct {
		shape    ShotShape
		expected int
	}{
		{Straight, -2},
		{Draw, 1},
		{Fade, 1},
		{Hook, -1},
		{Slice, -1},
	}

	for _, tt := range tests {
		if got := tt.shape.DifficultyModifier(); got != tt.expected {
			t.Errorf("%v.DifficultyModifier() = %v, want %v", tt.shape, got, tt.expected)
		}
	}
}

func TestShotShape_DrawEasierThanStraight(t *testing.T) {
	if Draw.DifficultyModifier() <= Straight.DifficultyModifier() {
		t.Error("Draw should be easier (higher modifier) than Straight")
	}
}

func TestShotShape_FadeEasierThanStraight(t *testing.T) {
	if Fade.DifficultyModifier() <= Straight.DifficultyModifier() {
		t.Error("Fade should be easier (higher modifier) than Straight")
	}
}

func TestShotShape_HookHarderThanDraw(t *testing.T) {
	if Hook.DifficultyModifier() >= Draw.DifficultyModifier() {
		t.Error("Hook should be harder (lower modifier) than Draw")
	}
}

func TestShotShape_SliceHarderThanFade(t *testing.T) {
	if Slice.DifficultyModifier() >= Fade.DifficultyModifier() {
		t.Error("Slice should be harder (lower modifier) than Fade")
	}
}

func TestDetermineActualShape_SuccessReturnsIntendedShape(t *testing.T) {
	shapes := []ShotShape{Straight, Draw, Fade, Hook, Slice}
	successOutcomes := []SkillCheckOutcome{Marginal, Good, Excellent, CriticalSuccess}

	for _, shape := range shapes {
		for _, outcome := range successOutcomes {
			result := SkillCheckResult{Outcome: outcome, Margin: 1}
			random := rand.New(rand.NewPCG(1, 2))

			shapeResult := DetermineActualShape(shape, result, random)

			if shapeResult.Actual != shape {
				t.Errorf("DetermineActualShape(%v, %v) Actual = %v, want %v",
					shape, outcome, shapeResult.Actual, shape)
			}
			if !shapeResult.Success {
				t.Errorf("DetermineActualShape(%v, %v) Success = false, want true",
					shape, outcome)
			}
			if shapeResult.Intended != shape {
				t.Errorf("DetermineActualShape(%v, %v) Intended = %v, want %v",
					shape, outcome, shapeResult.Intended, shape)
			}
		}
	}
}

func TestDetermineActualShape_FailedDrawBecomesHookOrStraight(t *testing.T) {
	result := SkillCheckResult{Outcome: Poor, Margin: -2}

	hookCount := 0
	straightCount := 0

	for seed := uint64(0); seed < 100; seed++ {
		random := rand.New(rand.NewPCG(seed, seed))
		shapeResult := DetermineActualShape(Draw, result, random)

		if !shapeResult.Success {
			if shapeResult.Actual == Hook {
				hookCount++
			} else if shapeResult.Actual == Straight {
				straightCount++
			} else {
				t.Errorf("Failed Draw should become Hook or Straight, got %v", shapeResult.Actual)
			}
		}
	}

	if hookCount < 60 {
		t.Errorf("Failed Draw should become Hook ~80%% of the time, got %d/100", hookCount)
	}
	if straightCount < 10 {
		t.Errorf("Failed Draw should become Straight ~20%% of the time, got %d/100", straightCount)
	}
}

func TestDetermineActualShape_FailedFadeBecomesSliceOrStraight(t *testing.T) {
	result := SkillCheckResult{Outcome: Poor, Margin: -2}

	sliceCount := 0
	straightCount := 0

	for seed := uint64(0); seed < 100; seed++ {
		random := rand.New(rand.NewPCG(seed, seed))
		shapeResult := DetermineActualShape(Fade, result, random)

		if !shapeResult.Success {
			if shapeResult.Actual == Slice {
				sliceCount++
			} else if shapeResult.Actual == Straight {
				straightCount++
			} else {
				t.Errorf("Failed Fade should become Slice or Straight, got %v", shapeResult.Actual)
			}
		}
	}

	if sliceCount < 60 {
		t.Errorf("Failed Fade should become Slice ~80%% of the time, got %d/100", sliceCount)
	}
	if straightCount < 10 {
		t.Errorf("Failed Fade should become Straight ~20%% of the time, got %d/100", straightCount)
	}
}

func TestDetermineActualShape_FailedStraightBecomesDrawOrFade(t *testing.T) {
	result := SkillCheckResult{Outcome: Poor, Margin: -2}

	drawCount := 0
	fadeCount := 0

	for seed := uint64(0); seed < 100; seed++ {
		random := rand.New(rand.NewPCG(seed, seed))
		shapeResult := DetermineActualShape(Straight, result, random)

		if !shapeResult.Success {
			if shapeResult.Actual == Draw {
				drawCount++
			} else if shapeResult.Actual == Fade {
				fadeCount++
			} else {
				t.Errorf("Failed Straight should become Draw or Fade, got %v", shapeResult.Actual)
			}
		}
	}

	if drawCount < 30 {
		t.Errorf("Failed Straight should become Draw ~50%% of the time, got %d/100", drawCount)
	}
	if fadeCount < 30 {
		t.Errorf("Failed Straight should become Fade ~50%% of the time, got %d/100", fadeCount)
	}
}

func TestDetermineActualShape_FailedHookBecomesHookOrDraw(t *testing.T) {
	result := SkillCheckResult{Outcome: Poor, Margin: -2}

	hookCount := 0
	drawCount := 0

	for seed := uint64(0); seed < 100; seed++ {
		random := rand.New(rand.NewPCG(seed, seed))
		shapeResult := DetermineActualShape(Hook, result, random)

		if !shapeResult.Success {
			if shapeResult.Actual == Hook {
				hookCount++
			} else if shapeResult.Actual == Draw {
				drawCount++
			} else {
				t.Errorf("Failed Hook should become Hook or Draw, got %v", shapeResult.Actual)
			}
		}
	}

	if hookCount < 50 {
		t.Errorf("Failed Hook should stay Hook ~70%% of the time, got %d/100", hookCount)
	}
}

func TestDetermineActualShape_FailedSliceBecomesSliceOrFade(t *testing.T) {
	result := SkillCheckResult{Outcome: Poor, Margin: -2}

	sliceCount := 0
	fadeCount := 0

	for seed := uint64(0); seed < 100; seed++ {
		random := rand.New(rand.NewPCG(seed, seed))
		shapeResult := DetermineActualShape(Slice, result, random)

		if !shapeResult.Success {
			if shapeResult.Actual == Slice {
				sliceCount++
			} else if shapeResult.Actual == Fade {
				fadeCount++
			} else {
				t.Errorf("Failed Slice should become Slice or Fade, got %v", shapeResult.Actual)
			}
		}
	}

	if sliceCount < 50 {
		t.Errorf("Failed Slice should stay Slice ~70%% of the time, got %d/100", sliceCount)
	}
}
