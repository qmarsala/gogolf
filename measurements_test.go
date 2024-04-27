package main

import (
	"testing"
)

func TestInchesToUnits(t *testing.T) {
	inches := float32(5)
	expectedUnits := float32(1)

	actual := Inch(inches).Units()
	if actual != Unit(expectedUnits) {
		t.Error("Expected 5 inches to equal 1 unit, but got", actual)
	}
}

func TestInchesToFeet(t *testing.T) {
	inches := float32(12)
	expectedFeet := float32(1)

	actual := Inch(inches).Feet()
	if actual != Foot(expectedFeet) {
		t.Error("Expected 12 inches to equal 1 feet, but got", actual)
	}
}

func TestInchesToYards(t *testing.T) {
	inches := float32(12)
	expectedYards := float32(0.33333334)

	actual := Inch(inches).Yards()
	if actual != Yard(expectedYards) {
		t.Error("Expected 12 inch to equal 0.33333334 yards, but got", actual)
	}
}

func TestUnitsToInches(t *testing.T) {
	units := float32(1)
	expectedInches := float32(5)

	actual := Unit(units).Inches()
	if actual != Inch(expectedInches) {
		t.Error("Expected 1 unit to equal 5 inches, but got", actual)
	}
}

func TestUnitsToFeet(t *testing.T) {
	units := float32(1)
	expectedFeet := 0.41666666

	actual := Unit(units).Feet()
	if actual != Foot(expectedFeet) {
		t.Error("Expected 1 unit to equal 0.41666666 feet, but got", actual)
	}
}

func TestUnitsToYard(t *testing.T) {
	units := float32(1)
	expectedYards := 0.13888888

	actual := Unit(units).Yards()
	if actual != Yard(expectedYards) {
		t.Error("Expected 1 unit to equal 0.13888888 yards, but got", actual)
	}
}

func TestYardsToUnits(t *testing.T) {
	yards := float32(10)
	expectedUnits := float32(72)

	actual := Yard(yards).Units()
	if actual != Unit(expectedUnits) {
		t.Error("Expected 10 yards to equal 72 units, but got", actual)
	}
}
