# Package Architecture Guide

## 📦 New Package Structure

```
westex/engines/economy/pkg/
├── entities/            ← Core business entities (WHAT exists)
│   ├── industry.go      → Industries that produce goods
│   ├── person.go        → People and PopulationSegments
│   ├── region.go        → Geographic container for all entities
│   ├── problem.go       → Economic needs/issues
│   └── resource.go      → Materials and commodities
│
├── market/              ← Economic interactions (WHAT happens)
│   ├── labor.go         → Employment transactions
│   └── trade.go         → Product buying/selling
│
└── core/                ← Simulation engine (HOW it runs)
    ├── engine.go        → Main simulation loop
    └── engine_test.go   → Tests
```

---

## 🎓 Why This Structure?

### 1. **Package Naming Rule**
```go
// Directory path and package name MUST match
📁 pkg/entities/
   └── person.go  → package entities  ✅

// Import and use:
import "westex/engines/economy/pkg/entities"
person := entities.NewPerson(...)  // Clear!
```

### 2. **Separation of Concerns**

| Package | Responsibility | Examples |
|---------|---------------|----------|

**A**: In Go, package name should match directory name.

```go
// ❌ BAD - Confusing!
📁 pkg/model/entities/
   industry.go → package model

import "path/to/entities"  // Says "entities"
model.Industry             // But uses "model" - CONFUSING!

// ✅ GOOD - Clear!
📁 pkg/entities/
   industry.go → package entities

import "path/to/entities"    // Says "entities"
entities.Industry            // Uses "entities" - MATCHES!
```

### Q: Where should PopulationSegment be defined?

**A**: In `person.go`, right before `Person` type.

**Reasoning**:
1. **Cohesion**: PopulationSegment is tightly coupled to Person
2. **Readability**: Readers see the dependency before it's used
3. **Go Convention**: Related types in same file

```go
// person.go
package entities

// PopulationSegment comes FIRST (dependency)
type PopulationSegment struct { ... }

// Person comes SECOND (uses PopulationSegment)
type Person struct {
    Segments []*PopulationSegment  // Uses above type
}
```

### Q: How to organize Problem and Labor?

**A**: They're in different packages based on their role:

- **Problem** → `entities/problem.go` (it's a core entity)
- **Labor** → `market/labor.go` (it's an interaction/transaction)

---

## 🔄 Migration Path

### Old Structure:
```
pkg/model/
├── entity.go      (Person, PopulationSegment)
├── region.go
├── problem.go
├── resource.go
└── entities/
    └── industry.go
```

### New Structure:
```
pkg/entities/        ← All core entities together
├── industry.go
├── person.go      (Person + PopulationSegment)
├── region.go
├── problem.go
└── resource.go

pkg/market/        ← All market interactions
├── labor.go
└── trade.go
```

---

## 💡 Best Practices

### 1. **One Concept Per File**
```go
// ✅ GOOD
person.go     → Person + PopulationSegment (related)
industry.go   → Industry only

// ❌ BAD
entities.go   → Person + Industry + Region (unrelated)
```

### 2. **Package Names**
- Use **singular** nouns: `domain`, not `domains`
- Use **lowercase**: `market`, not `Market`
- Keep it **short**: `pkg`, not `package`

### 3. **Import Organization**
```go
import (
    // Standard library first
    "fmt"
    
    // External packages
    "github.com/some/package"
    
    // Your packages (grouped by layer)
    "westex/engines/economy/pkg/entities"
    "westex/engines/economy/pkg/market"
)
```

---

## 🚀 Usage Examples

### Creating Entities
```go
import "westex/engines/economy/pkg/entities"

// Create a region
region := entities.NewRegion("Silicon Valley")

// Create a problem
foodProblem := entities.NewProblem("Food", "Need sustenance", 0.8)

// Create an industry (builder pattern)
industry := entities.CreateIndustry("FarmCo").
    SetupIndustry([]*entities.Problem{foodProblem}, nil, outputs).
    UpdateIndustryRates(200.0, 1.0, 10000.0)

// Create a person
person := entities.NewPerson("Alice", 500.0, 8.0)

// Create a population segment
segment := entities.NewPopulationSegment("Workers", 
    []*entities.Problem{foodProblem}, 100)
person.AddSegment(segment)
```

### Running Markets
```go
import (
    "westex/engines/economy/pkg/entities"
    "westex/engines/economy/pkg/market"
)

// Process labor market
laborLogs := market.ProcessLaborMarket(region, 10.0)

// Process product market
tradeLogs := market.ProcessProductMarket(region, 2.0)
```

---

## 📚 Further Reading

- [Effective Go - Package Names](https://go.dev/doc/effective_go#package-names)
- [Go Blog - Package Names](https://go.dev/blog/package-names)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)
