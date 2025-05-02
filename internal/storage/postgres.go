package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	log "github.com/sirupsen/logrus"
)

// Conn provides access to the database connection.
// It returns a pgxpool.Pool wrapped in the Connection interface.
//
// This function is declared as a variable
// so it can be overridden or mocked in tests.
var Conn = func() Connection {
	return activePool()
}

var pool *pgxpool.Pool

func activePool() *pgxpool.Pool {
	if pool == nil {
		panic("db is not initialized")
	}

	return pool
}

func closePool() {
	if pool != nil {
		pool.Close()
		pool = nil
	}
}

const DefaultStatementTimeout = 10 * time.Second

func Init(cfg *config.Config) (func(), error) {
	var err error

	pgxCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	pgxCfg.ConnConfig.Tracer = logrusTracer(tracelog.LogLevelWarn)
	pgxCfg.AfterConnect = compositeAfterConnect(
		registerTypes,
		setStatementTimeout(DefaultStatementTimeout),
	)

	pool, err = pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return closePool, nil
}

func logrusTracer(level tracelog.LogLevel) *tracelog.TraceLog {
	logger := logrus.NewLogger(log.StandardLogger())

	return &tracelog.TraceLog{
		Logger:   logger,
		LogLevel: level,
	}
}

type afterConnect func(ctx context.Context, conn *pgx.Conn) error

func compositeAfterConnect(funcs ...afterConnect) afterConnect {
	return func(ctx context.Context, conn *pgx.Conn) error {
		for _, f := range funcs {
			if err := f(ctx, conn); err != nil {
				return err
			}
		}

		return nil
	}
}

func registerTypes(_ context.Context, conn *pgx.Conn) error {
	decimal.Register(conn.TypeMap())
	return nil
}

func setStatementTimeout(timeout time.Duration) afterConnect {
	return func(ctx context.Context, conn *pgx.Conn) error {
		t := strconv.FormatFloat(timeout.Seconds(), 'f', -1, 64)
		_, err := conn.Exec(
			ctx,
			"SET statement_timeout = '"+t+"s'",
		)

		return err
	}
}
