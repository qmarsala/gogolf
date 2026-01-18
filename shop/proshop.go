package shop

import "gogolf"

type ProShop struct {
	Balls  []gogolf.Ball
	Gloves []gogolf.Glove
	Shoes  []gogolf.Shoes
}

func NewProShop() ProShop {
	return ProShop{
		Balls: []gogolf.Ball{
			{
				Name:          "Budget Ball",
				DistanceBonus: 0,
				SpinControl:   0.3,
				Cost:          20,
			},
			{
				Name:          "Standard Ball",
				DistanceBonus: 3,
				SpinControl:   0.5,
				Cost:          35,
			},
			{
				Name:          "Premium Ball",
				DistanceBonus: 5,
				SpinControl:   0.7,
				Cost:          50,
			},
			{
				Name:          "Pro V1",
				DistanceBonus: 8,
				SpinControl:   0.9,
				Cost:          75,
			},
		},
		Gloves: []gogolf.Glove{
			{
				Name:          "Basic Glove",
				AccuracyBonus: 0.02,
				Cost:          25,
			},
			{
				Name:          "Leather Pro",
				AccuracyBonus: 0.05,
				Cost:          45,
			},
			{
				Name:          "Precision Grip",
				AccuracyBonus: 0.08,
				Cost:          65,
			},
		},
		Shoes: []gogolf.Shoes{
			{
				Name:                "Casual Spikes",
				LiePenaltyReduction: 1,
				Cost:                30,
			},
			{
				Name:                "All-Terrain Pro",
				LiePenaltyReduction: 2,
				Cost:                55,
			},
			{
				Name:                "Tour Edition",
				LiePenaltyReduction: 3,
				Cost:                80,
			},
		},
	}
}

func (shop ProShop) PurchaseBall(golfer *gogolf.Golfer, ballName string) bool {
	var targetBall *gogolf.Ball
	for i := range shop.Balls {
		if shop.Balls[i].Name == ballName {
			targetBall = &shop.Balls[i]
			break
		}
	}

	if targetBall == nil {
		return false
	}

	if !golfer.SpendMoney(targetBall.Cost) {
		return false
	}

	golfer.EquipBall(targetBall)
	return true
}

func (shop ProShop) PurchaseGlove(golfer *gogolf.Golfer, gloveName string) bool {
	var targetGlove *gogolf.Glove
	for i := range shop.Gloves {
		if shop.Gloves[i].Name == gloveName {
			targetGlove = &shop.Gloves[i]
			break
		}
	}

	if targetGlove == nil {
		return false
	}

	if !golfer.SpendMoney(targetGlove.Cost) {
		return false
	}

	golfer.EquipGlove(targetGlove)
	return true
}

func (shop ProShop) PurchaseShoes(golfer *gogolf.Golfer, shoesName string) bool {
	var targetShoes *gogolf.Shoes
	for i := range shop.Shoes {
		if shop.Shoes[i].Name == shoesName {
			targetShoes = &shop.Shoes[i]
			break
		}
	}

	if targetShoes == nil {
		return false
	}

	if !golfer.SpendMoney(targetShoes.Cost) {
		return false
	}

	golfer.EquipShoes(targetShoes)
	return true
}
