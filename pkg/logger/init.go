package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetOutput(os.Stdout)

	levelStr := os.Getenv("LOG_LEVEL")

	level, err := log.ParseLevel(levelStr)
	if err != nil {
		level = log.InfoLevel
	}

	log.SetLevel(level)

	log.SetFormatter(&log.JSONFormatter{})
}
