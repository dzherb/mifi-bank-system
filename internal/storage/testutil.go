package storage

import (
	"context"
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func WithTempDB(cfg *config.Config, testRunner func() int) int {
	// Establish a connection to the main database.
	_, err := InitDP(cfg)
	if err != nil {
		log.Error(err)
		return 1
	}

	tm := testDBManager{}

	// Create a temporary database and switch Pool() to use it.
	name, err := tm.initTestDB()
	if err != nil {
		log.Error(err)
		return 1
	}

	// Schedule cleanup: drop the temporary database after tests finish.
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
