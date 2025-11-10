package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bismastr/cs-price-alert/internal/alerts"
	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/messaging"
	"github.com/bismastr/cs-price-alert/internal/repository"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()
	log.SetOutput(os.Stdout)
	cfg := config.Load()

	ctx := context.Background()
	crn := cron.New()

	db, err := db.NewDbClient(cfg)
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	repo := repository.New(db.PostgresPool)
	publisher, err := messaging.NewPublihser(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	alertsService := alerts.NewAlertService(repo, publisher)

	crn.AddFunc("@daily", func() {
		err := alertsService.DailyPriceSummary(ctx)
		if err != nil {
			log.Fatalf("Error creating daily summary : %v", err)
		}
	})

	fmt.Println("Running schedule cron...")
	crn.Start()

	select {}
}
