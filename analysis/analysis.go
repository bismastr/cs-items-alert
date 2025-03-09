package analysis

import (
	"context"
	"encoding/json"
	"log"
	"math"

	"github.com/bismastr/cs-price-alert/repository"
	"github.com/rabbitmq/amqp091-go"
)

type Analysis struct {
	repo *repository.Queries
}

func NewAnalysisService(repo *repository.Queries) *Analysis {
	return &Analysis{
		repo: repo,
	}
}

func (a *Analysis) PriceAnalysis(ctx context.Context, d amqp091.Delivery) error {
	log.Println("Running Price analysis...")
	var priceUpdate struct {
		ItemId int `json:"item_id"`
	}
	if err := json.Unmarshal(d.Body, &priceUpdate); err != nil {
		log.Printf("Error decoding message: %v", err)
		return err
	}

	priceHistory, err := a.repo.GetItemPrice(ctx, priceUpdate.ItemId)
	if err != nil {
		log.Printf("Error getting item price history: %v", err)
		return err
	}

	CalculateVolatility(&priceHistory)
	log.Printf("Item: %s%d, have a price changed %v percent from the last 5 hours", priceHistory.Name, priceHistory.ID, priceHistory.Volatility)

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
