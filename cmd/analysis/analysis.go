package main

import (
	"context"
	"log"

	"github.com/bismastr/cs-price-alert/analysis"
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

	repo := repository.New(db.Pool)
	consumer, err := messaging.NewConsumer()
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	var forever chan struct{}

	analysisService := analysis.NewAnalysisService(repo)
	err = consumer.PriceUpdateConsume(ctx, analysisService.PriceAnalysis)
	if err != nil {
		log.Printf("Err Consuming %v", err)
	}

	<-forever
}
