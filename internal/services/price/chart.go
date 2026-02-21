package price

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
)

type ChartPeriod string

const (
	ChartPeriod1D  ChartPeriod = "1d"
	ChartPeriod3D  ChartPeriod = "3d"
	ChartPeriod7D  ChartPeriod = "7d"
	ChartPeriod30D ChartPeriod = "30d"
	ChartPeriod1M  ChartPeriod = "1m"
	ChartPeriod3M  ChartPeriod = "3m"
	ChartPeriod6M  ChartPeriod = "6m"
)

type PriceChartResult struct {
	Timestamp time.Time `json:"timestamp"`
	Price     int32     `json:"price"`
	ChangePct float64   `json:"change_pct"`
}

func (s *PriceService) GetItemPriceChart(ctx context.Context, itemID int32, period ChartPeriod) ([]PriceChartResult, error) {

	switch period {
	case ChartPeriod1D:
		return s.GetItemPriceChartByHour(ctx, itemID, "1 day")
	case ChartPeriod3D:
		return s.GetItemPriceChartByHour(ctx, itemID, "3 days")
	case ChartPeriod7D:
		return s.GetItemPriceChartByDay(ctx, itemID, "7 days")
	case ChartPeriod30D:
		return s.GetItemPriceChartByDay(ctx, itemID, "30 days")
	case ChartPeriod1M:
		return s.GetItemPriceChartByDay(ctx, itemID, "1 month")
	case ChartPeriod3M:
		return s.GetItemPriceChartByDay(ctx, itemID, "3 months")
	case ChartPeriod6M:
		return s.GetItemPriceChartByDay(ctx, itemID, "6 months")
	default:
		return nil, fmt.Errorf("invalid period: %s, must be one of: 1d, 3d, 7d, 30d, 1m, 3m, 6m", period)
	}

}

func (s *PriceService) GetItemPriceChartByHour(ctx context.Context, itemId int32, interval string) ([]PriceChartResult, error) {
	rows, err := s.timescaleRepo.GetItemPriceChartByHour(ctx, timescale_repository.GetItemPriceChartByHourParams{
		ItemID:   itemId,
		Interval: interval,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []PriceChartResult{}, nil
	}

	baselinePrice := rows[0].ClosePrice

	results := make([]PriceChartResult, len(rows))
	for i, row := range rows {
		changePct := float64(row.ClosePrice-baselinePrice) / float64(baselinePrice) * 100
		changePct = math.Round(changePct*100) / 100

		results[i] = PriceChartResult{
			Timestamp: row.Bucket.Time,
			Price:     row.ClosePrice,
			ChangePct: changePct,
		}
	}
	return results, nil
}

func (s *PriceService) GetItemPriceChartByDay(ctx context.Context, itemId int32, interval string) ([]PriceChartResult, error) {
	rows, err := s.timescaleRepo.GetItemPriceChartByDay(ctx, timescale_repository.GetItemPriceChartByDayParams{
		ItemID:   itemId,
		Interval: interval,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []PriceChartResult{}, nil
	}

	baselinePrice := rows[0].ClosePrice

	results := make([]PriceChartResult, len(rows))
	for i, row := range rows {
		changePct := float64(row.ClosePrice-baselinePrice) / float64(baselinePrice) * 100
		changePct = math.Round(changePct*100) / 100

		results[i] = PriceChartResult{
			Timestamp: row.Bucket.Time,
			Price:     row.ClosePrice,
			ChangePct: changePct,
		}
	}
	return results, nil
}
