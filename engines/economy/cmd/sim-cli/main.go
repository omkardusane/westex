package main

import (
	"flag"
	"fmt"
	"log"

	"westex/engines/economy/pkg/config"
	"westex/engines/economy/pkg/core"
	"westex/engines/economy/pkg/entities"
	"westex/engines/economy/pkg/utils"
)

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "", "Path to YAML configuration file")
	flag.Parse()

	if *configFile != "" {
		// Run from YAML config
		runFromConfig(*configFile)
	} else {
		// Run with programmatic setup (default)
		runProgrammatic()
	}
}

// runFromConfig loads and runs simulation from a YAML configuration file
func runFromConfig(filepath string) {
	fmt.Println("=== Running simulation from config file ===")
	fmt.Printf("Loading: %s\n\n", filepath)

	// Load configuration
	cfg, err := config.LoadConfig(filepath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded config for: %s\n", cfg.Region.Name)
	fmt.Printf("  - %d problems defined\n", len(cfg.Problems))
	fmt.Printf("  - %d resources available\n", len(cfg.Resources))
	fmt.Printf("  - %d industries\n", len(cfg.Industries))
	fmt.Printf("  - Population: %d\n\n", cfg.Population.TotalSize)

	// Build region from config
	region, err := config.BuildRegionFromConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to build region: %v", err)
	}

	fmt.Printf("Region '%s' created successfully!\n", region.Name)
	fmt.Printf("  - Industries: %d\n", len(region.Industries))
	fmt.Printf("  - People: %d\n", len(region.People))
	fmt.Printf("  - Population Segments: %d\n\n", len(region.PopulationSegments))

	// Create engine with config parameters
	engine := core.NewEngineWithParams(
		region,
		cfg.Simulation.WagePerHour,
		cfg.Simulation.WeeksPerTick,
		cfg.Simulation.HoursPerWeek,
	)

	// Run simulation
	engine.Run(cfg.Simulation.Ticks)
}

// runProgrammatic runs simulation with programmatic setup
func runProgrammatic() {
	fmt.Println("=== Running simulation with programmatic setup ===")

	region := entities.NewRegion("Mumbai")

	// Define problems
	foodProblem := entities.NewProblem("Food", "Need for sustenance", 0.9)
	region.AddProblem(foodProblem)
	healthCareProblem := entities.NewProblem("Healthcare", "Need for medical services", 0.8)
	region.AddProblem(healthCareProblem)

	// Create resources
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
		// Probabilistically assign to workers segment
		if utils.ProbableChance(float32(workersPopulation.Size) / float32(generalPopulationSegment.Size)) {
			person.AddSegment(workersPopulation)
			workersCount++
		}
		region.AddPerson(person)
	}
	workersPopulation.UpdateSize(workersCount)

	// Update problem demands
	healthCareProblem.UpdateDemand(0.1)
	foodProblem.UpdateDemand(0.99)

	// Create and run engine
	engine := core.CreateNewEngine(region)
	engine.Run(3)
}
