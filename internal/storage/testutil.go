package storage

import (
	"context"
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
)

// WithTempDB sets up a temporary PostgreSQL database for testing.
//
// It connects to the main database using the given config, creates a temp database,
// switches the connection pool to use it, runs the tests, and then drops the temp database.
//
// Returns the exit code from testRunner or 1 on setup/cleanup failure.
func WithTempDB(cfg *config.Config, testRunner func() int) int {
	_, err := InitDP(cfg)
	if err != nil {
		log.Error(err)
		return 1
	}

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

// WithMigratedDB sets up and runs database migrations before executing the provided test runner.
//
// It establishes a database connection using the pgx connection pool, initializes the migration driver,
// applies all up migrations using golang-migrate, and then executes the given test runner.
//
// If any step fails (connection, migration, etc.), it logs the error and returns exit code 1.
//
// Returns the exit code from the test runner (typically passed to os.Exit).
func WithMigratedDB(testRunner func() int) int {
	db := stdlib.OpenDBFromPool(Pool())
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error(err)
		return 1
	}

	defer func(driver database.Driver) {
		err = driver.Close()
		if err != nil {
			log.Error(err)
		}
	}(driver)

	// Compute path relative to this file (always resolves to `.../storage/migrations`)
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Error("unable to determine current file path")
		return 1
	}
	migrationsPath := "file://" + filepath.Join(filepath.Dir(currentFile), "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		Pool().Config().ConnConfig.Database,
		driver,
	)
	if err != nil {
		log.Error(err)
		return 1
	}
	err = m.Up()
	if err != nil {
		log.Error(err)
		return 1
	}
	return testRunner()
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
