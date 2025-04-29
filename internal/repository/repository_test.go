package repository

import (
	"context"
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/models"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	code := storage.WithTempDB(cfg, func() int {
		return storage.WithMigratedDB(m.Run)
	})
	os.Exit(code)
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	tx, err := storage.Pool().Begin(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback(context.Background())

	now := time.Now().Add(-time.Second * 10)

	ur := NewUserRepository(tx)
	user, err := ur.Create(models.User{
		Email:    "test@test.com",
		Username: "test",
		Password: "test_pass",
	})
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Error("user ID is zero")
	}
	if user.Username != "test" {
		t.Errorf("expected username %q, got %q", user.Username, "test")
	}
	if user.Password != "test_pass" {
		t.Errorf("expected password %q, got %q", user.Password, "test")
	}
	if user.CreatedAt.Before(now) {
		t.Errorf("created_at %s is earlier than expected", user.CreatedAt)
	}
	if user.UpdatedAt.Before(now) {
		t.Errorf("updated_at %s is earlier than expected", user.UpdatedAt)
	}
}
