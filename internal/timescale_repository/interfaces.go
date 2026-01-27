package timescale_repository

import "context"

type Repository interface {
	Get24HourPricesChanges(ctx context.Context) ([]Get24HourPricesChangesRow, error)
	InsertPrice(ctx context.Context, params InsertPriceParams) error
	GetPriceChangesByItemIDs(ctx context.Context, arg GetPriceChangesByItemIDsParams) ([]GetPriceChangesByItemIDsRow, error)
	GetAllPriceChanges(ctx context.Context, arg GetAllPriceChangesParams) ([]GetAllPriceChangesRow, error)
	SearchPriceChangesByName(ctx context.Context, arg SearchPriceChangesByNameParams) ([]SearchPriceChangesByNameRow, error)
	CountSearchPriceChangesByName(ctx context.Context, query string) (int64, error)
	GetItemSparklineWeekly(ctx context.Context, itemID []int32) ([]GetItemSparklineWeeklyRow, error)
}
