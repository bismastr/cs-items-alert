package price

import (
	"context"
	"testing"

	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
	timescale_mocks "github.com/bismastr/cs-price-alert/internal/timescale_repository/mocks"
	"github.com/stretchr/testify/assert"
)

var mockStatsRow = timescale_repository.GetItemPriceStatsRow{
	High7d: 1500,
	Low7d:  900,
	High1m: 2000,
	Low1m:  800,
	High3m: 2500,
	Low3m:  700,
	High6m: 3000,
	Low6m:  600,
}

func TestGetItemPriceStats_Success(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceStats", ctx, int32(1)).Return(mockStatsRow, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceStats(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 4)

	assert.Equal(t, "7d", result[0].Interval)
	assert.Equal(t, "7 days", result[0].Label)
	assert.Equal(t, int32(1500), result[0].HighPrice)
	assert.Equal(t, int32(900), result[0].LowPrice)

	assert.Equal(t, "1m", result[1].Interval)
	assert.Equal(t, "1 month", result[1].Label)
	assert.Equal(t, int32(2000), result[1].HighPrice)
	assert.Equal(t, int32(800), result[1].LowPrice)

	assert.Equal(t, "3m", result[2].Interval)
	assert.Equal(t, "3 months", result[2].Label)
	assert.Equal(t, int32(2500), result[2].HighPrice)
	assert.Equal(t, int32(700), result[2].LowPrice)

	assert.Equal(t, "6m", result[3].Interval)
	assert.Equal(t, "6 months", result[3].Label)
	assert.Equal(t, int32(3000), result[3].HighPrice)
	assert.Equal(t, int32(600), result[3].LowPrice)

	mockTimescale.AssertExpectations(t)
}

func TestGetItemPriceStats_RepoError(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceStats", ctx, int32(1)).
		Return(timescale_repository.GetItemPriceStatsRow{}, assert.AnError)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceStats(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTimescale.AssertExpectations(t)
}

func TestGetItemPriceStats_ResultOrder(t *testing.T) {
	ctx := context.Background()
	mockTimescale := new(timescale_mocks.MockRepository)
	_, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceStats", ctx, int32(1)).Return(mockStatsRow, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceStats(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, []string{"7d", "1m", "3m", "6m"}, []string{
		result[0].Interval,
		result[1].Interval,
		result[2].Interval,
		result[3].Interval,
	})
}

func TestGetItemPriceStats_ZeroValues(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	// item with no price data returns zero values
	mockTimescale.On("GetItemPriceStats", ctx, int32(99)).
		Return(timescale_repository.GetItemPriceStatsRow{}, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceStats(ctx, 99)

	assert.NoError(t, err)
	assert.Len(t, result, 4)
	assert.Equal(t, int32(0), result[0].HighPrice)
	assert.Equal(t, int32(0), result[0].LowPrice)
	mockTimescale.AssertExpectations(t)
}
