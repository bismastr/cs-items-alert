package main

import (
	"log"

	app "github.com/bismastr/cs-price-alert/internal/app/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log.Println("Starting API server...")

	api, err := app.NewApiApp()
	if err != nil {
		log.Fatalf("Failed to create API app: %v", err)
	}
	defer api.Close()

	if err := api.Start(); err != nil {
		log.Fatalf("API server error: %v", err)
	}
}
