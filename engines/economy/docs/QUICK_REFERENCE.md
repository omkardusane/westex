# Quick Reference - Economy Simulation

## 🚀 Running the Simulation

```bash
# Run the main simulation
go run ./cmd/sim-cli

# Run tests
go test ./pkg/production -v
go test ./... -v  # All tests

# Run with config file (when ready)
# Edit main_config_example.go: rename mainConfigExample() to main()
# Comment out main() in main.go
go run ./cmd/sim-cli
```

---

## 📁 Project Structure

```
/economy
├── cmd/
│   └── sim-cli/
│       ├── main.go                    # Current active main
│       └── main_config_example.go     # YAML config example
├── pkg/
│   ├── core/
│   │   └── engine_new.go             # Main simulation engine
│   ├── entities/
│   │   ├── industry.go               # Industry with production history
│   │   ├── person.go                 # Person/worker
│   │   ├── problem.go                # Needs/problems
│   │   ├── resource.go               # Resources
│   │   └── region.go                 # Region container
│   ├── production/
│   │   ├── calculator.go             # Production calculations
│   │   ├── labor.go                  # Worker allocation & payments
│   │   ├── resources.go              # Resource consumption
│   │   └── production_test.go        # Tests (7/7 passing)
│   ├── config/
│   │   ├── config.go                 # YAML config structures
│   │   └── builder.go                # Build region from config
│   └── logging/
│       └── logger.go                 # Event logging
├── configs/
│   └── mumbai.yaml                   # Example config
└── docs/
    ├── CURRENT_PHASE.md              # What to work on next
    ├── PHASE_2_1_COMPLETION.md       # What we just finished
    ├── ROADMAP.md                    # Long-term plan
    └── START_HERE.md                 # Project overview
```

---

## 🔑 Key Concepts

### Simulation Flow (Each Tick)
1. **Production Phase** → Industries produce goods, pay workers
2. **Pricing Phase** → (Coming in Phase 2.2) Update prices
3. **Resource Regeneration** → Renewable resources replenish
4. **Product Market** → (Coming in Phase 3.1) People buy products

### Current Implementation Status

| Feature | Status | File |
|---------|--------|------|
| Production calculation | ✅ Complete | `production/calculator.go` |
| Labor payments | ✅ Complete | `production/labor.go` |
| Resource consumption | ✅ Complete | `production/resources.go` |
| Production history | ✅ Complete | `entities/industry.go` |
| Pricing system | ⏳ Next | To be created |
| Product market | ⏳ Next | To be created |
| Needs satisfaction | ⏳ Next | To be created |

---

## 💡 Common Tasks

### Adding a New Industry

```go
// In cmd/sim-cli/main.go
newIndustry := entities.CreateIndustry("Tech Industry").
    SetupIndustry(
        []*entities.Problem{techProblem},     // Problems it solves
        []*entities.Resource{electricity},     // Input resources
        []*entities.Resource{software},        // Output products
    ).
    UpdateLabor(float32(20)).                 // Workers needed
    SetInitialCapital(100000.0)               // Starting money

region.AddIndustry(newIndustry)
```

### Adding a New Resource

```go
// Create resource
water := entities.NewResource("Water", "liters")
water.Quantity = 50000              // Initial amount
water.IsFree = true                 // No cost to use
water.RegenerationRate = 1000       // Regenerates 1000/tick

region.AddResource(water)
```

### Adding a Population Segment

```go
segment := &entities.PopulationSegment{
    Name:     "Skilled Workers",
    Problems: []*entities.Problem{foodProblem, housingProblem},
    Size:     500,
}
region.AddPopulationSegment(segment)

// Create people in this segment
for i := 0; i < 500; i++ {
    person := entities.NewPerson(
        fmt.Sprintf("SkilledWorker-%d", i),
        100.0,  // Initial money
        40.0,   // Labor hours available
    )
    person.AddSegment(segment)
    region.AddPerson(person)
}
```

---

## 🧪 Testing

### Run Specific Tests
```bash
# Production tests
go test ./pkg/production -v

# Specific test
go test ./pkg/production -run TestPayWorkers -v
```

### Add New Test
```go
// In pkg/production/production_test.go
func TestMyNewFeature(t *testing.T) {
    // Arrange
    industry := entities.CreateIndustry("Test").UpdateLabor(5.0)
    
    // Act
    result := production.CalculateProduction(industry, 5.0, 40.0, 10.0)
    
    // Assert
    if result.UnitsProduced != 40.0 {
        t.Errorf("Expected 40 units, got %.2f", result.UnitsProduced)
    }
}
```

---

## 📊 Understanding the Output

### Production Phase Log
```
--- Agriculture Industry ---
  Allocated 4 workers (needs 4)           ← Worker allocation
  Production capacity: 100.0% (4/4)       ← % of capacity used
  💰 Paid $6400.00 in wages               ← Labor cost
  📉 Consumed 160.00 RawMaterial          ← Resources used
  ✅ Produced 160.00 Food                 ← Output
  📊 Total cost: $6560.00                 ← Production cost
      Labor: $6400.00                     ← Breakdown
      Resources: $160.00
      Per unit: $41.00                    ← Cost per unit (for pricing)
```

### Final Summary
```
🏭 INDUSTRIES:
  Agriculture Industry:
    Money: $30,800                        ← Cash remaining
    (Start: $50,000, Change: -19,200)     ← Spent on wages
    Products:
      - Food: 480.00 kg                   ← Inventory
    Production History: 3 records         ← Historical data
      Average cost/unit: $41.00           ← For pricing decisions
      Last cost/unit: $41.00
```

---

## 🎯 Current Parameters

### Simulation Settings
- **Wage per hour**: $10.00
- **Weeks per tick**: 4
- **Hours per week**: 40
- **Total hours per tick**: 160

### Industries
| Industry | Workers | Initial Capital | Input | Output |
|----------|---------|----------------|-------|--------|
| Agriculture | 4 | $50,000 | RawMaterial | Food |
| Healthcare | 10 | $80,000 | RawMaterial | Wellness, Medical |

### Population
- **Total**: 1,000 people
- **Workers**: ~191 (19.1%)
- **Initial money per person**: $50

### Resources
- **RawMaterial**: 10,000 units (non-regenerating)

---

## 🔧 Tuning the Economy

### Make industries more profitable
```go
// Increase initial capital
.SetInitialCapital(200000.0)

// Or reduce labor needs
.UpdateLabor(float32(2.0))
```

### Reduce unemployment
```go
// Increase labor needs
.UpdateLabor(float32(50.0))

// Or add more industries
```

### Make resources regenerate
```go
rawMaterial.RegenerationRate = 500  // +500 units per tick
```

---

## 📚 Next Steps

1. **Read**: `docs/CURRENT_PHASE.md` for next tasks
2. **Implement**: Pricing system (Phase 2.2)
3. **Implement**: Product market (Phase 3.1)
4. **Test**: Complete economic cycle

---

## 🆘 Troubleshooting

### "Industry cannot afford wages"
- Increase `SetInitialCapital()`
- Reduce number of workers needed
- Reduce wage rate

### "Insufficient resources"
- Increase initial resource quantity
- Add regeneration rate
- Reduce consumption (fewer workers)

### Tests failing
```bash
# Clean and rebuild
go clean
go test ./... -v
```

### Compilation errors
```bash
# Check for syntax errors
go build ./...

# Format code
go fmt ./...
```

---

**Happy coding!** 🚀

For detailed implementation guides, see:
- `docs/CURRENT_PHASE.md` - Current work
- `docs/ROADMAP.md` - Long-term vision
- `docs/PHASE_2_1_COMPLETION.md` - What's done
