package price

import (
	"context"
)

type GetItemPriceStatsResult struct {
	Interval  string `json:"interval"`
	Label     string `json:"label"`
	HighPrice int32  `json:"high_price"`
	LowPrice  int32  `json:"low_price"`
}

func (s *PriceService) GetItemPriceStats(ctx context.Context, itemID int32) ([]GetItemPriceStatsResult, error) {
	stats, err := s.timescaleRepo.GetItemPriceStats(ctx, itemID)
	if err != nil {
		return nil, err
	}

	results := []GetItemPriceStatsResult{
		{
			Interval:  "7d",
			Label:     "7 days",
			HighPrice: stats.High7d,
			LowPrice:  stats.Low7d,
		},
		{
			Interval:  "1m",
			Label:     "1 month",
			HighPrice: stats.High1m,
			LowPrice:  stats.Low1m,
		},
		{
			Interval:  "3m",
			Label:     "3 months",
			HighPrice: stats.High3m,
			LowPrice:  stats.Low3m,
		},
		{
			Interval:  "6m",
			Label:     "6 months",
			HighPrice: stats.High6m,
			LowPrice:  stats.Low6m,
		},
	}
	return results, nil
}
