package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	HTTPAddr           string
	Seed               bool
	StravaClientID     string
	StravaClientSecret string
	StravaRedirectURL  string
	StravaWebhookToken string
	WebBaseURL         string
}

func Load() Config {
	return Config{
		DatabaseURL:        getenv("DATABASE_URL", "postgres://pulse:pulse@localhost:5432/running_club?sslmode=disable"),
		JWTSecret:          getenv("JWT_SECRET", "dev-secret-change-me"),
		HTTPAddr:           getenv("HTTP_ADDR", ":8080"),
		Seed:               strings.EqualFold(getenv("SEED", "0"), "1") || strings.EqualFold(getenv("SEED", "0"), "true"),
		StravaClientID:     getenv("STRAVA_CLIENT_ID", ""),
		StravaClientSecret: getenv("STRAVA_CLIENT_SECRET", ""),
		StravaRedirectURL:  getenv("STRAVA_REDIRECT_URL", "http://localhost:8080/api/v1/integrations/strava/callback"),
		StravaWebhookToken: getenv("STRAVA_WEBHOOK_TOKEN", "dev-strava-webhook-token"),
		WebBaseURL:         getenv("WEB_BASE_URL", "http://localhost:5173"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
