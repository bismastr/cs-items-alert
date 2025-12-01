package price

import (
	"context"

	"github.com/bismastr/cs-price-alert/internal/repository"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
)

type TimescaleRepository interface {
	Get24HourPricesChanges(ctx context.Context) ([]timescale_repository.Get24HourPricesChangesRow, error)
	InsertPrice(ctx context.Context, params timescale_repository.InsertPriceParams) error
	GetPriceChangesByItemIDs(ctx context.Context, arg timescale_repository.GetPriceChangesByItemIDsParams) ([]timescale_repository.GetPriceChangesByItemIDsRow, error)
	GetAllPriceChanges(ctx context.Context, arg timescale_repository.GetAllPriceChangesParams) ([]timescale_repository.GetAllPriceChangesRow, error)
}

type PostgresRepository interface {
	GetItemByID(ctx context.Context, ids []int32) ([]repository.Item, error)
	SearchItemsByName(ctx context.Context, arg repository.SearchItemsByNameParams) ([]repository.SearchItemsByNameRow, error)
}

type PriceService struct {
	timescaleRepo TimescaleRepository
	postgresRepo  PostgresRepository
}

type GetPriceChange24HourResults struct {
	ItemId          int32
	ChangePct       float64
	Name            string
	OldSellPrice    int32
	LatestSellPrice int32
}

type InsertPriceParams struct {
	ItemID       int32
	SellPrice    int32
	SellListings int32
}

type Service interface {
	GetPriceChange24Hour(ctx context.Context) ([]GetPriceChange24HourResults, error)
	InsertPrice(ctx context.Context, params timescale_repository.InsertPriceParams) error
}

func NewPriceService(timescaleRepo TimescaleRepository,
	postgresRepo PostgresRepository) *PriceService {
	return &PriceService{
		timescaleRepo: timescaleRepo,
		postgresRepo:  postgresRepo,
	}
}

func NewPriceServiceWithRepos(
	timescaleRepo TimescaleRepository,
	postgresRepo PostgresRepository,
) *PriceService {
	return &PriceService{
		timescaleRepo: timescaleRepo,
		postgresRepo:  postgresRepo,
	}
}

func (s *PriceService) InsertPrice(ctx context.Context, item timescale_repository.InsertPriceParams) error {
	err := s.timescaleRepo.InsertPrice(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (s *PriceService) GetPriceChange24Hour(ctx context.Context) ([]GetPriceChange24HourResults, error) {
	priceChanges, err := s.timescaleRepo.Get24HourPricesChanges(ctx)
	if err != nil {
		return nil, err
	}

	priceMap := make(map[int32]GetPriceChange24HourResults)
	var itemdIds []int32
	for _, priceChange := range priceChanges {
		priceMap[priceChange.ItemID] = GetPriceChange24HourResults{
			ItemId:          priceChange.ItemID,
			ChangePct:       float64(priceChange.LatestSellPrice-priceChange.OldSellPrice) / float64(priceChange.OldSellPrice) * 100,
			OldSellPrice:    priceChange.OldSellPrice,
			LatestSellPrice: priceChange.LatestSellPrice,
		}
		itemdIds = append(itemdIds, priceChange.ItemID)
	}

	items, err := s.postgresRepo.GetItemByID(ctx, itemdIds)
	if err != nil {
		return nil, err
	}

	var result []GetPriceChange24HourResults
	for _, item := range items {
		priceChange, exists := priceMap[item.ID]
		if !exists {
			continue
		}

		result = append(result, GetPriceChange24HourResults{
			ItemId:          item.ID,
			ChangePct:       priceChange.ChangePct,
			Name:            item.Name,
			OldSellPrice:    priceChange.OldSellPrice,
			LatestSellPrice: priceChange.LatestSellPrice,
		})
	}

	return result, nil
}

type PriceChangeQueryParams struct {
	Query  string
	Limit  int32
	Offset int32
}

func (s *PriceService) GetSearchPriceChanges(ctx context.Context, params PriceChangeQueryParams) ([]GetPriceChange24HourResults, error) {

	if params.Query == "" {
		priceChanges, err := s.timescaleRepo.GetAllPriceChanges(ctx, timescale_repository.GetAllPriceChangesParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if err != nil {
			return nil, err
		}

		var itemsId []int32
		for _, priceChange := range priceChanges {
			itemsId = append(itemsId, priceChange.ItemID)
		}

		items, err := s.postgresRepo.GetItemByID(ctx, itemsId)
		if err != nil {
			return nil, err
		}

		var result []GetPriceChange24HourResults
		for i := 0; i < len(priceChanges); i++ {
			result = append(result, GetPriceChange24HourResults{
				ItemId:          priceChanges[i].ItemID,
				Name:            items[i].Name,
				ChangePct:       priceChanges[i].ChangePct,
				OldSellPrice:    priceChanges[i].OpenPrice,
				LatestSellPrice: priceChanges[i].ClosePrice,
			})
		}

		return result, nil
	}

	items, err := s.postgresRepo.SearchItemsByName(ctx,
		repository.SearchItemsByNameParams{
			Limit:  params.Limit,
			Name:   params.Query,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return nil, err
	}

	itemsId := make([]int32, 0, len(items))
	for _, item := range items {
		itemsId = append(itemsId, item.ID)
	}

	priceChanges, err := s.timescaleRepo.GetPriceChangesByItemIDs(ctx, timescale_repository.GetPriceChangesByItemIDsParams{
		ItemIds:    itemsId,
		MaxResults: params.Limit,
	})
	if err != nil {
		return nil, err
	}

	var result []GetPriceChange24HourResults
	for i := 0; i < len(priceChanges); i++ {
		result = append(result, GetPriceChange24HourResults{
			ItemId:          priceChanges[i].ItemID,
			Name:            items[i].Name,
			ChangePct:       priceChanges[i].ChangePct,
			OldSellPrice:    priceChanges[i].OpenPrice,
			LatestSellPrice: priceChanges[i].ClosePrice,
		})
	}

	return result, nil
}

func (s *PriceService) FormatPrice(cents int32) float64 {
	return float64(cents) / 100
}
