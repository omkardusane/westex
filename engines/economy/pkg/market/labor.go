package market

import (
	"fmt"
	"westex/engines/economy/pkg/entities"
)

// LaborTransaction represents a person working for an industry
type LaborTransaction struct {
	Person   *entities.Person
	Industry *entities.Industry
	Hours    float32
	Wage     float32 // Payment per hour
}

// ExecuteLaborTransaction processes a person renting their time to an industry
func ExecuteLaborTransaction(person *entities.Person, industry *entities.Industry, hours float32, wagePerHour float32) (bool, string) {
	// Check if person has enough labor hours
	if person.LaborHours < hours {
		return false, fmt.Sprintf("Person %s doesn't have enough labor hours (has %.2f, needs %.2f)",
			person.Name, person.LaborHours, hours)
	}

	totalWage := hours * wagePerHour

	// Check if industry can afford to pay
	if industry.Money < totalWage {
		return false, fmt.Sprintf("Industry %s cannot afford wage of %.2f (has %.2f)",
			industry.Name, totalWage, industry.Money)
	}

	// Execute transaction
	person.LaborHours -= hours
	person.Money += totalWage
	industry.Money -= totalWage

	return true, fmt.Sprintf("✓ %s worked %.2f hours for %s, earned %.2f",
		person.Name, hours, industry.Name, totalWage)
}

// ProcessLaborMarket simulates labor transactions in a region
func ProcessLaborMarket(region *entities.Region, wagePerHour float32) []string {
	logs := make([]string, 0)

	for _, industry := range region.Industries {
		laborNeeded := industry.LaborNeeded

		// Distribute labor among people
		for _, person := range region.People {
			if laborNeeded <= 0 {
				break
			}

			hoursToWork := laborNeeded
			if hoursToWork > person.LaborHours {
				hoursToWork = person.LaborHours
			}

			if hoursToWork > 0 {
				success, log := ExecuteLaborTransaction(person, industry, hoursToWork, wagePerHour)
				if success {
					logs = append(logs, log)
					laborNeeded -= hoursToWork
				}
			}
		}
	}

	return logs
}
