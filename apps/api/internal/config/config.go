package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	HTTPAddr    string
	Seed        bool
}

func Load() Config {
	return Config{
		DatabaseURL: getenv("DATABASE_URL", "postgres://pulse:pulse@localhost:5432/running_club?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET", "dev-secret-change-me"),
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		Seed:        strings.EqualFold(getenv("SEED", "0"), "1") || strings.EqualFold(getenv("SEED", "0"), "true"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
