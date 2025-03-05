// Package paths implements execute path service.
package paths

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"sbermortgagecalculator/internal/cache"
	"sbermortgagecalculator/internal/calculator"
	"sbermortgagecalculator/internal/models"
)

var requestIDCounter int64 = 0

// ExecuteLoanCalculation handler for mortgage calculation.
func ExecuteLoanCalculation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.LoanRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %v", err)
		writeJSONError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &request); err != nil {
		log.Printf("[ERROR] Invalid JSON format: %v", err)
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	aggregates, err := calculator.CalculateMortgageAggregates(request)
	if err != nil {
		log.Printf("[ERROR] Mortgage calculation failed: %v", err)
		writeJSONError(w, fmt.Sprintf("Calculation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	response := models.LoanResponse{
		Result: models.CalculationResult{
			Aggregates: aggregates,
			Params:     request.LoanParams,
			Program:    request.Program,
		},
	}

	requestID := atomic.AddInt64(&requestIDCounter, 1)
	cache.GetCache().Add(models.CachedLoan{
		ID:                int(requestID),
		CalculationResult: response.Result,
	})

	writeJSONResponse(w, response, http.StatusOK)
	log.Printf("[INFO] Calculation succeeded for Request ID: %d", requestID)
}
