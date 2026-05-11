package config

import (
	"os"
	"time"
)

type Config struct {
	Port           string
	DBDSN          string
	JWTSecret      string
	JWTAccessTTL   time.Duration
	JWTRefreshTTL  time.Duration
}

func Load() *Config {
	accessTTL, _ := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	refreshTTL, _ := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))
	return &Config{
		Port:          getEnv("PORT", "8001"),
		DBDSN:         getEnv("DB_DSN", ""),
		JWTSecret:     getEnv("JWT_SECRET", "change_me"),
		JWTAccessTTL:  accessTTL,
		JWTRefreshTTL: refreshTTL,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
