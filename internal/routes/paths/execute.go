// Package paths implements execute path service.
package paths

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"sbermortgagecalculator/internal/cache"
	"sbermortgagecalculator/internal/calculator"
	"sbermortgagecalculator/internal/models"
)

var ID = 0

// ExecuteLoanCalculation handler for mortgage calculation.
func ExecuteLoanCalculation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.LoanRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	aggregates, err := calculator.CalculateMortgageAggregates(request)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(models.LoanResponse{
		Result: models.CalculationResult{
			Aggregates: aggregates,
			Params:     request.LoanParams,
			Program:    request.Program,
		},
	})
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
	cache.GetCache().Add(models.CachedLoan{
		ID: ID,
		CalculationResult: models.CalculationResult{
			Aggregates: aggregates,
			Params:     request.LoanParams,
			Program:    request.Program,
		},
	})
	ID++
}
