package main

import (
	"log"

	"github.com/rasul07/alif-task/internal/config"
	"github.com/rasul07/alif-task/internal/handlers"
	"github.com/rasul07/alif-task/internal/service"
	"github.com/rasul07/alif-task/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := storage.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	walletService := service.NewWalletService(db)

	api := handlers.NewAPI(walletService)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := api.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
