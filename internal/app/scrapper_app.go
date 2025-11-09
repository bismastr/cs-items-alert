package app

import (
	"github.com/bismastr/cs-price-alert/internal/scrapper"
)

type ScrapperApp struct {
	scraper *scrapper.Scrapper
	config  scrapper.Config
}

func NewScraperApp() (*ScrapperApp, error) {
	config := scrapper.DefaultConfig()
	scraper := scrapper.NewScrapper(config)

	return &ScrapperApp{
		scraper: scraper,
		config:  config,
	}, nil
}

func (app *ScrapperApp) Start() error {
	return app.scraper.Start()
}
