package auth_test

import (
	"os"
	"testing"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/models"
	repo "github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/security"
	"github.com/dzherb/mifi-bank-system/internal/services/auth"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	cfg := config.Load()

	security.Init(cfg)

	_, err := storage.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := storage.RunTestsWithTempDB(func() int {
		return storage.RunTestsWithMigratedDB(m.Run)
	})
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	storage.TestWithTransaction(t)

	ur := repo.NewUserRepository()
	user, err := ur.Create(models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: "test_pass",
	})

	if err != nil {
		t.Fatal(err)
	}

	as := auth.NewService()

	got, err := as.Login(user.Email, "test_pass")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got.User != user {
		t.Errorf("got user %v, want %v", got.User, user)
		return
	}

	userID, err := security.ValidateToken(got.Token)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got.User.ID != userID {
		t.Errorf("got userID %v, want %v", userID, got.User.ID)
	}
}

func TestRegister(t *testing.T) {
	storage.TestWithTransaction(t)

	as := auth.NewService()

	res, err := as.Register("test@test.com", "user", "strongPass11")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	ur := repo.NewUserRepository()

	userFromDB, err := ur.Get(res.User.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if res.User != userFromDB {
		t.Errorf("got user %v, want %v", userFromDB, res.User)
	}

	userID, err := security.ValidateToken(res.Token)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if res.User.ID != userID {
		t.Errorf("got userID %v, want %v", userID, res.User.ID)
	}
}
