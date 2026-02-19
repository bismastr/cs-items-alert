package price

import (
	"context"
	"time"
)

type ChartPeriod string

const (
	ChartPeriod1D  ChartPeriod = "1d"
	ChartPeriod3D  ChartPeriod = "3d"
	ChartPeriod7D  ChartPeriod = "7d"
	ChartPeriod30D ChartPeriod = "30d"
)

type PriceChartResult struct {
	Timestamp time.Time `json:"timestamp"`
	Price     int32     `json:"price"`
}

func (s *PriceService) GetItemPriceChartByHour(ctx context.Context, itemID int32, period ChartPeriod) {

	switch period {
	case ChartPeriod1D:
	case ChartPeriod3D:
	case ChartPeriod7D:
	case ChartPeriod30D:
	}

}
