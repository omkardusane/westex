package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RegionConfig represents the complete configuration for a region
type RegionConfig struct {
	Region     RegionInfo           `yaml:"region"`
	Problems   []ProblemConfig      `yaml:"problems"`
	Resources  []ResourceConfig     `yaml:"resources"`
	Industries []IndustryConfig     `yaml:"industries"`
	Population PopulationConfig     `yaml:"population"`
	Simulation SimulationConfig     `yaml:"simulation"`
}

// RegionInfo contains basic region information
type RegionInfo struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ProblemConfig defines a problem/need in the economy
type ProblemConfig struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Demand      float32 `yaml:"demand"`      // 0.0 to 1.0 - what % of population needs this
	IsBasicNeed bool    `yaml:"basic_need"`  // true for survival needs, false for pleasures
}

// ResourceConfig defines a resource
type ResourceConfig struct {
	Name            string  `yaml:"name"`
	Unit            string  `yaml:"unit"`
	InitialQuantity float32 `yaml:"initial_quantity"`
	IsFree          bool    `yaml:"is_free"`          // true for land, water, etc.
	RegenerationRate float32 `yaml:"regeneration_rate"` // units per tick
}

// IndustryConfig defines an industry
type IndustryConfig struct {
	Name            string   `yaml:"name"`
	SolvesProblems  []string `yaml:"solves_problems"`  // Problem names
	InputResources  []string `yaml:"input_resources"`  // Resource names
	OutputResources []string `yaml:"output_resources"` // Resource names
	LaborNeeded     float32  `yaml:"labor_needed"`     // Number of workers
	InitialCapital  float32  `yaml:"initial_capital"`  // Starting money
}

// PopulationConfig defines population structure
type PopulationConfig struct {
	TotalSize int                       `yaml:"total_size"`
	Segments  []PopulationSegmentConfig `yaml:"segments"`
}

// PopulationSegmentConfig defines a population segment
type PopulationSegmentConfig struct {
	Name        string   `yaml:"name"`
	Percentage  float32  `yaml:"percentage"`   // % of total population
	HasProblems []string `yaml:"has_problems"` // Problem names
	InitialMoney float32 `yaml:"initial_money"` // Starting money per person
	LaborHours   float32 `yaml:"labor_hours"`   // Available hours per tick
}

// SimulationConfig defines simulation parameters
type SimulationConfig struct {
	Ticks                    int     `yaml:"ticks"`
	WeeksPerTick             int     `yaml:"weeks_per_tick"`
	HoursPerWeek             float32 `yaml:"hours_per_week"`
	WagePerHour              float32 `yaml:"wage_per_hour"`
	ProfitMargin             float32 `yaml:"profit_margin"`              // e.g., 0.10 for 10%
	ConsumptionFactorPerWeek float32 `yaml:"consumption_factor_per_week"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filepath string) (*RegionConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config RegionConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// validateConfig checks if the configuration is valid
func validateConfig(config *RegionConfig) error {
	if config.Region.Name == "" {
		return fmt.Errorf("region name is required")
	}

	if len(config.Problems) == 0 {
		return fmt.Errorf("at least one problem is required")
	}

	if len(config.Industries) == 0 {
		return fmt.Errorf("at least one industry is required")
	}

	if config.Population.TotalSize <= 0 {
		return fmt.Errorf("population size must be positive")
	}

	// Validate percentages sum to ~100%
	totalPercentage := float32(0)
	for _, segment := range config.Population.Segments {
		totalPercentage += segment.Percentage
	}
	if totalPercentage < 0.99 || totalPercentage > 1.01 {
		return fmt.Errorf("population segment percentages must sum to 1.0, got %.2f", totalPercentage)
	}

	return nil
}

// SaveConfig saves configuration to a YAML file
func SaveConfig(config *RegionConfig, filepath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
