package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	app "github.com/bismastr/cs-price-alert/internal/app/alert"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	c := cron.New()
	alertApp, err := app.NewAlertApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create alert app: %v", err)
	}

	defer alertApp.Close()
	_, err = c.AddFunc("0 0 * * *", func() {
		if err := alertApp.Start(); err != nil {
			log.Printf("Failed to start alert app: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron function: %v", err)
	}

	if err := alertApp.Start(); err != nil {
		log.Printf("Failed to start alert app: %v", err)
	}

	c.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}
