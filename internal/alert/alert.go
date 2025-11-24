package alert

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bismastr/cs-price-alert/internal/messaging"
	"github.com/bismastr/cs-price-alert/internal/price"
)

type AllertService struct {
	priceService *price.PriceService
	messaging    *messaging.Publisher
}

type PriceChangesAlert struct {
	ItemId          int32   `json:"item_id"`
	ChangePct       float64 `json:"change_pct"`
	Name            string  `json:"name"`
	AlertType       string  `json:"alert_type"`
	LatestSellPrice int32   `json:"latest_price"`
	OldSellPrice    int32   `json:"old_price"`
}

func NewAlertService(priceService *price.PriceService, messaging *messaging.Publisher) *AllertService {
	return &AllertService{
		priceService: priceService,
		messaging:    messaging,
	}
}

func (s *AllertService) Alert24Hour(ctx context.Context) error {
	priceChanges, err := s.priceService.GetPriceChange24Hour(ctx)
	if err != nil {
		return fmt.Errorf("error: get price change 24 hour: %w", err)
	}

	for _, price := range priceChanges {
		if price.ChangePct < 20.0 && price.ChangePct > -20.0 {
			continue
		}

		var alertType string
		if price.ChangePct < 0 {
			alertType = AlertTypeDecrease
		} else {
			alertType = AlertTypeIncrease
		}

		message, _ := json.Marshal(PriceChangesAlert{
			ItemId:          price.ItemId,
			ChangePct:       price.ChangePct,
			Name:            price.Name,
			AlertType:       alertType,
			LatestSellPrice: s.priceService.FormatPrice(price.LatestSellPrice),
			OldSellPrice:    s.priceService.FormatPrice(price.OldSellPrice),
		})

		s.messaging.Publish(ctx, messaging.QueueDiscordAlert, message)
	}

	return nil
}
