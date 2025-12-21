# 🛒 Consumption System - Implementation Plan

**Phase**: 3.1 - Product Market & Consumption  
**Status**: Ready to Implement  
**Estimated Time**: 4-6 hours

---

## 🎯 Overview

Implement a consumption system where people buy products to satisfy their needs, completing the economic cycle:

```
Production → Pay Wages → Workers Buy Products → Needs Satisfied → Revenue to Industries
```

---

## 📋 Implementation Tasks

### **Task 1: Needs-Based Purchasing Logic** (2 hours)

#### What to Build:
A market system where people buy products based on their needs and available money.

#### File: `pkg/market/product_market.go` (NEW)

```go
package market

import (
	"westex/engines/economy/pkg/entities"
)

// Purchase represents a completed transaction
type Purchase struct {
	PersonId     int
	PersonName   string
	IndustryName string
	ProductName  string
	ProblemSolved string  // Which need this satisfies
	Quantity     float32
	UnitPrice    float32
	TotalCost    float32
}

// MarketResult summarizes market activity
type MarketResult struct {
	Purchases        []Purchase
	TotalSpent       float32
	TotalRevenue     float32
	PeopleSatisfied  int
	PeopleUnsatisfied int
}

// ProcessProductMarket handles all purchases in one tick
func ProcessProductMarket(
	region *entities.Region,
	pricePerUnit float32,  // Temporary: will use industry.ProductPrice later
) *MarketResult {
	result := &MarketResult{
		Purchases: make([]Purchase, 0),
	}
	
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
			}
		}
	}
	
	return result
}

func findIndustryForProblem(region *entities.Region, problem *entities.Problem) *entities.Industry {
	for _, industry := range region.Industries {
		for _, p := range industry.OwnedProblems {
			if p.Name == problem.Name {
				return industry
			}
		}
	}
	return nil
}

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
		PersonName:    person.Name,
		IndustryName:  industry.Name,
		ProductName:   product.Name,
		ProblemSolved: need.Name,
		Quantity:      quantity,
		UnitPrice:     pricePerUnit,
		TotalCost:     cost,
	}
}
```

---

### **Task 2: Integrate into Engine** (1 hour)

#### File: `pkg/core/engine_new.go`

Add product market phase after production:

```go
func (e *Engine) processTick() {
	e.Logger.LogTick(e.CurrentTick)
	
	hoursAvailable := float32(e.WeeksPerTick) * e.HoursPerWeek
	
	// Phase 1: Production (includes labor payments)
	e.Logger.LogEvent("📦 PRODUCTION PHASE")
	e.processProductionPhase(hoursAvailable)
	
	// Phase 2: Product Market (NEW)
	e.Logger.LogEvent("\n🛒 PRODUCT MARKET PHASE")
	e.processProductMarket()
	
	// Phase 3: Resource regeneration
	e.Logger.LogEvent("\n🌱 RESOURCE REGENERATION")
	e.processResourceRegeneration()
}

func (e *Engine) processProductMarket() {
	// Temporary: use simple pricing (will be replaced with cost-plus pricing)
	pricePerUnit := float32(50.0) // Temporary fixed price
	
	result := market.ProcessProductMarket(e.Region, pricePerUnit)
	
	// Log summary
	e.Logger.LogEvent(fmt.Sprintf("💰 Total spent: $%.2f", result.TotalSpent))
	e.Logger.LogEvent(fmt.Sprintf("📊 Purchases made: %d", len(result.Purchases)))
	e.Logger.LogEvent(fmt.Sprintf("🏭 Industry revenue: $%.2f", result.TotalRevenue))
	
	// Log sample purchases (first 5)
	count := 0
	for _, purchase := range result.Purchases {
		if count >= 5 {
			e.Logger.LogEvent(fmt.Sprintf("   ... and %d more purchases", len(result.Purchases)-5))
			break
		}
		e.Logger.LogEvent(fmt.Sprintf("   🛍️  %s bought %.0f %s for $%.2f (solving %s)",
			purchase.PersonName, purchase.Quantity, purchase.ProductName,
			purchase.TotalCost, purchase.ProblemSolved))
		count++
	}
}
```

---

### **Task 3: Needs Satisfaction Tracking** (1.5 hours)

Track which needs are satisfied over time.

#### File: `pkg/entities/person.go`

Add tracking fields:

```go
type Person struct {
	Name       string
	Segments   []*PopulationSegment
	Money      float32
	LaborHours float32
	
	// Needs tracking (NEW)
	SatisfiedNeeds   map[string]int // Problem name → tick last satisfied
	UnsatisfiedTicks map[string]int // Problem name → consecutive ticks unsatisfied
}

func NewPerson(name string, initialMoney, laborHours float32) *Person {
	return &Person{
		Name:             name,
		Segments:         make([]*PopulationSegment, 0),
		Money:            initialMoney,
		LaborHours:       laborHours,
		SatisfiedNeeds:   make(map[string]int),
		UnsatisfiedTicks: make(map[string]int),
	}
}

// SatisfyNeed marks a need as satisfied this tick
func (p *Person) SatisfyNeed(problemName string, currentTick int) {
	p.SatisfiedNeeds[problemName] = currentTick
	delete(p.UnsatisfiedTicks, problemName)
}

// IncrementUnsatisfiedNeeds updates unsatisfied need counters
func (p *Person) IncrementUnsatisfiedNeeds(currentTick int) {
	allProblems := p.GetAllProblems()
	
	for _, problem := range allProblems {
		// Check if satisfied recently (this tick)
		if lastSatisfied, exists := p.SatisfiedNeeds[problem.Name]; exists {
			if lastSatisfied == currentTick {
				continue // Satisfied this tick
			}
		}
		
		// Increment unsatisfied counter
		p.UnsatisfiedTicks[problem.Name]++
	}
}

// GetUnsatisfiedNeeds returns problems not satisfied for N ticks
func (p *Person) GetUnsatisfiedNeeds(threshold int) []*Problem {
	unsatisfied := make([]*Problem, 0)
	allProblems := p.GetAllProblems()
	
	for _, problem := range allProblems {
		if ticks, exists := p.UnsatisfiedTicks[problem.Name]; exists {
			if ticks >= threshold {
				unsatisfied = append(unsatisfied, problem)
			}
		}
	}
	
	return unsatisfied
}
```

---

### **Task 4: Update Market to Track Satisfaction** (0.5 hours)

#### File: `pkg/market/product_market.go`

Update `attemptPurchase` to mark needs as satisfied:

```go
func attemptPurchase(
	person *entities.Person,
	industry *entities.Industry,
	need *entities.Problem,
	pricePerUnit float32,
	currentTick int,  // NEW parameter
) *Purchase {
	// ... existing purchase logic ...
	
	// Mark need as satisfied
	person.SatisfyNeed(need.Name, currentTick)
	
	return &Purchase{
		// ... existing fields ...
	}
}
```

And update the engine to call it:

```go
func (e *Engine) processProductMarket() {
	result := market.ProcessProductMarket(e.Region, pricePerUnit, e.CurrentTick)
	
	// After market, increment unsatisfied needs
	for _, person := range e.Region.People {
		person.IncrementUnsatisfiedNeeds(e.CurrentTick)
	}
	
	// ... logging ...
}
```

---

### **Task 5: Enhanced Logging & Summary** (1 hour)

#### Update Final Summary

Add needs satisfaction stats:

```go
func (e *Engine) printFinalSummary() {
	// ... existing summary ...
	
	// Needs satisfaction summary (NEW)
	fmt.Printf("\n🎯 NEEDS SATISFACTION:\n")
	
	satisfiedCount := 0
	unsatisfiedCount := 0
	
	for _, person := range e.Region.People {
		allProblems := person.GetAllProblems()
		for _, problem := range allProblems {
			if _, satisfied := person.SatisfiedNeeds[problem.Name]; satisfied {
				satisfiedCount++
			} else {
				unsatisfiedCount++
			}
		}
	}
	
	total := satisfiedCount + unsatisfiedCount
	satisfactionRate := float32(0)
	if total > 0 {
		satisfactionRate = float32(satisfiedCount) / float32(total) * 100
	}
	
	fmt.Printf("  Satisfaction Rate: %.1f%% (%d/%d needs met)\n", 
		satisfactionRate, satisfiedCount, total)
	
	// Show people with unsatisfied needs
	criticalCount := 0
	for _, person := range e.Region.People {
		critical := person.GetUnsatisfiedNeeds(3) // 3+ ticks unsatisfied
		if len(critical) > 0 {
			criticalCount++
		}
	}
	
	if criticalCount > 0 {
		fmt.Printf("  ⚠️  %d people have critical unsatisfied needs (3+ ticks)\n", criticalCount)
	}
}
```

---

## 🧪 Testing Plan

### Test 1: Basic Purchase
```go
func TestAttemptPurchase_Success(t *testing.T) {
	// Setup
	person := entities.NewPerson("TestPerson", 100.0, 8.0)
	problem := entities.NewProblem("Food", "Need food", 0.9)
	product := entities.NewResource("Food", "kg")
	product.Quantity = 10.0
	
	industry := entities.CreateIndustry("FoodCorp").
		SetupIndustry([]*entities.Problem{problem}, nil, []*entities.Resource{product}).
		SetInitialCapital(1000.0)
	
	// Act
	purchase := attemptPurchase(person, industry, problem, 50.0, 1)
	
	// Assert
	if purchase == nil {
		t.Fatal("Expected purchase to succeed")
	}
	if person.Money != 50.0 {
		t.Errorf("Expected person to have $50, got $%.2f", person.Money)
	}
	if industry.Money != 1050.0 {
		t.Errorf("Expected industry to have $1050, got $%.2f", industry.Money)
	}
}
```

### Test 2: Insufficient Funds
```go
func TestAttemptPurchase_InsufficientFunds(t *testing.T) {
	person := entities.NewPerson("PoorPerson", 10.0, 8.0)
	// ... setup ...
	
	purchase := attemptPurchase(person, industry, problem, 50.0, 1)
	
	if purchase != nil {
		t.Error("Expected purchase to fail due to insufficient funds")
	}
}
```

---

## 📊 Expected Behavior

### Tick 1:
```
📦 PRODUCTION PHASE
  Agriculture: Produced 160 Food, paid $6,400 wages
  Healthcare: Produced 160 Medical, paid $16,000 wages

🛒 PRODUCT MARKET PHASE
  💰 Total spent: $25,000
  📊 Purchases made: 500
  🏭 Industry revenue: $25,000
  🛍️  Person-1 bought 1 Food for $50 (solving Food)
  🛍️  Person-2 bought 1 Medical for $50 (solving Healthcare)
  ... and 498 more purchases
```

### Final Summary:
```
🎯 NEEDS SATISFACTION:
  Satisfaction Rate: 85.3% (853/1000 needs met)
  ⚠️  147 people have critical unsatisfied needs (3+ ticks)
```

---

## 🎯 Success Criteria

- [x] People can buy products
- [x] Money flows from people to industries
- [x] Products are consumed
- [x] Needs are tracked as satisfied/unsatisfied
- [x] Industries earn revenue
- [x] Complete economic cycle works

---

## 🔜 After Consumption

Once consumption is working, we can add:

1. **Dynamic Pricing** - Use production costs + markup
2. **Priority Purchasing** - Buy basic needs before luxuries
3. **Poverty Mechanics** - Track people who can't afford needs
4. **Demand-Based Production** - Industries produce based on demand

---

## 📝 Implementation Order

1. ✅ Create `pkg/market/product_market.go`
2. ✅ Add needs tracking to `Person`
3. ✅ Integrate market phase into engine
4. ✅ Add logging and summary
5. ✅ Write tests
6. ✅ Test with simulation
7. ✅ Balance parameters (prices, wages, etc.)

---

**Ready to start implementing?** Let's begin with Task 1! 🚀
