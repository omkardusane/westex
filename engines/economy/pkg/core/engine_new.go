package core

import (
	"fmt"
	"time"

	"westex/engines/economy/pkg/entities"
	"westex/engines/economy/pkg/logging"
	"westex/engines/economy/pkg/production"
)

// Engine is the core simulation engine
type Engine struct {
	Region       *entities.Region
	Logger       *logging.Logger
	CurrentTick  int
	WagePerHour  float32
	WeeksPerTick int
	HoursPerWeek float32
	InitialState *InitialState
}

// InitialState captures the starting state of the economy
type InitialState struct {
	IndustryMoney map[string]float32
	PersonMoney   map[string]float32
	TotalWealth   float32
}

// CreateNewEngine creates a new simulation engine with default parameters
func CreateNewEngine(region *entities.Region) *Engine {
	return NewEngineWithParams(region, 10.0, 4, 40.0)
}

// NewEngineWithParams creates a new simulation engine with custom parameters
func NewEngineWithParams(
	region *entities.Region,
	wagePerHour float32,
	weeksPerTick int,
	hoursPerWeek float32,
) *Engine {
	// Capture initial state
	initialState := &InitialState{
		IndustryMoney: make(map[string]float32),
		PersonMoney:   make(map[string]float32),
		TotalWealth:   0,
	}

	for _, ind := range region.Industries {
		initialState.IndustryMoney[ind.Name] = ind.Money
		initialState.TotalWealth += ind.Money
	}

	for _, p := range region.People {
		initialState.PersonMoney[p.Name] = p.Money
		initialState.TotalWealth += p.Money
	}

	return &Engine{
		Region:       region,
		Logger:       logging.NewLogger(true),
		CurrentTick:  0,
		WagePerHour:  wagePerHour,
		WeeksPerTick: weeksPerTick,
		HoursPerWeek: hoursPerWeek,
		InitialState: initialState,
	}
}

// Run executes the simulation for a given number of ticks
func (e *Engine) Run(ticks int) {
	fmt.Printf("\n🚀 Starting Economy Simulation for %d ticks...\n", ticks)
	fmt.Printf("Region: %s\n", e.Region.Name)
	fmt.Printf("Industries: %d, People: %d, Problems: %d\n",
		len(e.Region.Industries), len(e.Region.People), len(e.Region.Problems))
	fmt.Printf("Wage Rate: $%.2f/hour, Weeks/Tick: %d, Hours/Week: %.0f\n\n",
		e.WagePerHour, e.WeeksPerTick, e.HoursPerWeek)

	for i := 0; i < ticks; i++ {
		e.CurrentTick = i + 1
		e.processTick()
		time.Sleep(300 * time.Millisecond) // Slow down for readability
	}

	e.printFinalSummary()
}

// processTick handles one simulation tick
func (e *Engine) processTick() {
	e.Logger.LogTick(e.CurrentTick)

	// Calculate hours available this tick
	hoursAvailable := float32(e.WeeksPerTick) * e.HoursPerWeek

	// Phase 1: Production (includes labor payments)
	e.Logger.LogEvent("📦 PRODUCTION PHASE")
	e.processProductionPhase(hoursAvailable)

	// Phase 2: Resource regeneration
	e.Logger.LogEvent("\n🌱 RESOURCE REGENERATION")
	e.processResourceRegeneration()

	// Future phases:
	// Phase 3: Product market (people buy goods)
	// Phase 4: Needs satisfaction tracking
}

// processProductionPhase handles production and labor payments
func (e *Engine) processProductionPhase(hoursAvailable float32) {
	// Get available workers
	availableWorkers := e.getAvailableWorkers()
	e.Logger.LogEvent(fmt.Sprintf("Available workers: %d", len(availableWorkers)))

	totalWagesPaid := float32(0)
	totalUnitsProduced := float32(0)

	for _, industry := range e.Region.Industries {
		e.Logger.LogEvent(fmt.Sprintf("\n--- %s ---", industry.Name))

		// Allocate workers
		workers := production.AllocateWorkers(industry, availableWorkers)
		e.Logger.LogEvent(fmt.Sprintf("Allocated %d workers (needs %.0f)", len(workers), industry.LaborNeeded))

		if len(workers) == 0 {
			e.Logger.LogEvent("❌ No workers available")
			continue
		}

		// Calculate production
		result := production.CalculateProduction(
			industry,
			float32(len(workers)),
			hoursAvailable,
			e.WagePerHour,
		)

		e.Logger.LogEvent(fmt.Sprintf("Production capacity: %.1f%% (%.0f/%.0f workers)",
			(result.LaborUsed/industry.LaborNeeded)*100, result.LaborUsed, industry.LaborNeeded))

		// Pay workers FIRST (before production)
		payments, err := production.PayWorkers(
			industry,
			workers,
			hoursAvailable,
			e.WagePerHour,
		)

		if err != nil {
			e.Logger.LogEvent(fmt.Sprintf("❌ %s", err.Error()))
			continue
		}

		e.Logger.LogEvent(fmt.Sprintf("💰 Paid $%.2f in wages to %d workers", result.LaborCost, len(workers)))
		totalWagesPaid += result.LaborCost

		// Consume resources
		consumptions, err := production.ConsumeResources(industry, result.UnitsProduced)
		if err != nil {
			e.Logger.LogEvent(fmt.Sprintf("❌ Resource shortage: %s", err.Error()))
			// Refund workers since we can't produce
			for _, payment := range payments {
				for _, person := range e.Region.People {
					if person.Name == payment.PersonName {
						person.Money -= payment.TotalPaid
						industry.Money += payment.TotalPaid
						break
					}
				}
			}
			continue
		}

		// Log resource consumption
		for _, consumption := range consumptions {
			e.Logger.LogEvent(fmt.Sprintf("📉 Consumed %.2f %s (cost: $%.2f)",
				consumption.Quantity, consumption.ResourceName, consumption.Cost))
		}

		// Produce goods
		for _, product := range industry.OutputProducts {
			product.Add(result.UnitsProduced)
			e.Logger.LogEvent(fmt.Sprintf("✅ Produced %.2f %s (total: %.2f)",
				result.UnitsProduced, product.Name, product.Quantity))
			totalUnitsProduced += result.UnitsProduced
		}

		// Log costs
		e.Logger.LogEvent(fmt.Sprintf("📊 Total cost: $%.2f (Labor: $%.2f, Resources: $%.2f, Per unit: $%.2f)",
			result.TotalCost, result.LaborCost, result.ResourceCost, result.CostPerUnit))

		// Remove allocated workers from available pool
		availableWorkers = availableWorkers[len(workers):]
	}

	// Summary
	e.Logger.LogEvent(fmt.Sprintf("\n📈 PRODUCTION SUMMARY: %.2f units produced, $%.2f paid in wages",
		totalUnitsProduced, totalWagesPaid))

	unemployed := len(e.getAvailableWorkers()) - len(availableWorkers)
	if unemployed > 0 {
		e.Logger.LogEvent(fmt.Sprintf("⚠️  %d workers unemployed this tick", len(availableWorkers)))
	}
}

// processResourceRegeneration regenerates renewable resources
func (e *Engine) processResourceRegeneration() {
	production.RegenerateResources(e.Region.Resources)

	regenerated := 0
	for _, resource := range e.Region.Resources {
		if resource.RegenerationRate > 0 {
			e.Logger.LogEvent(fmt.Sprintf("🌿 %s regenerated +%.2f %s (total: %.2f)",
				resource.Name, resource.RegenerationRate, resource.Unit, resource.Quantity))
			regenerated++
		}
	}

	if regenerated == 0 {
		e.Logger.LogEvent("No renewable resources")
	}
}

// getAvailableWorkers returns all people in the "Workers" segment
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

// printFinalSummary prints statistics at the end of simulation
func (e *Engine) printFinalSummary() {
	fmt.Printf("\n\n" + "═══════════════════════════════════════\n")
	fmt.Printf("📊 FINAL SIMULATION SUMMARY\n")
	fmt.Printf("═══════════════════════════════════════\n\n")

	// Industry summary
	fmt.Printf("🏭 INDUSTRIES:\n")
	for _, industry := range e.Region.Industries {
		start := e.InitialState.IndustryMoney[industry.Name]
		change := industry.Money - start
		fmt.Printf("  %s:\n", industry.Name)
		fmt.Printf("    Money: $%.2f (Start: $%.2f, Change: %+.2f)\n", industry.Money, start, change)
		fmt.Printf("    Products:\n")
		for _, product := range industry.OutputProducts {
			fmt.Printf("      - %s: %.2f %s\n", product.Name, product.Quantity, product.Unit)
		}
	}

	// People summary
	fmt.Printf("\n👥 PEOPLE (showing first 5):\n")
	for i, person := range e.Region.People {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(e.Region.People)-5)
			break
		}
		start := e.InitialState.PersonMoney[person.Name]
		change := person.Money - start
		fmt.Printf("  %s: $%.2f (Start: $%.2f, Change: %+.2f)\n", person.Name, person.Money, start, change)
	}

	// Calculate total wealth
	totalWealth := float32(0.0)
	for _, person := range e.Region.People {
		totalWealth += person.Money
	}
	for _, industry := range e.Region.Industries {
		totalWealth += industry.Money
	}

	wealthChange := totalWealth - e.InitialState.TotalWealth

	fmt.Printf("\n💰 TOTAL WEALTH: $%.2f (Start: $%.2f, Change: %+.2f)\n", totalWealth, e.InitialState.TotalWealth, wealthChange)

	// Resource summary
	fmt.Printf("\n📦 RESOURCES:\n")
	for _, resource := range e.Region.Resources {
		status := ""
		if resource.IsFree {
			status = " (free resource)"
		}
		if resource.RegenerationRate > 0 {
			status += fmt.Sprintf(" (regenerates +%.0f/tick)", resource.RegenerationRate)
		}
		fmt.Printf("  %s: %.2f %s%s\n", resource.Name, resource.Quantity, resource.Unit, status)
	}

	fmt.Printf("\n✅ Simulation completed successfully!\n\n")
}
