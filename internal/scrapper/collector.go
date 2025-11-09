package scrapper

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func NewCollector(config Config) *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent(config.UserAgent),
	)

	c.Limit(
		&colly.LimitRule{
			Delay:       config.BaseDelay,
			RandomDelay: config.RandomDelay,
			Parallelism: 1,
			DomainGlob:  "*steamcommunity.*",
		},
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "application/json")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Cache-Control", "no-cache")
	})

	extensions.RandomUserAgent(c)
	return c
}
