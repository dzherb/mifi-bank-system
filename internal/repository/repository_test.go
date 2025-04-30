package repository

import (
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	_, err = storage.InitDP(cfg)
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
	storage.WithTransaction(t, func(tx pgx.Tx) {
		now := time.Now().Add(-time.Second * 10)
		ur := NewUserRepository(tx)

		userToCreate := testUser()
		user, err := ur.Create(userToCreate)
		if err != nil {
			t.Fatal(err)
		}

		if user.ID == 0 {
			t.Error("user ID is zero")
		}
		if user.Email != userToCreate.Email {
			t.Errorf("expected email %q, got %q", userToCreate.Email, user.Email)
		}
		if user.Username != userToCreate.Username {
			t.Errorf("expected username %q, got %q", userToCreate.Username, user.Username)
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
	storage.WithTransaction(t, func(tx pgx.Tx) {
		ur := NewUserRepository(tx)
		created, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		got, err := ur.Authenticate(created.Email, "test_pass")
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
	storage.WithTransaction(t, func(tx pgx.Tx) {
		ur := NewUserRepository(tx)
		created, err := ur.Create(testUser())
		if err != nil {
			t.Fatal(err)
		}

		_, err = ur.Authenticate(created.Email, "wrong_pass")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestUserRepositoryImpl_Get(t *testing.T) {
	storage.WithTransaction(t, func(tx pgx.Tx) {
		ur := NewUserRepository(tx)

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
	storage.WithTransaction(t, func(tx pgx.Tx) {
		now := time.Now().Add(-time.Second * 10)
		ar := NewAccountRepository(tx)
		ur := NewUserRepository(tx)

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
			t.Errorf("created_at %s is earlier than expected", account.CreatedAt)
		}
		if account.UpdatedAt.Before(now) {
			t.Errorf("updated_at %s is earlier than expected", account.UpdatedAt)
		}
	})
}

func TestAccountRepositoryImpl_Get(t *testing.T) {
	storage.WithTransaction(t, func(tx pgx.Tx) {
		ar := NewAccountRepository(tx)
		ur := NewUserRepository(tx)
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
