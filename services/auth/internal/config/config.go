package config

import (
	"os"
)

type Config struct {
	AppName   string
	Port      string
	DBURL     string
	JWTSecret string
}

func Load() (*Config, error) {
	return &Config{
		AppName:   getEnv("APP_NAME", "auth-service"),
		Port:      getEnv("PORT", "8080"),
		DBURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/auth_db?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "supersecret"),
	}, nil
}

func (c *Config) DatabaseURL() string {
	return c.DBURL
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
