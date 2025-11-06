package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleDB struct {
	Pool *pgxpool.Pool
}

func NewTimescaleDB() (*TimescaleDB, error) {
	dbUser := os.Getenv("TIMESCALE_DB_USERNAME")
	dbPassword := os.Getenv("TIMESCALE_DB_PASSWORD")
	dbName := os.Getenv("TIMESCALE_DB_NAME")
	dbHost := os.Getenv("TIMESCALE_DB_HOST")
	dbPort := os.Getenv("TIMESCALE_DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &TimescaleDB{
		Pool: pool,
	}, nil
}
