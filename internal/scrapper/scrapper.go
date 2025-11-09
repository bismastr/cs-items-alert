package scrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bismastr/cs-price-alert/internal/repository"
	"github.com/bismastr/cs-price-alert/internal/steam"
	"github.com/gocolly/colly"
)

type Scrapper struct {
	collector *colly.Collector
	config    Config
	repo      *repository.Queries
}

func NewScrapper(config Config, db repository.DBTX) *Scrapper {
	s := &Scrapper{
		collector: NewCollector(config),
		config:    config,
		repo:      repository.New(db),
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
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			_, err := s.repo.CreateItem(
				ctx,
				repository.CreateItemParams{
					Name:     item.Name,
					HashName: item.HashName,
				},
			)
			if err != nil {
				log.Printf("Error inserting item %s: %v ", item.HashName, err)
			}

			cancel()
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
