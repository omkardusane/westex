package entities

// PopulationSegment represents a group of people with shared characteristics
// This defines a category of people who face similar problems
// Examples: "Urban Workers", "Rural Farmers", "Students", "Retirees"
type PopulationSegment struct {
	Name     string
	Problems []*Problem // Problems this segment faces
	Size     int        // Number of people in this segment
}

// NewPopulationSegment creates a new population segment
func NewPopulationSegment(name string, problems []*Problem, size int) *PopulationSegment {
	return &PopulationSegment{
		Name:     name,
		Problems: problems,
		Size:     size,
	}
}

// Person represents an individual in the economy
type Person struct {
	Name       string
	Segments   []*PopulationSegment // A person can belong to multiple segments
	Money      float64              // Personal wealth
	LaborHours float64              // Available labor hours per time unit
}

// NewPerson creates a new Person instance
func NewPerson(name string, initialMoney, laborHours float64) *Person {
	return &Person{
		Name:       name,
		Segments:   make([]*PopulationSegment, 0),
		Money:      initialMoney,
		LaborHours: laborHours,
	}
}

// AddSegment adds a population segment to this person
func (p *Person) AddSegment(segment *PopulationSegment) {
	p.Segments = append(p.Segments, segment)
}

// GetAllProblems returns all unique problems from all segments
func (p *Person) GetAllProblems() []*Problem {
	problemMap := make(map[string]*Problem)
	for _, segment := range p.Segments {
		for _, problem := range segment.Problems {
			problemMap[problem.Name] = problem
		}
	}
	
	problems := make([]*Problem, 0, len(problemMap))
	for _, problem := range problemMap {
		problems = append(problems, problem)
	}
	return problems
}
