// Package paths implements execute path service.
package paths

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"sbermortgagecalculator/internal/calculator"
	"sbermortgagecalculator/internal/models"
	"sbermortgagecalculator/internal/utils"
)

var (
	requestIDCounter int64 = -1
	mu               sync.Mutex
)

// ExecuteLoanCalculation handler for mortgage calculation.
func ExecuteLoanCalculation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		utils.WriteJSONError(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var request models.LoanRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %v", err)
		utils.WriteJSONError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}()

	if err = json.Unmarshal(body, &request); err != nil {
		log.Printf("[ERROR] Invalid JSON format: %v", err)
		utils.WriteJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация запроса
	if validationErrors := utils.ValidateLoanRequest(request); len(validationErrors) > 0 {
		// Возвращаем первую ошибку валидации
		utils.WriteJSONError(w, validationErrors[0].Message, http.StatusBadRequest)
		return
	}

	aggregates, err := calculator.CalculateMortgageAggregates(request)
	if err != nil {
		log.Printf("[ERROR] Mortgage calculation failed: %v", err)
		utils.WriteJSONError(w, fmt.Sprintf("Calculation error: %s", err.Error()), http.StatusBadRequest)
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
	loanCache.Store(requestID, models.CachedLoan{
		ID:                requestID,
		CalculationResult: response.Result,
		CreatedAt:         time.Now(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] Failed to write response: %v", err)
	}
	log.Printf("[INFO] Calculation succeeded for Request ID: %d", requestID)
}
