package price

import (
	"context"

	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
)

type PriceRepository interface {
	Get24HourPricesChanges(ctx context.Context) ([]timescale_repository.Get24HourPricesChangesRow, error)
	InsertPrice(ctx context.Context, params timescale_repository.InsertPriceParams) error
}

type PriceService struct {
	repo PriceRepository
}

type GetPriceChange24Hour struct {
	ItemId    int32
	ChangePct float64
}

func NewPriceService(repo PriceRepository) *PriceService {
	return &PriceService{repo: repo}
}

func (s *PriceService) InsertItem(ctx context.Context, item timescale_repository.InsertPriceParams) error {
	err := s.repo.InsertPrice(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (s *PriceService) GetLatestPriceByHashName(ctx context.Context) ([]GetPriceChange24Hour, error) {
	prices, err := s.repo.Get24HourPricesChanges(ctx)
	if err != nil {
		return nil, err
	}

	var result []GetPriceChange24Hour
	for _, price := range prices {
		changePct := float64((price.LatestSellPrice - price.OldSellPrice)) / float64(price.OldSellPrice) * 100
		result = append(result, GetPriceChange24Hour{
			ItemId:    price.ItemID,
			ChangePct: changePct,
		})
	}

	return result, nil
}
