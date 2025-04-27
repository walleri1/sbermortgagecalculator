// Package main
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"sbermortgagecalculator/internal/middleware"
	"sbermortgagecalculator/internal/routes"
	"sbermortgagecalculator/internal/routes/paths"
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

	// Настройка логирования
	if config.Logging.File != "" {
		logFile, err := os.OpenFile(config.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Инициализация кэша
	paths.InitCache(config.Cache.MaxSize, config.Cache.TTL)
	log.Printf("Кэш инициализирован с параметрами: max_size=%d, ttl=%v", config.Cache.MaxSize, config.Cache.TTL)

	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	routes.SetupRoutes(r)

	// Настройка CORS
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins(config.CORS.AllowedOrigins),
		handlers.AllowedMethods(config.CORS.AllowedMethods),
		handlers.AllowedHeaders(config.CORS.AllowedHeaders),
	)

	// Настройка Swagger
	if config.Swagger.Enabled {
		swaggerHandler := http.FileServer(http.Dir("./docs/swagger"))
		r.PathPrefix(config.Swagger.Path).Handler(http.StripPrefix(config.Swagger.Path, swaggerHandler))
		log.Printf("Swagger UI доступен по адресу: %s", config.Swagger.Path)
	}

	address := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Сервер запущен на порту %s\n", address)
	srv := &http.Server{
		Addr:         address,
		Handler:      corsMiddleware(r),
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
