package entities

// Region represents a geographic/economic area containing all entities
type Region struct {
	Name       string
	Industries []*Industry
	People     []*Person
	Resources  []*Resource // Shared/available resources in the region
	Problems   []*Problem  // All problems present in the region
}

// NewRegion creates a new Region instance
func NewRegion(name string) *Region {
	return &Region{
		Name:       name,
		Industries: make([]*Industry, 0),
		People:     make([]*Person, 0),
		Resources:  make([]*Resource, 0),
		Problems:   make([]*Problem, 0),
	}
}

// AddIndustry adds an industry to the region
func (r *Region) AddIndustry(industry *Industry) {
	r.Industries = append(r.Industries, industry)
}

// AddPerson adds a person to the region
func (r *Region) AddPerson(person *Person) {
	r.People = append(r.People, person)
}

// AddResource adds a resource to the region
func (r *Region) AddResource(resource *Resource) {
	r.Resources = append(r.Resources, resource)
}

// AddProblem adds a problem to the region
func (r *Region) AddProblem(problem *Problem) {
	r.Problems = append(r.Problems, problem)
}

// GetResource finds a resource by name
func (r *Region) GetResource(name string) *Resource {
	for _, resource := range r.Resources {
		if resource.Name == name {
			return resource
		}
	}
	return nil
}

// GetProblem finds a problem by name
func (r *Region) GetProblem(name string) *Problem {
	for _, problem := range r.Problems {
		if problem.Name == name {
			return problem
		}
	}
	return nil
}
