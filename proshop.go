package gogolf

// ProShop represents the equipment shop where players can purchase items
type ProShop struct {
	Balls  []Ball
	Gloves []Glove
	Shoes  []Shoes
}

// NewProShop creates a ProShop with starter inventory
func NewProShop() ProShop {
	return ProShop{
		Balls: []Ball{
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
		Gloves: []Glove{
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
		Shoes: []Shoes{
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

// PurchaseBall attempts to purchase and equip a ball by name
// Returns true if successful, false if insufficient funds or item not found
func (shop ProShop) PurchaseBall(golfer *Golfer, ballName string) bool {
	// Find the ball in inventory
	var targetBall *Ball
	for i := range shop.Balls {
		if shop.Balls[i].Name == ballName {
			targetBall = &shop.Balls[i]
			break
		}
	}

	if targetBall == nil {
		return false // Ball not found
	}

	// Attempt to spend money
	if !golfer.SpendMoney(targetBall.Cost) {
		return false // Insufficient funds
	}

	// Equip the ball
	golfer.EquipBall(targetBall)
	return true
}

// PurchaseGlove attempts to purchase and equip a glove by name
// Returns true if successful, false if insufficient funds or item not found
func (shop ProShop) PurchaseGlove(golfer *Golfer, gloveName string) bool {
	// Find the glove in inventory
	var targetGlove *Glove
	for i := range shop.Gloves {
		if shop.Gloves[i].Name == gloveName {
			targetGlove = &shop.Gloves[i]
			break
		}
	}

	if targetGlove == nil {
		return false // Glove not found
	}

	// Attempt to spend money
	if !golfer.SpendMoney(targetGlove.Cost) {
		return false // Insufficient funds
	}

	// Equip the glove
	golfer.EquipGlove(targetGlove)
	return true
}

// PurchaseShoes attempts to purchase and equip shoes by name
// Returns true if successful, false if insufficient funds or item not found
func (shop ProShop) PurchaseShoes(golfer *Golfer, shoesName string) bool {
	// Find the shoes in inventory
	var targetShoes *Shoes
	for i := range shop.Shoes {
		if shop.Shoes[i].Name == shoesName {
			targetShoes = &shop.Shoes[i]
			break
		}
	}

	if targetShoes == nil {
		return false // Shoes not found
	}

	// Attempt to spend money
	if !golfer.SpendMoney(targetShoes.Cost) {
		return false // Insufficient funds
	}

	// Equip the shoes
	golfer.EquipShoes(targetShoes)
	return true
}
