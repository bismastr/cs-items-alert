package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var (
    baseUrl = "https://steamcommunity.com/market/search/render/?count=100&search_descriptions=0&sort_column=popular&sort_dir=desc&norender=1&category_730_Type=tag_CSGO_Type_WeaponCase&category_730_ItemSet%5B%5D=any&category_730_ProPlayer%5B%5D=any&category_730_StickerCapsule%5B%5D=any&category_730_Tournament%5B%5D=any&category_730_TournamentTeam%5B%5D=any&category_730_Type%5B%5D=tag_CSGO_Type_WeaponCase&category_730_Weapon%5B%5D=any&appid=730"
)
func main() {
    c := defaultCollector()

    c.OnResponse(func(r *colly.Response) {
        json.Marshal(r.Body)
    })

    for start := 0; start <= 90; start += 100 {
        url := fmt.Sprintf("%s&start=%d", baseUrl, start)
        c.Visit(url)
    }
}


func defaultCollector() *colly.Collector {
    c := colly.NewCollector(
        colly.MaxDepth(1),
    )

    c.Limit(
        &colly.LimitRule{
            RandomDelay: 2 * time.Second,
            Parallelism: 2,
            DomainGlob:  "*steamcommunity.*",
        },
    )

    c.OnError(func(r *colly.Response, err error) {
        fmt.Println("error scrapper")
    })

    extensions.RandomUserAgent(c)

    return c
}
