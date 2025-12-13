package main

import (
	"fmt"
	"westex/engines/economy/pkg/core"
	"westex/engines/economy/pkg/entities"
	"westex/engines/economy/pkg/utils"
)

func main() {
	sim2()
}

func sim2() {
	region := entities.NewRegion("Mumbai")

	// Define problems
	foodProblem := entities.NewProblem("Food", "Need for sustenance", 0.9)
	region.AddProblem(foodProblem)
	healthCareProblem := entities.NewProblem("Healthcare", "Need for medical services", 0.8)
	region.AddProblem(healthCareProblem)
	// educationProblem := entities.NewProblem("Education", "Need for learning and knowledge", 0.6)
	// entertainmentProblem := entities.NewProblem("Entertainment", "Need for leisure and fun", 0.2)
	// region.AddProblem(educationProblem)
	// region.AddProblem(entertainmentProblem)
	// Create industries using builder pattern

	rawMaterial := entities.NewResource("RawMaterial", "units")
	rawMaterial.Quantity = 10000 // Initial supply
	region.AddResource(rawMaterial)

	// Industry 1: Food Production (solves food problem)
	foodProduct := entities.NewResource("Food", "kg")
	foodIndustry := entities.CreateIndustry("Agriculture Industry").
		SetupIndustry([]*entities.Problem{foodProblem}, []*entities.Resource{rawMaterial}, []*entities.Resource{foodProduct}).
		UpdateLabor(float32(4.0)).
		SetInitialCapital(50000.0) // Starting capital for wages
	region.AddIndustry(foodIndustry)

	// Industry 2: Healthcare Services (solves healthcare problem)
	wellnessServices := entities.NewResource("Wellness", "visits")
	healthcareServices := entities.NewResource("Medical", "treatments")
	healthcareIndustry := entities.CreateIndustry("Health Industry").
		SetupIndustry([]*entities.Problem{healthCareProblem}, []*entities.Resource{rawMaterial}, []*entities.Resource{wellnessServices, healthcareServices}).
		UpdateLabor(float32(10)).
		SetInitialCapital(80000.0) // Starting capital for wages
	region.AddIndustry(healthcareIndustry)

	// Create population segments
	workersPopulation := &entities.PopulationSegment{
		Name:     "Workers",
		Problems: []*entities.Problem{},
		Size:     200,
	}
	generalPopulationSegment := &entities.PopulationSegment{
		Name:     "General Population",
		Problems: []*entities.Problem{foodProblem, healthCareProblem},
		Size:     1000,
	}
	region.AddPopulationSegment(workersPopulation)
	region.AddPopulationSegment(generalPopulationSegment)
	// Create 1000 people
	workersCount := 0
	for i := 1; i <= generalPopulationSegment.Size; i++ {
		person := entities.NewPerson(fmt.Sprintf("Person-%d", i), 50.0, 8.0)
		person.AddSegment(generalPopulationSegment)
		// probabilistically assign to workers segment
		if utils.ProbableChance(float32(workersPopulation.Size) / float32(generalPopulationSegment.Size)) {
			person.AddSegment(workersPopulation)
			workersCount++
		}
		region.AddPerson(person)
	}
	workersPopulation.UpdateSize(workersCount)

	/*
		we have 200 workers and a population of 1000 people
		Industries will use labor from the 200 workers
	*/

	/*
		Goals of industries: to produce food for 1000 people and healthcare for 1000 people
	*/
	// Create and run the simulation engine
	healthCareProblem.UpdateDemand(0.1)
	foodProblem.UpdateDemand(0.99)
	engine := core.CreateNewEngine(region)
	engine.Run(3)
	// set production parameters
	// engine.
}

func sim1() {
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
	foodProduct := entities.NewResource("Food", "kg")
	foodIndustry := entities.CreateIndustry("Farm Ind").
		SetupIndustry([]*entities.Problem{foodProblem}, nil, []*entities.Resource{foodProduct}).
		UpdateIndustryRates(200.0, 1.0, 10000.0)
	region.AddIndustry(foodIndustry)

	// Industry 2: Entertainment (solves entertainment problem)
	entertainmentProduct := entities.NewResource("Entertainment", "hours")
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
	engine := core.NewEngineWithParams(
		region,
		10.0, // Wage per hour: $10
		4,    // Weeks per tick
		40.0, // Hours per week
	)

	// Run for 3 ticks
	engine.Run(3)
}
