package price

import (
	"context"
	"testing"

	"github.com/bismastr/cs-price-alert/internal/repository"
	postgres_mocks "github.com/bismastr/cs-price-alert/internal/repository/mocks"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
	timescale_mocks "github.com/bismastr/cs-price-alert/internal/timescale_repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPriceChange24Hour_Success(t *testing.T) {
	ctx := context.Background()

	mockTimescaleRepo := new(timescale_mocks.MockRepository)
	mockTimescaleRepo.On("Get24HourPricesChanges", mock.Anything).Return([]timescale_repository.Get24HourPricesChangesRow{
		{ItemID: 1, LatestSellPrice: 1200, OldSellPrice: 1000},
		{ItemID: 2, LatestSellPrice: 800, OldSellPrice: 1000},
	}, nil)

	mockPostgresRepo := new(postgres_mocks.MockRepository)
	mockPostgresRepo.On("GetItemByID", mock.Anything, []int32{1, 2}).Return([]repository.Item{
		{ID: 1, HashName: "item1"},
		{ID: 2, HashName: "item2"},
	}, nil)

	service := NewPriceService(mockTimescaleRepo, mockPostgresRepo)

	result, err := service.GetPriceChange24Hour(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.EqualValues(t, 20.0, result[0].ChangePct)
	assert.EqualValues(t, -20.0, result[1].ChangePct)

	mockTimescaleRepo.AssertExpectations(t)
}

func TestGetPriceChange24Hour_Error(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(timescale_mocks.MockRepository)
	mockRepo.On("Get24HourPricesChanges", mock.Anything).Return(
		[]timescale_repository.Get24HourPricesChangesRow(nil),
		assert.AnError,
	)

	mockPostgresRepo := new(postgres_mocks.MockRepository)
	// GetItemByID is not called if Get24HourPricesChanges fails

	service := NewPriceService(mockRepo, mockPostgresRepo)

	result, err := service.GetPriceChange24Hour(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestInsertItem_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(timescale_mocks.MockRepository)
	mockPostgresRepo := new(postgres_mocks.MockRepository)
	params := timescale_repository.InsertPriceParams{
		ItemID:       1,
		SellPrice:    100,
		SellListings: 50,
	}

	mockRepo.On("InsertPrice", ctx, params).Return(nil)

	service := NewPriceService(mockRepo, mockPostgresRepo)

	err := service.InsertPrice(ctx, params)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestInsertItem_Error(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(timescale_mocks.MockRepository)
	params := timescale_repository.InsertPriceParams{
		ItemID:       1,
		SellPrice:    100,
		SellListings: 50,
	}

	mockRepo.On("InsertPrice", ctx, params).Return(assert.AnError)

	mockPostgresRepo := new(postgres_mocks.MockRepository)

	service := NewPriceService(mockRepo, mockPostgresRepo)

	err := service.InsertPrice(ctx, params)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSearchPriceChanges_Success(t *testing.T) {
	ctx := context.Background()

	mockTimescaleRepo := new(timescale_mocks.MockRepository)
	mockPostgresRepo := new(postgres_mocks.MockRepository)

	searchParams := repository.SearchItemsByNameParams{
		Limit: 10,
		Name:  "item",
	}

	mockPostgresRepo.On("SearchItemsByName", ctx, searchParams).Return([]repository.SearchItemsByNameRow{
		{ID: 1, Name: "item1"},
		{ID: 2, Name: "item2"},
	}, nil)

	priceChangesParams := timescale_repository.GetPriceChangesByItemIDsParams{
		ItemIds:    []int32{1, 2},
		MaxResults: 10,
	}

	mockTimescaleRepo.On("GetPriceChangesByItemIDs", ctx, priceChangesParams).Return([]timescale_repository.GetPriceChangesByItemIDsRow{
		{ItemID: 1, OpenPrice: 1000, ClosePrice: 1200, ChangePct: 20.0},
		{ItemID: 2, OpenPrice: 1000, ClosePrice: 800, ChangePct: -20.0},
	}, nil)

	mockPostgresRepo.On("SearchItemsCount", mock.Anything, "item").Return(int64(2), nil)

	service := NewPriceService(mockTimescaleRepo, mockPostgresRepo)

	result, totalCount, err := service.GetSearchPriceChanges(ctx, PriceChangeQueryParams{
		Limit: 10,
		Query: "item",
	})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.EqualValues(t, 20.0, result[0].ChangePct)
	assert.EqualValues(t, -20.0, result[1].ChangePct)
	assert.EqualValues(t, 2, totalCount)

	mockPostgresRepo.AssertExpectations(t)
	mockTimescaleRepo.AssertExpectations(t)
}

func TestSearchPriceChangesWithoutQuery_Success(t *testing.T) {
	ctx := context.Background()

	mockTimescaleRepo := new(timescale_mocks.MockRepository)
	mockPostgresRepo := new(postgres_mocks.MockRepository)

	mockPostgresRepo.On("GetAllItemsCount", mock.Anything).Return(int64(2), nil)

	priceChangesParams := timescale_repository.GetAllPriceChangesParams{
		Limit:  2,
		Offset: 0,
	}

	mockTimescaleRepo.On("GetAllPriceChanges", ctx, priceChangesParams).Return([]timescale_repository.GetAllPriceChangesRow{
		{ItemID: 1, OpenPrice: 1000, ClosePrice: 1200, ChangePct: 20.0},
		{ItemID: 2, OpenPrice: 1000, ClosePrice: 800, ChangePct: -20.0},
	}, nil)

	mockPostgresRepo.On("GetItemByID", mock.Anything, []int32{1, 2}).Return([]repository.Item{
		{ID: 1, Name: "item1"},
		{ID: 2, Name: "item2"},
	}, nil)

	service := NewPriceService(mockTimescaleRepo, mockPostgresRepo)
	result, totalCount, err := service.GetSearchPriceChanges(ctx, PriceChangeQueryParams{
		Query:  "",
		Limit:  2,
		Offset: 0,
	})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.EqualValues(t, 20.0, result[0].ChangePct)
	assert.EqualValues(t, -20.0, result[1].ChangePct)
	assert.EqualValues(t, 2, totalCount)

	mockPostgresRepo.AssertExpectations(t)
	mockTimescaleRepo.AssertExpectations(t)
}
