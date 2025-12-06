package entities

// Industry represents a business entity that produces goods/services
type Industry struct {
	Name            string
	OwnedProblems   []*Problem  // Problems this industry solves (1-2 problems)
	InputResources  []*Resource // Resources needed for production
	OutputProducts  []*Resource // Products produced
	LaborNeeded     float32     // Hours of labor needed per time unit
	ConsumptionRate float32     // Rate at which input resources are consumed per unit labor week
	ProductionRate  float32     // Rate at which output products are produced per unit labor hour
	Money           float32     // Money owned by the industry
	LaborEmployed   float32     // Number of laborers employed per tick
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
func (i *Industry) UpdateIndustryRates(laborNeeded, consumptionRate, productionRate float32) *Industry {
	i.LaborNeeded = laborNeeded
	i.ConsumptionRate = consumptionRate
	i.ProductionRate = productionRate
	return i
}

func (i *Industry) UpdateLabor(laborNeeded float32) *Industry {
	i.LaborNeeded = laborNeeded
	return i
}

func (i *Industry) UpdateConsumptionRate(consumptionRate float32) {
	i.ConsumptionRate = consumptionRate
}

func (i *Industry) UpdateProductionrate(productionRate float32) {
	i.ProductionRate = productionRate
}

// UpdateIndustryMoney updates the industry's cash balance
func (i *Industry) UpdateIndustryMoney(amount float32) *Industry {
	i.Money += amount
	return i
}

// SetInitialCapital sets the starting capital for the industry
func (i *Industry) SetInitialCapital(amount float32) *Industry {
	i.Money = amount
	return i
}
