package paths

import (
	"encoding/json"
	"log"
	"net/http"
	"sbermortgagecalculator/internal/models"
	"sync"
)

var loanCache sync.Map

// getLoansFromSyncMap bypasses sync.Map and returns the CachedLoan slice.
func getLoansFromSyncMap(m *sync.Map) []models.CachedLoan {
	var loans []models.CachedLoan

	m.Range(func(_, value any) bool {
		if loan, ok := value.(models.CachedLoan); ok {
			loans = append(loans, loan)
		} else {
			log.Printf("Warning: unexpected type in sync.Map: %T\n", value)
		}
		return true
	})

	return loans
}

// writeJSONResponse writes a JSON response with the specified status code.
func writeJSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[ERROR] Failed to write JSON response: %v", err)
	}
}

// writeJSONError writes a JSON error response with the specified status code.
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	writeJSONResponse(w, map[string]string{"error": message}, statusCode)
}
