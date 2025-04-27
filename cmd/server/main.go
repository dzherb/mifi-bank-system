package main

import "github.com/dzherb/mifi-bank-system/pkg/logger"

func main() {
	log := logger.New()
	log.Info("starting server")
}
