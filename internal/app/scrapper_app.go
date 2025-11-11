package app

import (
	"context"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/scrapper"
)

type ScrapperApp struct {
	ctx     context.Context
	scraper *scrapper.Scrapper
	config  scrapper.Config
}

func NewScraperApp(ctx context.Context) (*ScrapperApp, error) {
	cfg := config.Load()

	dbClient, err := db.NewDbClient(cfg)
	if err != nil {
		return nil, err
	}

	config := scrapper.DefaultConfig()
	scraper := scrapper.NewScrapper(ctx, config, dbClient)

	return &ScrapperApp{
		ctx:     ctx,
		scraper: scraper,
		config:  config,
	}, nil
}

func (app *ScrapperApp) Start() error {
	return app.scraper.Start()
}
