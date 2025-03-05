// Package paths implements cache path service.
package paths

import (
	"log"
	"net/http"
)

// GetCachedLoans handler for getting the cache of all calculations.
func GetCachedLoans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	cachedLoans := getLoansFromSyncMap(&loanCache)
	if len(cachedLoans) == 0 {
		log.Println("[INFO] Cache is empty, no loans to retrieve")
		writeJSONError(w, "empty cache", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, cachedLoans, http.StatusOK)
	log.Printf("[INFO] Successfully returned %d cached loan(s)", len(cachedLoans))
}
