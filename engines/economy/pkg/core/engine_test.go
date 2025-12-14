package core

import (
	"testing"
	"westex/engines/economy/pkg/entities"
)

func TestCreateNewEngine(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")

	// Act
	engine := CreateNewEngine(region)

	// Assert
	if engine == nil {
		t.Fatal("Expected engine to be created, got nil")
	}

	if engine.Region != region {
		t.Error("Expected engine to have the provided region")
	}

	if engine.CurrentTick != 0 {
		t.Errorf("Expected CurrentTick to be 0, got %d", engine.CurrentTick)
	}

	if engine.WagePerHour != 10.0 {
		t.Errorf("Expected default WagePerHour to be 10.0, got %.2f", engine.WagePerHour)
	}

	if engine.WeeksPerTick != 4 {
		t.Errorf("Expected default WeeksPerTick to be 4, got %d", engine.WeeksPerTick)
	}

	if engine.HoursPerWeek != 40.0 {
		t.Errorf("Expected default HoursPerWeek to be 40.0, got %.2f", engine.HoursPerWeek)
	}

	if engine.Logger == nil {
		t.Error("Expected Logger to be initialized")
	}

	if engine.InitialState == nil {
		t.Error("Expected InitialState to be initialized")
	}
}

func TestNewEngineWithParams(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")
	wagePerHour := float32(15.0)
	weeksPerTick := 2
	hoursPerWeek := float32(35.0)

	// Act
	engine := NewEngineWithParams(region, wagePerHour, weeksPerTick, hoursPerWeek)

	// Assert
	if engine == nil {
		t.Fatal("Expected engine to be created, got nil")
	}

	if engine.WagePerHour != wagePerHour {
		t.Errorf("Expected WagePerHour to be %.2f, got %.2f", wagePerHour, engine.WagePerHour)
	}

	if engine.WeeksPerTick != weeksPerTick {
		t.Errorf("Expected WeeksPerTick to be %d, got %d", weeksPerTick, engine.WeeksPerTick)
	}

	if engine.HoursPerWeek != hoursPerWeek {
		t.Errorf("Expected HoursPerWeek to be %.2f, got %.2f", hoursPerWeek, engine.HoursPerWeek)
	}
}

func TestInitialState_CapturesIndustryMoney(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")
	industry := entities.CreateIndustry("TestIndustry").SetInitialCapital(5000.0)
	region.AddIndustry(industry)

	// Act
	engine := CreateNewEngine(region)

	// Assert
	if engine.InitialState.IndustryMoney["TestIndustry"] != 5000.0 {
		t.Errorf("Expected initial industry money to be 5000.0, got %.2f",
			engine.InitialState.IndustryMoney["TestIndustry"])
	}

	if engine.InitialState.TotalWealth != 5000.0 {
		t.Errorf("Expected total wealth to be 5000.0, got %.2f", engine.InitialState.TotalWealth)
	}
}

func TestInitialState_CapturesPersonMoney(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")
	person := entities.NewPerson("TestPerson", 100.0, 8.0)
	region.AddPerson(person)

	// Act
	engine := CreateNewEngine(region)

	// Assert
	if engine.InitialState.PersonMoney["TestPerson"] != 100.0 {
		t.Errorf("Expected initial person money to be 100.0, got %.2f",
			engine.InitialState.PersonMoney["TestPerson"])
	}

	if engine.InitialState.TotalWealth != 100.0 {
		t.Errorf("Expected total wealth to be 100.0, got %.2f", engine.InitialState.TotalWealth)
	}
}

func TestInitialState_CapturesTotalWealth(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")

	industry1 := entities.CreateIndustry("Industry1").SetInitialCapital(5000.0)
	industry2 := entities.CreateIndustry("Industry2").SetInitialCapital(3000.0)
	region.AddIndustry(industry1)
	region.AddIndustry(industry2)

	person1 := entities.NewPerson("Person1", 100.0, 8.0)
	person2 := entities.NewPerson("Person2", 200.0, 8.0)
	region.AddPerson(person1)
	region.AddPerson(person2)

	// Act
	engine := CreateNewEngine(region)

	// Assert
	expectedTotal := float32(5000.0 + 3000.0 + 100.0 + 200.0)
	if engine.InitialState.TotalWealth != expectedTotal {
		t.Errorf("Expected total wealth to be %.2f, got %.2f",
			expectedTotal, engine.InitialState.TotalWealth)
	}
}

func TestGetAvailableWorkers(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")

	workersSegment := &entities.PopulationSegment{
		Name:     "Workers",
		Problems: []*entities.Problem{},
		Size:     10,
	}
	region.AddPopulationSegment(workersSegment)

	// Create 10 workers
	for i := 0; i < 10; i++ {
		person := entities.NewPerson("Worker", 50.0, 8.0)
		person.AddSegment(workersSegment)
		region.AddPerson(person)
	}

	// Create 5 non-workers
	otherSegment := &entities.PopulationSegment{
		Name:     "Other",
		Problems: []*entities.Problem{},
		Size:     5,
	}
	region.AddPopulationSegment(otherSegment)

	for i := 0; i < 5; i++ {
		person := entities.NewPerson("Other", 50.0, 8.0)
		person.AddSegment(otherSegment)
		region.AddPerson(person)
	}

	engine := CreateNewEngine(region)

	// Act
	workers := engine.getAvailableWorkers()

	// Assert
	if len(workers) != 10 {
		t.Errorf("Expected 10 workers, got %d", len(workers))
	}
}

func TestGetAvailableWorkers_NoWorkersSegment(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")

	otherSegment := &entities.PopulationSegment{
		Name:     "Other",
		Problems: []*entities.Problem{},
		Size:     5,
	}
	region.AddPopulationSegment(otherSegment)

	for i := 0; i < 5; i++ {
		person := entities.NewPerson("Other", 50.0, 8.0)
		person.AddSegment(otherSegment)
		region.AddPerson(person)
	}

	engine := CreateNewEngine(region)

	// Act
	workers := engine.getAvailableWorkers()

	// Assert
	if len(workers) != 0 {
		t.Errorf("Expected 0 workers when no Workers segment exists, got %d", len(workers))
	}
}

func TestEngine_ProcessTick_DoesNotPanic(t *testing.T) {
	// Arrange
	region := entities.NewRegion("TestRegion")

	// Create a minimal setup
	problem := entities.NewProblem("Food", "Need food", 0.9)
	region.AddProblem(problem)

	resource := entities.NewResource("RawMaterial", "units")
	resource.Quantity = 1000
	region.AddResource(resource)

	product := entities.NewResource("Food", "kg")

	industry := entities.CreateIndustry("TestIndustry").
		SetupIndustry([]*entities.Problem{problem}, []*entities.Resource{resource}, []*entities.Resource{product}).
		UpdateLabor(2.0).
		SetInitialCapital(10000.0)
	region.AddIndustry(industry)

	workersSegment := &entities.PopulationSegment{
		Name:     "Workers",
		Problems: []*entities.Problem{},
		Size:     5,
	}
	region.AddPopulationSegment(workersSegment)

	for i := 0; i < 5; i++ {
		person := entities.NewPerson("Worker", 50.0, 8.0)
		person.AddSegment(workersSegment)
		region.AddPerson(person)
	}

	engine := CreateNewEngine(region)

	// Act & Assert - should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processTick panicked: %v", r)
		}
	}()

	engine.processTick()
}
