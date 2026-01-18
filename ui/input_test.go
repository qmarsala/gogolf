package ui

import (
	"math"
	"testing"
	"time"
)

// Test NewPowerMeter creates a power meter
func TestNewPowerMeter(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	if pm == nil {
		t.Fatal("NewPowerMeter returned nil")
	}

	if pm.renderer == nil {
		t.Error("PowerMeter renderer is nil")
	}

	if pm.maxPower != 1.0 {
		t.Errorf("PowerMeter maxPower = %f, want 1.0", pm.maxPower)
	}

	if pm.maxTime != 2*time.Second {
		t.Errorf("PowerMeter maxTime = %v, want 2s", pm.maxTime)
	}

	if pm.sweetSpotStart != 0.75 {
		t.Errorf("PowerMeter sweetSpotStart = %f, want 0.75", pm.sweetSpotStart)
	}

	if pm.sweetSpotEnd != 0.85 {
		t.Errorf("PowerMeter sweetSpotEnd = %f, want 0.85", pm.sweetSpotEnd)
	}
}

// Test calculatePower returns correct power values with sweet spot
func TestPowerMeter_calculatePower(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	tests := []struct {
		elapsed  time.Duration
		expected float64
		desc     string
	}{
		{0, 0.0, "Start = 0%"},
		{500 * time.Millisecond, 0.317, "25% time -> 31.7% power"},
		{1000 * time.Millisecond, 0.633, "50% time -> 63.3% power"},
		{1499 * time.Millisecond, 0.949, "Just before zone = ~95%"},
		{1500 * time.Millisecond, 1.0, "Zone start = 100%"},
		{1600 * time.Millisecond, 1.0, "Zone middle = 100%"},
		{1700 * time.Millisecond, 1.0, "Zone end = 100%"},
		{1850 * time.Millisecond, 0.725, "Midway overshoot"},
		{2000 * time.Millisecond, 0.50, "Max overshoot = 50%"},
		{3000 * time.Millisecond, 0.50, "Beyond max = 50%"},
	}

	for _, tt := range tests {
		power := pm.calculatePower(tt.elapsed)
		if math.Abs(power-tt.expected) > 0.01 {
			t.Errorf("%s: calculatePower(%v) = %.3f, want %.3f", tt.desc, tt.elapsed, power, tt.expected)
		}
	}
}

// Test drawMeterBar creates correct bar representation with sweet spot
func TestPowerMeter_drawMeterBar(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	tests := []struct {
		elapsed  time.Duration
		stopped  bool
		expected string
		desc     string
	}{
		{
			0,
			false,
			"[               ( )  ]",
			"Start - empty bar with zone",
		},
		{
			500 * time.Millisecond,
			false,
			"[=====          ( )  ]",
			"25% time - before zone",
		},
		{
			1499 * time.Millisecond,
			false,
			"[============== ( )  ]",
			"Just before zone (not stopped)",
		},
		{
			1500 * time.Millisecond,
			true,
			"[===============( )  ]",
			"Zone start (stopped at boundary)",
		},
		{
			1600 * time.Millisecond,
			true,
			"[===============(o)  ]",
			"In zone (stopped) - 100%",
		},
		{
			1700 * time.Millisecond,
			true,
			"[===============(=)  ]",
			"Zone end (stopped)",
		},
		{
			2000 * time.Millisecond,
			true,
			"[===============(=)==]",
			"Max time (stopped) - overshoot",
		},
	}

	for _, tt := range tests {
		bar := pm.drawMeterBar(tt.elapsed, tt.stopped)
		if bar != tt.expected {
			t.Errorf("%s: drawMeterBar(%v, %v) = %q, want %q", tt.desc, tt.elapsed, tt.stopped, bar, tt.expected)
		}
	}
}

// Test drawMeterBar doesn't panic with edge cases
func TestPowerMeter_drawMeterBar_EdgeCases(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("drawMeterBar panicked: %v", r)
		}
	}()

	pm.drawMeterBar(0, false)
	pm.drawMeterBar(0, true)
	pm.drawMeterBar(3*time.Second, false)
	pm.drawMeterBar(3*time.Second, true)
}

func TestPowerMeter_SetClubDistance(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	pm.SetClubDistance(180)

	if pm.clubMaxDistance != 180 {
		t.Errorf("clubMaxDistance = %f, want 180", pm.clubMaxDistance)
	}
}

func TestPowerMeter_CalculateProjectedDistance(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)
	pm.SetClubDistance(200)

	tests := []struct {
		power            float64
		expectedDistance float64
		desc             string
	}{
		{0.0, 0.0, "0% power = 0 yards"},
		{0.5, 100.0, "50% power = 100 yards"},
		{1.0, 200.0, "100% power = 200 yards"},
		{0.75, 150.0, "75% power = 150 yards"},
	}

	for _, tt := range tests {
		distance := pm.calculateProjectedDistance(tt.power)
		if math.Abs(distance-tt.expectedDistance) > 0.01 {
			t.Errorf("%s: calculateProjectedDistance(%.2f) = %.2f, want %.2f",
				tt.desc, tt.power, distance, tt.expectedDistance)
		}
	}
}

func TestPowerMeter_FormatDistanceDisplay(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)
	pm.SetClubDistance(180)

	display := pm.formatDistanceDisplay(0.5, 90)
	if display == "" {
		t.Error("formatDistanceDisplay returned empty string")
	}

	if !containsSubstring(display, "90") {
		t.Errorf("display should contain projected distance, got: %s", display)
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (containsSubstring(s[1:], substr) || s[:len(substr)] == substr))
}

func TestPowerMeter_SetPuttingMode(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	pm.SetPuttingMode(15)

	if !pm.isPutting {
		t.Error("SetPuttingMode should set isPutting to true")
	}
	if pm.puttDistanceFeet != 15 {
		t.Errorf("puttDistanceFeet = %f, want 15", pm.puttDistanceFeet)
	}
}

func TestPowerMeter_PuttingModeScalesDistance(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	pm.SetPuttingMode(20)

	expectedMax := 30.0
	if pm.clubMaxDistance != expectedMax {
		t.Errorf("clubMaxDistance = %f, want %f (20 + 50%% buffer)", pm.clubMaxDistance, expectedMax)
	}
}

func TestPowerMeter_PuttingModeMinimumDistance(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	pm.SetPuttingMode(3)

	expectedMin := 10.0
	if pm.clubMaxDistance < expectedMin {
		t.Errorf("clubMaxDistance = %f, want at least %f (minimum)", pm.clubMaxDistance, expectedMin)
	}
}

func TestPowerMeter_PuttingModeDisplaysFeet(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)
	pm.SetPuttingMode(15)

	display := pm.formatDistanceDisplay(0.5, 11.25)

	if !containsSubstring(display, "feet") {
		t.Errorf("putting mode should display feet, got: %s", display)
	}
}

func TestPowerMeter_ClearPuttingMode(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	pm.SetPuttingMode(15)
	pm.ClearPuttingMode()

	if pm.isPutting {
		t.Error("ClearPuttingMode should set isPutting to false")
	}
}

func TestPowerMeter_SetPuttingModeWithClubDistance(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	pm.SetPuttingModeWithClubDistance(15, 40)

	if pm.putterMaxYards != 40 {
		t.Errorf("putterMaxYards = %f, want 40", pm.putterMaxYards)
	}
	if !pm.isPutting {
		t.Error("isPutting should be true")
	}
}
