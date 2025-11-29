package core

import (
	"fmt"
	"time"
	"westex/engines/economy/pkg/entities"
	"westex/engines/economy/pkg/logging"
	"westex/engines/economy/pkg/market"
)

// Engine is the core simulation engine
type Engine struct {
	Region         *entities.Region
	Logger         *logging.Logger
	CurrentTick    int
	WagePerHour    float64 // Standard wage rate
	PricePerUnit   float64 // Standard price per unit of product
	ProductionRate float64 // How much each industry produces per tick
	InitialState   *InitialState
}

// InitialState captures the starting state of the economy
type InitialState struct {
	IndustryMoney map[string]float64
	PersonMoney   map[string]float64
	TotalWealth   float64
}

// NewEngine creates a new simulation engine
func NewEngine(region *entities.Region, wagePerHour, pricePerUnit, productionRate float64) *Engine {
	// Capture initial state
	initialState := &InitialState{
		IndustryMoney: make(map[string]float64),
		PersonMoney:   make(map[string]float64),
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
		Region:         region,
		Logger:         logging.NewLogger(true),
		CurrentTick:    0,
		WagePerHour:    wagePerHour,
		PricePerUnit:   pricePerUnit,
		ProductionRate: productionRate,
		InitialState:   initialState,
	}
}

// Run executes the simulation for a given number of ticks
func (e *Engine) Run(ticks int) {
	fmt.Printf("\n🚀 Starting Economy Simulation for %d ticks...\n", ticks)
	fmt.Printf("Region: %s\n", e.Region.Name)
	fmt.Printf("Industries: %d, People: %d, Problems: %d\n",
		len(e.Region.Industries), len(e.Region.People), len(e.Region.Problems))
	fmt.Printf("Initial Total Wealth: $%.2f\n\n", e.InitialState.TotalWealth)

	for i := 0; i < ticks; i++ {
		e.CurrentTick = i + 1
		e.processTick()
		time.Sleep(500 * time.Millisecond) // Slow down for readability
	}

	e.printFinalSummary()
}

// processTick handles one simulation tick
func (e *Engine) processTick() {
	e.Logger.LogTick(e.CurrentTick)

	// Phase 1: Industries produce goods
	e.Logger.LogEvent("📦 PRODUCTION PHASE")
	productionLogs := e.processProduction()
	e.Logger.LogEvents(productionLogs)

	// Phase 2: Labor market - people work for industries
	e.Logger.LogEvent("\n💼 LABOR MARKET PHASE")
	laborLogs := market.ProcessLaborMarket(e.Region, e.WagePerHour)
	e.Logger.LogEvents(laborLogs)

	// Phase 3: Product market - people buy products
	e.Logger.LogEvent("\n🛒 PRODUCT MARKET PHASE")
	tradeLogs := market.ProcessProductMarket(e.Region, e.PricePerUnit)
	// tradeLogsSummary := make([]string, 0)
	// for _, tradeLogItem := range tradeLogs {
	// 	tradeLogItem
	// }
	// e.Logger.LogEvents(tradeLogsSummary)
	e.Logger.LogEvents(tradeLogs)

	// Phase 4: Reset labor hours for next tick
	e.resetLaborHours()
}

// processProduction simulates industries producing goods
func (e *Engine) processProduction() []string {
	logs := make([]string, 0)

	for _, industry := range e.Region.Industries {
		for _, product := range industry.OutputProducts {
			product.Add(e.ProductionRate)
			logs = append(logs, fmt.Sprintf("✓ %s produced %.2f units of %s",
				industry.Name, e.ProductionRate, product.Name))
		}
	}

	return logs
}

// resetLaborHours resets everyone's labor hours for the next tick
func (e *Engine) resetLaborHours() {
	for _, person := range e.Region.People {
		person.LaborHours = 8.0 // Standard 8-hour workday
	}
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
	totalWealth := 0.0
	for _, person := range e.Region.People {
		totalWealth += person.Money
	}
	for _, industry := range e.Region.Industries {
		totalWealth += industry.Money
	}

	wealthChange := totalWealth - e.InitialState.TotalWealth

	fmt.Printf("\n💰 TOTAL WEALTH: $%.2f (Start: $%.2f, Change: %+.2f)\n", totalWealth, e.InitialState.TotalWealth, wealthChange)
	fmt.Printf("\n✅ Simulation completed successfully!\n\n")
}
