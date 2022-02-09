package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Configuration represents application configuration for serve action.
type Configuration struct {
	Database DatabaseConfig
	Minio    MinioConfig
}

type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"storage.ops.boobsrate.com"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"golangbackend"`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:"142701foobar"`
	Bucket    string `env:"MINIO_BUCKET" required:"true"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"true"`
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" required:"false" envDefault:"132.226.58.135"`
	Port     string `env:"DATABASE_PORT" required:"false" envDefault:"5432"`
	User     string `env:"DATABASE_USER" required:"false" envDefault:"golangbackend"`
	Password string `env:"DATABASE_PASSWORD" required:"false" envDefault:"Assimilate142701"`
	Name     string `env:"DATABASE_NAME" required:"false" envDefault:"tits"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.Name)
}

// LoadConfiguration returns a new application configuration parsed from environment variables.
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
