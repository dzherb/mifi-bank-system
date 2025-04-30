package storage

import (
	"context"
	"github.com/dzherb/mifi-bank-system/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	cfg := config.Load()

	_, err := Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	code := WithTempDB(m.Run)
	os.Exit(code)
}

func TestTempDB(t *testing.T) {
	testPool := Pool()

	dbName := testPool.Config().ConnConfig.Database
	if !strings.Contains(dbName, ".go_test") {
		t.Error("unexpected database name")
	}

	err := testPool.Ping(context.Background())
	if err != nil {
		t.Error(err)
	}
}
