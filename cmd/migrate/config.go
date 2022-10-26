package main

import (
	"github.com/caarlos0/env/v6"
)

// Configuration represents application configuration for serve action.
type Configuration struct {
	DatabaseDSN   string `env:"connection_dsn"`
	MigrationsDir string `env:"MIGRATIONS_DIR" envDefault:"migrations/"`
}

// LoadConfiguration returns a new application configuration parsed from environment variables.
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
