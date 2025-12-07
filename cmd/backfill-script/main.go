package main

import (
	"context"
	"log"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/bismastr/cs-price-alert/internal/db"
	"github.com/joho/godotenv"
)

type Item struct {
	Id   int64
	Name string
}

func main() {
	godotenv.Load()
	cfg := config.Load()
	ctx := context.Background()

	dbClient, err := db.NewDbClient(cfg)
	if err != nil {
		log.Printf("Failed to create db client: %v", err)
	}
	defer dbClient.PostgresPool.Close()
	defer dbClient.TimescalePool.Close()

	rows, err := dbClient.PostgresPool.Query(ctx, `SELECT id, name FROM items`)
	if err != nil {
		log.Printf("Failed to get rows: %v", err)
	}

	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Id, &item.Name); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	const batchSize = 50
	totalUpdated := 0

	for i := 0; i < len(items); i += batchSize {

		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		itemsBatch := items[i:end]
		tx, err := dbClient.TimescalePool.Begin(ctx)
		if err != nil {
			log.Printf("Failed to start transaction for batch %d-%d: %v", i, end, err)
			continue
		}

		batchUpdated := 0
		hasError := false
		for _, b := range itemsBatch {
			result, err := tx.Exec(ctx,
				`UPDATE prices SET item_name = $1 WHERE item_id = $2`,
				b.Name, b.Id,
			)
			if err != nil {
				log.Printf("Failed to update item_id %d: %v", b.Id, err)
				hasError = true
				break
			}

			rowsAffected := result.RowsAffected()
			batchUpdated += int(rowsAffected)
		}

		if hasError {
			tx.Rollback(ctx)
			log.Printf("❌ Rolled back batch %d-%d due to errors", i, end)
			continue
		}

		// Commit the transaction
		if err := tx.Commit(ctx); err != nil {
			log.Printf("Failed to commit batch %d-%d: %v", i, end, err)
			continue
		}

		totalUpdated += batchUpdated

		log.Printf("Progress: %d/%d items processed | %d price records updated in this batch | Total: %d",
			end, len(items), batchUpdated, totalUpdated)
	}

	log.Printf("Backfill complete! Total price records updated: %d", totalUpdated)
}
