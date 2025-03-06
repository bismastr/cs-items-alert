package repository

import "context"

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
`

func (q *Queries) InsertItem(ctx context.Context, item InsertItem) error {
	_, err := q.db.Exec(ctx, insertItem, item.Name, item.HashName, item.SellPrice, item.SellListings)
	return err
}
