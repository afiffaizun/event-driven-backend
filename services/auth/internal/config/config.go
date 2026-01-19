package config

import "github.com/caarlos0/env/v10"

type Config struct {
	AppName string `env:"APP_NAME" envDefault:"auth-service"`
	Port    string `env:"PORT" envDefault:"8080"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// future
	// DBURL string `env:"DATABASE_URL"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
