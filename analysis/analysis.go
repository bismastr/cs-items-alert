package analysis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/bismastr/cs-price-alert/bot"
	"github.com/bismastr/cs-price-alert/messaging"
	"github.com/bismastr/cs-price-alert/repository"
)

type Analysis struct {
	repo     *repository.Queries
	bot      *bot.Bot
	consumer *messaging.Consumer
}

func NewAnalysisService(repo *repository.Queries, bot *bot.Bot, consumer *messaging.Consumer) *Analysis {
	return &Analysis{
		repo:     repo,
		bot:      bot,
		consumer: consumer,
	}
}

func (a *Analysis) PriceAnalysis(ctx context.Context) error {
	log.Println("Running Price analysis...")
	msgs, close, err := a.consumer.PriceUpdateConsume("price_updates")
	if err != nil {
		return err
	}

	defer close()

	var priceUpdate struct {
		ItemId int `json:"item_id"`
	}

	for d := range msgs {
		if err := json.Unmarshal(d.Body, &priceUpdate); err != nil {
			log.Printf("Error decoding message: %v", err)
			return err
		}

		dailySummary, err := a.repo.GetDailySummaryByItem(ctx, priceUpdate.ItemId)
		if err != nil {
			log.Printf("Error getting item price history: %v", err)
			return err
		}

		content := fmt.Sprintf("Item: %d, have a price changed %v with max price: %d and min price: %d", dailySummary.ItemId, dailySummary.ChangePct, dailySummary.MaxPrice, dailySummary.MinPrice)
		log.Println(content)
		a.bot.SendMessageToChannel("1276782792876888075", content)
	}

	return nil
}

func CalculateVolatility(i *repository.ItemWithPriceHistory) {
	if len(i.PriceHistory) < 2 {
		i.Volatility = 0
		return
	}

	sum := 0.0
	prev := float64(i.PriceHistory[0].SellPrice)
	for _, entry := range i.PriceHistory[1:] {
		current := float64(entry.SellPrice)
		sum += math.Abs((current - prev) / prev)
		prev = current
	}

	i.Volatility = sum / float64(len(i.PriceHistory)-1) * 100
}
