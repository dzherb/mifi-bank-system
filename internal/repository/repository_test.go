package repo_test

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/shopspring/decimal"
)

func TestMain(m *testing.M) {
	cfg := config.Load()

	_, err := storage.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := storage.WithTempDB(func() int {
		return storage.WithMigratedDB(m.Run)
	})
	os.Exit(code)
}

func testUser() models.User {
	return models.User{
		Email:    "test@test.com",
		Username: "test",
		Password: "test_pass",
	}
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	storage.WithTransaction(t, func() {
		now := time.Now().Add(-time.Second * 10)
		ur := repo.NewUserRepository()

		userToCreate := testUser()

		user, err := ur.Create(userToCreate)
		if err != nil {
			t.Fatal(err)
		}

		if user.ID == 0 {
			t.Error("user ID is zero")
		}

		if user.Email != userToCreate.Email {
			t.Errorf(
				"expected email %q, got %q",
				userToCreate.Email,
				user.Email,
			)
		}

		if user.Username != userToCreate.Username {
			t.Errorf(
				"expected username %q, got %q",
				userToCreate.Username,
				user.Username,
			)
		}

		if user.Password == userToCreate.Password {
			t.Errorf("expected password to be hashed")
		}

		if user.CreatedAt.Before(now) {
			t.Errorf("created_at %s is earlier than expected", user.CreatedAt)
		}

		if user.UpdatedAt.Before(now) {
			t.Errorf("updated_at %s is earlier than expected", user.UpdatedAt)
		}
	})
}

func TestUserRepositoryImpl_Authenticate(t *testing.T) {
	storage.WithTransaction(t, func() {
		ur := repo.NewUserRepository()

		created, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		got, err := ur.GetByCredentials(created.Email, "test_pass")
		if err != nil {
			t.Error(err)
			return
		}

		if got != created {
			t.Errorf("expected user to be %v, got %v", created, got)
		}
	})
}

func TestUserRepositoryImpl_Authenticate2(t *testing.T) {
	storage.WithTransaction(t, func() {
		ur := repo.NewUserRepository()

		created, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		_, err = ur.GetByCredentials(created.Email, "wrong_pass")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestUserRepositoryImpl_Get(t *testing.T) {
	storage.WithTransaction(t, func() {
		ur := repo.NewUserRepository()

		created, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		got, err := ur.Get(created.ID)

		if err != nil {
			t.Fatal(err)
		}

		if got != created {
			t.Errorf("expected user %+v, got %+v", created, got)
		}
	})
}

func TestAccountRepositoryImpl_Create(t *testing.T) {
	storage.WithTransaction(t, func() {
		now := time.Now().Add(-time.Second * 10)
		ar := repo.NewAccountRepository()
		ur := repo.NewUserRepository()

		user, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		account, err := ar.Create(models.Account{
			UserID: user.ID,
		})
		if err != nil {
			t.Fatal(err)
		}

		if account.ID == 0 {
			t.Error("account ID is zero")
		}

		if account.UserID != user.ID {
			t.Errorf("expected user ID %d, got %d", user.ID, account.UserID)
		}

		if !account.Balance.IsZero() {
			t.Errorf("expected balance 0, got %q", account.Balance)
		}

		if account.CreatedAt.Before(now) {
			t.Errorf(
				"created_at %s is earlier than expected",
				account.CreatedAt,
			)
		}

		if account.UpdatedAt.Before(now) {
			t.Errorf(
				"updated_at %s is earlier than expected",
				account.UpdatedAt,
			)
		}
	})
}

func TestAccountRepositoryImpl_Get(t *testing.T) {
	storage.WithTransaction(t, func() {
		ar := repo.NewAccountRepository()
		ur := repo.NewUserRepository()

		user, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		created, err := ar.Create(models.Account{
			UserID: user.ID,
		})
		if err != nil {
			t.Fatal(err)
		}

		got, err := ar.Get(created.ID)
		if err != nil {
			return
		}

		if !reflect.DeepEqual(got, created) {
			t.Errorf("expected account %+v, got %+v", created, got)
		}
	})
}

func TestAccountRepositoryImpl_Update(t *testing.T) {
	storage.WithTransaction(t, func() {
		ar := repo.NewAccountRepository()
		ur := repo.NewUserRepository()

		user, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		created, err := ar.Create(models.Account{
			UserID: user.ID,
		})
		if err != nil {
			t.Fatal(err)
		}

		updated := models.Account{
			ID:      created.ID,
			Balance: created.Balance.Add(decimal.NewFromInt(100)),
		}

		updated, err = ar.Update(updated)
		if err != nil {
			t.Fatal(err)
		}

		if !updated.Balance.Equal(decimal.NewFromInt(100)) {
			t.Errorf("expected balance to be 100, got %v", updated.Balance)
		}
	})
}
