# Commit: Add YAML Configuration System

## Summary
Implemented a complete configuration system for defining economic regions via YAML files. This allows easy experimentation with different scenarios without code changes.

## Changes

### New Files
- `pkg/config/config.go` - YAML loading, validation, and saving
- `pkg/config/builder.go` - Converts YAML config to entity structures
- `pkg/config/config_test.go` - Unit tests for configuration system
- `configs/mumbai.yaml` - Complete Mumbai region configuration
- `docs/CONFIGURATION.md` - Configuration system documentation
- `docs/CONFIG_SUMMARY.md` - Quick reference summary
- `cmd/sim-cli/main_config_example.go` - Example usage

### Modified Files
- `go.mod` - Added `gopkg.in/yaml.v3` dependency
- `pkg/entities/problem.go` - Added `IsBasicNeed` field
- `pkg/entities/resource.go` - Added `IsFree` and `RegenerationRate` fields
- `pkg/entities/industry.go` - Added `SetInitialCapital()` method

## Features

### Configuration Structure
```yaml
region:          # Basic region info
problems:        # Economic needs (food, healthcare, etc.)
resources:       # Materials and free resources (land, water)
industries:      # Production entities
population:      # People and segments
simulation:      # Simulation parameters
```

### Validation
- Region name required
- At least one problem and industry
- Population size > 0
- Segment percentages sum to 1.0
- Industries reference valid problems/resources

### Benefits
- ✅ Easy scenario creation
- ✅ No recompilation needed
- ✅ Version-controllable scenarios
- ✅ Self-documenting YAML
- ✅ Automatic validation

## Usage

```go
// Load from YAML
cfg, _ := config.LoadConfig("configs/mumbai.yaml")
region, _ := config.BuildRegionFromConfig(cfg)
engine := core.CreateNewEngine(region)
engine.Run(cfg.Simulation.Ticks)
```

## Testing
```bash
go mod tidy
go test ./pkg/config -v
go run ./cmd/sim-cli/main_config_example.go
```

## Next Steps
Phase 2.1: Production with Labor Payments
- Create `pkg/production/` package
- Extract production calculations
- Implement immediate wage payments
- Add resource consumption
- Track production costs
