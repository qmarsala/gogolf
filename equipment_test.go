package gogolf

import "testing"

// Test Ball struct has required fields
func TestBall_Fields(t *testing.T) {
	ball := Ball{
		Name:          "Pro V1",
		DistanceBonus: 5.0,
		SpinControl:   0.8,
		Cost:          50,
	}

	if ball.Name != "Pro V1" {
		t.Errorf("Ball name = %s, want Pro V1", ball.Name)
	}
	if ball.DistanceBonus != 5.0 {
		t.Errorf("Ball DistanceBonus = %f, want 5.0", ball.DistanceBonus)
	}
	if ball.SpinControl != 0.8 {
		t.Errorf("Ball SpinControl = %f, want 0.8", ball.SpinControl)
	}
	if ball.Cost != 50 {
		t.Errorf("Ball Cost = %d, want 50", ball.Cost)
	}
}

// Test Glove struct has required fields
func TestGlove_Fields(t *testing.T) {
	glove := Glove{
		Name:          "Leather Pro",
		AccuracyBonus: 0.05,
		Cost:          30,
	}

	if glove.Name != "Leather Pro" {
		t.Errorf("Glove name = %s, want Leather Pro", glove.Name)
	}
	if glove.AccuracyBonus != 0.05 {
		t.Errorf("Glove AccuracyBonus = %f, want 0.05", glove.AccuracyBonus)
	}
	if glove.Cost != 30 {
		t.Errorf("Glove Cost = %d, want 30", glove.Cost)
	}
}

// Test Shoes struct has required fields
func TestShoes_Fields(t *testing.T) {
	shoes := Shoes{
		Name:               "Spiked Comfort",
		LiePenaltyReduction: 1,
		Cost:               40,
	}

	if shoes.Name != "Spiked Comfort" {
		t.Errorf("Shoes name = %s, want Spiked Comfort", shoes.Name)
	}
	if shoes.LiePenaltyReduction != 1 {
		t.Errorf("Shoes LiePenaltyReduction = %d, want 1", shoes.LiePenaltyReduction)
	}
	if shoes.Cost != 40 {
		t.Errorf("Shoes Cost = %d, want 40", shoes.Cost)
	}
}

// Test Golfer has equipment fields
func TestGolfer_EquipmentFields(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	// Check that equipment fields exist (may be nil initially)
	_ = golfer.Ball
	_ = golfer.Glove
	_ = golfer.Shoes
}

// Test Golfer can equip a ball
func TestGolfer_EquipBall(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	ball := &Ball{
		Name:          "Distance King",
		DistanceBonus: 10.0,
		SpinControl:   0.5,
		Cost:          75,
	}

	golfer.EquipBall(ball)

	if golfer.Ball == nil {
		t.Fatal("Golfer ball is nil after EquipBall")
	}
	if golfer.Ball.Name != "Distance King" {
		t.Errorf("Equipped ball name = %s, want Distance King", golfer.Ball.Name)
	}
}

// Test Golfer can equip a glove
func TestGolfer_EquipGlove(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	glove := &Glove{
		Name:          "Precision Grip",
		AccuracyBonus: 0.1,
		Cost:          50,
	}

	golfer.EquipGlove(glove)

	if golfer.Glove == nil {
		t.Fatal("Golfer glove is nil after EquipGlove")
	}
	if golfer.Glove.Name != "Precision Grip" {
		t.Errorf("Equipped glove name = %s, want Precision Grip", golfer.Glove.Name)
	}
}

// Test Golfer can equip shoes
func TestGolfer_EquipShoes(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	shoes := &Shoes{
		Name:               "All-Terrain Pro",
		LiePenaltyReduction: 2,
		Cost:               60,
	}

	golfer.EquipShoes(shoes)

	if golfer.Shoes == nil {
		t.Fatal("Golfer shoes are nil after EquipShoes")
	}
	if golfer.Shoes.Name != "All-Terrain Pro" {
		t.Errorf("Equipped shoes name = %s, want All-Terrain Pro", golfer.Shoes.Name)
	}
}

// Test equipment can be changed
func TestGolfer_ChangeEquipment(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	ball1 := &Ball{Name: "Ball 1", DistanceBonus: 5, SpinControl: 0.5, Cost: 30}
	ball2 := &Ball{Name: "Ball 2", DistanceBonus: 8, SpinControl: 0.7, Cost: 50}

	golfer.EquipBall(ball1)
	if golfer.Ball.Name != "Ball 1" {
		t.Errorf("First ball = %s, want Ball 1", golfer.Ball.Name)
	}

	golfer.EquipBall(ball2)
	if golfer.Ball.Name != "Ball 2" {
		t.Errorf("Second ball = %s, want Ball 2", golfer.Ball.Name)
	}
}

// Test GetEquippedBall returns the equipped ball
func TestGolfer_GetEquippedBall(t *testing.T) {
	golfer := NewGolfer("TestPlayer")
	ball := &Ball{Name: "Test Ball", DistanceBonus: 5, SpinControl: 0.6, Cost: 40}

	golfer.EquipBall(ball)

	equipped := golfer.GetEquippedBall()
	if equipped == nil {
		t.Fatal("GetEquippedBall returned nil")
	}
	if equipped.Name != "Test Ball" {
		t.Errorf("GetEquippedBall name = %s, want Test Ball", equipped.Name)
	}
}

// Test GetEquippedBall returns nil when no ball equipped
func TestGolfer_GetEquippedBall_NoEquipment(t *testing.T) {
	golfer := NewGolfer("TestPlayer")

	equipped := golfer.GetEquippedBall()
	if equipped != nil {
		t.Errorf("GetEquippedBall = %+v, want nil (no ball equipped)", equipped)
	}
}
