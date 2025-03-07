// Package main
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"sbermortgagecalculator/internal/middleware"
	"sbermortgagecalculator/internal/routes"
	"sbermortgagecalculator/internal/utils"
)

func main() {
	configPath := flag.String("config", "config.yml", "The path to the configuration file")
	flag.Parse()
	if configPath == nil {
		log.Fatalln("Error set config path")
	}
	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error load config server: %v", err)
	}

	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	routes.SetupRoutes(r)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	address := fmt.Sprintf(":%d", config.Port)
	log.Printf("The server is running on the port %s\n", address)
	srv := &http.Server{
		Addr:         address,
		Handler:      corsMiddleware(r),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
