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

	code := storage.WithTempDB(func() int {
		return storage.WithMigratedDB(m.Run)
	})
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	storage.WithTransaction(t, func() {
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
	})
}
