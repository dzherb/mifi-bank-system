package config

import (
	"os"
)

type Config struct {
	ServerHost string
	ServerPort string
}

func Load() (*Config, error) {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		ServerHost: host,
		ServerPort: port,
	}, nil
}
