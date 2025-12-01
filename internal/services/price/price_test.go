package price

import (
	"context"
	"testing"

	"github.com/bismastr/cs-price-alert/internal/repository"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTimescaleRepo struct {
	mock.Mock
}

type MockPostgresRepo struct {
	mock.Mock
}

func (m *MockPostgresRepo) GetItemByID(ctx context.Context, ids []int32) ([]repository.Item, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]repository.Item), args.Error(1)
}

func (m *MockPostgresRepo) SearchItemsByName(ctx context.Context, arg repository.SearchItemsByNameParams) ([]repository.SearchItemsByNameRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repository.SearchItemsByNameRow), args.Error(1)
}

func (m *MockPostgresRepo) SearchItemsCount(ctx context.Context, name string) (int64, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostgresRepo) GetAllItemsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTimescaleRepo) Get24HourPricesChanges(ctx context.Context) ([]timescale_repository.Get24HourPricesChangesRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]timescale_repository.Get24HourPricesChangesRow), args.Error(1)
}

func (m *MockTimescaleRepo) InsertPrice(ctx context.Context, params timescale_repository.InsertPriceParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockTimescaleRepo) GetPriceChangesByItemIDs(ctx context.Context, arg timescale_repository.GetPriceChangesByItemIDsParams) ([]timescale_repository.GetPriceChangesByItemIDsRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]timescale_repository.GetPriceChangesByItemIDsRow), args.Error(1)
}

func (m *MockTimescaleRepo) GetAllPriceChanges(ctx context.Context, arg timescale_repository.GetAllPriceChangesParams) ([]timescale_repository.GetAllPriceChangesRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]timescale_repository.GetAllPriceChangesRow), args.Error(1)
}

func TestGetPriceChange24Hour_Success(t *testing.T) {
	ctx := context.Background()

	mockTimescaleRepo := new(MockTimescaleRepo)
	mockTimescaleRepo.On("Get24HourPricesChanges", mock.Anything).Return([]timescale_repository.Get24HourPricesChangesRow{
		{ItemID: 1, LatestSellPrice: 1200, OldSellPrice: 1000},
		{ItemID: 2, LatestSellPrice: 800, OldSellPrice: 1000},
	}, nil)

	mockPostgresRepo := new(MockPostgresRepo)
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

	mockRepo := new(MockTimescaleRepo)
	mockRepo.On("Get24HourPricesChanges", mock.Anything).Return(
		[]timescale_repository.Get24HourPricesChangesRow(nil),
		assert.AnError,
	)

	mockPostgresRepo := new(MockPostgresRepo)
	mockPostgresRepo.On("GetItemByID", mock.Anything, []int32{1, 2}).Return([]repository.Item{
		{ID: 1, HashName: "item1"},
		{ID: 2, HashName: "item2"},
	}, nil)

	service := NewPriceService(mockRepo, mockPostgresRepo)

	result, err := service.GetPriceChange24Hour(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestInsertItem_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockTimescaleRepo)
	mockPostgresRepo := new(MockPostgresRepo)
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

	mockRepo := new(MockTimescaleRepo)
	params := timescale_repository.InsertPriceParams{
		ItemID:       1,
		SellPrice:    100,
		SellListings: 50,
	}

	mockRepo.On("InsertPrice", ctx, params).Return(assert.AnError)

	mockPostgresRepo := new(MockPostgresRepo)

	service := NewPriceService(mockRepo, mockPostgresRepo)

	err := service.InsertPrice(ctx, params)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSearchPriceChanges_Success(t *testing.T) {
	ctx := context.Background()

	mockTimescaleRepo := new(MockTimescaleRepo)
	mockPostgresRepo := new(MockPostgresRepo)

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
