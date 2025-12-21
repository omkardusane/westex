package entities

var resourceIDCounter = 0

// Resource represents a material or commodity that can be consumed or produced
type Resource struct {
	ID               int
	Name             string
	Quantity         float32 // Can change over time
	Unit             string  // e.g., "kg", "liters", "units"
	IsFree           bool    // true for government-controlled resources (land, water, minerals)
	RegenerationRate float32 // units regenerated per tick (e.g., forests regrow)
}

// NewResource creates a new Resource instance
func NewResource(name string, unit string) *Resource {
	resourceIDCounter++
	return &Resource{
		ID:       resourceIDCounter,
		Name:     name,
		Quantity: 0,
		Unit:     unit,
	}
}

// Add increases the resource quantity
func (r *Resource) Add(amount float32) {
	r.Quantity += amount
}

// Consume decreases the resource quantity
// Returns true if successful, false if insufficient quantity
func (r *Resource) Consume(amount float32) bool {
	if r.Quantity >= amount {
		r.Quantity -= amount
		return true
	}
	return false
}
