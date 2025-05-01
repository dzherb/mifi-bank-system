package storage

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

// WithTempDB sets up a temporary PostgreSQL database for testing.
//
// It connects to the main database using the given config,
// creates a temp database, switches the connection pool
// to use it, runs the tests, and then drops the temp database.
//
// Returns the exit code from testRunner or 1 on setup/cleanup failure.
func WithTempDB(testRunner func() int) int {
	tm := testDBManager{}

	name, err := tm.initTestDB()
	if err != nil {
		log.Error(err)
		return 1
	}

	defer func() {
		err = tm.dropTestDB(name)
		if err != nil {
			log.Error(err)
		}
		// Close the connection to the original database.
		ClosePool()
	}()

	// Run the tests.
	return testRunner()
}

// WithMigratedDB applies all up migrations, runs the provided test runner,
// and then rolls back all migrations.
//
// If any step fails, the error is logged and the function returns exit code 1.
//
// Returns the exit code from the test runner (typically passed to os.Exit).
func WithMigratedDB(testRunner func() int) int {
	m, err := migrator()
	if err != nil {
		log.Error(err)
		return 1
	}

	err = m.Up()
	if err != nil {
		log.Error(err)
		return 1
	}

	defer func(m *migrate.Migrate) {
		err = m.Down()
		if err != nil { // coverage-ignore
			log.Error(err)
		}

		err, err2 := m.Close()
		if err != nil { // coverage-ignore
			log.Error(err)
		}

		if err2 != nil { // coverage-ignore
			log.Error(err2)
		}
	}(m)

	return testRunner()
}

// WithTransaction runs the test function
// within a rolled-back transaction.
// Fails the test immediately if beginning
// or rolling back the transaction fails.
func WithTransaction(
	t *testing.T,
	test func(),
) {
	tx, err := Pool().Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	// Create a savepoint to prevent any function
	// from performing a real commit.
	// This ensures changes are rolled back after the test.
	nestedTx, err := tx.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	prevConn := Conn

	// Mock the Conn function
	// to return the nested transaction.
	Conn = func() Connection {
		return nestedTx
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Fatal(err)
		}

		Conn = prevConn
	}(tx, t.Context())

	test()
}

func migrator() (*migrate.Migrate, error) {
	db := stdlib.OpenDBFromPool(Pool())

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	// Compute path relative to this file (resolves to `.../storage/migrations`)
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("unable to determine current file path")
	}

	migrationsPath := "file://" + filepath.Join(
		filepath.Dir(currentFile),
		"migrations",
	)

	return migrate.NewWithDatabaseInstance(
		migrationsPath,
		Pool().Config().ConnConfig.Database,
		driver,
	)
}

type testDBManager struct {
	originalPool *pgxpool.Pool
}

func (tm *testDBManager) initTestDB() (string, error) {
	name := uuid.NewString() + ".go_test"

	err := createDB(name, "template1")
	if err != nil {
		return "", err
	}

	newCfg := pool.Config()
	newCfg.ConnConfig.Database = name

	tm.originalPool = pool

	pool, err = pgxpool.NewWithConfig(context.Background(), newCfg)

	if err != nil {
		return "", err
	}

	return name, nil
}

func (tm *testDBManager) dropTestDB(name string) error {
	// Close the connection to the temporary database.
	ClosePool()

	pool = tm.originalPool

	err := dropDB(name)
	if err != nil {
		return err
	}

	tm.originalPool = nil

	return nil
}
