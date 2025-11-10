package db

import (
	"context"
	"fmt"
	"time"

	"github.com/bismastr/cs-price-alert/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	PostgresPool  *pgxpool.Pool
	TimescalePool *pgxpool.Pool
}

func NewDbClient(cfg *config.Config) (*Db, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	timescalePool, err := createPool(ctx, cfg.TimescaleDB)
	if err != nil {
		return nil, err
	}

	postgresPool, err := createPool(ctx, cfg.Database)
	if err != nil {
		postgresPool.Close()
		return nil, err
	}

	return &Db{
		PostgresPool:  postgresPool,
		TimescalePool: timescalePool,
	}, nil
}

func createPool(ctx context.Context, dbCfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Database,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: parse config failed: %w", dbCfg.Database, err)
	}

	poolConfig.MaxConns = dbCfg.MaxConns
	poolConfig.MinConns = dbCfg.MinConns
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: create pool failed: %w", dbCfg.Database, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: ping failed: %w", dbCfg.Database, err)
	}

	return pool, nil
}
