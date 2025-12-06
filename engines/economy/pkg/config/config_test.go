package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configYAML := `
region:
  name: "Test Region"
  description: "A test region"

problems:
  - name: "Food"
    description: "Need for food"
    demand: 0.9
    basic_need: true

resources:
  - name: "Land"
    unit: "acres"
    initial_quantity: 1000
    is_free: true
    regeneration_rate: 0

industries:
  - name: "Farm"
    solves_problems:
      - "Food"
    input_resources:
      - "Land"
    output_resources:
      - "Food"
    labor_needed: 10
    initial_capital: 5000

population:
  total_size: 100
  segments:
    - name: "Workers"
      percentage: 1.0
      has_problems:
        - "Food"
      initial_money: 50
      labor_hours: 8

simulation:
  ticks: 5
  weeks_per_tick: 4
  hours_per_week: 40
  wage_per_hour: 10.0
  profit_margin: 0.10
  consumption_factor_per_week: 1.0
`

	// Write to temp file
	tmpfile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configYAML)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load config
	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify
	if config.Region.Name != "Test Region" {
		t.Errorf("Expected region name 'Test Region', got '%s'", config.Region.Name)
	}

	if len(config.Problems) != 1 {
		t.Errorf("Expected 1 problem, got %d", len(config.Problems))
	}

	if len(config.Industries) != 1 {
		t.Errorf("Expected 1 industry, got %d", len(config.Industries))
	}

	if config.Population.TotalSize != 100 {
		t.Errorf("Expected population 100, got %d", config.Population.TotalSize)
	}
}

func TestBuildRegionFromConfig(t *testing.T) {
	config := &RegionConfig{
		Region: RegionInfo{
			Name:        "Test",
			Description: "Test region",
		},
		Problems: []ProblemConfig{
			{
				Name:        "Food",
				Description: "Need food",
				Demand:      0.9,
				IsBasicNeed: true,
			},
		},
		Resources: []ResourceConfig{
			{
				Name:            "Land",
				Unit:            "acres",
				InitialQuantity: 1000,
				IsFree:          true,
			},
		},
		Industries: []IndustryConfig{
			{
				Name:            "Farm",
				SolvesProblems:  []string{"Food"},
				InputResources:  []string{"Land"},
				OutputResources: []string{"Food"},
				LaborNeeded:     10,
				InitialCapital:  5000,
			},
		},
		Population: PopulationConfig{
			TotalSize: 100,
			Segments: []PopulationSegmentConfig{
				{
					Name:         "Workers",
					Percentage:   1.0,
					HasProblems:  []string{"Food"},
					InitialMoney: 50,
					LaborHours:   8,
				},
			},
		},
	}

	region, err := BuildRegionFromConfig(config)
	if err != nil {
		t.Fatalf("Failed to build region: %v", err)
	}

	if region.Name != "Test" {
		t.Errorf("Expected region name 'Test', got '%s'", region.Name)
	}

	if len(region.Industries) != 1 {
		t.Errorf("Expected 1 industry, got %d", len(region.Industries))
	}

	if len(region.People) != 100 {
		t.Errorf("Expected 100 people, got %d", len(region.People))
	}
}
