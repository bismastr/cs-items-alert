package scrapper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bismastr/cs-price-alert/internal/steam"
	"github.com/gocolly/colly"
)

type Scrapper struct {
	collector *colly.Collector
	config    Config
}

func NewScrapper(config Config) *Scrapper {
	s := &Scrapper{
		collector: NewCollector(config),
		config:    config,
	}

	s.setupHandlers()
	return s
}

func (s *Scrapper) setupHandlers() {
	s.collector.OnResponse(func(r *colly.Response) {
		var response steam.SteamSearchResponse
		err := json.Unmarshal(r.Body, &response)
		if err != nil {
			log.Printf("Error unmarshalling response: %v", err)
			return
		}

		for _, item := range response.Results {
			log.Printf("Item: %s, Price: %s", item.Name, item.SellPriceText)
		}
	})
}

func (s *Scrapper) Start() error {
	totalPages := (s.config.TotalCount + s.config.PageSize - 1) / s.config.PageSize

	for page := 0; page < totalPages; page++ {
		start := page * s.config.PageSize

		pageUrl := fmt.Sprintf("%s&start=%d", s.config.BaseUrl, start)
		err := s.collector.Visit(pageUrl)
		if err != nil {
			log.Printf("Error visiting page %d: %v", page, err)
		}
	}

	return nil
}
