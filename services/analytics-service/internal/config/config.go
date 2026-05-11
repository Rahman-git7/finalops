package config

import (
	"os"
	"time"
)

type Config struct {
	Port              string
	JWTSecret         string
	WorkoutServiceURL string
	RedisURL          string
	CacheTTL          time.Duration
}

func Load() *Config {
	cacheTTL, _ := time.ParseDuration(getEnv("CACHE_TTL", "5m"))
	return &Config{
		Port:              getEnv("PORT", "8004"),
		JWTSecret:         getEnv("JWT_SECRET", "change_me"),
		WorkoutServiceURL: getEnv("WORKOUT_SERVICE_URL", "http://localhost:8003"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379"),
		CacheTTL:          cacheTTL,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
