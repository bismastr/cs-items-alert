package price

import (
	"context"

	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
)

type Service interface {
	GetPriceChange24Hour(ctx context.Context) ([]GetPriceChange24HourResults, error)
	InsertPrice(ctx context.Context, params timescale_repository.InsertPriceParams) error
	GetSearchPriceChanges(ctx context.Context, params PriceChangeQueryParams) ([]GetPriceChange24HourResults, int, error)
	GetItemPriceChart(ctx context.Context, itemID int32, period ChartPeriod) ([]PriceChartResult, error)
}
