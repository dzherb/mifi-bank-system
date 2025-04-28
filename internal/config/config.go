package config

import (
	"errors"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	ServerHost  string
	ServerPort  string
	PostgresURL string
}

func Load() (*Config, error) {
	loadDotenv()

	postgresConn, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		return nil, errors.New("no POSTGRES_URL environment variable set")
	}

	return &Config{
		ServerHost:  envOrDefault("SERVER_HOST", "localhost"),
		ServerPort:  envOrDefault("SERVER_PORT", "8080"),
		PostgresURL: postgresConn,
	}, nil
}

func loadDotenv() {
	tryFiles := []string{".env", "../../.env"}

	for _, file := range tryFiles {
		err := godotenv.Load(file)
		if err != nil {
			continue
		}
		return
	}

	log.Warn("no .env file found")
}

// Simple helper function to read an environment or return a default value
func envOrDefault(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
