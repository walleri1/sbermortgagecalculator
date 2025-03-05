package paths

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSONResponse writes a JSON response with the specified status code.
func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
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
