package ui

import (
	"fmt"
	"gogolf"
	"math/rand"
	"time"
)

// PowerMeter manages the spacebar-based power input
type PowerMeter struct {
	renderer         *Renderer
	maxPower         float64
	maxTime          time.Duration
	sweetSpotStart   float64 // Sweet spot starts at 75% of bar
	sweetSpotEnd     float64 // Sweet spot ends at 85% of bar
	clubMaxDistance  float64 // Max distance of current club in yards (or feet for putting)
	isPutting        bool
	puttDistanceFeet float64
	putterMaxYards   float64 // Putter's max distance in yards for power conversion
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
				projectedDist := pm.calculateProjectedDistance(power)

				// Draw power meter with projected distance
				pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
				meterBar := pm.drawMeterBar(elapsed, false)
				fmt.Printf("Power: %s %s   ", meterBar, pm.formatDistanceDisplay(power, projectedDist))

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
	finalDist := pm.calculateProjectedDistance(finalPower)

	// Show final position with marker
	pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
	meterBar := pm.drawMeterBar(finalElapsed, true)
	fmt.Printf("Power: %s %s   ", meterBar, pm.formatDistanceDisplay(finalPower, finalDist))

	// Brief pause to show result
	time.Sleep(500 * time.Millisecond)

	// Clear the meter display
	pm.renderer.Terminal.MoveCursor(meterRow, panel.X+2)
	fmt.Print("                                                        ")

	if pm.isPutting && pm.putterMaxYards > 0 {
		selectedFeet := finalDist
		selectedYards := selectedFeet / 3.0
		return selectedYards / pm.putterMaxYards
	}

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

// SetClubDistance sets the max distance of the current club
func (pm *PowerMeter) SetClubDistance(distance float64) {
	pm.clubMaxDistance = distance
	pm.isPutting = false
}

// SetPuttingMode configures the meter for putting with auto-scaled distance
// clubMaxDistanceYards is the putter's max distance in yards (used for power calculation)
func (pm *PowerMeter) SetPuttingMode(distanceFeet float64) {
	pm.isPutting = true
	pm.puttDistanceFeet = distanceFeet

	maxDistance := distanceFeet * 1.5
	if maxDistance < 10 {
		maxDistance = 10
	}
	pm.clubMaxDistance = maxDistance
}

// SetPuttingModeWithClubDistance configures putting mode and stores club distance for power conversion
func (pm *PowerMeter) SetPuttingModeWithClubDistance(distanceFeet float64, clubMaxDistanceYards float64) {
	pm.SetPuttingMode(distanceFeet)
	pm.putterMaxYards = clubMaxDistanceYards
}

// ClearPuttingMode resets the meter to normal mode
func (pm *PowerMeter) ClearPuttingMode() {
	pm.isPutting = false
	pm.puttDistanceFeet = 0
}

// calculateProjectedDistance returns the projected shot distance based on power
func (pm *PowerMeter) calculateProjectedDistance(power float64) float64 {
	return pm.clubMaxDistance * power
}

// formatDistanceDisplay creates a string showing power percentage and projected distance
func (pm *PowerMeter) formatDistanceDisplay(power float64, projectedDistance float64) string {
	if pm.isPutting {
		return fmt.Sprintf("%.0f%% | %.0f feet", power*100, projectedDistance)
	}
	return fmt.Sprintf("%.0f%% | %.0f yards", power*100, projectedDistance)
}

// ShotShapeSelector manages shot shape selection in the game UI
type ShotShapeSelector struct {
	renderer *Renderer
}

// NewShotShapeSelector creates a shot shape selector
func NewShotShapeSelector(renderer *Renderer) *ShotShapeSelector {
	return &ShotShapeSelector{renderer: renderer}
}

// SelectShotShape displays shape options and returns selected shape
// Default is Straight if user just presses Enter or space
func (s *ShotShapeSelector) SelectShotShape() gogolf.ShotShape {
	panel := s.renderer.Layout.LeftPanel
	row := panel.Height - 5

	s.renderer.Terminal.MoveCursor(row, panel.X+2)
	fmt.Print("Shot shape: [1]Straight [2]Draw [3]Fade [4]Hook [5]Slice")
	s.renderer.Terminal.MoveCursor(row+1, panel.X+2)
	fmt.Print("Press 1-5 or Enter for Straight:                        ")

	key := s.waitForShapeKey()

	s.renderer.Terminal.MoveCursor(row, panel.X+2)
	fmt.Print("                                                        ")
	s.renderer.Terminal.MoveCursor(row+1, panel.X+2)
	fmt.Print("                                                        ")

	switch key {
	case '2':
		return gogolf.Draw
	case '3':
		return gogolf.Fade
	case '4':
		return gogolf.Hook
	case '5':
		return gogolf.Slice
	default:
		return gogolf.Straight
	}
}

// waitForShapeKey waits for a valid shape selection key
func (s *ShotShapeSelector) waitForShapeKey() byte {
	for {
		key := readSingleKey()
		if key == ' ' || key == '\r' || key == '\n' || key == '1' {
			return '1'
		}
		if key >= '2' && key <= '5' {
			return key
		}
	}
}

// DiceRoller displays an animated dice rolling effect
type DiceRoller struct {
	renderer *Renderer
}

// NewDiceRoller creates a dice roller UI component
func NewDiceRoller(renderer *Renderer) *DiceRoller {
	return &DiceRoller{renderer: renderer}
}

// ShowRoll displays an animated dice roll, stopping each die one at a time
func (dr *DiceRoller) ShowRoll(finalRolls []int, targetNumber int) {
	if len(finalRolls) != 3 {
		return
	}

	panel := dr.renderer.Layout.LeftPanel
	row := panel.Height - 6

	dr.renderer.Terminal.MoveCursor(row-1, panel.X+2)
	fmt.Printf("Target: %d                              ", targetNumber)

	stopped := [3]bool{false, false, false}
	displayed := [3]int{1, 1, 1}

	animationDuration := 1500 * time.Millisecond
	stopInterval := animationDuration / 4
	frameInterval := 50 * time.Millisecond

	startTime := time.Now()
	ticker := time.NewTicker(frameInterval)
	defer ticker.Stop()

	for {
		elapsed := time.Since(startTime)

		if elapsed > stopInterval && !stopped[0] {
			stopped[0] = true
			displayed[0] = finalRolls[0]
		}
		if elapsed > stopInterval*2 && !stopped[1] {
			stopped[1] = true
			displayed[1] = finalRolls[1]
		}
		if elapsed > stopInterval*3 && !stopped[2] {
			stopped[2] = true
			displayed[2] = finalRolls[2]
		}

		for i := 0; i < 3; i++ {
			if !stopped[i] {
				displayed[i] = rand.Intn(6) + 1
			}
		}

		dr.renderer.Terminal.MoveCursor(row, panel.X+2)
		dr.drawDice(displayed, stopped)

		if stopped[0] && stopped[1] && stopped[2] {
			break
		}

		<-ticker.C
	}

	total := finalRolls[0] + finalRolls[1] + finalRolls[2]
	dr.renderer.Terminal.MoveCursor(row+1, panel.X+2)
	fmt.Printf("Total: %d (Target: %d)                  ", total, targetNumber)

	time.Sleep(500 * time.Millisecond)
}

// drawDice renders the three dice with their current values
func (dr *DiceRoller) drawDice(values [3]int, stopped [3]bool) {
	fmt.Print("Dice: ")
	for i, val := range values {
		if stopped[i] {
			fmt.Printf("%s[%d]%s ", colorBrightWhite, val, colorReset)
		} else {
			fmt.Printf("%s[%d]%s ", colorDim, val, colorReset)
		}
	}
	fmt.Print("      ")
}

// ClearDiceDisplay clears the dice roll display area
func (dr *DiceRoller) ClearDiceDisplay() {
	panel := dr.renderer.Layout.LeftPanel
	row := panel.Height - 6
	dr.renderer.Terminal.MoveCursor(row-1, panel.X+2)
	fmt.Print("                                        ")
	dr.renderer.Terminal.MoveCursor(row, panel.X+2)
	fmt.Print("                                        ")
	dr.renderer.Terminal.MoveCursor(row+1, panel.X+2)
	fmt.Print("                                        ")
}
