package price

import (
	"context"
	"sync"

	"github.com/bismastr/cs-price-alert/internal/repository"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
)

type PriceService struct {
	timescaleRepo timescale_repository.Repository
	postgresRepo  repository.Repository
}

type GetPriceChange24HourResults struct {
	ItemId          int32   `json:"item_id"`
	ChangePct       float64 `json:"change_pct"`
	Name            string  `json:"name"`
	OldSellPrice    int32   `json:"old_sell_price"`
	LatestSellPrice int32   `json:"latest_sell_price"`
	IconUrl         string  `json:"icon_url"`
	Sparkline       []int32 `json:"sparkline,omitempty"`
}

type InsertPriceParams struct {
	ItemID       int32
	SellPrice    int32
	SellListings int32
}

type PriceChangeQueryParams struct {
	Query  string
	Limit  int32
	Offset int32
	SortBy string
}

func NewPriceService(timescaleRepo timescale_repository.Repository,
	postgresRepo repository.Repository) *PriceService {
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
			IconUrl:         item.IconUrl.String,
		})
	}

	return result, nil
}

func (s *PriceService) getEmptyQueryResults(ctx context.Context, params PriceChangeQueryParams) ([]GetPriceChange24HourResults, int, error) {
	var wg sync.WaitGroup
	var priceChanges []timescale_repository.GetAllPriceChangesRow
	var itemsCount int64
	var errPrice, errCount error

	wg.Add(2)
	go func() {
		defer wg.Done()
		priceChanges, errPrice = s.timescaleRepo.GetAllPriceChanges(ctx, timescale_repository.GetAllPriceChangesParams{
			Limit:  params.Limit,
			Offset: params.Offset,
			SortBy: params.SortBy,
		})
	}()

	go func() {
		defer wg.Done()
		itemsCount, errCount = s.postgresRepo.GetAllItemsCount(ctx)
	}()

	wg.Wait()

	if errPrice != nil {
		return nil, 0, errPrice
	}

	if errCount != nil {
		return nil, 0, errCount
	}

	var itemIds []int32
	for _, price := range priceChanges {
		itemIds = append(itemIds, price.ItemID)
	}

	items, err := s.postgresRepo.GetItemByID(ctx, itemIds)
	if err != nil {
		return nil, 0, err
	}

	itemMap := make(map[int32]repository.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	sparklineMap, err := s.GetItemSparklineWeekly(ctx, itemIds)
	if err != nil {
		return nil, 0, err
	}

	var result []GetPriceChange24HourResults
	for _, priceChange := range priceChanges {
		item, exists := itemMap[priceChange.ItemID]
		if !exists {
			continue
		}

		result = append(result, GetPriceChange24HourResults{
			ItemId:          item.ID,
			Name:            item.Name,
			ChangePct:       priceChange.ChangePct,
			OldSellPrice:    priceChange.OpenPrice,
			LatestSellPrice: priceChange.ClosePrice,
			IconUrl:         item.IconUrl.String,
			Sparkline:       sparklineMap[item.ID],
		})
	}

	return result, int(itemsCount), nil
}

func (s *PriceService) getSearchQueryResults(ctx context.Context, params PriceChangeQueryParams) ([]GetPriceChange24HourResults, int, error) {
	var wg sync.WaitGroup
	var priceChanges []timescale_repository.SearchPriceChangesByNameRow
	var totalCount int64
	var errPrice, errCount error

	wg.Add(2)
	go func() {
		defer wg.Done()
		priceChanges, errPrice = s.timescaleRepo.SearchPriceChangesByName(ctx, timescale_repository.SearchPriceChangesByNameParams{
			Limit:  params.Limit,
			Offset: params.Offset,
			Query:  params.Query,
			SortBy: params.SortBy,
		})
	}()

	go func() {
		defer wg.Done()
		totalCount, errCount = s.timescaleRepo.CountSearchPriceChangesByName(ctx, params.Query)
	}()

	wg.Wait()

	if errCount != nil {
		return nil, 0, errCount
	}

	if errPrice != nil {
		return nil, 0, errPrice
	}

	var itemIds []int32
	for _, price := range priceChanges {
		itemIds = append(itemIds, price.ItemID)
	}

	items, err := s.postgresRepo.GetItemByID(ctx, itemIds)
	if err != nil {
		return nil, 0, err
	}

	itemMap := make(map[int32]repository.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	sparklineMap, err := s.GetItemSparklineWeekly(ctx, itemIds)
	if err != nil {
		return nil, 0, err
	}

	var result []GetPriceChange24HourResults
	for _, priceChange := range priceChanges {
		item, exists := itemMap[priceChange.ItemID]
		if !exists {
			continue
		}

		result = append(result, GetPriceChange24HourResults{
			ItemId:          item.ID,
			Name:            item.Name,
			ChangePct:       priceChange.ChangePct,
			OldSellPrice:    priceChange.OpenPrice,
			LatestSellPrice: priceChange.ClosePrice,
			IconUrl:         item.IconUrl.String,
			Sparkline:       sparklineMap[item.ID],
		})
	}

	return result, int(totalCount), nil
}

func (s *PriceService) GetSearchPriceChanges(ctx context.Context, params PriceChangeQueryParams) ([]GetPriceChange24HourResults, int, error) {
	if params.Query == "" {
		return s.getEmptyQueryResults(ctx, params)
	}
	return s.getSearchQueryResults(ctx, params)
}

func (s *PriceService) GetItemSparklineWeekly(ctx context.Context, itemID []int32) (map[int32][]int32, error) {
	sparklines, err := s.timescaleRepo.GetItemSparklineWeekly(ctx, itemID)
	if err != nil {
		return nil, err
	}

	result := make(map[int32][]int32)
	for _, sparkline := range sparklines {
		result[sparkline.ItemID] = sparkline.Sparkline
	}

	return result, nil
}

func (s *PriceService) FormatPrice(cents int32) float64 {
	return float64(cents) / 100
}
