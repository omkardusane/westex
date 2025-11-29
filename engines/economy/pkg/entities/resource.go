package entities

// Resource represents a material or commodity that can be consumed or produced
type Resource struct {
	Name     string
	Quantity float64 // Can change over time
	Unit     string  // e.g., "kg", "liters", "units"
}

// NewResource creates a new Resource instance
func NewResource(name string, quantity float64, unit string) *Resource {
	return &Resource{
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
	}
}

// Add increases the resource quantity
func (r *Resource) Add(amount float64) {
	r.Quantity += amount
}

// Consume decreases the resource quantity
// Returns true if successful, false if insufficient quantity
func (r *Resource) Consume(amount float64) bool {
	if r.Quantity >= amount {
		r.Quantity -= amount
		return true
	}
	return false
}
