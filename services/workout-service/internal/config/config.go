package config

import "os"

type Config struct {
	Port               string
	DBDSN              string
	JWTSecret          string
	ExerciseServiceURL string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8003"),
		DBDSN:              getEnv("DB_DSN", ""),
		JWTSecret:          getEnv("JWT_SECRET", "change_me"),
		ExerciseServiceURL: getEnv("EXERCISE_SERVICE_URL", "http://localhost:8002"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
