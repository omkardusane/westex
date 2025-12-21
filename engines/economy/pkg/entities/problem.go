package entities

var problemIDCounter = 0

// Problem represents a high-level need or issue in the economy
// Examples: food, water, entertainment, civil-infra
type Problem struct {
	ID          int
	Name        string
	Description string
	Severity    float32 // 0.0 to 1.0, how critical this problem is
	Demand      float32 // Calculated demand based on population sentiments
	IsBasicNeed bool    // true for survival needs (food, water), false for pleasures (entertainment)
}

// NewProblem creates a new Problem instance
func NewProblem(name, description string, severity float32) *Problem {
	problemIDCounter++
	return &Problem{
		ID:          problemIDCounter,
		Name:        name,
		Description: description,
		Severity:    severity,
		Demand:      0.5,
	}
}

func (p *Problem) getName() string {
	return p.Name
}

func (p *Problem) UpdateDemand(demand float32) {
	p.Demand = demand
}
