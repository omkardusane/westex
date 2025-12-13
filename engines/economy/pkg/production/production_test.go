package production

import (
	"testing"
	"westex/engines/economy/pkg/entities"
)

func TestCalculateProduction(t *testing.T) {
	// Create test industry
	industry := entities.CreateIndustry("TestCorp").
		UpdateLabor(10.0) // Needs 10 workers

	// Test with sufficient labor
	result := CalculateProduction(industry, 10.0, 40.0, 10.0)

	if result.LaborUsed != 10.0 {
		t.Errorf("Expected 10 workers used, got %.2f", result.LaborUsed)
	}

	expectedCost := float32(10.0 * 10.0 * 40.0) // workers * wage * hours
	if result.LaborCost != expectedCost {
		t.Errorf("Expected labor cost %.2f, got %.2f", expectedCost, result.LaborCost)
	}

	// Full capacity: 10/10 workers * 40 hours = 40 units
	if result.UnitsProduced != 40.0 {
		t.Errorf("Expected 40 units, got %.2f", result.UnitsProduced)
	}
}

func TestCalculateProduction_InsufficientLabor(t *testing.T) {
	industry := entities.CreateIndustry("TestCorp").
		UpdateLabor(10.0)

	// Only 5 workers available
	result := CalculateProduction(industry, 5.0, 40.0, 10.0)

	if result.LaborUsed != 5.0 {
		t.Errorf("Expected 5 workers used, got %.2f", result.LaborUsed)
	}

	// Production should be half of full capacity: 5/10 * 40 = 20 units
	expectedProduction := float32(20.0)
	if result.UnitsProduced != expectedProduction {
		t.Errorf("Expected %.2f units, got %.2f", expectedProduction, result.UnitsProduced)
	}
}

func TestPayWorkers(t *testing.T) {
	industry := entities.CreateIndustry("TestCorp").
		SetInitialCapital(10000.0)

	workers := []*entities.Person{
		entities.NewPerson("Alice", 100.0, 8.0),
		entities.NewPerson("Bob", 100.0, 8.0),
	}

	payments, err := PayWorkers(industry, workers, 40.0, 10.0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(payments) != 2 {
		t.Errorf("Expected 2 payments, got %d", len(payments))
	}

	// Each worker should receive 40 hours * $10/hour = $400
	expectedPay := float32(400.0)
	if payments[0].TotalPaid != expectedPay {
		t.Errorf("Expected payment %.2f, got %.2f", expectedPay, payments[0].TotalPaid)
	}

	// Workers should have received money
	if workers[0].Money != 500.0 { // 100 + 400
		t.Errorf("Expected worker money 500, got %.2f", workers[0].Money)
	}

	// Industry should have paid out
	expectedIndustryMoney := float32(10000.0 - 800.0) // 2 workers * 400
	if industry.Money != expectedIndustryMoney {
		t.Errorf("Expected industry money %.2f, got %.2f", expectedIndustryMoney, industry.Money)
	}
}

func TestPayWorkers_InsufficientFunds(t *testing.T) {
	industry := entities.CreateIndustry("TestCorp").
		SetInitialCapital(100.0) // Not enough money

	workers := []*entities.Person{
		entities.NewPerson("Alice", 100.0, 8.0),
	}

	_, err := PayWorkers(industry, workers, 40.0, 10.0)
	if err == nil {
		t.Error("Expected error for insufficient funds")
	}
}

func TestAllocateWorkers(t *testing.T) {
	industry := entities.CreateIndustry("TestCorp").
		UpdateLabor(5.0) // Needs 5 workers

	workers := []*entities.Person{
		entities.NewPerson("Alice", 100.0, 8.0),
		entities.NewPerson("Bob", 100.0, 8.0),
		entities.NewPerson("Charlie", 100.0, 8.0),
		entities.NewPerson("David", 100.0, 8.0),
		entities.NewPerson("Eve", 100.0, 8.0),
		entities.NewPerson("Frank", 100.0, 8.0),
		entities.NewPerson("Grace", 100.0, 8.0),
	}

	allocated := AllocateWorkers(industry, workers)

	if len(allocated) != 5 {
		t.Errorf("Expected 5 workers allocated, got %d", len(allocated))
	}
}

func TestConsumeResources(t *testing.T) {
	// Create resources
	rawMaterial := entities.NewResource("RawMaterial", "units")
	rawMaterial.Quantity = 100.0

	land := entities.NewResource("Land", "acres")
	land.Quantity = 50.0
	land.IsFree = true

	industry := entities.CreateIndustry("TestCorp")
	industry.InputResources = []*entities.Resource{rawMaterial, land}

	// Consume resources for 10 units of production
	consumptions, err := ConsumeResources(industry, 10.0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(consumptions) != 2 {
		t.Errorf("Expected 2 consumptions, got %d", len(consumptions))
	}

	// RawMaterial should be consumed
	if rawMaterial.Quantity != 90.0 {
		t.Errorf("Expected 90 units remaining, got %.2f", rawMaterial.Quantity)
	}

	// Land should be consumed (but free)
	if land.Quantity != 40.0 {
		t.Errorf("Expected 40 acres remaining, got %.2f", land.Quantity)
	}

	// Free resource should have 0 cost
	for _, c := range consumptions {
		if c.ResourceName == "Land" && c.Cost != 0 {
			t.Errorf("Expected land cost 0, got %.2f", c.Cost)
		}
	}
}

func TestConsumeResources_Insufficient(t *testing.T) {
	rawMaterial := entities.NewResource("RawMaterial", "units")
	rawMaterial.Quantity = 5.0 // Not enough

	industry := entities.CreateIndustry("TestCorp")
	industry.InputResources = []*entities.Resource{rawMaterial}

	_, err := ConsumeResources(industry, 10.0)
	if err == nil {
		t.Error("Expected error for insufficient resources")
	}
}
