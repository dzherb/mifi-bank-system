package storage

import (
	"context"
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitDP(cfg *config.Config) (func(), error) {
	var err error
	pool, err = pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return ClosePool, nil
}

func Pool() *pgxpool.Pool {
	if pool == nil {
		panic("db is not initialized")
	}
	return pool
}

func ClosePool() {
	if pool != nil {
		pool.Close()
		pool = nil
	}
}
