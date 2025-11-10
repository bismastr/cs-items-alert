package main

import (
	"context"
	"log"
	"os"

	"github.com/bismastr/cs-price-alert/internal/analysis"
	"github.com/bismastr/cs-price-alert/internal/bot"
	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/messaging"
	"github.com/bismastr/cs-price-alert/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	cfg := config.Load()

	db, err := db.NewDbClient(cfg)
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	bot := bot.NewBot()

	repo := repository.New(db.PostgresPool)
	consumer, err := messaging.NewConsumer(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	analysisService := analysis.NewAnalysisService(repo, bot, consumer)
	analysisService.PriceAnalysis(ctx)

	select {}
}
