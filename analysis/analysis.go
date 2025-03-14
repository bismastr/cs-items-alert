package analysis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	alertsMap, err := a.alertsRealTime(ctx)

	msgs, close, err := a.consumer.Consume("price_updates")
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

		if v, ok := alertsMap[priceUpdate.ItemId]; ok {
			dailySummary, err := a.repo.GetDailySummaryByItem(ctx, priceUpdate.ItemId)
			if err != nil {
				log.Printf("Error getting item price history: %v", err)
				return err
			}
			change := fmt.Sprintf("%.2f%%", dailySummary.ChangePct)

			report := fmt.Sprintf("Hi <@%v> Real Time Alert when price is above  %v\n", v.DiscordId, v.Threshold)
			report += "Open  | Close | Change  |\n"
			report += "------------------------	\n"
			report += fmt.Sprintf("$%4d | $%4d | %7s | ",
				dailySummary.OpeningPrice,
				dailySummary.ClosingPrice,
				change,
			)

			a.bot.SendMessageToChannel("1276782792876888075", report)
		}
	}

	return nil
}
