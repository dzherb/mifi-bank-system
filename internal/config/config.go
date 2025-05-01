package config

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerHost     string
	ServerPort     string
	DatabaseURL    string
	SecretKey      string
	AccessTokenTTL time.Duration
}

func Load() *Config {
	loadDotenv()

	return &Config{
		ServerHost:     envOrDefault("SERVER_HOST", "localhost"),
		ServerPort:     envOrDefault("SERVER_PORT", "8080"),
		DatabaseURL:    envOrPanic("DATABASE_URL"),
		SecretKey:      envOrPanic("SECRET_KEY"),
		AccessTokenTTL: accessTokenTTL(),
	}
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

func envOrDefault(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func envOrPanic(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("environment variable %s not set", key))
}

func accessTokenTTL() time.Duration {
	const defaultTTL = time.Minute * 60

	ttl, ok := os.LookupEnv("ACCESS_TOKEN_TTL_IN_MINUTES")
	if !ok {
		return defaultTTL
	}

	tokenTTL, err := strconv.Atoi(ttl)
	if err != nil {
		return defaultTTL
	}

	return time.Duration(tokenTTL) * time.Minute
}
