package production

import (
	"fmt"
	"westex/engines/economy/pkg/entities"
)

// ResourceConsumption tracks resources used in production
type ResourceConsumption struct {
	ResourceName string
	Quantity     float32
	Cost         float32
}

// ConsumeResources deducts input resources needed for production
func ConsumeResources(
	industry *entities.Industry,
	unitsToProdu float32,
) ([]ResourceConsumption, error) {
	consumptions := make([]ResourceConsumption, 0)

	// For each input resource
	for _, input := range industry.InputResources {
		// Calculate how much needed
		// Simplified: 1 unit of input → 1 unit of output
		needed := unitsToProdu

		// Check availability
		if input.Quantity < needed {
			return nil, fmt.Errorf("insufficient %s: need %.2f, have %.2f",
				input.Name, needed, input.Quantity)
		}

		// Consume
		success := input.Consume(needed)
		if !success {
			return nil, fmt.Errorf("failed to consume %s", input.Name)
		}

		// Calculate cost
		costPerUnit := float32(1.0) // Default cost

		// Free resources have no cost
		if input.IsFree {
			costPerUnit = 0
		}

		consumptions = append(consumptions, ResourceConsumption{
			ResourceName: input.Name,
			Quantity:     needed,
			Cost:         needed * costPerUnit,
		})
	}

	return consumptions, nil
}

// RegenerateResources adds regeneration to renewable resources
func RegenerateResources(resources []*entities.Resource) {
	for _, resource := range resources {
		if resource.RegenerationRate > 0 {
			resource.Add(resource.RegenerationRate)
		}
	}
}
