package repository

import "context"

type InsertItem struct {
	Name         string
	HashName     string
	SellPrice    int
	SellListings int
}

const insertItem = `
INSERT INTO items (name, hash_name, sell_price, sell_listings)
VALUES ($1, $2, $3, $4)
ON CONFLICT (hash_name) 
DO UPDATE SET sell_price = $3, sell_listings = $4
`

func (q *Queries) InsertItem(ctx context.Context, item InsertItem) error {
	_, err := q.db.Exec(ctx, insertItem, item.Name, item.HashName, item.SellPrice, item.SellListings)
	return err
}
