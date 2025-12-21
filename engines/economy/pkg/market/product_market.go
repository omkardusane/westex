package market

import (
	"westex/engines/economy/pkg/entities"
)

// Purchase represents a completed transaction
type Purchase struct {
	PersonID      int
	PersonName    string
	IndustryID    int
	IndustryName  string
	ProductID     int
	ProductName   string
	ProblemID     int
	ProblemSolved string
	Quantity      float32
	UnitPrice     float32
	TotalCost     float32
}

// MarketResult summarizes market activity for one tick
type MarketResult struct {
	Purchases         []Purchase
	TotalSpent        float32
	TotalRevenue      float32
	PeopleSatisfied   int
	PeopleUnsatisfied int
}

// ProcessProductMarket handles all purchases in one tick
func ProcessProductMarket(
	region *entities.Region,
	pricePerUnit float32,
) *MarketResult {
	result := &MarketResult{
		Purchases: make([]Purchase, 0),
	}

	satisfiedPeople := make(map[int]bool) // Track people who bought something

	// For each person
	for _, person := range region.People {
		// Get their needs (from all segments)
		needs := person.GetAllProblems()

		// Try to satisfy each need
		for _, need := range needs {
			// Find industries that solve this need
			industry := findIndustryForProblem(region, need)
			if industry == nil {
				continue
			}

			// Try to buy product
			purchase := attemptPurchase(person, industry, need, pricePerUnit)
			if purchase != nil {
				result.Purchases = append(result.Purchases, *purchase)
				result.TotalSpent += purchase.TotalCost
				result.TotalRevenue += purchase.TotalCost
				satisfiedPeople[person.ID] = true
			}
		}
	}

	// Count satisfied vs unsatisfied people
	result.PeopleSatisfied = len(satisfiedPeople)
	result.PeopleUnsatisfied = len(region.People) - result.PeopleSatisfied

	return result
}

// findIndustryForProblem finds the first industry that solves a given problem
func findIndustryForProblem(region *entities.Region, problem *entities.Problem) *entities.Industry {
	for _, industry := range region.Industries {
		for _, p := range industry.OwnedProblems {
			if p.ID == problem.ID {
				return industry
			}
		}
	}
	return nil
}

// attemptPurchase tries to make a purchase for a person
func attemptPurchase(
	person *entities.Person,
	industry *entities.Industry,
	need *entities.Problem,
	pricePerUnit float32,
) *Purchase {
	// Check if industry has products
	if len(industry.OutputProducts) == 0 {
		return nil
	}

	product := industry.OutputProducts[0] // Simplified: use first product

	// Check if product available
	if product.Quantity < 1.0 {
		return nil
	}

	// Check if person can afford
	if person.Money < pricePerUnit {
		return nil
	}

	// Make purchase
	quantity := float32(1.0) // Buy 1 unit
	cost := pricePerUnit * quantity

	// Transfer money
	person.Money -= cost
	industry.Money += cost

	// Transfer product
	product.Consume(quantity)

	return &Purchase{
		PersonID:      person.ID,
		PersonName:    person.Name,
		IndustryID:    industry.ID,
		IndustryName:  industry.Name,
		ProductID:     product.ID,
		ProductName:   product.Name,
		ProblemID:     need.ID,
		ProblemSolved: need.Name,
		Quantity:      quantity,
		UnitPrice:     pricePerUnit,
		TotalCost:     cost,
	}
}
