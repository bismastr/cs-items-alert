package alerts

import (
	"context"
	"encoding/json"
	"github.com/bismastr/cs-price-alert/internal/messaging"
	"github.com/bismastr/cs-price-alert/internal/repository"
)

type AlertService struct {
	repo      *repository.Queries
	publisher *messaging.Publisher
}

func NewAlertService(repo *repository.Queries, publihser *messaging.Publisher) *AlertService {
	return &AlertService{
		repo:      repo,
		publisher: publihser,
	}
}

type NotificationPriceSummary struct {
	ItemId        int     `json:"item_id"`
	ItemName      string  `json:"item_name"`
	AvgPrice      float64 `json:"avg_price"`
	MaxPrice      int     `json:"max_price"`
	MinPrice      int     `json:"min_price"`
	OpeningPrice  int     `json:"opening_price"`
	CloseingPrice int     `json:"closing_price"`
	ChangePct     float64 `json:"change_pct"`
	DiscordId     int64   `json:"discord_id"`
}

func (a *AlertService) DailyPriceSummary(ctx context.Context) error {
	daily, err := a.repo.GetDailyAlertSchedule(ctx)
	if err != nil {
		return err
	}

	for _, d := range daily {
		summary, err := a.repo.GetDailySummaryByItem(ctx, d.ItemId)
		if err != nil {
			return err
		}

		itemData, err := a.repo.GetItemData(ctx, d.ItemId)
		if err != nil {
			return err
		}

		notification := NotificationPriceSummary{
			ItemId:        summary.ItemId,
			ItemName:      itemData.Name,
			AvgPrice:      summary.AvgPrice,
			MaxPrice:      summary.MaxPrice,
			MinPrice:      summary.MinPrice,
			OpeningPrice:  summary.OpeningPrice,
			CloseingPrice: summary.ClosingPrice,
			ChangePct:     summary.ChangePct,
			DiscordId:     d.DiscordId,
		}

		message, err := json.Marshal(notification)
		if err != nil {
			return err
		}
		a.publisher.Publish("notification_price_alerts", message)
	}

	return nil
}
