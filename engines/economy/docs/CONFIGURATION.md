# Configuration System

## Overview

The economy simulation now supports YAML configuration files for defining regions, industries, population, and simulation parameters.

## Quick Start

### 1. Create a configuration file

See `configs/mumbai.yaml` for a complete example.

### 2. Load and run simulation

```go
package main

import (
    "fmt"
    "westex/engines/economy/pkg/config"
    "westex/engines/economy/pkg/core"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("configs/mumbai.yaml")
    if err != nil {
        panic(err)
    }
    
    // Build region from config
    region, err := config.BuildRegionFromConfig(cfg)
    if err != nil {
        panic(err)
    }
    
    // Create engine
    engine := core.CreateNewEngine(region)
    
    // Run simulation
    engine.Run(cfg.Simulation.Ticks)
}
```

## Configuration Structure

### Region
```yaml
region:
  name: "Mumbai"
  description: "A bustling metropolitan economy"
```

### Problems (Needs)
```yaml
problems:
  - name: "Food"
    description: "Need for sustenance"
    demand: 0.99          # 99% of population needs this
    basic_need: true      # Survival need vs pleasure
```

- **demand**: 0.0 to 1.0, percentage of population that needs this
- **basic_need**: `true` for survival (food, water), `false` for pleasures (entertainment)

### Resources
```yaml
resources:
  - name: "Land"
    unit: "acres"
    initial_quantity: 5000
    is_free: true              # Government-controlled resource
    regeneration_rate: 0       # Units regenerated per tick
```

- **is_free**: `true` for land, water, minerals (allocated by government)
- **regeneration_rate**: How much regenerates each tick (e.g., forests regrow)

### Industries
```yaml
industries:
  - name: "Agriculture Industry"
    solves_problems:
      - "Food"                 # Problem names this industry addresses
    input_resources:
      - "RawMaterial"          # Resources consumed
      - "Land"
      - "Water"
    output_resources:
      - "Food"                 # Products produced
    labor_needed: 50           # Number of workers required
    initial_capital: 50000     # Starting money
```

### Population
```yaml
population:
  total_size: 1000
  segments:
    - name: "Workers"
      percentage: 0.20         # 20% of total (200 people)
      has_problems: []
      initial_money: 100       # Starting money per person
      labor_hours: 8           # Hours available per tick
      
    - name: "General Population"
      percentage: 0.80         # 80% of total (800 people)
      has_problems:
        - "Food"
        - "Healthcare"
      initial_money: 50
      labor_hours: 0           # Non-workers
```

**Important**: Segment percentages must sum to 1.0 (100%)

### Simulation Parameters
```yaml
simulation:
  ticks: 10                           # Number of simulation ticks
  weeks_per_tick: 4                   # How many weeks each tick represents
  hours_per_week: 40                  # Working hours per week
  wage_per_hour: 10.0                 # Hourly wage rate
  profit_margin: 0.10                 # 10% markup on production costs
  consumption_factor_per_week: 1.0    # Consumption rate
```

## Creating New Scenarios

### Example: Small Village
```yaml
region:
  name: "Small Village"
  description: "A simple agricultural community"

problems:
  - name: "Food"
    description: "Basic sustenance"
    demand: 1.0
    basic_need: true

resources:
  - name: "Land"
    unit: "acres"
    initial_quantity: 100
    is_free: true
    regeneration_rate: 0

industries:
  - name: "Farm"
    solves_problems: ["Food"]
    input_resources: ["Land"]
    output_resources: ["Food"]
    labor_needed: 5
    initial_capital: 1000

population:
  total_size: 50
  segments:
    - name: "Farmers"
      percentage: 1.0
      has_problems: ["Food"]
      initial_money: 20
      labor_hours: 8

simulation:
  ticks: 20
  weeks_per_tick: 1
  hours_per_week: 40
  wage_per_hour: 5.0
  profit_margin: 0.05
  consumption_factor_per_week: 1.0
```

## Validation

The config system validates:
- ✅ Region name is not empty
- ✅ At least one problem exists
- ✅ At least one industry exists
- ✅ Population size is positive
- ✅ Segment percentages sum to ~1.0
- ✅ Industries reference valid problems
- ✅ Industries reference valid resources

## Testing

Run tests:
```bash
go test ./pkg/config -v
```

## Benefits

1. **Easy experimentation**: Change parameters without recompiling
2. **Scenario management**: Save different economic scenarios
3. **Version control**: Track changes to scenarios
4. **Sharing**: Share configurations with others
5. **Documentation**: YAML is self-documenting

## Next Steps

After configuration is working:
1. Implement production phase with labor payments
2. Add resource consumption
3. Implement pricing system
4. Add trade/consumption phase
