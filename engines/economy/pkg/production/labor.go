package production

import (
	"fmt"
	"westex/engines/economy/pkg/entities"
)

// LaborPayment represents a wage payment to a worker
type LaborPayment struct {
	PersonName   string
	IndustryName string
	HoursWorked  float32
	WageRate     float32
	TotalPaid    float32
}

// PayWorkers distributes wages to workers employed by an industry
func PayWorkers(
	industry *entities.Industry,
	workers []*entities.Person,
	hoursPerWorker float32,
	wageRate float32,
) ([]LaborPayment, error) {
	payments := make([]LaborPayment, 0)
	totalWages := float32(0)

	// Calculate total wages needed
	for range workers {
		wages := hoursPerWorker * wageRate
		totalWages += wages
	}

	// Check if industry can afford
	if industry.Money < totalWages {
		return nil, fmt.Errorf("industry %s cannot afford wages: needs %.2f, has %.2f",
			industry.Name, totalWages, industry.Money)
	}

	// Pay each worker
	for _, worker := range workers {
		wages := hoursPerWorker * wageRate

		// Deduct from industry
		industry.Money -= wages

		// Pay worker
		worker.Money += wages

		// Record payment
		payments = append(payments, LaborPayment{
			PersonName:   worker.Name,
			IndustryName: industry.Name,
			HoursWorked:  hoursPerWorker,
			WageRate:     wageRate,
			TotalPaid:    wages,
		})
	}

	return payments, nil
}

// AllocateWorkers assigns workers to an industry based on labor needs
func AllocateWorkers(
	industry *entities.Industry,
	availableWorkers []*entities.Person,
) []*entities.Person {
	needed := int(industry.LaborNeeded)
	available := len(availableWorkers)

	// Take minimum of needed and available
	count := needed
	if available < needed {
		count = available
	}

	if count <= 0 {
		return []*entities.Person{}
	}

	return availableWorkers[:count]
}
