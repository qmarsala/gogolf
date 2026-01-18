package ui

import (
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
}

// Test calculatePower returns correct power values
func TestPowerMeter_calculatePower(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	tests := []struct {
		elapsed  time.Duration
		expected float64
	}{
		{0, 0.0},
		{500 * time.Millisecond, 0.25},
		{1 * time.Second, 0.5},
		{1500 * time.Millisecond, 0.75},
		{2 * time.Second, 1.0},
		{3 * time.Second, 1.0}, // Should cap at max
	}

	for _, tt := range tests {
		power := pm.calculatePower(tt.elapsed)
		if power < tt.expected-0.01 || power > tt.expected+0.01 {
			t.Errorf("calculatePower(%v) = %f, want %f", tt.elapsed, power, tt.expected)
		}
	}
}

// Test drawMeterBar creates correct bar representation
func TestPowerMeter_drawMeterBar(t *testing.T) {
	renderer := NewRenderer()
	pm := NewPowerMeter(renderer)

	tests := []struct {
		power    float64
		expected string
	}{
		{0.0, "[                    ]"},
		{0.25, "[=====               ]"},
		{0.5, "[==========          ]"},
		{0.75, "[===============     ]"},
		{1.0, "[====================]"},
	}

	for _, tt := range tests {
		bar := pm.drawMeterBar(tt.power)
		if bar != tt.expected {
			t.Errorf("drawMeterBar(%f) = %q, want %q", tt.power, bar, tt.expected)
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

	pm.drawMeterBar(-0.1)
	pm.drawMeterBar(1.5)
	pm.drawMeterBar(0.0)
	pm.drawMeterBar(1.0)
}
