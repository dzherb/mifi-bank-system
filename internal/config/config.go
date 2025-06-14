package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ServerHost     string
	ServerPort     string
	DatabaseURL    string
	SecretKey      string
	AccessTokenTTL time.Duration

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

func Load() *Config {
	loadDotenv()

	return &Config{
		ServerHost:     envOrDefault("SERVER_HOST", "localhost"),
		ServerPort:     envOrDefault("SERVER_PORT", "8080"),
		DatabaseURL:    envOrPanic("DATABASE_URL"),
		SecretKey:      envOrPanic("SECRET_KEY"),
		AccessTokenTTL: accessTokenTTL(),
		SMTPHost:       envOrPanic("SMTP_HOST"),
		SMTPPort:       envAndConvert("SMTP_PORT", strconv.Atoi),
		SMTPUsername:   envOrPanic("SMTP_USERNAME"),
		SMTPPassword:   envOrPanic("SMTP_PASSWORD"),
	}
}

func loadDotenv() {
	tryFiles := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

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

func envAndConvert[T any](key string, convert func(val string) (T, error)) T {
	val := envOrPanic(key)

	res, err := convert(val)
	if err != nil {
		panic(err)
	}

	return res
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
