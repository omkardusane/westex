package main

import (
	"fmt"
	"westex/engines/economy/pkg/core"
	"westex/engines/economy/pkg/entities"
)

func main() {
	// Create a new region
	region := entities.NewRegion("Silicon Valley")

	// Define problems
	foodProblem := entities.NewProblem("Food", "Need for sustenance", 0.8)
	waterProblem := entities.NewProblem("Water", "Need for clean water", 0.9)
	entertainmentProblem := entities.NewProblem("Entertainment", "Need for leisure and fun", 0.5)

	region.AddProblem(foodProblem)
	region.AddProblem(waterProblem)
	region.AddProblem(entertainmentProblem)

	// Create industries using builder pattern
	// Industry 1: Food Production (solves food problem)
	foodProduct := entities.NewResource("Food", 100.0, "kg")
	foodIndustry := entities.CreateIndustry("Farm Ind").
		SetupIndustry([]*entities.Problem{foodProblem}, nil, []*entities.Resource{foodProduct}).
		UpdateIndustryRates(200.0, 1.0, 10000.0)
	region.AddIndustry(foodIndustry)

	// Industry 2: Entertainment (solves entertainment problem)
	entertainmentProduct := entities.NewResource("Entertainment", 50.0, "hours")
	entertainmentIndustry := entities.CreateIndustry("FunZone").
		SetupIndustry([]*entities.Problem{entertainmentProblem}, nil, []*entities.Resource{entertainmentProduct}).
		UpdateIndustryRates(150.0, 1.0, 8000.0)
	region.AddIndustry(entertainmentIndustry)

	// Create population segments
	workersSegment := &entities.PopulationSegment{
		Name:     "Workers",
		Problems: []*entities.Problem{foodProblem, waterProblem, entertainmentProblem},
		Size:     20,
	}

	// Create 100 people
	for i := 1; i <= 20; i++ {
		person := entities.NewPerson(fmt.Sprintf("Person-%d", i), 50.0, 8.0)
		person.AddSegment(workersSegment)
		region.AddPerson(person)
	}

	// Create and run the simulation engine
	engine := core.NewEngine(
		region,
		10.0, // Wage per hour: $10
		2.0,  // Price per unit: $2
		50.0, // Production rate: 50 units per tick
	)

	// Run for 10 ticks
	engine.Run(3)
}
