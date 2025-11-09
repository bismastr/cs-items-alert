package app

import (
	"context"

	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/scrapper"
)

type ScrapperApp struct {
	scraper *scrapper.Scrapper
	config  scrapper.Config
}

func NewScraperApp() (*ScrapperApp, error) {
	dbClient, err := db.NewDbClient()
	if err != nil {
		return nil, err
	}

	config := scrapper.DefaultConfig()
	scraper := scrapper.NewScrapper(config, dbClient.Pool)

	return &ScrapperApp{
		scraper: scraper,
		config:  config,
	}, nil
}

func (app *ScrapperApp) Start(ctx context.Context) error {
	return app.scraper.Start(ctx)
}
