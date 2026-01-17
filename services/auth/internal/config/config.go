package config

import "os"

type Config struct {
	ServiceName string
	Port        string
}

func Load() *Config {
	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "auth-service"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
