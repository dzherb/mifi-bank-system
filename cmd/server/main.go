package main

import (
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/security"
	"github.com/dzherb/mifi-bank-system/internal/server"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/dzherb/mifi-bank-system/pkg/logger"
	log "github.com/sirupsen/logrus"
)

func main() { // coverage-ignore
	logger.Init()

	cfg := config.Load()

	security.Init(cfg)

	closeConn, err := storage.Init(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	defer closeConn()

	err = server.Start(cfg)
	if err != nil {
		log.Error(err)
	}
}
