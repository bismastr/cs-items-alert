package repository

import "context"

type Repository interface {
	CreateItem(ctx context.Context, arg CreateItemParams) (Item, error)
	GetAllItemsCount(ctx context.Context) (int64, error)
	GetItemByHashName(ctx context.Context, hashName string) (GetItemByHashNameRow, error)
	GetItemByID(ctx context.Context, ids []int32) ([]Item, error)
	GetItems(ctx context.Context) ([]Item, error)
	SearchItemsByName(ctx context.Context, arg SearchItemsByNameParams) ([]SearchItemsByNameRow, error)
	SearchItemsCount(ctx context.Context, name string) (int64, error)
}
