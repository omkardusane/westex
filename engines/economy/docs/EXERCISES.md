# 🛠️ Practical Exercises Workbook

This workbook contains hands-on exercises to practice Go concepts. Each exercise builds on your economy simulation engine.

---

## 📋 How to Use This Workbook

1. Read the concept explanation
2. Try to implement the solution yourself
3. Check the provided solution
4. Run tests to verify
5. Experiment and modify

---

## Exercise 1: Understanding Exported vs Unexported

### Concept
In Go, capitalization determines visibility:
- **Capitalized** = Exported (public) - accessible from other packages
- **lowercase** = unexported (private) - only accessible within the package

### Task
Add wealth calculation to `Region` with proper encapsulation.

### Implementation

**File:** `pkg/entities/region.go`

```go
// Add these methods to the Region struct

// calculateIndustryWealth is unexported - internal helper
func calculateIndustryWealth(industries []*Industry) float64 {
    total := 0.0
    for _, industry := range industries {
        total += industry.Money
    }
    return total
}

// calculatePeopleWealth is unexported - internal helper
func calculatePeopleWealth(people []*Person) float64 {
    total := 0.0
    for _, person := range people {
        total += person.Money
    }
    return total
}

// GetTotalWealth is exported - public API
func (r *Region) GetTotalWealth() float64 {
    industryWealth := calculateIndustryWealth(r.Industries)
    peopleWealth := calculatePeopleWealth(r.People)
    return industryWealth + peopleWealth
}

// GetWealthDistribution is exported - returns breakdown
func (r *Region) GetWealthDistribution() map[string]float64 {
    return map[string]float64{
        "industries": calculateIndustryWealth(r.Industries),
        "people":     calculatePeopleWealth(r.People),
        "total":      r.GetTotalWealth(),
    }
}
```

### Test

**File:** `pkg/entities/region_test.go`

```go
package entities

import "testing"

func TestRegion_GetTotalWealth(t *testing.T) {
    region := NewRegion("Test")
    
    // Add industry with money
    industry := CreateIndustry("TestCorp").UpdateIndustryMoney(1000.0)
    region.AddIndustry(industry)
    
    // Add person with money
    person := NewPerson("Alice", 500.0, 8.0)
    region.AddPerson(person)
    
    expected := 1500.0
    actual := region.GetTotalWealth()
    
    if actual != expected {
        t.Errorf("expected %.2f, got %.2f", expected, actual)
    }
}
```

### Run Test
```bash
go test ./pkg/model -v -run TestRegion_GetTotalWealth
```

### Questions to Think About
1. Why are the helper functions unexported?
2. What happens if you try to call `calculateIndustryWealth` from `pkg/core`?
3. How does this help with maintainability?

---

## Exercise 2: Pointer vs Value Receivers

### Concept
- **Value receiver** (`func (r Region)`) - receives a copy, changes don't persist
- **Pointer receiver** (`func (r *Region)`) - receives the original, changes persist

### Task
Add methods to `Person` that demonstrate both receiver types.

### Implementation

**File:** `pkg/entities/person.go`

```go
// Add these methods to Person

// GetFullInfo returns formatted info (value receiver - read-only)
func (p Person) GetFullInfo() string {
    return fmt.Sprintf("%s (Money: $%.2f, Labor: %.2f hrs)", 
        p.Name, p.Money, p.LaborHours)
}

// CanAfford checks affordability (value receiver - read-only)
func (p Person) CanAfford(amount float64) bool {
    return p.Money >= amount
}

// HasLaborAvailable checks labor availability (value receiver)
func (p Person) HasLaborAvailable(hours float64) bool {
    return p.LaborHours >= hours
}

// EarnMoney adds money (pointer receiver - modifies)
func (p *Person) EarnMoney(amount float64) {
    p.Money += amount
}

// SpendMoney deducts money with validation (pointer receiver)
func (p *Person) SpendMoney(amount float64) error {
    if !p.CanAfford(amount) {
        return fmt.Errorf("%s cannot afford $%.2f (has $%.2f)", 
            p.Name, amount, p.Money)
    }
    p.Money -= amount
    return nil
}

// WorkHours deducts labor hours (pointer receiver)
func (p *Person) WorkHours(hours float64) error {
    if !p.HasLaborAvailable(hours) {
        return fmt.Errorf("%s doesn't have %.2f hours available (has %.2f)", 
            p.Name, hours, p.LaborHours)
    }
    p.LaborHours -= hours
    return nil
}

// ResetLaborHours sets labor hours to default (pointer receiver)
func (p *Person) ResetLaborHours(hours float64) {
    p.LaborHours = hours
}
```

### Test

**File:** `pkg/entities/person_test.go`

```go
package entities

import "testing"

func TestPerson_ValueReceivers(t *testing.T) {
    person := NewPerson("Alice", 100.0, 8.0)
    
    // Test read-only methods
    if !person.CanAfford(50.0) {
        t.Error("should be able to afford 50")
    }
    
    if person.CanAfford(150.0) {
        t.Error("should not be able to afford 150")
    }
    
    if !person.HasLaborAvailable(5.0) {
        t.Error("should have 5 hours available")
    }
}

func TestPerson_PointerReceivers(t *testing.T) {
    person := NewPerson("Bob", 100.0, 8.0)
    
    // Test money operations
    person.EarnMoney(50.0)
    if person.Money != 150.0 {
        t.Errorf("expected 150.0, got %.2f", person.Money)
    }
    
    err := person.SpendMoney(30.0)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if person.Money != 120.0 {
        t.Errorf("expected 120.0, got %.2f", person.Money)
    }
    
    // Test insufficient funds
    err = person.SpendMoney(200.0)
    if err == nil {
        t.Error("expected error for insufficient funds")
    }
    
    // Test labor operations
    err = person.WorkHours(3.0)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if person.LaborHours != 5.0 {
        t.Errorf("expected 5.0 hours, got %.2f", person.LaborHours)
    }
    
    // Reset
    person.ResetLaborHours(8.0)
    if person.LaborHours != 8.0 {
        t.Errorf("expected 8.0 hours after reset, got %.2f", person.LaborHours)
    }
}
```

### Experiment
Try changing `SpendMoney` to a value receiver:
```go
func (p Person) SpendMoney(amount float64) error {
    p.Money -= amount  // This won't persist!
    return nil
}
```

Run the test - it will fail! This demonstrates why we need pointer receivers for mutations.

---

## Exercise 3: Custom Error Types

### Concept
Custom errors provide more context and allow type-based error handling.

### Task
Create domain-specific errors for the economy simulation.

### Implementation

**File:** `pkg/entities/errors.go` (new file)

```go
package entities

import "fmt"

// InsufficientFundsError represents a lack of money
type InsufficientFundsError struct {
    Entity    string
    Available float64
    Required  float64
}

func (e *InsufficientFundsError) Error() string {
    return fmt.Sprintf("%s has insufficient funds: available=%.2f, required=%.2f", 
        e.Entity, e.Available, e.Required)
}

// InsufficientLaborError represents a lack of labor hours
type InsufficientLaborError struct {
    PersonName string
    Available  float64
    Required   float64
}

func (e *InsufficientLaborError) Error() string {
    return fmt.Sprintf("%s has insufficient labor: available=%.2f hrs, required=%.2f hrs", 
        e.PersonName, e.Available, e.Required)
}

// InsufficientResourceError represents a lack of resources
type InsufficientResourceError struct {
    ResourceName string
    Available    float64
    Required     float64
}

func (e *InsufficientResourceError) Error() string {
    return fmt.Sprintf("insufficient %s: available=%.2f, required=%.2f", 
        e.ResourceName, e.Available, e.Required)
}

// NotFoundError represents an entity not found
type NotFoundError struct {
    EntityType string
    EntityID   string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s not found: %s", e.EntityType, e.EntityID)
}
```

### Update Person Methods

**File:** `pkg/entities/person.go`

```go
// Update SpendMoney to use custom error
func (p *Person) SpendMoney(amount float64) error {
    if !p.CanAfford(amount) {
        return &InsufficientFundsError{
            Entity:    p.Name,
            Available: p.Money,
            Required:  amount,
        }
    }
    p.Money -= amount
    return nil
}

// Update WorkHours to use custom error
func (p *Person) WorkHours(hours float64) error {
    if !p.HasLaborAvailable(hours) {
        return &InsufficientLaborError{
            PersonName: p.Name,
            Available:  p.LaborHours,
            Required:   hours,
        }
    }
    p.LaborHours -= hours
    return nil
}
```

### Update Resource Methods

**File:** `pkg/entities/resource.go`

```go
// Update Consume to return error instead of bool
func (r *Resource) Consume(amount float64) error {
    if r.Quantity < amount {
        return &InsufficientResourceError{
            ResourceName: r.Name,
            Available:    r.Quantity,
            Required:     amount,
        }
    }
    r.Quantity -= amount
    return nil
}
```

### Test Error Handling

**File:** `pkg/entities/errors_test.go`

```go
package entities

import (
    "errors"
    "testing"
)

func TestInsufficientFundsError(t *testing.T) {
    person := NewPerson("Alice", 50.0, 8.0)
    
    err := person.SpendMoney(100.0)
    
    // Check error occurred
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    
    // Check error type
    var insufficientFunds *InsufficientFundsError
    if !errors.As(err, &insufficientFunds) {
        t.Fatal("expected InsufficientFundsError")
    }
    
    // Check error details
    if insufficientFunds.Available != 50.0 {
        t.Errorf("expected available=50.0, got %.2f", insufficientFunds.Available)
    }
    
    if insufficientFunds.Required != 100.0 {
        t.Errorf("expected required=100.0, got %.2f", insufficientFunds.Required)
    }
}

func TestInsufficientResourceError(t *testing.T) {
    resource := NewResource("Food", 10.0, "kg")
    
    err := resource.Consume(20.0)
    
    var insufficientResource *InsufficientResourceError
    if !errors.As(err, &insufficientResource) {
        t.Fatal("expected InsufficientResourceError")
    }
    
    if insufficientResource.ResourceName != "Food" {
        t.Errorf("expected resource=Food, got %s", insufficientResource.ResourceName)
    }
}
```

### Advanced: Error Handling in Interactions

**File:** `pkg/interaction/trade.go`

Update to use custom errors and handle them:

```go
func ExecuteTradeTransaction(buyer *entities.Person, seller *entities.Industry, productName string, quantity float64, pricePerUnit float64) (bool, string) {
    // Find product
    var product *entities.Resource
    for _, p := range seller.OutputProducts {
        if p.Name == productName {
            product = p
            break
        }
    }
    
    if product == nil {
        return false, fmt.Sprintf("Industry %s doesn't produce %s", seller.Name, productName)
    }
    
    totalPrice := quantity * pricePerUnit
    
    // Try to spend money
    if err := buyer.SpendMoney(totalPrice); err != nil {
        // Handle specific error type
        var insufficientFunds *entities.InsufficientFundsError
        if errors.As(err, &insufficientFunds) {
            return false, fmt.Sprintf("❌ %s", err.Error())
        }
        return false, err.Error()
    }
    
    // Try to consume product
    if err := product.Consume(quantity); err != nil {
        // Refund the buyer
        buyer.EarnMoney(totalPrice)
        return false, fmt.Sprintf("❌ %s", err.Error())
    }
    
    // Complete transaction
    seller.Money += totalPrice
    
    return true, fmt.Sprintf("✓ %s bought %.2f %s from %s for $%.2f", 
        buyer.Name, quantity, productName, seller.Name, totalPrice)
}
```

---

## Exercise 4: Interfaces for Flexibility

### Concept
Interfaces define behavior contracts, enabling polymorphism and testability.

### Task
Create interfaces for the simulation engine components.

### Implementation

**File:** `pkg/core/interfaces.go` (new file)

```go
package core

import "westex/engines/economy/pkg/entities"

// SimulationEngine defines the contract for any simulation engine
type SimulationEngine interface {
    Run(ticks int)
    GetRegion() *entities.Region
    GetCurrentTick() int
    ProcessTick()
}

// ProductionSystem handles production logic
type ProductionSystem interface {
    ProduceGoods(industries []*entities.Industry, rate float64) []string
}

// MarketSystem handles market transactions
type MarketSystem interface {
    ProcessLaborMarket(region *entities.Region, wageRate float64) []string
    ProcessProductMarket(region *entities.Region, priceRate float64) []string
}

// Verify Engine implements SimulationEngine at compile time
var _ SimulationEngine = (*Engine)(nil)
```

### Update Engine

**File:** `pkg/core/engine.go`

```go
// Add these methods to satisfy the interface

func (e *Engine) GetRegion() *entities.Region {
    return e.Region
}

func (e *Engine) GetCurrentTick() int {
    return e.CurrentTick
}

func (e *Engine) ProcessTick() {
    e.processTick()
}
```

### Create Alternative Implementation

**File:** `pkg/core/simple_engine.go` (new file)

```go
package core

import (
    "fmt"
    "westex/engines/economy/pkg/entities"
)

// SimpleEngine is a minimal simulation engine (no logging)
type SimpleEngine struct {
    Region      *entities.Region
    CurrentTick int
}

func NewSimpleEngine(region *entities.Region) *SimpleEngine {
    return &SimpleEngine{
        Region:      region,
        CurrentTick: 0,
    }
}

func (e *SimpleEngine) Run(ticks int) {
}

func (e *SimpleEngine) ProcessTick() {
    // Minimal tick processing
    fmt.Printf("Tick %d\n", e.CurrentTick)
}

// Verify SimpleEngine also implements SimulationEngine
var _ SimulationEngine = (*SimpleEngine)(nil)
```

### Use Interface in Tests

**File:** `pkg/core/interface_test.go`

```go
package core

import (
    "simulation-engine/pkg/model"
    "testing"
)

// Test that works with ANY SimulationEngine implementation
func testEngineInterface(t *testing.T, engine SimulationEngine) {
    engine.Run(3)
    
    if engine.GetCurrentTick() != 3 {
        t.Errorf("expected tick=3, got %d", engine.GetCurrentTick())
    }
    
    region := engine.GetRegion()
    if region == nil {
        t.Error("region should not be nil")
    }
}

func TestEngine_Interface(t *testing.T) {
    region := model.NewRegion("Test")
    engine := NewEngine(region, nil, 10.0, 2.0, 50.0)
    testEngineInterface(t, engine)
}

func TestSimpleEngine_Interface(t *testing.T) {
    region := model.NewRegion("Test")
    engine := NewSimpleEngine(region)
    testEngineInterface(t, engine)
}
```

### Benefits
1. **Polymorphism**: Use different engine implementations interchangeably
2. **Testing**: Easy to create mock implementations
3. **Flexibility**: Add new implementations without changing existing code

---

## Exercise 5: Dependency Injection with Interfaces

### Concept
Inject dependencies as interfaces, not concrete types.

### Task
Refactor the engine to inject all dependencies.

### Implementation

**File:** `pkg/logging/interface.go` (new file)

```go
package logging

// Logger defines logging behavior
type Logger interface {
    LogTick(tick int)
    LogEvent(message string)
    LogEvents(messages []string)
    LogSummary(title string, data map[string]interface{})
    LogError(err error)
}

// Verify our logger implements the interface
var _ Logger = (*Logger)(nil)
```

**File:** `pkg/logging/console_logger.go`

Rename `logger.go` to `console_logger.go` and update:

```go
package logging

// ConsoleLogger implements Logger for console output
type ConsoleLogger struct {
    enabled bool
}

func NewConsoleLogger(enabled bool) *ConsoleLogger {
    return &ConsoleLogger{enabled: enabled}
}

// ... rest of the methods stay the same
```

**File:** `pkg/logging/mock_logger.go` (new file for testing)

```go
package logging

// MockLogger implements Logger for testing
type MockLogger struct {
    Events []string
    Ticks  []int
}

func NewMockLogger() *MockLogger {
    return &MockLogger{
        Events: make([]string, 0),
        Ticks:  make([]int, 0),
    }
}

func (m *MockLogger) LogTick(tick int) {
    m.Ticks = append(m.Ticks, tick)
}

func (m *MockLogger) LogEvent(message string) {
    m.Events = append(m.Events, message)
}

func (m *MockLogger) LogEvents(messages []string) {
    m.Events = append(m.Events, messages...)
}

func (m *MockLogger) LogSummary(title string, data map[string]interface{}) {
    // Store summary if needed
}

func (m *MockLogger) LogError(err error) {
    m.Events = append(m.Events, "ERROR: "+err.Error())
}
```

### Update Engine Constructor

**File:** `pkg/core/engine.go`

```go
// Update to accept Logger interface
func NewEngine(
    region *entities.Region,
    logger logging.Logger,  // Interface, not concrete type!
    wagePerHour,
    pricePerUnit,
    productionRate float64,
) *Engine {
    // If no logger provided, use a default
    if logger == nil {
        logger = logging.NewConsoleLogger(true)
    }
    
    return &Engine{
        Region:         region,
        Logger:         logger,
        CurrentTick:    0,
        WagePerHour:    wagePerHour,
        PricePerUnit:   pricePerUnit,
        ProductionRate: productionRate,
    }
}
```

### Update CLI

**File:** `cmd/sim-cli/main.go`

```go
import (
    "westex/engines/economy/pkg/core"
    "westex/engines/economy/pkg/logging"
    "westex/engines/economy/pkg/entities"
)

func main() {
    region := setupRegion()
    
    // Inject console logger
    logger := logging.NewConsoleLogger(true)
    
    engine := core.NewEngine(region, logger, 10.0, 2.0, 50.0)
    engine.Run(10)
}
```

### Test with Mock Logger

**File:** `pkg/core/engine_di_test.go`

```go
package core

import (
    "westex/engines/economy/pkg/logging"
    "westex/engines/economy/pkg/entities"
    "testing"
)

func TestEngine_WithMockLogger(t *testing.T) {
    region := model.NewRegion("Test")
    mockLogger := logging.NewMockLogger()
    
    engine := NewEngine(region, mockLogger, 10.0, 2.0, 50.0)
    engine.Run(3)
    
    // Verify logging happened
    if len(mockLogger.Ticks) != 3 {
        t.Errorf("expected 3 ticks logged, got %d", len(mockLogger.Ticks))
    }
    
    if mockLogger.Ticks[0] != 1 {
        t.Errorf("expected first tick=1, got %d", mockLogger.Ticks[0])
    }
}
```

---

## 🎯 Summary

You've now practiced:
1. ✅ Exported vs unexported (encapsulation)
2. ✅ Pointer vs value receivers
3. ✅ Custom error types
4. ✅ Interfaces for abstraction
5. ✅ Dependency injection

## 📝 Next Steps

1. Implement all exercises in order
2. Run the tests to verify
3. Experiment with variations
4. Move on to `LEARNING_PATH.md` for more advanced topics

## 🧪 Running All Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/model -v
```

---

**Remember:** The best way to learn is by doing. Don't just read the solutions - type them out, break them, fix them, and understand why they work!
