# Quick Start Guide

## 🎯 Goal
Get the economy simulation engine running in 5 minutes!

## Prerequisites

### Install Go
1. Download Go from [golang.org/dl](https://golang.org/dl/)
2. Install for Windows (choose the `.msi` installer)
3. Verify installation:
   ```powershell
   go version
   ```
   You should see something like: `go version go1.21.x windows/amd64`

## 🚀 Running the Simulation

### Option 1: Run Directly (Recommended for Development)

```powershell
# Navigate to the project directory
cd d:\code4\westex\engines\economy

# Run the simulation
go run ./cmd/sim-cli
```

### Option 2: Build and Run

```powershell
# Build the executable
go build -o sim-cli.exe ./cmd/sim-cli

# Run it
.\sim-cli.exe
```

## 🧪 Running Tests

```powershell
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📊 What You'll See

When you run the simulation, you'll see:

1. **Initialization**: Region setup with industries and people
2. **Tick-by-Tick Logs**: 
   - Production phase (industries creating goods)
   - Labor market (people working for wages)
   - Product market (people buying goods)
3. **Final Summary**: Economic statistics

Example output:
```
🚀 Starting Economy Simulation for 10 ticks...
Region: Silicon Valley
Industries: 2, People: 100, Problems: 3

========== TICK 1 [10:30:45] ==========
  📦 PRODUCTION PHASE
  ✓ FarmCo produced 50.00 Food
  ✓ FunZone produced 50.00 Entertainment

  💼 LABOR MARKET PHASE
  ✓ Person-1 worked 2.00 hours for FarmCo, earned 20.00
  ✓ Person-2 worked 2.00 hours for FarmCo, earned 20.00
  ...

  🛒 PRODUCT MARKET PHASE
  ✓ Person-1 bought 8.00 Food from FarmCo for 16.00
  ✓ Person-2 bought 5.00 Entertainment from FunZone for 10.00
  ...
```

## 🎮 Customizing the Simulation

Edit `cmd/sim-cli/main.go` to customize:

### Change Number of People
```go
for i := 1; i <= 200; i++ {  // Change 100 to 200
    person := model.NewPerson(fmt.Sprintf("Person-%d", i), 500.0, 8.0)
    // ...
}
```

### Add More Industries
```go
waterIndustry := model.NewIndustry("AquaCorp", []*model.Problem{waterProblem}, 100.0, 5000.0)
waterProduct := model.NewResource("Water", 200.0, "liters")
waterIndustry.AddOutputProduct(waterProduct)
region.AddIndustry(waterIndustry)
```

### Adjust Economic Parameters
```go
engine := core.NewEngine(
    region,
    15.0,  // Wage per hour (was 10.0)
    3.0,   // Price per unit (was 2.0)
    100.0, // Production rate (was 50.0)
)
```

### Change Simulation Duration
```go
engine.Run(20)  // Run for 20 ticks instead of 10
```

## 🔧 Troubleshooting

### "go: command not found"
- Go is not installed or not in your PATH
- Solution: Install Go and restart your terminal

### Import errors
```powershell
# Clean and rebuild
go clean -cache
go mod tidy
go run ./cmd/sim-cli
```

### Module issues
```powershell
# Verify go.mod exists
cat go.mod

# Should show:
# module simulation-engine
# go 1.21
```

## 📁 Project Structure Quick Reference

```
/economy
├── cmd/sim-cli/main.go       ← Start here to customize
├── pkg/
│   ├── core/engine.go        ← Simulation logic
│   ├── model/                ← Data structures
│   │   ├── region.go
│   │   ├── entity.go
│   │   ├── problem.go
│   │   └── resource.go
│   ├── interaction/          ← Business logic
│   │   ├── labor.go
│   │   └── trade.go
│   └── logging/logger.go     ← Output formatting
└── go.mod                    ← Module definition
```

## 🎓 Next Steps

1. **Experiment**: Modify parameters and see how the economy changes
2. **Add Features**: 
   - New industries
   - More complex production chains
   - Resource dependencies
3. **Learn**: Read `ARCHITECTURE.md` for deeper understanding
4. **Extend**: Add persistence, APIs, or a web interface

## 💡 Common Modifications

### Make Industries Consume Resources
```go
// In pkg/core/engine.go, modify processProduction()
func (e *Engine) processProduction() []string {
    logs := make([]string, 0)
    
    for _, industry := range e.Region.Industries {
        // Check if industry has required inputs
        canProduce := true
        for _, input := range industry.InputResources {
            if input.Quantity < 10.0 {
                canProduce = false
                break
            }
        }
        
        if canProduce {
            // Consume inputs
            for _, input := range industry.InputResources {
                input.Consume(10.0)
            }
            
            // Produce outputs
            for _, product := range industry.OutputProducts {
                product.Add(e.ProductionRate)
                logs = append(logs, fmt.Sprintf("✓ %s produced %.2f %s", 
                    industry.Name, e.ProductionRate, product.Name))
            }
        }
    }
    
    return logs
}
```

### Add Random Events
```go
import "math/rand"

// In processTick(), add:
if rand.Float64() < 0.1 {  // 10% chance
    e.Logger.LogEvent("🌪️  Natural disaster! Production halved this tick.")
    e.ProductionRate *= 0.5
}
```

## 🤝 Getting Help

- Read the `README.md` for overview
- Check `ARCHITECTURE.md` for design patterns
- Review test files for examples
- Experiment and learn by doing!

---

Happy simulating! 🎉
