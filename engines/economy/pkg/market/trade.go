package market

import (
	"fmt"
	"westex/engines/economy/pkg/entities"
)

// TradeTransaction represents a purchase of products
type TradeTransaction struct {
	Buyer    *entities.Person
	Seller   *entities.Industry
	Product  *entities.Resource
	Quantity float64
	Price    float64 // Total price
}

// ExecuteTradeTransaction processes a person buying products from an industry
func ExecuteTradeTransaction(buyer *entities.Person, seller *entities.Industry, productName string, quantity float64, pricePerUnit float64) (bool, string) {
	// Find the product in industry's output
	var product *entities.Resource
	for _, p := range seller.OutputProducts {
		if p.Name == productName {
			product = p
			break
		}
	}

	if product == nil {
		return false, fmt.Sprintf("Industry %s doesn't produce %s", seller.Name, productName)
	}

	// Check if industry has enough product
	if product.Quantity < quantity {
		return false, fmt.Sprintf("Industry %s doesn't have enough %s (has %.2f, needs %.2f)",
			seller.Name, productName, product.Quantity, quantity)
	}

	totalPrice := quantity * pricePerUnit

	// Check if buyer can afford it
	if buyer.Money < totalPrice {
		return false, fmt.Sprintf("Person %s cannot afford %.2f %s (costs %.2f, has %.2f)",
			buyer.Name, quantity, productName, totalPrice, buyer.Money)
	}

	// Execute transaction
	buyer.Money -= totalPrice
	seller.Money += totalPrice
	product.Quantity -= quantity

	return true, fmt.Sprintf("✓ %s bought %.2f %s from %s for %.2f",
		buyer.Name, quantity, productName, seller.Name, totalPrice)
}

// ProcessProductMarket simulates people buying products to solve their problems
func ProcessProductMarket(region *entities.Region, pricePerUnit float64) []string {
	logs := make([]string, 0)
	var summary map[string]float64 = make(map[string]float64)
	// For each person, try to buy products that solve their problems
	for _, person := range region.People {
		problems := person.GetAllProblems()
		summary["transactions.people"] += 1
		for _, problem := range problems {
			// Find industries that solve this problem
			for _, industry := range region.Industries {
				for _, ownedProblem := range industry.OwnedProblems {
					if ownedProblem.Name == problem.Name {
						// Try to buy products from this industry
						for _, product := range industry.OutputProducts {
							// Buy a small amount based on problem severity
							quantityToBuy := problem.Severity * 10.0

							success, log := ExecuteTradeTransaction(person, industry, product.Name, quantityToBuy, pricePerUnit)

							summary[product.Name+".qty"] += quantityToBuy
							summary[product.Name+".cost"] += pricePerUnit * quantityToBuy
							summary["transactions.count"] += 1
							if success {
								logs = append(logs, log)
							}
						}
					}
				}
			}
		}
	}
	// convert summary into a log statement and append to log
	summaryLog := fmt.Sprintf("Market Summary - People: %.0f, Transactions: %.0f",
		summary["transactions.people"], summary["transactions.count"])
	for key, value := range summary {
		if key != "transactions.people" && key != "transactions.count" {
			summaryLog += fmt.Sprintf(" | %s: %.2f", key, value)
		}
	}
	logs = append(make([]string, 0), summaryLog)

	return logs
}
