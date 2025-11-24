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

func (m *MockTimescaleRepo) Get24HourPricesChanges(ctx context.Context) ([]timescale_repository.Get24HourPricesChangesRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]timescale_repository.Get24HourPricesChangesRow), args.Error(1)
}

func (m *MockTimescaleRepo) InsertPrice(ctx context.Context, params timescale_repository.InsertPriceParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
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

	service := NewPriceServiceWithRepos(mockTimescaleRepo, mockPostgresRepo)

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

	service := NewPriceServiceWithRepos(mockRepo, mockPostgresRepo)

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

	service := NewPriceServiceWithRepos(mockRepo, mockPostgresRepo)

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

	service := NewPriceServiceWithRepos(mockRepo, mockPostgresRepo)

	err := service.InsertPrice(ctx, params)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
