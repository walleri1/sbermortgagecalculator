// The routes package implements execute path service.
package paths

import (
	"encoding/json"
	"log"
	"net/http"
	"sbermortgagecalculator/internal/models"
)

// ExecuteLoanCalculation handler for mortgage calculation.
func ExecuteLoanCalculation(w http.ResponseWriter, r *http.Request) {
	var request models.LoanRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	// TODO: release logic path /execute
	// Mock
	response := models.LoanResponse{
		Result: models.CalculationResult{
			Params:  request.Params,
			Program: request.Program,
			Aggregates: models.Aggregates{
				Rate:            8,
				LoanSum:         4000000,
				MonthlyPayment:  33458,
				Overpayment:     4029920,
				LastPaymentDate: "2044-02-18",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Unmarshell %v", err)
	}
	w.Write(data)
}
