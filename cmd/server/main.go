package main

import (
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/server"
	"github.com/dzherb/mifi-bank-system/internal/storage"
	"github.com/dzherb/mifi-bank-system/pkg/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	closeConn, err := storage.InitDP(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	defer closeConn()

	log.Info("connected to postgres")

	err = server.Start(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
}
