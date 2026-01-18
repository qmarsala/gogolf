package ui

import "testing"

// Test NewRenderer creates a renderer
func TestNewRenderer(t *testing.T) {
	renderer := NewRenderer()

	if renderer == nil {
		t.Fatal("NewRenderer returned nil")
	}

	if renderer.Terminal == nil {
		t.Error("Renderer Terminal is nil")
	}

	if renderer.Layout == nil {
		t.Error("Renderer Layout is nil")
	}
}

// Test Render doesn't panic with minimal state
func TestRenderer_Render_MinimalState(t *testing.T) {
	renderer := NewRenderer()

	state := GameState{
		PlayerName:   "TestPlayer",
		Money:        100,
		HoleNumber:   1,
		TotalHoles:   9,
		Par:          4,
		HoleDistance: 350,
		BallLie:      "Tee",
		DistanceToHole: 350,
		PromptMsg:    "Enter club selection:",
		Skills:       make(map[string]SkillDisplay),
		Abilities:    make(map[string]AbilityDisplay),
	}

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Render panicked: %v", r)
		}
	}()

	renderer.Render(state)
}

// Test Render with full state
func TestRenderer_Render_FullState(t *testing.T) {
	renderer := NewRenderer()

	state := GameState{
		PlayerName:        "TestPlayer",
		Money:             150,
		HoleNumber:        1,
		TotalHoles:        9,
		Par:               4,
		HoleDistance:      350,
		BallLie:           "Fairway",
		BallLieDifficulty: 0,
		DistanceToHole:    150,
		BallLocationX:     200,
		BallLocationY:     50,
		HoleLocationX:     350,
		HoleLocationY:     0,
		TotalStrokes:      1,
		ScoreToPar:        0,
		StrokesThisHole:   1,
		PromptMsg:         "Enter club selection:",
		LastShot: &ShotDisplay{
			ClubName:    "Driver",
			Outcome:     "Good",
			Margin:      2,
			Description: "Good contact. Ball flights well.",
			Rotation:    5.2,
			RotationDir: "right",
			Power:       0.95,
			Distance:    265,
			XPEarned:    2,
			LevelUps:    []string{},
		},
		Skills: map[string]SkillDisplay{
			"Driver": {Name: "Driver", Level: 1, Value: 5, CurrentXP: 2, XPForNext: 10},
		},
		Abilities: map[string]AbilityDisplay{
			"Strength": {Name: "Strength", Level: 1, Value: 5, CurrentXP: 2, XPForNext: 10},
		},
		Equipment: EquipmentDisplay{
			BallName:  "Standard Ball",
			BallBonus: "+3 dist",
		},
	}

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Render panicked: %v", r)
		}
	}()

	renderer.Render(state)
}

// Test RenderBorder doesn't panic
func TestRenderer_RenderBorder(t *testing.T) {
	renderer := NewRenderer()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderBorder panicked: %v", r)
		}
	}()

	renderer.RenderBorder()
}

// Test RenderLeftPanel doesn't panic
func TestRenderer_RenderLeftPanel(t *testing.T) {
	renderer := NewRenderer()

	state := GameState{
		HoleNumber:     1,
		Par:            4,
		HoleDistance:   350,
		BallLie:        "Tee",
		DistanceToHole: 350,
		PromptMsg:      "Enter club:",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderLeftPanel panicked: %v", r)
		}
	}()

	renderer.RenderLeftPanel(state)
}

// Test RenderRightPanel doesn't panic
func TestRenderer_RenderRightPanel(t *testing.T) {
	renderer := NewRenderer()

	state := GameState{
		PlayerName:     "TestPlayer",
		Money:          100,
		TotalStrokes:   0,
		ScoreToPar:     0,
		StrokesThisHole: 0,
		HoleNumber:     1,
		TotalHoles:     9,
		Skills:         make(map[string]SkillDisplay),
		Abilities:      make(map[string]AbilityDisplay),
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderRightPanel panicked: %v", r)
		}
	}()

	renderer.RenderRightPanel(state)
}

// Test printInPanel doesn't panic with long text
func TestRenderer_printInPanel_LongText(t *testing.T) {
	renderer := NewRenderer()
	panel := renderer.Layout.LeftPanel

	longText := "This is a very long text that should be truncated to fit within the panel width without causing any issues or panics"

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printInPanel panicked with long text: %v", r)
		}
	}()

	renderer.printInPanel(panel, 1, longText, false)
}

// Test printInPanel with centered text
func TestRenderer_printInPanel_Centered(t *testing.T) {
	renderer := NewRenderer()
	panel := renderer.Layout.LeftPanel

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printInPanel panicked with centered text: %v", r)
		}
	}()

	renderer.printInPanel(panel, 1, "Centered Text", true)
}

func TestShotDisplay_HasTargetNumber(t *testing.T) {
	shot := ShotDisplay{
		ClubName:     "Driver",
		Outcome:      "Good",
		Margin:       2,
		TargetNumber: 12,
		DiceRolls:    []int{3, 4, 3},
	}

	if shot.TargetNumber != 12 {
		t.Errorf("ShotDisplay.TargetNumber = %d, want 12", shot.TargetNumber)
	}
}
