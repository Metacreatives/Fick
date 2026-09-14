package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	Address     string
}

func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	address := os.Getenv("FICK_ADDR")

	if address == "" {
		address = "127.0.0.1:3000"
	}

	return Config{
		DatabaseURL: databaseURL,
		Address:     address,
	}, nil
}
