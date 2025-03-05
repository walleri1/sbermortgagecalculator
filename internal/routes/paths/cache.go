// Package paths implements cache path service.
package paths

import (
	"encoding/json"
	"log"
	"net/http"

	"sbermortgagecalculator/internal/cache"
)

// GetCachedLoans handler for getting the cache of all calculations.
func GetCachedLoans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	cachedLoans := cache.GetCache().GetSortedLoans()
	if len(cachedLoans) == 0 {
		http.Error(w, `{"error": empty cach}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(cachedLoans)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
