package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL   string
	Address       string
	SecureCookies bool
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

	secureCookies := false

	if value := os.Getenv(
		"FICK_SECURE_COOKIES",
	); value != "" {
		parsed, err := strconv.ParseBool(value)

		if err != nil {
			return Config{},
				fmt.Errorf(
					"parse FICK_SECURE_COOKIES: %w",
					err,
				)
		}

		secureCookies = parsed
	}

	return Config{
		DatabaseURL:   databaseURL,
		Address:       address,
		SecureCookies: secureCookies,
	}, nil
}
