package storage_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	cfg := config.Load()

	_, err := storage.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := storage.WithTempDB(m.Run)
	os.Exit(code)
}

func TestTempDB(t *testing.T) {
	testPool := storage.Pool()

	dbName := testPool.Config().ConnConfig.Database
	if !strings.Contains(dbName, ".go_test") {
		t.Error("unexpected database name")
	}

	err := testPool.Ping(t.Context())
	if err != nil {
		t.Error(err)
	}
}
