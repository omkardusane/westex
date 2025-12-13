package main

import (
	"fmt"
	"log"

	"westex/engines/economy/pkg/config"
	"westex/engines/economy/pkg/core"
)

// mainConfigExample demonstrates loading configuration from YAML
// To use this, rename to main() and comment out the main() in main.go
func mainConfigExample() {
	// Example 1: Load from config file
	runFromConfig()

	// Example 2: Programmatic setup (your existing code)
	// sim2()
}

func runFromConfig() {
	fmt.Println("=== Running simulation from config file ===\n")

	// Load configuration
	cfg, err := config.LoadConfig("configs/mumbai.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded config for: %s\n", cfg.Region.Name)
	fmt.Printf("  - %d problems defined\n", len(cfg.Problems))
	fmt.Printf("  - %d resources available\n", len(cfg.Resources))
	fmt.Printf("  - %d industries\n", len(cfg.Industries))
	fmt.Printf("  - Population: %d\n\n", cfg.Population.TotalSize)

	// Build region from config
	region, err := config.BuildRegionFromConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to build region: %v", err)
	}

	fmt.Printf("Region '%s' created successfully!\n", region.Name)
	fmt.Printf("  - Industries: %d\n", len(region.Industries))
	fmt.Printf("  - People: %d\n", len(region.People))
	fmt.Printf("  - Population Segments: %d\n\n", len(region.PopulationSegments))

	// Create and run engine
	engine := core.CreateNewEngine(region)

	// Run simulation with configured parameters
	engine.Run(cfg.Simulation.Ticks)
}
