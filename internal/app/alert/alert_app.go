package app

import (
	"context"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/bismastr/cs-price-alert/internal/messaging"
	"github.com/bismastr/cs-price-alert/internal/services/alert"
	"github.com/bismastr/cs-price-alert/internal/services/price"
)

type AlertApp struct {
	ctx          context.Context
	alertService *alert.AllertService
	db           *db.Db
}

func NewAlertApp(ctx context.Context) (*AlertApp, error) {
	cfg := config.Load()

	dbClient, err := db.NewDbClient(cfg)
	if err != nil {
		return nil, err
	}

	priceService := price.NewPriceService(dbClient)

	messagingPublisher, err := messaging.NewPublisher(cfg)
	if err != nil {
		return nil, err
	}

	alertService := alert.NewAlertService(priceService, messagingPublisher)

	return &AlertApp{
		ctx:          ctx,
		alertService: alertService}, nil
}

func (app *AlertApp) Start() error {
	return app.alertService.Alert24Hour(app.ctx)
}

func (app *AlertApp) Close() {
	app.db.PostgresPool.Close()
	app.db.TimescalePool.Close()
}
