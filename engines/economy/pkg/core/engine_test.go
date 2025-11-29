package core

import (
	"westex/engines/economy/pkg/entities"
	"testing"
)

func TestEngineCreation(t *testing.T) {
	region := entities.NewRegion("Test Region")
	engine := NewEngine(region, 10.0, 2.0, 50.0)

	if engine == nil {
		t.Fatal("Engine should not be nil")
	}

	if engine.Region.Name != "Test Region" {
		t.Errorf("Expected region name 'Test Region', got '%s'", engine.Region.Name)
	}

	if engine.WagePerHour != 10.0 {
		t.Errorf("Expected wage per hour 10.0, got %.2f", engine.WagePerHour)
	}
}

func TestSimulationRuns(t *testing.T) {
	// Create a minimal region
	region := entities.NewRegion("Test Region")

	// Add a problem
	foodProblem := entities.NewProblem("Food", "Test food problem", 0.5)
	region.AddProblem(foodProblem)

	// Add an industry
	product := entities.NewResource("TestProduct", 10.0, "units")
	industry := entities.CreateIndustry("TestCorp").
		SetupIndustry([]*entities.Problem{foodProblem}, nil, []*entities.Resource{product}).
		UpdateIndustryRates(10.0, 1.0, 1000.0)
	region.AddIndustry(industry)

	// Add a person
	person := entities.NewPerson("TestPerson", 100.0, 8.0)
	region.AddPerson(person)

	// Create engine with logging disabled
	engine := NewEngine(region, 5.0, 1.0, 10.0)
	engine.Logger = engine.Logger // Keep logger but we could disable it

	// Run simulation - should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Simulation panicked: %v", r)
		}
	}()

	// Run for just 2 ticks to keep test fast
	engine.Run(2)

	// Verify tick count
	if engine.CurrentTick != 2 {
		t.Errorf("Expected current tick to be 2, got %d", engine.CurrentTick)
	}
}

func TestProductionPhase(t *testing.T) {
	region := entities.NewRegion("Test Region")
	
	problem := entities.NewProblem("Test", "Test problem", 0.5)
	product := entities.NewResource("TestProduct", 100.0, "units")
	industry := entities.CreateIndustry("TestCorp").
		SetupIndustry([]*entities.Problem{problem}, nil, []*entities.Resource{product}).
		UpdateIndustryRates(10.0, 1.0, 1000.0)
	region.AddIndustry(industry)

	engine := NewEngine(region, 5.0, 1.0, 50.0)
	
	initialQuantity := product.Quantity
	logs := engine.processProduction()

	if len(logs) != 1 {
		t.Errorf("Expected 1 production log, got %d", len(logs))
	}

	expectedQuantity := initialQuantity + engine.ProductionRate
	if product.Quantity != expectedQuantity {
		t.Errorf("Expected product quantity %.2f, got %.2f", expectedQuantity, product.Quantity)
	}
}
