package main

import (
	"context"
	"log"
	"os"

	"github.com/bismastr/cs-price-alert/analysis"
	"github.com/bismastr/cs-price-alert/bot"
	"github.com/bismastr/cs-price-alert/db"
	"github.com/bismastr/cs-price-alert/messaging"
	"github.com/bismastr/cs-price-alert/repository"
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
