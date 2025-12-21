# 🎯 Current Phase: Pricing & Product Market

**Phase**: 2.2 & 3.1 - Pricing System and Product Market  
**Duration**: Week 2 (7 days)  
**Status**: Ready to Start  
**Previous Phase**: ✅ Phase 2.1 COMPLETED (see PHASE_2_1_COMPLETION.md)

---

## 📋 Overview

**Goal**: Implement pricing based on production costs and enable people to buy products, completing the economic cycle.

**Why this matters**: This closes the loop - workers earn wages, then spend money on products, which gives industries revenue to pay more wages.

---

## ✅ Phase 2.1 Completion Summary

All tasks from Phase 2.1 are **COMPLETE**:
- ✅ Production calculation extracted and tested
- ✅ Labor payments implemented (immediate, before production)
- ✅ Resource consumption with validation
- ✅ Production history tracking for cost analysis
- ✅ Initial industry funding
- ✅ All tests passing (7/7)

See `PHASE_2_1_COMPLETION.md` for full details.

---

## 🎯 Phase 2.2: Pricing System (Days 1-3)

### Task 1: Cost-Plus Pricing (Day 1)
**Priority**: High  
**Estimated Time**: 3 hours

#### What to do:
Implement automatic pricing based on production costs + 10% markup.

#### Implementation:

**File**: `pkg/pricing/calculator.go` (new file)

```go
package pricing

import "westex/engines/economy/pkg/entities"

// PricingStrategy defines how prices are calculated
type PricingStrategy string

const (
	CostPlus    PricingStrategy = "cost_plus"    // Cost + fixed markup
	MarketBased PricingStrategy = "market_based" // Supply/demand
)

// PriceCalculator calculates product prices
type PriceCalculator struct {
	Strategy     PricingStrategy
	MarkupRate   float32 // e.g., 0.10 for 10%
	MinimumPrice float32 // Floor price
}

// NewPriceCalculator creates a calculator with cost-plus strategy
func NewPriceCalculator(markupRate float32) *PriceCalculator {
	return &PriceCalculator{
		Strategy:     CostPlus,
		MarkupRate:   markupRate,
		MinimumPrice: 1.0,
	}
}

// CalculatePrice determines the selling price for a product
func (pc *PriceCalculator) CalculatePrice(industry *entities.Industry) float32 {
	switch pc.Strategy {
	case CostPlus:
		return pc.calculateCostPlusPrice(industry)
	default:
		return pc.MinimumPrice
	}
}

func (pc *PriceCalculator) calculateCostPlusPrice(industry *entities.Industry) float32 {
	// Get average production cost
	avgCost := industry.GetAverageCostPerUnit()
	
	// If no history, use last cost or minimum
	if avgCost == 0 {
		avgCost = industry.GetLastProductionCost()
	}
	if avgCost == 0 {
		return pc.MinimumPrice
	}
	
	// Apply markup
	price := avgCost * (1.0 + pc.MarkupRate)
	
	// Enforce minimum
	if price < pc.MinimumPrice {
		price = pc.MinimumPrice
	}
	
	return price
}
```

**File**: `pkg/entities/industry.go` (add field)

```go
type Industry struct {
	// ... existing fields ...
	ProductPrice float32 // Current selling price per unit
}
```

---

### Task 2: Price Updates (Day 2)
**Priority**: High  
**Estimated Time**: 2 hours

#### What to do:
Update prices each tick based on latest production costs.

#### Implementation:

**File**: `pkg/core/engine_new.go`

Add new phase after production:

```go
func (e *Engine) processTick() {
	e.Logger.LogTick(e.CurrentTick)
	
	hoursAvailable := float32(e.WeeksPerTick) * e.HoursPerWeek
	
	// Phase 1: Production (includes labor payments)
	e.Logger.LogEvent("📦 PRODUCTION PHASE")
	e.processProductionPhase(hoursAvailable)
	
	// Phase 2: Pricing (update based on costs)
	e.Logger.LogEvent("\n💰 PRICING PHASE")
	e.processPricingPhase()
	
	// Phase 3: Resource regeneration
	e.Logger.LogEvent("\n🌱 RESOURCE REGENERATION")
	e.processResourceRegeneration()
}

func (e *Engine) processPricingPhase() {
	priceCalc := pricing.NewPriceCalculator(0.10) // 10% markup
	
	for _, industry := range e.Region.Industries {
		newPrice := priceCalc.CalculatePrice(industry)
		oldPrice := industry.ProductPrice
		industry.ProductPrice = newPrice
		
		e.Logger.LogEvent(fmt.Sprintf("💵 %s: $%.2f → $%.2f (based on avg cost $%.2f)",
			industry.Name, oldPrice, newPrice, industry.GetAverageCostPerUnit()))
	}
}
```

---

### Task 3: Price Display (Day 3)
**Priority**: Medium  
**Estimated Time**: 1 hour

#### What to do:
Show prices in logs and final summary.

#### Implementation:

Update production logging to show price:

```go
e.Logger.LogEvent(fmt.Sprintf("📊 Total cost: $%.2f, Price: $%.2f (%.0f%% markup)",
	result.TotalCost, industry.ProductPrice, 
	((industry.ProductPrice/result.CostPerUnit)-1)*100))
```

Update final summary:

```go
fmt.Printf("    Current Price: $%.2f/unit\n", industry.ProductPrice)
fmt.Printf("    Potential Revenue: $%.2f (if all sold)\n", 
	industry.ProductPrice * totalUnitsProduced)
```

---

## 🛒 Phase 3.1: Product Market (Days 4-7)

### Task 4: Needs-Based Purchasing (Day 4-5)
**Priority**: Critical  
**Estimated Time**: 4 hours

#### What to do:
People buy products to satisfy their needs.

#### Implementation:

**File**: `pkg/market/product_market.go` (new file)

```go
package market

import (
	"westex/engines/economy/pkg/entities"
)

// Purchase represents a transaction
type Purchase struct {
	PersonName   string
	IndustryName string
	ProductName  string
	Quantity     float32
	TotalCost    float32
}

// ProcessProductMarket handles people buying products
func ProcessProductMarket(
	region *entities.Region,
) []Purchase {
	purchases := make([]Purchase, 0)
	
	// For each person
	for _, person := range region.People {
		// Identify their needs
		needs := getPersonNeeds(person)
		
		// For each need, try to buy products
		for _, need := range needs {
			// Find industries that solve this need
			for _, industry := range region.Industries {
				if solvesProblem(industry, need) {
					// Try to buy
					purchase := attemptPurchase(person, industry, need)
					if purchase != nil {
						purchases = append(purchases, *purchase)
					}
				}
			}
		}
	}
	
	return purchases
}

func getPersonNeeds(person *entities.Person) []*entities.Problem {
	needs := make([]*entities.Problem, 0)
	for _, segment := range person.Segments {
		needs = append(needs, segment.Problems...)
	}
	return needs
}

func solvesProblem(industry *entities.Industry, problem *entities.Problem) bool {
	for _, p := range industry.OwnedProblems {
		if p.Name == problem.Name {
			return true
		}
	}
	return false
}

func attemptPurchase(
	person *entities.Person,
	industry *entities.Industry,
	need *entities.Problem,
) *Purchase {
	// Check if industry has products
	if len(industry.OutputProducts) == 0 {
		return nil
	}
	
	product := industry.OutputProducts[0] // Simplified: take first product
	price := industry.ProductPrice
	
	// Check if person can afford
	if person.Money < price {
		return nil
	}
	
	// Check if product available
	if product.Quantity < 1.0 {
		return nil
	}
	
	// Make purchase
	quantity := float32(1.0) // Buy 1 unit
	cost := price * quantity
	
	person.Money -= cost
	industry.Money += cost
	product.Consume(quantity)
	
	return &Purchase{
		PersonName:   person.Name,
		IndustryName: industry.Name,
		ProductName:  product.Name,
		Quantity:     quantity,
		TotalCost:    cost,
	}
}
```

---

### Task 5: Needs Satisfaction Tracking (Day 6)
**Priority**: Medium  
**Estimated Time**: 3 hours

#### What to do:
Track which needs are satisfied and which aren't.

#### Implementation:

**File**: `pkg/entities/person.go` (add fields)

```go
type Person struct {
	// ... existing fields ...
	SatisfiedNeeds   map[string]int // Problem name → ticks since satisfied
	UnsatisfiedNeeds map[string]int // Problem name → ticks unsatisfied
}

func (p *Person) SatisfyNeed(problemName string) {
	if p.SatisfiedNeeds == nil {
		p.SatisfiedNeeds = make(map[string]int)
	}
	p.SatisfiedNeeds[problemName] = 0
	delete(p.UnsatisfiedNeeds, problemName)
}

func (p *Person) IncrementUnsatisfiedNeeds() {
	if p.UnsatisfiedNeeds == nil {
		p.UnsatisfiedNeeds = make(map[string]int)
	}
	
	// For each problem in segments
	for _, segment := range p.Segments {
		for _, problem := range segment.Problems {
			if _, satisfied := p.SatisfiedNeeds[problem.Name]; !satisfied {
				p.UnsatisfiedNeeds[problem.Name]++
			}
		}
	}
}
```

---

### Task 6: Integration & Testing (Day 7)
**Priority**: High  
**Estimated Time**: 3 hours

#### What to do:
Integrate product market into engine and test complete cycle.

#### Implementation:

**File**: `pkg/core/engine_new.go`

```go
func (e *Engine) processTick() {
	// ... existing phases ...
	
	// Phase 4: Product Market
	e.Logger.LogEvent("\n🛒 PRODUCT MARKET PHASE")
	e.processProductMarket()
	
	// Phase 5: Needs tracking
	e.updateNeedsSatisfaction()
}

func (e *Engine) processProductMarket() {
	purchases := market.ProcessProductMarket(e.Region)
	
	totalSpent := float32(0)
	for _, purchase := range purchases {
		totalSpent += purchase.TotalCost
		e.Logger.LogEvent(fmt.Sprintf("🛍️  %s bought %.0f %s for $%.2f from %s",
			purchase.PersonName, purchase.Quantity, purchase.ProductName,
			purchase.TotalCost, purchase.IndustryName))
	}
	
	e.Logger.LogEvent(fmt.Sprintf("\n📊 MARKET SUMMARY: %d purchases, $%.2f spent",
		len(purchases), totalSpent))
}
```

---

## 📊 Success Criteria

By end of Week 2, you should have:

✅ **Pricing System** where:
- Prices automatically calculated from production costs
- 10% markup applied
- Prices update each tick

✅ **Product Market** where:
- People buy products based on needs
- Money flows: Industry → Workers → Industry
- Products are consumed

✅ **Complete Economic Cycle**:
```
Tick 1: Produce → Pay Workers → Update Prices → Workers Buy → Needs Satisfied
Tick 2: Produce (with revenue) → Pay Workers → ...
```

✅ **Sustainable Economy**:
- Industries earn revenue from sales
- Revenue covers future wage costs
- Workers can afford basic needs

---

## 🐛 Common Issues & Solutions

### Issue 1: Industries run out of products
**Solution**: Increase production rates or reduce population

### Issue 2: People can't afford products
**Solution**: Lower prices or increase wages

### Issue 3: Industries lose money despite sales
**Solution**: Check markup rate (should be > 0.10 for profit)

---

## 🔜 After This Phase

**Phase 4**: Survival & Mortality
- Poverty mechanics
- Death from unsatisfied basic needs
- Population dynamics

---

**Ready to implement pricing and product market!** 🚀
