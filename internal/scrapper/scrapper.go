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

	//handle errors
	s.collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v \nError: %v", r.Request.URL, r, err)
	})
}

func (s *Scrapper) Start(ctx context.Context) error {
	totalPages := (s.config.TotalCount + s.config.PageSize - 1) / s.config.PageSize
	var rateLimitHits int

	for page := 0; page < totalPages; page++ {
		start := page * s.config.PageSize

		//handling rate limit
		if rateLimitHits > 0 {
			if rateLimitHits > s.config.MaxRateLimitHits {
				return fmt.Errorf("exceeded maximum rate limit hits: %d", rateLimitHits)
			}

			backoff := time.Duration(rateLimitHits*rateLimitHits) * s.config.BaseDelay
			log.Printf("Rate limit hit detected. Backing off for %v before retrying...", backoff)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		pageUrl := fmt.Sprintf("%s&start=%d", s.config.BaseUrl, start)
		err := s.visitWithRetry(ctx, pageUrl)
		if err != nil {
			log.Printf("Error visiting page %d: %v", page, err)
			page--
			rateLimitHits++

			continue
		}
	}

	return nil
}

func (s *Scrapper) visitWithRetry(ctx context.Context, url string) error {
	for attempt := 1; attempt < s.config.MaxRetries; attempt++ {
		err := s.collector.Visit(url)
		if err == nil {
			return nil //success!
		}

		if err.Error() == "Too Many Requests" {
			return err //need to handle rate limit
		}

		if attempt < s.config.MaxRetries {
			backoff := time.Duration(attempt*attempt) * s.config.BaseDelay

			select {
			case <-time.After(backoff):
				// this will continue to next attempt
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("all retry attempts failed")
}
