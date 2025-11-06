package main

import (
	"context"
	"github.com/bismastr/cs-price-alert/internal/analysis"
	"github.com/bismastr/cs-price-alert/internal/bot"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/messaging"
	"github.com/bismastr/cs-price-alert/internal/repository"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	db, err := db.NewDbClient()
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	bot := bot.NewBot()

	repo := repository.New(db.Pool)
	consumer, err := messaging.NewConsumer(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	analysisService := analysis.NewAnalysisService(repo, bot, consumer)
	analysisService.PriceAnalysis(ctx)

	select {}
}
