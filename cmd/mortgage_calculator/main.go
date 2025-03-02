// Package main
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sbermortgagecalculator/internal/routes"
	"sbermortgagecalculator/internal/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	routes.SetupRoutes(r)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	address := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("The server is running on the port %s\n", address)
	if err := http.ListenAndServe(address, corsMiddleware(r)); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
