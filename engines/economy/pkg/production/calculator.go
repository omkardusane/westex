package production

import "westex/engines/economy/pkg/entities"

// ProductionResult contains the outcome of production calculation
type ProductionResult struct {
	UnitsProduced float32
	LaborUsed     float32
	LaborCost     float32
	ResourceCost  float32
	TotalCost     float32
	CostPerUnit   float32
}

// CalculateProduction determines how much can be produced given constraints
func CalculateProduction(
	industry *entities.Industry,
	availableLabor float32,
	availableHours float32,
	wageRate float32,
) *ProductionResult {
	result := &ProductionResult{}

	// Calculate labor utilization
	laborNeeded := industry.LaborNeeded
	laborUsed := min(availableLabor, laborNeeded)
	result.LaborUsed = laborUsed

	// Calculate production capacity (what % of full capacity)
	productionRate := laborUsed / laborNeeded
	if laborNeeded == 0 {
		productionRate = 0
	}

	// Units produced: production rate × available hours
	// Simplified: 1 unit per hour of effective labor
	result.UnitsProduced = productionRate * availableHours

	// Calculate costs
	result.LaborCost = laborUsed * wageRate * availableHours
	result.ResourceCost = calculateResourceCost(industry, result.UnitsProduced)
	result.TotalCost = result.LaborCost + result.ResourceCost

	if result.UnitsProduced > 0 {
		result.CostPerUnit = result.TotalCost / result.UnitsProduced
	}

	return result
}

// calculateResourceCost estimates the cost of resources consumed
func calculateResourceCost(industry *entities.Industry, unitsProduced float32) float32 {
	totalCost := float32(0)

	// Simplified: each input resource costs 1.0 per unit consumed
	// In future, this will use actual market prices
	for _, input := range industry.InputResources {
		// Assume 1:1 ratio: 1 unit of input → 1 unit of output
		unitsNeeded := unitsProduced
		costPerUnit := float32(1.0) // Default cost

		// Free resources (land, water) have no cost
		if input.IsFree {
			costPerUnit = 0
		}

		totalCost += unitsNeeded * costPerUnit
	}

	return totalCost
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
