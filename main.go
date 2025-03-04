package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bismastr/cs-price-alert/db"
	"github.com/bismastr/cs-price-alert/items"
	"github.com/bismastr/cs-price-alert/repository"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/joho/godotenv"
)

var (
	baseUrl = "https://steamcommunity.com/market/search/render/?count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&norender=1&category_730_Type=tag_CSGO_Type_WeaponCase&category_730_ItemSet%5B%5D=any&category_730_ProPlayer%5B%5D=any&category_730_StickerCapsule%5B%5D=any&category_730_Tournament%5B%5D=any&category_730_TournamentTeam%5B%5D=any&category_730_Type%5B%5D=tag_CSGO_Type_WeaponCase&category_730_Weapon%5B%5D=any&appid=730"
)

func main() {
	godotenv.Load()
	db, err := db.NewDbClient()
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	repo := repository.New(db.Pool)
	ctx := context.Background()

	c := defaultCollector(1 * time.Second)

	var result items.SteamSearchResponse
	c.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &result)
		if err != nil {
			panic(err)
		}

		for _, item := range result.Results {
			err := repo.InsertItem(ctx, repository.InsertItem{
				Name:         item.Name,
				HashName:     item.HashName,
				SellPrice:    item.SellPrice,
				SellListings: item.SellListings,
			})
			if err != nil {
				panic(err)
			}
		}
	})

	for start := 100; start <= 400; start += 100 {
		url := fmt.Sprintf("%s&start=%d", baseUrl, start)
		c.Visit(url)
	}
}

func defaultCollector(delay time.Duration) *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	c.Limit(
		&colly.LimitRule{
			Delay:       delay,
			Parallelism: 1,
			DomainGlob:  "*steamcommunity.*",
		},
	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.Request.Body)
		panic(err.Error())
	})

	extensions.RandomUserAgent(c)

	return c
}
