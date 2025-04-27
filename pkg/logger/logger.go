package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func New() *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)

	levelStr := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetLevel(logrus.InfoLevel)

	return log
}
