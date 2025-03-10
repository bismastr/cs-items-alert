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
    INSERT INTO items (name, hash_name, sell_price, sell_listings)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (hash_name) 
    DO UPDATE SET
        sell_price = EXCLUDED.sell_price,
        sell_listings = EXCLUDED.sell_listings
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

const getItemPrice = `
WITH current_item AS (
	SELECT id , name
	FROM items
	WHERE id = $1
),
price_history AS (
	SELECT sell_price, sell_listings, created_at
	FROM prices
	WHERE item_id = $1
	AND created_at >= NOW() - INTERVAL '5 hours'
	ORDER BY created_at DESC
)
SELECT 
	ci.id,
	ci.name,
	ph.sell_price,
	ph.sell_listings,
	ph.created_at
FROM current_item ci
LEFT JOIN price_history ph ON true
ORDER BY ph.created_at DESC
`

type ItemWithPriceHistory struct {
	ID           int
	Name         string
	CurrentPrice int
	PriceHistory []PriceEntry
	Volatility   float64
}

type PriceEntry struct {
	SellPrice   int
	SellListing int
	CreatedAt   pgtype.Timestamp
}

func (q *Queries) GetItemPrice(ctx context.Context, itemId int) (ItemWithPriceHistory, error) {
	rows, err := q.db.Query(ctx, getItemPrice, itemId)
	if err != nil {
		return ItemWithPriceHistory{}, err
	}
	defer rows.Close()

	var item ItemWithPriceHistory
	var hasData bool

	for rows.Next() {
		var entry PriceEntry
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&entry.SellPrice,
			&entry.SellListing,
			&entry.CreatedAt,
		)
		if err != nil {
			return ItemWithPriceHistory{}, err
		}

		if !hasData {
			item.CurrentPrice = entry.SellPrice
			hasData = true
		}

		item.PriceHistory = append(item.PriceHistory, entry)
	}

	return item, rows.Err()
}
