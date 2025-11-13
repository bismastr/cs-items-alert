package price

import (
	"context"
	"testing"

	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPriceRepo struct {
	mock.Mock
}

func (m *MockPriceRepo) Get24HourPricesChanges(ctx context.Context) ([]timescale_repository.Get24HourPricesChangesRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]timescale_repository.Get24HourPricesChangesRow), args.Error(1)
}

func (m *MockPriceRepo) InsertPrice(ctx context.Context, params timescale_repository.InsertPriceParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func TestGetLatestPriceByHashName_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockPriceRepo)
	mockRepo.On("Get24HourPricesChanges", mock.Anything).Return([]timescale_repository.Get24HourPricesChangesRow{
		{ItemID: 1, LatestSellPrice: 1200, OldSellPrice: 1000},
		{ItemID: 2, LatestSellPrice: 800, OldSellPrice: 1000},
	}, nil)

	service := NewPriceService(mockRepo)

	result, err := service.GetLatestPriceByHashName(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.EqualValues(t, 20.0, result[0].ChangePct)
	assert.EqualValues(t, -20.0, result[1].ChangePct)

	mockRepo.AssertExpectations(t)
}

func TestGetLatestPriceByHashName_Error(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockPriceRepo)
	mockRepo.On("Get24HourPricesChanges", mock.Anything).Return(
		[]timescale_repository.Get24HourPricesChangesRow(nil),
		assert.AnError,
	)

	service := NewPriceService(mockRepo)

	result, err := service.GetLatestPriceByHashName(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestInsertItem_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockPriceRepo)
	params := timescale_repository.InsertPriceParams{
		ItemID:       1,
		SellPrice:    100,
		SellListings: 50,
	}

	mockRepo.On("InsertPrice", ctx, params).Return(nil)

	service := NewPriceService(mockRepo)

	err := service.InsertItem(ctx, params)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestInsertItem_Error(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockPriceRepo)
	params := timescale_repository.InsertPriceParams{
		ItemID:       1,
		SellPrice:    100,
		SellListings: 50,
	}

	mockRepo.On("InsertPrice", ctx, params).Return(assert.AnError)

	service := NewPriceService(mockRepo)

	err := service.InsertItem(ctx, params)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
