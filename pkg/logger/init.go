package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
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
