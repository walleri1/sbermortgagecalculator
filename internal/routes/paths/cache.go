// The paths package implements cache path service
package paths

import (
	"encoding/json"
	"net/http"
	"sbermortgagecalculator/internal/models"
)

// GetCachedLoans handler for getting the cache of all calculations
func GetCachedLoans(w http.ResponseWriter, r *http.Request) {
	// TODO: release logic path /cache
	// Mock
	cachedLoans := []models.CachedLoan{
		{
			ID: 0,
			Result: models.CalculationResult{
				Params: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{
					Salary: true,
				},
				Aggregates: models.Aggregates{
					Rate:            8,
					LoanSum:         4000000,
					MonthlyPayment:  33458,
					Overpayment:     4029920,
					LastPaymentDate: "2044-02-18",
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cachedLoans)
}
