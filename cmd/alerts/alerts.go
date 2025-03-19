package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bismastr/cs-price-alert/alerts"
	"github.com/bismastr/cs-price-alert/db"
	"github.com/bismastr/cs-price-alert/messaging"
	"github.com/bismastr/cs-price-alert/repository"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	crn := cron.New()

	db, err := db.NewDbClient()
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	repo := repository.New(db.Pool)
	publisher, err := messaging.NewPublihser(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	alertsService := alerts.NewAlertService(repo, publisher)

	crn.AddFunc("@daily", func() {
		alertsService.DailyPriceSummary(ctx)
	})

	fmt.Println("Running schedule cron...")
	crn.Start()

	select {}
}
