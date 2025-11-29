package entities

// Problem represents a high-level need or issue in the economy
// Examples: food, water, entertainment, civil-infra
type Problem struct {
	Name        string
	Description string
	Severity    float64 // 0.0 to 1.0, how critical this problem is
}

// NewProblem creates a new Problem instance
func NewProblem(name, description string, severity float64) *Problem {
	return &Problem{
		Name:        name,
		Description: description,
		Severity:    severity,
	}
}
