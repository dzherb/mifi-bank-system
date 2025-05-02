package storage

import (
	"context"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	log "github.com/sirupsen/logrus"
)

var pool *pgxpool.Pool

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

var Conn = func() Connection { // can be mocked in tests
	return Pool()
}

func Init(cfg *config.Config) (func(), error) {
	var err error

	pgxCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	pgxCfg.ConnConfig.Tracer = logrusTracer(tracelog.LogLevelWarn)
	pgxCfg.AfterConnect = registerTypes

	pool, err = pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return ClosePool, nil
}

func registerTypes(_ context.Context, conn *pgx.Conn) error {
	decimal.Register(conn.TypeMap())
	return nil
}

func logrusTracer(level tracelog.LogLevel) *tracelog.TraceLog {
	logger := logrus.NewLogger(log.StandardLogger())

	return &tracelog.TraceLog{
		Logger:   logger,
		LogLevel: level,
	}
}
