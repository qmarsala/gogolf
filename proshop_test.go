package gogolf

import (
	"testing"
)

func TestProShop_HasInventory(t *testing.T) {
	shop := ProShop{
		Balls:  []Ball{},
		Gloves: []Glove{},
		Shoes:  []Shoes{},
	}

	if shop.Balls == nil {
		t.Error("ProShop Balls is nil")
	}
	if shop.Gloves == nil {
		t.Error("ProShop Gloves is nil")
	}
	if shop.Shoes == nil {
		t.Error("ProShop Shoes is nil")
	}
}

func TestNewProShop(t *testing.T) {
	shop := NewProShop()

	if len(shop.Balls) == 0 {
		t.Error("NewProShop has no balls in inventory")
	}
	if len(shop.Gloves) == 0 {
		t.Error("NewProShop has no gloves in inventory")
	}
	if len(shop.Shoes) == 0 {
		t.Error("NewProShop has no shoes in inventory")
	}
}

func TestNewProShop_BallVariety(t *testing.T) {
	shop := NewProShop()

	if len(shop.Balls) < 3 {
		t.Errorf("ProShop has %d balls, want at least 3 options", len(shop.Balls))
	}

	foundCheap := false
	foundExpensive := false
	for _, ball := range shop.Balls {
		if ball.Cost <= 30 {
			foundCheap = true
		}
		if ball.Cost >= 50 {
			foundExpensive = true
		}
	}

	if !foundCheap {
		t.Error("ProShop should have at least one affordable ball (<= 30)")
	}
	if !foundExpensive {
		t.Error("ProShop should have at least one premium ball (>= 50)")
	}
}

func TestNewProShop_GloveVariety(t *testing.T) {
	shop := NewProShop()

	if len(shop.Gloves) < 2 {
		t.Errorf("ProShop has %d gloves, want at least 2 options", len(shop.Gloves))
	}
}

func TestNewProShop_ShoeVariety(t *testing.T) {
	shop := NewProShop()

	if len(shop.Shoes) < 2 {
		t.Errorf("ProShop has %d shoes, want at least 2 options", len(shop.Shoes))
	}
}

func TestProShop_PurchaseBall_Success(t *testing.T) {
	shop := NewProShop()
	golfer := NewGolfer("TestPlayer")
	golfer.Money = 100

	var targetBall *Ball
	for i := range shop.Balls {
		if shop.Balls[i].Cost <= 100 {
			targetBall = &shop.Balls[i]
			break
		}
	}

	if targetBall == nil {
		t.Fatal("No affordable ball found in shop")
	}

	initialMoney := golfer.Money
	success := shop.PurchaseBall(&golfer, targetBall.Name)

	if !success {
		t.Error("PurchaseBall failed, want success")
	}

	expectedMoney := initialMoney - targetBall.Cost
	if golfer.Money != expectedMoney {
		t.Errorf("After purchase, money = %d, want %d", golfer.Money, expectedMoney)
	}

	if golfer.Ball == nil {
		t.Fatal("Ball not equipped after purchase")
	}
	if golfer.Ball.Name != targetBall.Name {
		t.Errorf("Equipped ball = %s, want %s", golfer.Ball.Name, targetBall.Name)
	}
}

func TestProShop_PurchaseBall_InsufficientFunds(t *testing.T) {
	shop := NewProShop()
	golfer := NewGolfer("TestPlayer")
	golfer.Money = 10

	var expensiveBall *Ball
	for i := range shop.Balls {
		if shop.Balls[i].Cost > 10 {
			expensiveBall = &shop.Balls[i]
			break
		}
	}

	if expensiveBall == nil {
		t.Skip("No expensive ball found in shop")
	}

	initialMoney := golfer.Money
	success := shop.PurchaseBall(&golfer, expensiveBall.Name)

	if success {
		t.Error("PurchaseBall succeeded with insufficient funds, want failure")
	}

	if golfer.Money != initialMoney {
		t.Errorf("Money changed after failed purchase: %d, want %d", golfer.Money, initialMoney)
	}
}

func TestProShop_PurchaseBall_NotFound(t *testing.T) {
	shop := NewProShop()
	golfer := NewGolfer("TestPlayer")
	golfer.Money = 1000

	success := shop.PurchaseBall(&golfer, "NonExistent Ball")

	if success {
		t.Error("PurchaseBall succeeded for non-existent ball, want failure")
	}
}

func TestProShop_PurchaseGlove_Success(t *testing.T) {
	shop := NewProShop()
	golfer := NewGolfer("TestPlayer")
	golfer.Money = 100

	var targetGlove *Glove
	for i := range shop.Gloves {
		if shop.Gloves[i].Cost <= 100 {
			targetGlove = &shop.Gloves[i]
			break
		}
	}

	if targetGlove == nil {
		t.Fatal("No affordable glove found in shop")
	}

	success := shop.PurchaseGlove(&golfer, targetGlove.Name)

	if !success {
		t.Error("PurchaseGlove failed, want success")
	}

	if golfer.Glove == nil {
		t.Fatal("Glove not equipped after purchase")
	}
}

func TestProShop_PurchaseShoes_Success(t *testing.T) {
	shop := NewProShop()
	golfer := NewGolfer("TestPlayer")
	golfer.Money = 100

	var targetShoes *Shoes
	for i := range shop.Shoes {
		if shop.Shoes[i].Cost <= 100 {
			targetShoes = &shop.Shoes[i]
			break
		}
	}

	if targetShoes == nil {
		t.Fatal("No affordable shoes found in shop")
	}

	success := shop.PurchaseShoes(&golfer, targetShoes.Name)

	if !success {
		t.Error("PurchaseShoes failed, want success")
	}

	if golfer.Shoes == nil {
		t.Fatal("Shoes not equipped after purchase")
	}
}
