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
	ctx       context.Context
	collector *colly.Collector
	config    Config
	repo      *repository.Queries
}

func NewScrapper(ctx context.Context, config Config, db repository.DBTX) *Scrapper {
	s := &Scrapper{
		ctx:       ctx,
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
			ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)

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

func (s *Scrapper) Start() error {
	totalPages := (s.config.TotalCount + s.config.PageSize - 1) / s.config.PageSize
	var rateLimitHits int

	for page := 0; page < totalPages; page++ {
		start := page * s.config.PageSize
		log.Printf("Starting to scrape page %d (start=%d)", page, start)
		//handling rate limit
		if rateLimitHits > 0 {
			if rateLimitHits > s.config.MaxRateLimitHits {
				return fmt.Errorf("exceeded maximum rate limit hits: %d", rateLimitHits)
			}

			backoff := time.Duration(rateLimitHits*rateLimitHits) * s.config.BaseDelay
			log.Printf("Rate limit hit detected. Backing off for %v before retrying...", backoff)
			log.Printf("Recreating collector")
			s.recreateCollector()

			select {
			case <-time.After(backoff):
				log.Printf("Backoff completed, retrying page %d", page)
			case <-s.ctx.Done():
				return s.ctx.Err()
			}
		}

		pageUrl := fmt.Sprintf("%s&start=%d", s.config.BaseUrl, start)
		err := s.visitWithRetry(pageUrl)
		if err != nil {
			log.Printf("Error visiting page %d: %v", page, err)
			rateLimitHits++
			page--

			continue
		}

		//reset rate limit hits on successful visit
		rateLimitHits = 0
	}

	return nil
}

func (s *Scrapper) visitWithRetry(url string) error {
	for attempt := 1; attempt <= s.config.MaxRetries; attempt++ {
		err := s.collector.Visit(url)
		if err == nil {
			return nil //success!
		}

		if err.Error() == "Too Many Requests" {
			return err //need to handle rate limit out of this function
		}

		if attempt < s.config.MaxRetries {
			log.Printf("Error visiting attempt %d: %v", attempt, err)
			backoff := time.Duration(attempt*attempt) * s.config.BaseDelay

			select {
			case <-time.After(backoff):
			case <-s.ctx.Done():
				return s.ctx.Err()
			}
		}
	}

	return fmt.Errorf("all retry attempts failed")
}

func (s *Scrapper) recreateCollector() {
	s.collector = NewCollector(s.config)
	s.setupHandlers()
}
