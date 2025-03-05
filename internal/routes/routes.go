// Package routes implements all service paths.
package routes

import (
	"github.com/gorilla/mux"

	"sbermortgagecalculator/internal/routes/paths"
)

// SetupRoutes sets handlers for paths.
func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/execute", paths.ExecuteLoanCalculation).Methods("POST")
	router.HandleFunc("/cache", paths.GetCachedLoans).Methods("GET")
}
