package gogolf

import (
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
