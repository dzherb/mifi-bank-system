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

	_, err = storage.InitDP(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := storage.WithTempDB(func() int {
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

	userToCreate := models.User{
		Email:    "test@test.com",
		Username: "test",
		Password: "test_pass",
	}
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
	if user.Password != userToCreate.Password {
		t.Errorf("expected password %q, got %q", userToCreate.Password, user.Password)
	}
	if user.CreatedAt.Before(now) {
		t.Errorf("created_at %s is earlier than expected", user.CreatedAt)
	}
	if user.UpdatedAt.Before(now) {
		t.Errorf("updated_at %s is earlier than expected", user.UpdatedAt)
	}
}

func TestUserRepositoryImpl_Get(t *testing.T) {
	tx, err := storage.Pool().Begin(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback(context.Background())

	ur := NewUserRepository(tx)

	user := models.User{
		Email:    "test@test.com",
		Username: "test",
		Password: "test_pass",
	}
	created, err := ur.Create(user)
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
}
