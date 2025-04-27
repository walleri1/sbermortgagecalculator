// Package paths implements cache path service.
package paths

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"sbermortgagecalculator/internal/models"
	"sbermortgagecalculator/internal/utils"
)

var (
	loanCache *utils.Cache
)

// InitCache инициализирует кэш с заданными параметрами
func InitCache(maxSize int, ttl time.Duration) {
	loanCache = utils.NewCache(maxSize, ttl)
}

// GetCachedLoans handler for retrieving cached loan calculations.
func GetCachedLoans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var cachedLoans []models.CachedLoan
	loanCache.Range(func(key int64, value models.CachedLoan) bool {
		cachedLoans = append(cachedLoans, value)
		return true
	})

	if len(cachedLoans) == 0 {
		utils.WriteJSONError(w, "empty cache", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cachedLoans); err != nil {
		log.Printf("[ERROR] Failed to write response: %v", err)
	}
}

// GetCachedResult handler for retrieving a specific cached loan calculation by ID.
func GetCachedResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/cache/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSONError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	value, ok := loanCache.Load(id)
	if !ok {
		utils.WriteJSONError(w, "Loan calculation not found", http.StatusNotFound)
		return
	}

	cachedLoan := value

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cachedLoan); err != nil {
		log.Printf("[ERROR] Failed to write response: %v", err)
	}
}
