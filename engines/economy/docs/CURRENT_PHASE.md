# 🎯 Current Phase: Production with Labor Payments

**Phase**: 2.1 - Production Phase  
**Duration**: Week 1 (7 days)  
**Status**: In Progress  

---

## 📋 Overview

**Goal**: Complete the production phase where industries:
1. Hire workers and pay wages immediately
2. Consume input resources
3. Produce output goods
4. Track production costs

**Why this matters**: This creates the foundation of the economic cycle. Workers earn money, which they'll later spend on goods.

---

## ✅ Tasks Breakdown

### **Task 1: Extract Production Calculation** (Day 1)
**Priority**: High  
**Estimated Time**: 2 hours

#### What to do:
Create a clean, testable function for production calculations.

#### Implementation:

**File**: `pkg/production/calculator.go` (new file)

```go
package production

import "westex/engines/economy/pkg/entities"

// ProductionResult contains the outcome of production calculation
type ProductionResult struct {
    UnitsProduced    float32
    LaborUsed        float32
    LaborCost        float32
    ResourceCost     float32
    TotalCost        float32
    CostPerUnit      float32
}

// CalculateProduction determines how much can be produced given constraints
func CalculateProduction(
    industry *entities.Industry,
    availableLabor float32,
    availableHours float32,
    wageRate float32,
) *ProductionResult {
    result := &ProductionResult{}
    
    // Calculate labor utilization
    laborNeeded := industry.LaborNeeded
    laborUsed := min(availableLabor, laborNeeded)
    result.LaborUsed = laborUsed
    
    // Calculate production capacity
    productionRate := laborUsed / laborNeeded
    hoursUsed := productionRate * availableHours
    
    // Units produced (simplified: 1 unit per hour of effective labor)
    result.UnitsProduced = hoursUsed
    
    // Calculate costs
    result.LaborCost = laborUsed * wageRate * availableHours
    result.ResourceCost = calculateResourceCost(industry)
    result.TotalCost = result.LaborCost + result.ResourceCost
    
    if result.UnitsProduced > 0 {
        result.CostPerUnit = result.TotalCost / result.UnitsProduced
    }
    
    return result
}

func calculateResourceCost(industry *entities.Industry) float32 {
    // TODO: Calculate based on input resources consumed
    // For now, return 0 (will implement in Task 3)
    return 0
}

func min(a, b float32) float32 {
    if a < b {
        return a
    }
    return b
}
```

#### Testing:

**File**: `pkg/production/calculator_test.go`

```go
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
    
    expectedCost := 10.0 * 10.0 * 40.0 // workers * wage * hours
    if result.LaborCost != expectedCost {
        t.Errorf("Expected labor cost %.2f, got %.2f", expectedCost, result.LaborCost)
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
    
    // Production should be half of full capacity
    expectedProduction := 5.0 / 10.0 * 40.0 // 20 units
    if result.UnitsProduced != expectedProduction {
        t.Errorf("Expected %.2f units, got %.2f", expectedProduction, result.UnitsProduced)
    }
}
```

---

### **Task 2: Implement Labor Payments** (Day 2)
**Priority**: Critical  
**Estimated Time**: 3 hours

#### What to do:
During production, immediately pay workers their wages.

#### Implementation:

**File**: `pkg/production/labor.go` (new file)

```go
package production

import (
    "fmt"
    "westex/engines/economy/pkg/entities"
)

// LaborPayment represents a wage payment to a worker
type LaborPayment struct {
    PersonName   string
    IndustryName string
    HoursWorked  float32
    WageRate     float32
    TotalPaid    float32
}

// PayWorkers distributes wages to workers employed by an industry
func PayWorkers(
    industry *entities.Industry,
    workers []*entities.Person,
    hoursPerWorker float32,
    wageRate float32,
) ([]LaborPayment, error) {
    payments := make([]LaborPayment, 0)
    totalWages := float32(0)
    
    // Calculate total wages needed
    for _, worker := range workers {
        wages := hoursPerWorker * wageRate
        totalWages += wages
    }
    
    // Check if industry can afford
    if industry.Money < totalWages {
        return nil, fmt.Errorf("industry %s cannot afford wages: needs %.2f, has %.2f",
            industry.Name, totalWages, industry.Money)
    }
    
    // Pay each worker
    for _, worker := range workers {
        wages := hoursPerWorker * wageRate
        
        // Deduct from industry
        industry.Money -= wages
        
        // Pay worker
        worker.Money += wages
        
        // Record payment
        payments = append(payments, LaborPayment{
            PersonName:   worker.Name,
            IndustryName: industry.Name,
            HoursWorked:  hoursPerWorker,
            WageRate:     wageRate,
            TotalPaid:    wages,
        })
    }
    
    return payments, nil
}

// AllocateWorkers assigns workers to an industry based on labor needs
func AllocateWorkers(
    industry *entities.Industry,
    availableWorkers []*entities.Person,
) []*entities.Person {
    needed := int(industry.LaborNeeded)
    available := len(availableWorkers)
    
    // Take minimum of needed and available
    count := needed
    if available < needed {
        count = available
    }
    
    return availableWorkers[:count]
}
```

#### Update Engine:

**File**: `pkg/core/engine.go`

```go
import "westex/engines/economy/pkg/production"

func (e *Engine) processProduction(hours float32) []ProducedGoods {
    pGoodsList := []ProducedGoods{}
    
    // Get available workers (from worker population segment)
    availableWorkers := e.getAvailableWorkers()
    
    for _, industry := range e.Region.Industries {
        // Allocate workers to this industry
        workers := production.AllocateWorkers(industry, availableWorkers)
        
        // Calculate production
        result := production.CalculateProduction(
            industry,
            float32(len(workers)),
            hours,
            SimConf.WagePerHour,
        )
        
        // Pay workers immediately
        payments, err := production.PayWorkers(
            industry,
            workers,
            hours,
            SimConf.WagePerHour,
        )
        
        if err != nil {
            e.Logger.LogEvent(fmt.Sprintf("❌ %s", err.Error()))
            continue
        }
        
        // Log payments
        e.Logger.LogEvent(fmt.Sprintf("💰 %s paid %.2f in wages to %d workers",
            industry.Name, result.LaborCost, len(workers)))
        
        // Produce goods
        for _, product := range industry.OutputProducts {
            product.Add(result.UnitsProduced)
            
            pGoodsList = append(pGoodsList, ProducedGoods{
                IndustryName:     industry.Name,
                ProductName:      product.Name,
                Quantity:         result.UnitsProduced,
                ProductionRate:   result.LaborUsed / industry.LaborNeeded,
                ProductionTarget: 0, // Will calculate from demand
            })
        }
        
        // Remove allocated workers from available pool
        availableWorkers = availableWorkers[len(workers):]
    }
    
    return pGoodsList
}

func (e *Engine) getAvailableWorkers() []*entities.Person {
    workers := make([]*entities.Person, 0)
    
    // Find worker population segment
    for _, segment := range e.Region.PopulationSegments {
        if segment.Name == "Workers" {
            // Get all people in this segment
            for _, person := range e.Region.People {
                for _, personSegment := range person.Segments {
                    if personSegment.Name == segment.Name {
                        workers = append(workers, person)
                        break
                    }
                }
            }
            break
        }
    }
    
    return workers
}
```

---

### **Task 3: Resource Consumption** (Day 3-4)
**Priority**: High  
**Estimated Time**: 4 hours

#### What to do:
Industries consume input resources when producing.

#### Implementation:

**File**: `pkg/production/resources.go` (new file)

```go
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
    resourceCostPerUnit map[string]float32,
) ([]ResourceConsumption, error) {
    consumptions := make([]ResourceConsumption, 0)
    
    // For each input resource
    for _, input := range industry.InputResources {
        // Calculate how much needed (simplified: 1 input per 1 output)
        needed := unitsToProdu
        
        // Check availability
        if input.Quantity < needed {
            return nil, fmt.Errorf("insufficient %s: need %.2f, have %.2f",
                input.Name, needed, input.Quantity)
        }
        
        // Consume
        err := input.Consume(needed)
        if err != nil {
            return nil, err
        }
        
        // Calculate cost
        costPerUnit := resourceCostPerUnit[input.Name]
        if costPerUnit == 0 {
            costPerUnit = 1.0 // Default cost
        }
        
        consumptions = append(consumptions, ResourceConsumption{
            ResourceName: input.Name,
            Quantity:     needed,
            Cost:         needed * costPerUnit,
        })
    }
    
    return consumptions, nil
}
```

#### Update Calculator:

**File**: `pkg/production/calculator.go`

```go
func calculateResourceCost(industry *entities.Industry) float32 {
    totalCost := float32(0)
    
    // Simplified: assume each input resource costs 1.0 per unit
    for _, input := range industry.InputResources {
        // In future, this will be actual market price
        costPerUnit := float32(1.0)
        totalCost += costPerUnit
    }
    
    return totalCost
}
```

---

### **Task 4: Track Production Costs** (Day 5)
**Priority**: Medium  
**Estimated Time**: 2 hours

#### What to do:
Store historical production costs for pricing calculations.

#### Implementation:

**File**: `pkg/entities/industry.go`

Add to Industry struct:
```go
type Industry struct {
    // ... existing fields ...
    
    ProductionHistory []ProductionRecord
}

type ProductionRecord struct {
    Tick            int
    UnitsProduced   float32
    TotalCost       float32
    CostPerUnit     float32
    LaborCost       float32
    ResourceCost    float32
}

func (i *Industry) RecordProduction(record ProductionRecord) {
    i.ProductionHistory = append(i.ProductionHistory, record)
    
    // Keep only last 10 records
    if len(i.ProductionHistory) > 10 {
        i.ProductionHistory = i.ProductionHistory[1:]
    }
}

func (i *Industry) GetAverageCostPerUnit() float32 {
    if len(i.ProductionHistory) == 0 {
        return 0
    }
    
    total := float32(0)
    for _, record := range i.ProductionHistory {
        total += record.CostPerUnit
    }
    
    return total / float32(len(i.ProductionHistory))
}
```

---

### **Task 5: Initial Industry Funding** (Day 6)
**Priority**: Critical  
**Estimated Time**: 1 hour

#### What to do:
Industries need starting capital to pay first round of wages.

#### Implementation:

**File**: `pkg/entities/industry.go`

```go
func (i *Industry) SetInitialCapital(amount float32) *Industry {
    i.Money = amount
    return i
}
```

**File**: `cmd/sim-cli/main.go`

```go
foodIndustry := entities.CreateIndustry("Agriculture Industry").
    SetupIndustry([]*entities.Problem{foodProblem}, inputs, outputs).
    UpdateLabor(float32(4.0)).
    SetInitialCapital(10000.0) // Starting capital
```

---

### **Task 6: Testing & Balancing** (Day 7)
**Priority**: High  
**Estimated Time**: 3 hours

#### What to do:
Test the complete production cycle and balance parameters.

#### Test Scenarios:

1. **Sufficient Labor & Resources**
   - Expected: Full production, all workers paid
   
2. **Insufficient Labor**
   - Expected: Reduced production, some workers unemployed
   
3. **Insufficient Resources**
   - Expected: Production halts, error logged
   
4. **Industry Bankruptcy**
   - Expected: Cannot pay wages, production stops

#### Parameters to Balance:
- Initial industry capital
- Wage rates
- Labor requirements per industry
- Resource consumption rates
- Production rates

---

## 📊 Success Criteria

By end of Week 1, you should have:

✅ **Working production phase** where:
- Industries hire workers
- Workers get paid immediately
- Resources are consumed
- Goods are produced
- Costs are tracked

✅ **Clean code** with:
- Extracted, testable functions
- Proper error handling
- Clear logging

✅ **Balanced parameters** where:
- Industries can afford wages
- Production rates are realistic
- Economy doesn't immediately collapse

---

## 🐛 Common Issues & Solutions

### Issue 1: Industries run out of money
**Solution**: Increase initial capital or reduce wage rates

### Issue 2: No workers available
**Solution**: Increase worker population segment size

### Issue 3: Resource depletion too fast
**Solution**: Increase initial resource quantities or reduce consumption rate

---

## 📝 Notes

- Keep production logic separate from engine logic
- Use dependency injection for testability
- Log everything for debugging
- Don't worry about optimization yet - focus on correctness

---

## 🔜 Next Week Preview

**Week 2** will focus on:
- Pricing system (cost + 10%)
- Product market (people buying goods)
- Needs satisfaction tracking

---

**Let's build this step by step. Start with Task 1 today!** 🚀
