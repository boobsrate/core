package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Configuration represents application configuration for serve action.
type Configuration struct {
	Host          string `env:"DATABASE_HOST" required:"true" envDefault:"132.226.58.135"`
	Port          int    `env:"DATABASE_PORT" required:"true" envDefault:"5432"`
	User          string `env:"DATABASE_USER" required:"true" envDefault:"golangbackend"`
	Password      string `env:"DATABASE_PASSWORD" required:"true" envDefault:"Assimilate142701"`
	Name          string `env:"DATABASE_NAME" required:"true" envDefault:"tits"`
	MigrationsDir string `env:"MIGRATIONS_DIR" envDefault:"migrations/"`
}

func (d *Configuration) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.Name)
}

// LoadConfiguration returns a new application configuration parsed from environment variables.
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
