package ui

import (
	"fmt"
	"time"
)

// PowerMeter manages the spacebar-based power input
type PowerMeter struct {
	renderer       *Renderer
	maxPower       float64
	maxTime        time.Duration
	sweetSpotStart float64 // Sweet spot starts at 75% of bar
	sweetSpotEnd   float64 // Sweet spot ends at 85% of bar
}

// NewPowerMeter creates a power meter with default settings
func NewPowerMeter(renderer *Renderer) *PowerMeter {
	return &PowerMeter{
		renderer:       renderer,
		maxPower:       1.0,
		maxTime:        2 * time.Second, // 2 seconds for max time
		sweetSpotStart: 0.75,             // Sweet spot at 75%-85%
		sweetSpotEnd:   0.85,
	}
}

// GetPower displays a power meter and waits for two spacebar presses
// Returns power value between 0.0 and 1.0 based on time between presses
func (pm *PowerMeter) GetPower() float64 {
	panel := pm.renderer.Layout.LeftPanel
	meterRow := panel.Height - 3

	// Display initial instruction
	pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
	fmt.Print("Press SPACE to start power meter...                    ")

	// Wait for first spacebar press
	pm.waitForSpacebar()

	// Start timer and show meter
	startTime := time.Now()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	stopChan := make(chan bool, 1)
	powerChan := make(chan float64, 1)

	// Goroutine to wait for second spacebar press
	go func() {
		pm.waitForSpacebar()
		stopChan <- true
	}()

	// Goroutine to update meter display
	var finalElapsed time.Duration
	go func() {
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(startTime)
				power := pm.calculatePower(elapsed)

				// Draw power meter (not stopped yet)
				pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
				meterBar := pm.drawMeterBar(elapsed, false)
				fmt.Printf("Power: %s %.0f%%   ", meterBar, power*100)

				// Check if max time reached
				if elapsed >= pm.maxTime {
					finalElapsed = elapsed
					powerChan <- pm.calculatePower(elapsed)
					return
				}
			case <-stopChan:
				elapsed := time.Since(startTime)
				finalElapsed = elapsed
				powerChan <- pm.calculatePower(elapsed)
				return
			}
		}
	}()

	// Wait for power value
	finalPower := <-powerChan

	// Show final position with marker
	pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
	meterBar := pm.drawMeterBar(finalElapsed, true)
	fmt.Printf("Power: %s %.0f%%   ", meterBar, finalPower*100)

	// Brief pause to show result
	time.Sleep(500 * time.Millisecond)

	// Clear the meter display
	pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
	fmt.Print("                                                        ")

	return finalPower
}

// calculatePower converts elapsed time to power value (0.0 to 1.0)
// with sweet spot mechanics: <75% = linear to 95%, 75-85% = 100%, >85% = decay to 50%
func (pm *PowerMeter) calculatePower(elapsed time.Duration) float64 {
	timeRatio := float64(elapsed) / float64(pm.maxTime)
	if timeRatio > 1.0 {
		timeRatio = 1.0
	}

	// Before sweet spot: linear increase to 95%
	if timeRatio < pm.sweetSpotStart {
		return (timeRatio / pm.sweetSpotStart) * 0.95
	}

	// In sweet spot: 100% power
	if timeRatio <= pm.sweetSpotEnd {
		return 1.0
	}

	// Past sweet spot: decay from 95% to 50%
	overshoot := (timeRatio - pm.sweetSpotEnd) / (1.0 - pm.sweetSpotEnd)
	return 0.95 - (overshoot * 0.45)
}

// drawMeterBar creates a visual power bar with sweet spot zone
func (pm *PowerMeter) drawMeterBar(elapsed time.Duration, stopped bool) string {
	barWidth := 20
	timeRatio := float64(elapsed) / float64(pm.maxTime)
	if timeRatio > 1.0 {
		timeRatio = 1.0
	}

	sweetStart := int(pm.sweetSpotStart * float64(barWidth))
	sweetEnd := int(pm.sweetSpotEnd * float64(barWidth))
	currentPos := int(timeRatio * float64(barWidth))

	bar := "["
	for i := 0; i < barWidth; i++ {
		// Sweet spot zone boundaries
		if i == sweetStart {
			bar += "("
			continue
		}
		if i == sweetEnd {
			bar += ")"
			continue
		}

		// Current position marker (only when stopped)
		if stopped && i == currentPos {
			bar += "o"
			continue
		}

		// Fill bar up to current position
		if i < currentPos {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"

	return bar
}
