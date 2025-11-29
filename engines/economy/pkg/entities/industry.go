package entities

// Industry represents a business entity that produces goods/services
type Industry struct {
	Name            string
	OwnedProblems   []*Problem   // Problems this industry solves (1-2 problems)
	InputResources  []*Resource  // Resources needed for production
	OutputProducts  []*Resource  // Products produced
	LaborNeeded     float64      // Hours of labor needed per time unit
	ConsumptionRate float64      // Rate at which input resources are consumed per unit labor week
	ProductionRate  float64      // Rate at which output products are produced per unit labor week
	Money           float64      // Money owned by the industry
}

// CreateIndustry sets up the industry with name and returns a new Industry instance
func CreateIndustry(name string) *Industry {
	return &Industry{
		Name:           name,
		OwnedProblems:  make([]*Problem, 0),
		InputResources: make([]*Resource, 0),
		OutputProducts: make([]*Resource, 0),
		Money:          0,
	}
}

// SetupIndustry sets OwnedProblems, InputResources, OutputProducts
func (i *Industry) SetupIndustry(problems []*Problem, inputs []*Resource, outputs []*Resource) *Industry {
	i.OwnedProblems = problems
	i.InputResources = inputs
	i.OutputProducts = outputs
	return i
}

// UpdateIndustryRates sets LaborNeeded, ConsumptionRate, ProductionRate
func (i *Industry) UpdateIndustryRates(laborNeeded, consumptionRate, productionRate float64) *Industry {
	i.LaborNeeded = laborNeeded
	i.ConsumptionRate = consumptionRate
	i.ProductionRate = productionRate
	return i
}

// UpdateIndustryMoney updates the industry's cash balance
func (i *Industry) UpdateIndustryMoney(amount float64) *Industry {
	i.Money += amount
	return i
}
