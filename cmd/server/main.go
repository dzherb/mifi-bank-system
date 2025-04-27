package main

import (
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/internal/server"
	"github.com/dzherb/mifi-bank-system/pkg/logger"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load the config: %v", err)
		return
	}

	err = server.Start(cfg)
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}
