package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type InsertItem struct {
	Id           int
	Name         string
	HashName     string
	SellPrice    int
	SellListings int
}

const insertItem = `
WITH inserted_item AS (
    INSERT INTO items (name, hash_name)
    VALUES ($1, $2)
    ON CONFLICT (hash_name) 
   	DO UPDATE SET name = EXCLUDED.name 
    RETURNING id 
)
INSERT INTO prices (item_id, sell_price, sell_listings)
SELECT id, $3, $4 FROM inserted_item
RETURNING item_id
`

func (q *Queries) InsertItem(ctx context.Context, arg InsertItem) (int, error) {
	var id int
	err := q.db.QueryRow(ctx, insertItem, arg.Name, arg.HashName, arg.SellPrice, arg.SellListings).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

type GetItem struct {
	Name     string
	HashName string
}

const getItem = `
SELECT name, hash_name FROM items 
WHERE id = $1          
`

func (q *Queries) GetItemData(ctx context.Context, itemId int) (GetItem, error) {
	row := q.db.QueryRow(ctx, getItem, itemId)
	var i GetItem
	err := row.Scan(
		&i.Name,
		&i.HashName,
	)

	return i, err
}

const getDailySummaryByItem = `
SELECT 
	item_id, 
	dps.bucket , 
	dps.avg_price , 
	dps.max_price, 
	dps.min_price, 
	dps.opening_price,
	dps.closing_price,
	dps.data_points,
	(closing_price - opening_price)::FLOAT / opening_price * 100 AS change_pct 
	FROM daily_price_summary dps 
WHERE dps.item_id = $1
ORDER BY bucket DESC  
LIMIT 1          
`

type GetDailySummaryByItem struct {
	ItemId       int
	Bucket       pgtype.Timestamp
	AvgPrice     float64
	MaxPrice     int
	MinPrice     int
	OpeningPrice int
	ClosingPrice int
	DataPoints   int
	ChangePct    float64
}

func (q *Queries) GetDailySummaryByItem(ctx context.Context, itemId int) (GetDailySummaryByItem, error) {
	row := q.db.QueryRow(ctx, getDailySummaryByItem, itemId)
	var i GetDailySummaryByItem
	err := row.Scan(
		&i.ItemId,
		&i.Bucket,
		&i.AvgPrice,
		&i.MaxPrice,
		&i.MinPrice,
		&i.OpeningPrice,
		&i.ClosingPrice,
		&i.DataPoints,
		&i.ChangePct,
	)

	return i, err
}
