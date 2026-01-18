package ui

import (
	"fmt"
	"time"
)

// PowerMeter manages the spacebar-based power input
type PowerMeter struct {
	renderer *Renderer
	maxPower float64
	maxTime  time.Duration
}

// NewPowerMeter creates a power meter with default settings
func NewPowerMeter(renderer *Renderer) *PowerMeter {
	return &PowerMeter{
		renderer: renderer,
		maxPower: 1.0,
		maxTime:  2 * time.Second, // 2 seconds for full power
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
	go func() {
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(startTime)
				power := pm.calculatePower(elapsed)

				// Draw power meter
				pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
				meterBar := pm.drawMeterBar(power)
				fmt.Printf("Power: %s %.0f%%   ", meterBar, power*100)

				// Check if max time reached
				if elapsed >= pm.maxTime {
					powerChan <- pm.maxPower
					return
				}
			case <-stopChan:
				elapsed := time.Since(startTime)
				power := pm.calculatePower(elapsed)
				powerChan <- power
				return
			}
		}
	}()

	// Wait for power value
	finalPower := <-powerChan

	// Clear the meter display
	pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
	fmt.Print("                                                        ")

	return finalPower
}

// calculatePower converts elapsed time to power value (0.0 to 1.0)
func (pm *PowerMeter) calculatePower(elapsed time.Duration) float64 {
	if elapsed >= pm.maxTime {
		return pm.maxPower
	}

	ratio := float64(elapsed) / float64(pm.maxTime)
	return ratio * pm.maxPower
}

// drawMeterBar creates a visual power bar
func (pm *PowerMeter) drawMeterBar(power float64) string {
	barWidth := 20
	filledWidth := int(power * float64(barWidth))

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filledWidth {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"

	return bar
}
