package config

import "github.com/caarlos0/env/v10"

type Config struct {
	AppName string `env:"APP_NAME" envDefault:"auth-service"`
	Port    string `env:"PORT" envDefault:"8080"`

	DBhost string `env:"DB_HOST" envDefault:"localhost"`
	DBPort string `env:"DB_PORT" envDefault:"5432"`
	DBUser string `env:"DB_USER" envDefault:"authuser"`
	DBPass string `env:"DB_PASS" envDefault:"authpass"`
	DBName string `env:"DB_NAME" envDefault:"authdb"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPass + "@" + c.DBhost + ":" + c.DBPort + "/" + c.DBName
}
