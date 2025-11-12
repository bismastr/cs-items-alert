package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bismastr/cs-price-alert/internal/app"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	c := cron.New()
	scrapper, err := app.NewScraperApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create scrapper app: %v", err)
	}

	defer scrapper.Close()
	_, err = c.AddFunc("0 * * * *", func() {
		if err := scrapper.Start(); err != nil {
			log.Printf("Failed to start scrapper app: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron function: %v", err)
	}

	c.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}
