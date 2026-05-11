package config

import "os"

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8002"),
		DBDSN:     getEnv("DB_DSN", ""),
		JWTSecret: getEnv("JWT_SECRET", "change_me"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
