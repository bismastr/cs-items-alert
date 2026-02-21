package price

import (
	"context"
	"testing"
	"time"

	postgres_mocks "github.com/bismastr/cs-price-alert/internal/repository/mocks"
	"github.com/bismastr/cs-price-alert/internal/timescale_repository"
	timescale_mocks "github.com/bismastr/cs-price-alert/internal/timescale_repository/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockBucket = pgtype.Timestamptz{Time: time.Now(), Valid: true}

// baseline=1000, second=1100 → +10%, third=900 → -10%
var mockHourRows = []timescale_repository.GetItemPriceChartByHourRow{
	{Bucket: mockBucket, OpenPrice: 1000, ClosePrice: 1000, SellListings: 50, ChangePct: 0},
	{Bucket: mockBucket, OpenPrice: 1000, ClosePrice: 1100, SellListings: 40, ChangePct: 10.0},
	{Bucket: mockBucket, OpenPrice: 1100, ClosePrice: 900, SellListings: 60, ChangePct: -10.0},
}

var mockDayRows = []timescale_repository.GetItemPriceChartByDayRow{
	{Bucket: mockBucket, OpenPrice: 1000, ClosePrice: 1100, SellListings: 50, ChangePct: 10.0},
	{Bucket: mockBucket, OpenPrice: 1100, ClosePrice: 1050, SellListings: 40, ChangePct: -4.5},
}

func newMocks() (*timescale_mocks.MockRepository, *postgres_mocks.MockRepository) {
	return new(timescale_mocks.MockRepository), new(postgres_mocks.MockRepository)
}

func TestGetItemPriceChart_1D(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceChartByHour", ctx, timescale_repository.GetItemPriceChartByHourParams{
		ItemID:   int32(1),
		Interval: "1 day",
	}).Return(mockHourRows, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceChart(ctx, 1, ChartPeriod1D)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, mockBucket.Time, result[0].Timestamp)
	// baseline row: price=1000, change_pct=0
	assert.Equal(t, int32(1000), result[0].Price)
	assert.Equal(t, float64(0), result[0].ChangePct)
	// second row: (1100-1000)/1000*100 = 10.00
	assert.Equal(t, int32(1100), result[1].Price)
	assert.Equal(t, float64(10), result[1].ChangePct)
	// third row: (900-1000)/1000*100 = -10.00
	assert.Equal(t, int32(900), result[2].Price)
	assert.Equal(t, float64(-10), result[2].ChangePct)
	mockTimescale.AssertExpectations(t)
}

func TestGetItemPriceChart_3D(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceChartByHour", ctx, timescale_repository.GetItemPriceChartByHourParams{
		ItemID:   int32(1),
		Interval: "3 days",
	}).Return(mockHourRows, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceChart(ctx, 1, ChartPeriod3D)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	// baseline row change_pct is always 0
	assert.Equal(t, float64(0), result[0].ChangePct)
	assert.Equal(t, float64(10), result[1].ChangePct)
	mockTimescale.AssertExpectations(t)
}

func TestGetItemPriceChart_7D(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceChartByDay", ctx, timescale_repository.GetItemPriceChartByDayParams{
		ItemID:   int32(1),
		Interval: "7 days",
	}).Return(mockDayRows, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceChart(ctx, 1, ChartPeriod7D)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int32(1100), result[0].Price)
	// ChangePct not computed in GetItemPriceChartByDay, always 0
	assert.Equal(t, float64(0), result[0].ChangePct)
	assert.Equal(t, float64(0), result[1].ChangePct)
	mockTimescale.AssertExpectations(t)
}

func TestGetItemPriceChart_30D(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceChartByDay", ctx, timescale_repository.GetItemPriceChartByDayParams{
		ItemID:   int32(1),
		Interval: "30 days",
	}).Return(mockDayRows, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceChart(ctx, 1, ChartPeriod30D)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int32(1100), result[0].Price)
	// ChangePct not computed in GetItemPriceChartByDay, always 0
	assert.Equal(t, float64(0), result[0].ChangePct)
	mockTimescale.AssertExpectations(t)
}

func TestGetItemPriceChart_InvalidPeriod(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceChart(ctx, 1, ChartPeriod("invalid"))

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTimescale.AssertNotCalled(t, "GetItemPriceChartByDay", mock.Anything, mock.Anything)
	mockTimescale.AssertNotCalled(t, "GetItemPriceChartByHour", mock.Anything, mock.Anything)
}

func TestGetItemPriceChart_RepoError(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceChartByHour", ctx, timescale_repository.GetItemPriceChartByHourParams{
		ItemID:   int32(1),
		Interval: "1 day",
	}).Return([]timescale_repository.GetItemPriceChartByHourRow(nil), assert.AnError)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceChart(ctx, 1, ChartPeriod1D)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTimescale.AssertExpectations(t)
}

func TestGetItemPriceChart_EmptyResult(t *testing.T) {
	ctx := context.Background()
	mockTimescale, mockPostgres := newMocks()

	mockTimescale.On("GetItemPriceChartByHour", ctx, timescale_repository.GetItemPriceChartByHourParams{
		ItemID:   int32(1),
		Interval: "1 day",
	}).Return([]timescale_repository.GetItemPriceChartByHourRow{}, nil)

	service := NewPriceService(mockTimescale, mockPostgres)
	result, err := service.GetItemPriceChart(ctx, 1, ChartPeriod1D)

	// empty rows should return empty slice without panicking
	assert.NoError(t, err)
	assert.Empty(t, result)
	mockTimescale.AssertExpectations(t)
}
