package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bismastr/cs-price-alert/db"
	messaaging "github.com/bismastr/cs-price-alert/messaging"
	"github.com/bismastr/cs-price-alert/price"
	"github.com/bismastr/cs-price-alert/repository"
	"github.com/bismastr/cs-price-alert/steam"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var (
	baseUrl = "https://steamcommunity.com/market/search/render/?count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&norender=1&category_730_Type=tag_CSGO_Type_WeaponCase&category_730_ItemSet%5B%5D=any&category_730_ProPlayer%5B%5D=any&category_730_StickerCapsule%5B%5D=any&category_730_Tournament%5B%5D=any&category_730_TournamentTeam%5B%5D=any&category_730_Type%5B%5D=tag_CSGO_Type_WeaponCase&category_730_Weapon%5B%5D=any&appid=730"
)

func main() {
	godotenv.Load()
	log.SetOutput(os.Stdout)

	crn := cron.New()

	db, err := db.NewDbClient()
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	repo := repository.New(db.Pool)
	publisher, err := messaaging.NewPublihser(os.Getenv("RMQ_URL"))
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	priceService := price.NewPriceService(repo, publisher)

	_, err = crn.AddFunc("@hourly", func() {
		scrapper(context.Background(), priceService)
	})

	if err != nil {
		log.Fatalln("cannot run cron")
	}

	crn.Start()
	log.Println("Cron scheduler started")

	select {}
}

func scrapper(ctx context.Context, priceService *price.PriceService) {
	log.Println("Starting scrape job...")
	c := defaultCollector(5*time.Second, 2*time.Second) // Increased base delay to 5s + 2s jitter

	var result steam.SteamSearchResponse
	c.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &result)
		if err != nil {
			log.Printf("Error unmarshalling response: %v", err)
			return
		}

		// Check if we got empty results, which might indicate rate limiting
		if len(result.Results) == 0 {
			log.Printf("Received empty results for URL %s, possible rate limiting", r.Request.URL)
			// Treat empty results as an error and retry
			r.Request.Retry()
			return
		}

		for _, item := range result.Results {
			insertItem := repository.InsertItem{
				Name:         item.Name,
				HashName:     item.HashName,
				SellPrice:    item.SellPrice,
				SellListings: item.SellListings,
			}

			err := priceService.InsertItem(ctx, insertItem)
			if err != nil {
				log.Printf("Error inserting item: %v", err)
			}
		}
	})

	// Add a global rate limit counter
	var rateLimitHits int

	for start := 0; start <= 428; start += 10 {
		select {
		case <-ctx.Done():
			log.Println("Scraping cancelled")
			return
		default:
			// If we've hit rate limits multiple times, add an adaptive delay
			if rateLimitHits > 0 {
				delay := time.Duration(rateLimitHits*30) * time.Second
				if delay > 5*time.Minute {
					delay = 5 * time.Minute // Cap at 5 minutes
				}
				log.Printf("Adding adaptive delay of %v due to rate limiting", delay)
				time.Sleep(delay)
			}

			url := fmt.Sprintf("%s&start=%d", baseUrl, start)
			if err := c.Visit(url); err != nil {
				log.Printf("Visit error: %v", err)
				// Increment rate limit counter if it looks like a rate limit error
				if err.Error() == "Forbidden" || err.Error() == "Too Many Requests" {
					rateLimitHits++
					log.Printf("Rate limit detected, hits: %d", rateLimitHits)
					// Add a longer delay before continuing
					time.Sleep(2 * time.Minute)
					// Retry this same URL in the next iteration
					start -= 10
				}
			}
		}
	}
}

func defaultCollector(baseDelay, randomDelay time.Duration) *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	c.Limit(
		&colly.LimitRule{
			Delay:       baseDelay,
			RandomDelay: randomDelay, // Adds jitter
			Parallelism: 1,
			DomainGlob:  "*steamcommunity.*",
		},
	)

	c.OnError(func(r *colly.Response, err error) {
		retries := 0
		if r.Request.Ctx != nil {
			if val := r.Request.Ctx.GetAny("retries"); val != nil {
				retries = val.(int)
			}
		}

		if retries < 5 {
			retries++
			log.Printf("Retrying (%d) URL %s: %v", retries, r.Request.URL, err)
			r.Request.Ctx.Put("retries", retries)

			// Exponential backoff with jitter
			backoff := time.Duration(60*(retries*retries)) * time.Second
			log.Printf("Waiting %v before retry", backoff)
			time.Sleep(backoff)
			r.Request.Retry()
		} else {
			log.Printf("Failed URL %s after 5 retries: %v", r.Request.URL, err)
		}
	})

	extensions.RandomUserAgent(c)
	return c
}
