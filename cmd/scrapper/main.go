package main

import (
	"log"

	"github.com/bismastr/cs-price-alert/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	scrapper, err := app.NewScraperApp()
	if err != nil {
		log.Fatalf("Failed to create scrapper app: %v", err)
	}

	if err := scrapper.Start(); err != nil {
		log.Fatalf("Failed to start scrapper app: %v", err)
	}
}
