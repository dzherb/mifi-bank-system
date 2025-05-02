package storage_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/models"
	repo "github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	cfg := config.Load()

	_, err := storage.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := storage.RunTestsWithTempDB(m.Run)
	os.Exit(code)
}

func TestTempDB(t *testing.T) {
	testPool := storage.ActivePool()

	dbName := testPool.Config().ConnConfig.Database
	if !strings.Contains(dbName, ".go_test") {
		t.Error("unexpected database name")
	}

	err := testPool.Ping(t.Context())
	if err != nil {
		t.Error(err)
	}
}

func TestWithTransactionHelperCommits(t *testing.T) {
	storage.TestWithMigratedDB(t)
	storage.TestWithTransaction(t)

	var user models.User

	err := storage.WithTransaction(
		t.Context(),
		storage.Conn(),
		func(tx pgx.Tx) error {
			var err error

			ur := repo.NewUserRepositoryFromTx(tx)
			user, err = ur.Create(models.User{
				Email:    "test@test.com",
				Username: "test",
				Password: "test",
			})

			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	ur := repo.NewUserRepository()

	_, err = ur.Get(user.ID)
	if err != nil {
		t.Errorf("expected to retrieve the user, got error: %v", err)
	}
}

func TestWithTransactionHelperRollbacks(t *testing.T) {
	storage.TestWithMigratedDB(t)
	storage.TestWithTransaction(t)

	var user models.User

	err := storage.WithTransaction(
		t.Context(),
		storage.Conn(),
		func(tx pgx.Tx) error {
			var err error

			ur := repo.NewUserRepositoryFromTx(tx)
			user, err = ur.Create(models.User{
				Email:    "test@test.com",
				Username: "test",
				Password: "test",
			})

			if err != nil {
				return err
			}

			// user with the same email should return an error
			_, err = ur.Create(models.User{
				Email:    "test@test.com",
				Username: "test2",
				Password: "test",
			})

			if err != nil {
				return err
			}

			return nil
		},
	)

	if err == nil {
		t.Error("expected an error")
	}

	ur := repo.NewUserRepository()

	_, err = ur.Get(user.ID)
	if err == nil {
		t.Error("transaction hasn't rollback, expected an error")
	}
}
